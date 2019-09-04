package twitter

import (
	"github.com/ChimeraCoder/anaconda"
)

type View struct {
	Client *anaconda.TwitterApi
}

func (v *View) Show(text string) error {
	_, err := v.Client.PostTweet(text, nil)
	return err
}

func (v *View) ReplyToTweet(text string, toTweetIDStr string) error {
	_, err := PostReply(v.Client, text, toTweetIDStr)
	return err
}

func NewView(accessToken, accessTokenSecret, consumerKey, consumerSecret string) *View {
	return &View{
		Client: anaconda.NewTwitterApiWithCredentials(
			accessToken,
			accessTokenSecret,
			consumerKey,
			consumerSecret,
		),
	}
}
