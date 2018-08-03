package models

type Message struct {
	Id        string `json:"_id,omitempty"`
	ChannelId string `json:"rid"`
	Text      string `json:"msg"`
	Timestamp string `json:"ts,omitempty"`
	User      User   `json:"u,omitempty"`

	Alias     string `json:"alias,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Emoji     string `json:"emoji,omitempty"`
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