package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type TweetType int

const (
	Tweet TweetType = iota + 1
	Reply
	QuoteTweet
	ReplyWithQuote
)

func toUser(anacondaUser *anaconda.User) *model.TwitterUser {
	return &model.TwitterUser{
		ID:         anacondaUser.Id,
		ScreenName: anacondaUser.ScreenName,
	}
}
