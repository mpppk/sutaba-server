package model

import (
	"strconv"
)

type Message struct {
	ID                  int64
	User                User
	Text                string
	QuoteMessage        *Message
	InReplyToStatusID   int64
	InReplyToUserID     int64
	InReplyToScreenName string
	MediaURLs           []string
}

func (t *Message) GetIDStr() string {
	return strconv.FormatInt(t.ID, 10)
}

func (t *Message) IsReply() bool {
	return t.InReplyToUserID != 0
}

func (t *Message) HasQuoteTweet() bool {
	return t.QuoteMessage != nil
}
