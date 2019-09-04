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
	var medias []*model.MessageMedia
	for _, mediaURL := range tweet.MediaURLs {
		medias = append(medias, model.NewMessageMedia(mediaURL))
	}

	message := &model.Message{
		ID:     tweet.ID,
		User:   tweet.User,
		Text:   tweet.Text,
		Medias: medias,
	}

	if tweet.InReplyToUserID != 0 {
		message.RepliedUser = &model.User{
			ID:   tweet.InReplyToUserID,
			Name: tweet.InReplyToScreenName,
		}
	}

	if tweet.InReplyToStatusID != 0 {
		message.RepliedMessageID = tweet.InReplyToStatusID
	}

	if tweet.QuoteTweet != nil {
		message.ReferencedMessage = r.NewMessage(tweet.QuoteTweet)
	}

	key := strconv.FormatInt(message.ID, 10) // FIXME consider other SNS like slack, mastdon, etc...
	r.tweetMap[key] = tweet
	return message
}
