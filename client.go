package rocketchatgo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gopackage/ddp"
	"math/rand"
	"time"
)

type Session struct {
	ddp   *ddp.Client
	state *State
}

func (s *Session) Close() {
	s.ddp.Close()
}

func (s *Session) Login(username string, email string, password string) error {
	// clean the state
	s.state = NewState()
	digest := sha256.Sum256([]byte(password))

	loginResult, err := s.ddp.Call("login", ddpLoginRequest{
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
	result, err := s.ddp.Call("rooms/get")
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
		err := s.ddp.Sub("stream-room-messages", room.ID, true)
		if err != nil {
			panic(err)
		}
	}

	channel := make(chan *Message, 100)
	s.ddp.CollectionByName("stream-room-messages").AddUpdateListener(&messageExtractor{"update", channel})

	for message := range channel {
		room := s.state.GetChannelById(message.ChannelID)
		fmt.Printf("%s: %s - %s", room.Name, message.Sender.Username, message.Message)
	}
}

func NewClient(host string, port int, tls bool) (*Session, error) {

	rand.Seed(time.Now().UTC().UnixNano())
	protocol := "ws"
	if tls {
		protocol = "wss"
	}

	hostString := fmt.Sprintf("%s://%s:%d/websocket", protocol, host, port)
	ddpClient := ddp.NewClient(hostString, "http://"+host)
	ddpClient.SetSocketLogActive(true)

	if err := ddpClient.Connect(); err != nil {
		return nil, err
	}

	return &Session{
		ddp: ddpClient,
	}, nil
}

type messageExtractor struct {
	operation string
	channel   chan *Message
}

func (u messageExtractor) CollectionUpdate(collection, operation, id string, doc ddp.Update) {
	if u.operation == operation {
		jsonBytes, err := json.Marshal(doc["args"].([]interface{})[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(jsonBytes))
		message := &Message{}
		err = json.Unmarshal(jsonBytes, message)
		if err != nil {
			fmt.Println(err)
			return
		}

		u.channel <- message
	}
}
