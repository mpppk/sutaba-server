package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/sutaba-server/pkg/domain/model"
	"github.com/mpppk/sutaba-server/pkg/util"
)

type TweetType int

const (
	Tweet TweetType = iota + 1
	Reply
	QuoteTweet
	ReplyWithQuote
)

type User struct {
	Client *anaconda.TwitterApi
	ID     int64
}

func NewUser(
	accessToken, accessTokenSecret, consumerKey, consumerSecret string,
	id int64, keyword string, isErrorReporter bool, tweetType TweetType) *User {
	return &User{
		Client: anaconda.NewTwitterApiWithCredentials(
			accessToken,
			accessTokenSecret,
			consumerKey,
			consumerSecret,
		),
		ID: id,
	}
}

func (u *User) PostQuoteTweet(text string, quotedTweetIDStr, quotedTweetUserID string) (*model.Tweet, error) {
	return PostQuoteTweet(u.Client, text, quotedTweetIDStr, quotedTweetUserID)
}

func (u *User) PostReply(text, toTweetIDStr string, toScreenNames []string) (*model.Tweet, error) {
	return PostReply(u.Client, text, toTweetIDStr, toScreenNames)
}

func (u *User) PostReplyWithQuote(text string, quotedTweetIDStr, quotedTweetUserScreenName, toTweetIDStr string, toScreenNames []string) (*model.Tweet, error) {
	return PostReplyWithQuote(u.Client, text, quotedTweetIDStr, quotedTweetUserScreenName, toTweetIDStr, toScreenNames)
}

func (u *User) PostErrorTweet(notifyText, sorryText, toSorryTweetIDStr, toSorryUserScreenName string) {
	if _, err := u.Client.PostTweet(notifyText, nil); err != nil {
		util.LogPrintlnInOneLine("failed to tweet error notify message", err)
	}
	if _, err := u.PostReply(sorryText, toSorryTweetIDStr, []string{toSorryUserScreenName}); err != nil {
		util.LogPrintlnInOneLine("failed to tweet sorry message", err)
	}
}

func toUser(anacondaUser *anaconda.User) *model.TwitterUser {
	return &model.TwitterUser{
		ID:         anacondaUser.Id,
		ScreenName: anacondaUser.ScreenName,
	}
}
