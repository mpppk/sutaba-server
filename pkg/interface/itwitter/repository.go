package itwitter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mpppk/sutaba-server/pkg/util"
	"golang.org/x/xerrors"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type Twitter struct {
	tweetMap map[string]*Tweet
}

func NewTwitter() *Twitter {
	return &Twitter{
		tweetMap: map[string]*Tweet{},
	}
}

func (r *Twitter) NewMessage(tweet *Tweet) *model.Message {
	message := &model.Message{
		ID:       tweet.ID,
		User:     tweet.User,
		Text:     tweet.Text,
		MediaNum: len(tweet.MediaURLs),
	}

	if tweet.QuoteTweet != nil {
		message.ReferencedMessage = r.NewMessage(tweet.QuoteTweet)
	}

	r.tweetMap[message.GetIDStr()] = tweet
	return message
}

func (r *Twitter) RetrieveTweetFromMessage(message *model.Message) (*Tweet, bool) {
	tweet, ok := r.tweetMap[message.GetIDStr()]
	return tweet, ok
}

func DownloadMediaFromTweet(tweet *Tweet, retryNum, retryInterval int) ([]byte, error) {
	mediaURL, ok := tweet.GetFirstMediaURL()
	if !ok {
		return nil, xerrors.Errorf("tweet has no media: %#v", tweet)
	}
	mediaUrl, err := url.Parse(mediaURL)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse media url(%s): %w", mediaURL, err)
	}

	mediaUrlPaths := strings.Split(mediaUrl.Path, "/")
	if len(mediaUrlPaths) == 0 {
		return nil, xerrors.Errorf("invalid mediaUrl: %s", mediaURL)
	}

	cnt := 0
	for {
		bytes, err := util.DownloadFile(mediaURL)
		if err != nil {
			if cnt >= retryNum {
				return nil, xerrors.Errorf("failed to download file from %s (retired %d times): %w", mediaURL, retryNum, err)
			}

			fmt.Println(xerrors.Errorf("failed to download file from %s: %w", mediaURL, err))
			time.Sleep(time.Duration(retryInterval) * time.Second)
			cnt++
			continue
		}
		return bytes, nil
	}
}
