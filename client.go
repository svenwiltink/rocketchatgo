package rocketchatgo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/svenwiltink/ddpgo"
	"net/url"
	"log"
)

type Session struct {
	ddp   *ddpgo.Client
	state *State

	messageChan chan *Message
}

func (s *Session) Close() {
	//s.ddp.Close()
}

func (s *Session) Login(username string, email string, password string) error {
	// clean the state
	s.state = NewState()
	digest := sha256.Sum256([]byte(password))

	loginResult, err := s.ddp.CallMethod("login", ddpLoginRequest{
		User:     ddpUser{Email: email, Username: username},
		Password: ddpPassword{Digest: hex.EncodeToString(digest[:]), Algorithm: "sha-256"}})

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
	channels, err := s.GetChannels()
	if err != nil {
		return err
	}

	s.state.AddRooms(channels...)
	s.startEventListener()
	return nil
}

func (s *Session) GetChannels() ([]*Room, error) {
	result, err := s.ddp.CallMethod("rooms/get")
	if err != nil {
		return nil, err
	}

	jsonArray, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	channelArray := make([]*Room, 0)
	err = json.Unmarshal(jsonArray, &channelArray)

	return channelArray, err
}

func (s *Session) startEventListener() {
	for _, room := range s.state.channelMap {
		_, err := s.ddp.Subscribe("stream-room-messages", room.ID, true)
		if err != nil {
			panic(err)
		}
	}

	s.ddp.GetCollectionByName("stream-room-messages").AddChangedEventHandler(s.OnRoomMessage)

	for message := range s.messageChan {
		if message.Sender.ID == s.state.UserID {
			continue
		}

		s.SendMessage(message.ChannelID, message.Message)
	}
}

func (s Session) SendMessage(channelID string, message string) error {
	_, err := s.ddp.CallMethod("sendMessage", struct {
		Message   string `json:"msg"`
		ChannelID string `json:"rid"`
	}{
		Message:   message,
		ChannelID: channelID,
	})

	return err
}

func (s Session) OnRoomMessage(event ddpgo.CollectionChangedEvent) {
	jsonBytes, err := json.Marshal(event.Fields.Args)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(jsonBytes))
	messageList := make([]*Message, 0)
	err = json.Unmarshal(jsonBytes, &messageList)
	if err != nil {
		log.Println(err)
		return
	}

	for _, message := range messageList {
		s.messageChan <- message
	}
}

func NewClient(host string) (*Session, error) {

	ddpClient := ddpgo.NewClient(url.URL{Host: host})

	if err := ddpClient.Connect(); err != nil {
		return nil, err
	}

	return &Session{
		ddp:         ddpClient,
		messageChan: make(chan *Message, 100),
	}, nil
}
