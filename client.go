package rocketchatgo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/svenwiltink/ddpgo"
	"log"
	"net/url"
	"strings"
	"sync"
)

type Session struct {
	ddp   *ddpgo.Client
	state *State

	eventHandlers     map[string][]EventHandler
	eventHandlerMutex sync.RWMutex
}

func (s *Session) Close() {
	s.ddp.Close()
}

// NotifyRoom notifies the room of a given event
func (s *Session) NotifyRoom(roomID string, params ...interface{}) (err error) {
	_, err = s.ddp.CallMethod("stream-notify-room", params...)
	return
}

// IsTyping enables/disables the '... is typing' event for the given room and user
func (s *Session) IsTyping(roomID, username string, flag bool) (err error) {
	err = s.NotifyRoom(roomID+"/typing", username, flag)
	return
}

func (s *Session) AddHandler(i interface{}) error {
	handler, err := getHandlerFromInterface(i)
	if err != nil {
		return err
	}

	s.eventHandlerMutex.Lock()
	defer s.eventHandlerMutex.Unlock()
	handlerType := handler.Type()

	slice, exists := s.eventHandlers[handlerType]

	if !exists {
		slice = make([]EventHandler, 0)
	}

	s.eventHandlers[handlerType] = append(slice, handler)
	return nil
}

func (s *Session) Login(username string, email string, password string) error {
	// clean the state
	s.state = NewState()
	digest := sha256.Sum256([]byte(password))

	loginResult, err := s.ddp.Login(ddpgo.Credentials{
		User: ddpgo.User{
			Username: username,
			Email: email,
		},
		Password: ddpgo.Password{
			Digest: hex.EncodeToString(digest[:]),
			Algorithm: "sha-256",
		},
	})

	if err != nil {
		return err
	}

	jsonString, err := json.Marshal(loginResult)
	if err != nil {
		return err
	}

	loginResultStruct := &ddpLoginResponse{}
	err = json.Unmarshal(jsonString, loginResultStruct)
	if err != nil {
		return err
	}

	s.state.UserID = loginResultStruct.ID
	s.updateChannels()
	s.startEventListener()
	return nil
}

func (s *Session) SendMessage(channelId string, message string) error {
	err := s.SendCustomMessage(Message{
		Message:   message,
		ChannelID: channelId,
	})

	return err
}

// SendCustomMessage sends a customizable message using the full REST-model (which strangely also applies to realtime)
func (s *Session) SendCustomMessage(message Message) error {
	_, err := s.ddp.CallMethod("sendMessage", message)
	return err
}

func (s *Session) GetChannels() ([]*Room, error) {
	result, err := s.ddp.CallMethod("rooms/get", map[string]int{
		"$date": 0,
	})
	if err != nil {
		return nil, err
	}

	jsonArray, err := json.Marshal(result.(map[string]interface{})["update"])
	if err != nil {
		return nil, err
	}

	channelArray := make([]*Room, 0)
	err = json.Unmarshal(jsonArray, &channelArray)

	return channelArray, err
}

func (s *Session) GetChannelByName(name string) *Room {
	room := s.state.GetChannelByName(name)

	if room != nil {
		return room
	}

	s.updateChannels()

	return s.state.GetChannelByName(name)
}

func (s *Session) GetChannelById(channelId string) *Room {
	room := s.state.GetChannelById(channelId)

	if room != nil {
		return room
	}

	s.updateChannels()

	return s.state.GetChannelById(channelId)
}

func (s *Session) updateChannels() {
	channels, _ := s.GetChannels()
	s.state.SetRooms(channels...)
}

func (s *Session) SetChannels(channels []*Room) {
	s.state.SetRooms(channels...)
}

func (s *Session) startEventListener() {
	if len(s.state.channelMap) == 0 {
		panic("THIS WILL SEGFAULT, MAKE SURE THE BOT IS IN SOME CHANNELS!")
	}

	for _, room := range s.state.channelMap {
		_, err := s.ddp.Subscribe("stream-room-messages", room.ID, true)
		if err != nil {
			panic(err)
		}
	}

	s.ddp.Subscribe("stream-notify-user", s.state.UserID+"/rooms-changed", true)

	s.ddp.GetCollectionByName("stream-notify-user").AddChangedEventHandler(s.onUserChange)
	s.ddp.GetCollectionByName("stream-room-messages").AddChangedEventHandler(s.onRoomMessage)
}

// pass the event to all eventhandlers
func (s *Session) handleEvent(event Event) {
	eventType := event.GetType()
	s.eventHandlerMutex.RLock()
	defer s.eventHandlerMutex.RUnlock()

	handlers, exists := s.eventHandlers[eventType]
	if !exists {
		return
	}

	for _, handler := range handlers {
		handler.Handle(s, event)
	}
}

func (s *Session) onUserChange(event ddpgo.CollectionChangedEvent) {
	eventType := strings.Split(event.Fields.EventName, "/")[1]
	switch eventType {
	case "rooms-changed":
		{
			channel := &Room{}
			jsonBytes, err := json.Marshal(event.Fields.Args[1])
			if err != nil {
				log.Println(err)
				return
			}

			err = json.Unmarshal(jsonBytes, channel)
			if err != nil {
				log.Println(err)
				return
			}

			channelEvent := event.Fields.Args[0].(string)

			switch channelEvent {
			case "inserted":
				{
					event := &ChannelJoinEvent{
						Channel: channel,
					}

					s.state.AddRoom(channel)
					s.handleEvent(event)

				}
			case "removed":
				{
					event := &ChannelLeaveEvent{
						Channel: channel,
					}

					s.state.RemoveRoom(channel)
					s.handleEvent(event)
				}
			}
		}
	}
}

func (s *Session) onRoomMessage(event ddpgo.CollectionChangedEvent) {
	jsonBytes, err := json.Marshal(event.Fields.Args)
	if err != nil {
		log.Println(err)
		return
	}

	messageList := make([]*Message, 0)
	err = json.Unmarshal(jsonBytes, &messageList)
	if err != nil {
		log.Println(err)
		return
	}

	for _, message := range messageList {
		log.Println(message)
		event := MessageCreateEvent{Message: message}
		s.handleEvent(&event)
	}
}

func (s *Session) GetUserID() string {
	return s.state.UserID
}

func NewClient(host string, ssl bool) (*Session, error) {

	scheme := "wss"
	if !ssl {
		scheme = "ws"
	}

	ddpClient := ddpgo.NewClient(url.URL{Host: host, Scheme: scheme})

	if err := ddpClient.Connect(); err != nil {
		return nil, err
	}

	return &Session{
		ddp:           ddpClient,
		eventHandlers: make(map[string][]EventHandler),
	}, nil
}
