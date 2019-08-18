package twitter

import (
	"log"

	"github.com/ChimeraCoder/anaconda"
)

type TweetType int

const (
	Tweet TweetType = iota + 1
	Reply
	QuoteTweet
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

func (u *User) LogAndPostErrorTweet(text string, err error) {
	log.Println(err)
	if _, err := u.Client.PostTweet(text, nil); err != nil {
		log.Println("failed to tweet error message", err)
	}
}
