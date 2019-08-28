package twitter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

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

func toQuoteTweet(text string, quotedTweet *anaconda.Tweet) string {
	tweetIdStr := quotedTweet.IdStr
	return text + " " + BuildTweetUrl(
		quotedTweet.User.ScreenName,
		tweetIdStr,
	)
}

func PostQuoteTweet(api *anaconda.TwitterApi, text string, quotedTweet *anaconda.Tweet) (anaconda.Tweet, error) {
	return api.PostTweet(toQuoteTweet(text, quotedTweet), nil)
}

func PostReply(api *anaconda.TwitterApi, text, toTweetIDStr string, toScreenNames []string) (anaconda.Tweet, error) {
	v := url.Values{}
	v.Set("in_reply_to_status_id", toTweetIDStr)
	var mentions []string
	for _, toScreenName := range toScreenNames {
		mentions = append(mentions, "@"+toScreenName)
	}
	newText := fmt.Sprintf("%s\n%s", strings.Join(mentions, " "), text)
	return api.PostTweet(newText, v)
}

func PostReplyWithQuote(
	api *anaconda.TwitterApi,
	text string,
	quotedTweet *anaconda.Tweet,
	toTweetIDStr string,
	toScreenNames []string,
) (anaconda.Tweet, error) {
	return PostReply(api, toQuoteTweet(text, quotedTweet), toTweetIDStr, toScreenNames)
}
