package model

import (
	"strconv"
	"strings"
)

type Message struct {
	ID                int64
	User              User
	Text              string
	ReferencedMessage *Message
	MediaNum          int
}

func (m *Message) GetIDStr() string {
	return strconv.FormatInt(m.ID, 10)
}

func (m *Message) HasMessageReference() bool {
	return m.ReferencedMessage != nil
}

func (m *Message) HasKeyWord(word string) bool {
	return strings.Contains(m.Text, word)
}
