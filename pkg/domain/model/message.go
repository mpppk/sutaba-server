package model

import (
	"strconv"
	"strings"
)

type MessageMedia struct {
	url string
}

func NewMessageMedia(url string) *MessageMedia {
	return &MessageMedia{
		url: url,
	}
}

func (m *MessageMedia) GetURL() string {
	return m.url
}

type Message struct {
	ID                int64
	User              User
	Text              string
	ReferencedMessage *Message
	RepliedMessageID  int64
	RepliedUser       *User
	Medias            []*MessageMedia
}

func (m *Message) GetIDStr() string {
	return strconv.FormatInt(m.ID, 10)
}

func (m *Message) IsReply() bool {
	return m.RepliedMessageID != 0 || m.RepliedUser != nil
}

func (m *Message) HasMessageReference() bool {
	return m.ReferencedMessage != nil
}

func (m *Message) GetFirstMedia() (*MessageMedia, bool) {
	if len(m.Medias) == 0 {
		return nil, false
	}

	return m.Medias[0], true
}

func (m *Message) HasKeyWord(word string) bool {
	return strings.Contains(m.Text, word)
}
