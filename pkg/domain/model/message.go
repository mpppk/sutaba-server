package model

import (
	"strconv"
	"strings"
)

type MessageID int64

func (m MessageID) ToString() string {
	return strconv.FormatInt(int64(m), 10)
}

type MessageText string

func (m MessageText) HasKeyword(keyword string) bool {
	return strings.Contains(string(m), keyword)
}

type Message struct {
	ID                MessageID
	User              User
	Text              MessageText
	ReferencedMessage *Message
	MediaNum          int
}

func (m *Message) GetIDStr() string {
	return m.ID.ToString()
}

func (m *Message) HasMessageReference() bool {
	return m.ReferencedMessage != nil
}

func (m *Message) HasKeyWord(keyword string) bool {
	return m.Text.HasKeyword(keyword)
}
