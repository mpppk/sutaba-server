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
	ID                  int64
	User                User
	Text                string
	QuoteMessage        *Message
	InReplyToStatusID   int64
	InReplyToUserID     int64
	InReplyToScreenName string
	Medias              []*MessageMedia
}

func (m *Message) GetIDStr() string {
	return strconv.FormatInt(m.ID, 10)
}

func (m *Message) IsReply() bool {
	return m.InReplyToUserID != 0
}

func (m *Message) HasMessageRefference() bool {
	return m.QuoteMessage != nil
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
