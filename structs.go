package rocketchatgo

import (
	"strconv"
	"time"
)

type ddpLoginResponse struct {
	ID string `json:"id"`
}

type Room struct {
	ID           string   `json:"_id"`
	Type         RoomType `json:"t"`
	CreationDate JsonDate `json:"ts"`
	Name         string   `json:"name"`
	LastMessage  JsonDate `json:"lm"`
	MessageCount int      `json:"msg"`
	CanLeave     bool     `json:"cl"`
	ReadOnly     bool     `json:"ro"`
	Usernames    []string `json:"usernames"`
	Owner        *User    `json:"u"`
}

type RoomType string

const (
	RoomTypeChannel RoomType = "c"
	RoomTypeDirect  RoomType = "d"
)

type User struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type Message struct {
	ID           string   `json:"_id"`
	Type         string   `json:"t,omitempty"`
	CreationDate JsonDate `json:"ts,omitempty"`
	Message      string   `json:"msg"`
	Url          []string `json:"url,omitempty"`
	ExpireAt     JsonDate `json:"url,omitempty"`
	Mentions     []*User  `json:"mentions,omitempty"`
	Sender       *User    `json:"u,omitempty"`
	ChannelID    string   `json:"rid,omitempty"`
}

type JsonDate struct {
	Time jsonTime `json:"$date"`
}

type jsonTime time.Time

func (t jsonTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *jsonTime) UnmarshalJSON(s []byte) (err error) {
	r := string(s)

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}

	*(*time.Time)(t) = time.Unix(0, q*int64(time.Millisecond)).UTC()
	return
}
