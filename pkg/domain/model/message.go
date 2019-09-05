package model

import (
	"strconv"
	"strings"
)

type MessageID int64
type MessageText string

type Message struct {
	ID                MessageID
	User              User
	Text              MessageText
	ReferencedMessage *Message
	MediaNum          int
}

func (m *Message) GetIDStr() string {
	return strconv.FormatInt(int64(m.ID), 10)
}

func (m *Message) HasMessageReference() bool {
	return m.ReferencedMessage != nil
}

func (m *Message) HasKeyWord(word string) bool {
	return strings.Contains(string(m.Text), word)
}
