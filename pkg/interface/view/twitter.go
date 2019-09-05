package view

import "github.com/mpppk/sutaba-server/pkg/domain/model"

type TwitterView interface {
	Show(text string) error
	ReplyToTweet(text string, toTweetID model.MessageID) error
}
