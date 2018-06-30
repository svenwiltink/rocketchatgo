package rocketchatgo

import "sync"

type State struct {
	UserID    string
	UserMap   map[string]*User
	userMutex sync.RWMutex

	channelMap   map[string]*Room
	channelMutex sync.RWMutex
}

func NewState() *State {
	return &State{
		UserMap:    make(map[string]*User),
		channelMap: make(map[string]*Room),
	}
}

func (s *State) AddRooms(rooms ...*Room) {
	s.channelMutex.Lock()
	defer s.channelMutex.Unlock()

	for _, room := range rooms {
		s.channelMap[room.ID] = room
	}
}

func (s *State) SetRooms(rooms ...*Room) {
	s.channelMutex.Lock()
	defer s.channelMutex.Unlock()

	s.channelMap = make(map[string]*Room)

	for _, room := range rooms {
		s.channelMap[room.ID] = room
	}
}

func (s *State) AddRoom(room *Room) {
	s.channelMutex.Lock()
	defer s.channelMutex.Unlock()

	s.channelMap[room.ID] = room
}

func (s *State) GetChannelById(channelId string) *Room {
	s.channelMutex.RLock()
	defer s.channelMutex.RUnlock()

	if channel, ok := s.channelMap[channelId]; ok {
		return channel
	}

	return nil
}

func (s *State) GetChannelByName(name string) *Room {
	s.channelMutex.RLock()
	defer s.channelMutex.RUnlock()

	for _, room := range s.channelMap {
		if room.Name == name {
			return room
		}
	}

	return nil
}

func (s *State) GetUserById(userId string) *User {
	s.userMutex.RLock()
	defer s.userMutex.RUnlock()

	if user, ok := s.UserMap[userId]; ok {
		return user
	}

	return nil
}

func (s *State) GetUserByName(name string) *User {
	s.userMutex.RLock()
	defer s.userMutex.RUnlock()

	for _, user := range s.UserMap {
		if user.Username == name {
			return user
		}
	}

	return nil
}
