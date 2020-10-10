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

var debugKeyword = "--debug"

func (m MessageText) HasKeyword(keyword string) bool {
	return strings.Contains(string(m), keyword)
}

type MessageEvent struct {
	TargetUserID UserID
	IsShared     bool
	Message      *Message
}

type Message struct {
	ID                MessageID
	User              User
	Text              MessageText
	ReferencedMessage *Message
	MediaNum          int
	ReplyUserID       UserID
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

func (m *Message) IsRepliedTo(user *User) bool {
	return m.ReplyUserID == user.ID || strings.Contains(string(m.Text), "@"+string(user.Name))
}

func (m *Message) IsDebugMode() bool {
	return m.Text.HasKeyword(debugKeyword)
}
