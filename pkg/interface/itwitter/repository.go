package itwitter

import (
	"strconv"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type Twitter struct {
	tweetMap map[string]*Tweet
}

func NewTwitter() *Twitter {
	return &Twitter{
		tweetMap: map[string]*Tweet{},
	}
}

func (r *Twitter) NewMessage(tweet *Tweet) *model.Message {
	message := &model.Message{
		ID:                  tweet.ID,
		User:                tweet.User,
		Text:                tweet.Text,
		InReplyToStatusID:   tweet.InReplyToStatusID,
		InReplyToUserID:     tweet.InReplyToUserID,
		InReplyToScreenName: tweet.InReplyToScreenName,
		MediaURLs:           tweet.MediaURLs,
	}

	if tweet.QuoteTweet != nil {
		message.QuoteMessage = r.NewMessage(tweet.QuoteTweet)
	}

	key := strconv.FormatInt(message.ID, 10) // FIXME consider other SNS like slack, mastdon, etc...
	r.tweetMap[key] = tweet
	return message
}
