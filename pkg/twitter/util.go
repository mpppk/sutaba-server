package twitter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

func BuildTweetUrl(userName, id string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", userName, id)
}

func CreateCRCToken(crcToken, consumerSecret string) string {
	mac := hmac.New(sha256.New, []byte(consumerSecret))
	mac.Write([]byte(crcToken))
	return "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func PostQuoteTweet(api *anaconda.TwitterApi, text string, quotedTweet *anaconda.Tweet) (anaconda.Tweet, error) {
	tweetIdStr := quotedTweet.IdStr
	tweetText := text + " " + BuildTweetUrl(
		quotedTweet.User.ScreenName,
		tweetIdStr,
	)
	return api.PostTweet(tweetText, nil)
}
