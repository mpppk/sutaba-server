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

func toUser(anacondaUser *anaconda.User) *model.User {
	return &model.User{
		ID:   model.UserID(anacondaUser.Id),
		Name: model.UserName(anacondaUser.ScreenName),
	}
}
