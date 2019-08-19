package twitter

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
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
	Client          *anaconda.TwitterApi
	ID              int64
	TargetKeyword   string
	IsErrorReporter bool
	TweetType       TweetType
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
		ID:              id,
		TargetKeyword:   keyword,
		IsErrorReporter: isErrorReporter,
		TweetType:       tweetType,
	}
}

func (u *User) PostByTweetType(text string, tweet *anaconda.Tweet) (anaconda.Tweet, error) {
	switch u.TweetType {
	case Tweet:
		return u.Client.PostTweet(text, nil)
	case Reply:
		return u.PostReply(text, tweet.User.ScreenName, tweet.IdStr)
	case QuoteTweet:
		return u.PostQuoteTweet(text, tweet)
	case ReplyWithQuote:
		return u.PostReplyWithQuote(text, tweet, tweet.User.ScreenName, tweet.IdStr)
	}
	return anaconda.Tweet{}, fmt.Errorf("unknown TweetType: %v", u.TweetType)
}

func (u *User) PostQuoteTweet(text string, quoteTweet *anaconda.Tweet) (anaconda.Tweet, error) {
	return PostQuoteTweet(u.Client, text, quoteTweet)
}

func (u *User) PostReply(text, toScreenName, toTweetIDStr string) (anaconda.Tweet, error) {
	return PostReply(u.Client, text, toScreenName, toTweetIDStr)
}

func (u *User) PostReplyWithQuote(text string, quoteTweet *anaconda.Tweet, toScreenName, toTweetIDStr string) (anaconda.Tweet, error) {
	return PostReplyWithQuote(u.Client, text, quoteTweet, toScreenName, toTweetIDStr)
}

func (u *User) PostErrorTweet(notifyText, sorryText, toSorryUserScreenName, toSorryTweetIDStr string) {
	if _, err := u.Client.PostTweet(notifyText, nil); err != nil {
		util.LogPrintlnInOneLine("failed to tweet error notify message", err)
	}
	if _, err := u.PostReply(sorryText, toSorryUserScreenName, toSorryTweetIDStr); err != nil {
		util.LogPrintlnInOneLine("failed to tweet sorry message", err)
	}
}
