package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type View struct {
	Client *anaconda.TwitterApi
}

func (v *View) Show(text string) error {
	_, err := v.Client.PostTweet(text, nil)
	return err
}

func (v *View) ReplyToTweet(text string, toTweetID model.MessageID) error {
	_, err := PostReply(v.Client, text, toTweetID)
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
