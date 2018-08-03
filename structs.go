package rocketchatgo

import (
	"strconv"
	"time"
)

type ddpLoginRequest struct {
	User     ddpUser     `json:"user"`
	Password ddpPassword `json:"password"`
}

type ddpUser struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

type ddpPassword struct {
	Digest    string `json:"digest"`
	Algorithm string `json:"algorithm"`
}

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
	ID           string `json:"_id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Token        string `json:"token,omitempty"`
	TokenExpires int64  `json:"tokenExpires,omitempty"`
}

type Message struct {
	ChannelID string `json:"rid"`
	Message   string `json:"msg"`

	// Optional fields, or for receiver only
	ID          string       `json:"_id,omitempty"`
	Timestamp   string       `json:"ts,omitempty"`
	Sender      User         `json:"u,omitempty"`
	Alias       string       `json:"alias,omitempty"`
	Avatar      string       `json:"avatar,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Emoji       string       `json:"emoji,omitempty"`
}

type Attachment struct {
	Color         string            `json:"color,omitempty"`
	Text          string            `json:"text,omitempty"`
	Timestamp     string            `json:"ts,omitempty"`
	ThumbnailURL  string            `json:"thumb_url,omitempty"`
	MessageURL    string            `json:"message_link,omitempty"`
	Collapsed     bool              `json:"collapsed,omitempty"`
	AuthorName    string            `json:"author_name,omitempty"`
	AuthorURL     string            `json:"author_link,omitempty"`
	AuthorIcon    string            `json:"author_icon,omitempty"`
	Title         string            `json:"title,omitempty"`
	TitleURL      string            `json:"title_link,omitempty"`
	TitleDownload bool              `json:"title_link_download,omitempty"`
	ImageURL      string            `json:"image_url,omitempty"`
	AudioURL      string            `json:"audio_url,omitempty"`
	VideoURL      string            `json:"video_url,omitempty"`
	Fields        []AttachmentField `json:"fields,omitempty"`
}

type AttachmentField struct {
	Short bool   `json:"short,omitempty"`
	Title string `json:"title"`
	Value string `json:"value"`
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
