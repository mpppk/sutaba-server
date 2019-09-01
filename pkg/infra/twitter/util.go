package twitter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/mpppk/sutaba-server/pkg/domain/model"

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

func toQuoteTweet(text string, quotedTweetIDStr, quotedTweetUserScreenName string) string {
	return text + " " + BuildTweetUrl(
		quotedTweetUserScreenName,
		quotedTweetIDStr,
	)
}

func PostQuoteTweet(api *anaconda.TwitterApi, text string, quotedTweetIDStr, quotedTweetUserScreenName string) (*model.Tweet, error) {
	anacondaTweet, err := api.PostTweet(toQuoteTweet(text, quotedTweetIDStr, quotedTweetUserScreenName), nil)
	if err != nil {
		return nil, err
	}
	return ToTweet(&anacondaTweet), nil
}

func PostReply(api *anaconda.TwitterApi, text, toTweetIDStr string, toScreenNames []string) (*model.Tweet, error) {
	v := url.Values{}
	v.Set("in_reply_to_status_id", toTweetIDStr)
	var mentions []string
	for _, toScreenName := range toScreenNames {
		mentions = append(mentions, "@"+toScreenName)
	}
	newText := fmt.Sprintf("%s\n%s", strings.Join(mentions, " "), text)
	anacondaTweet, err := api.PostTweet(newText, v)
	if err != nil {
		return nil, err
	}
	return ToTweet(&anacondaTweet), nil
}

func PostReplyWithQuote(
	api *anaconda.TwitterApi,
	text string,
	quotedTweetIDStr string,
	quotedTweetUserScreenName string,
	toTweetIDStr string,
	toScreenNames []string,
) (*model.Tweet, error) {
	quotedTweetText := toQuoteTweet(text, quotedTweetIDStr, quotedTweetUserScreenName)
	return PostReply(api, quotedTweetText, toTweetIDStr, toScreenNames)

}
