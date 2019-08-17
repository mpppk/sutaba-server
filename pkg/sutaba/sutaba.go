package sutaba

import (
	"strings"

	"github.com/mpppk/sutaba-server/pkg/twitter"
)

type ImagePredictResponse struct {
	Pred       string `json:"pred"`
	Confidence string `json:"confidence"`
}

func IsTargetTweetCreateEvents(events *twitter.TweetCreateEvents, toUserId int64, keyword string) bool {
	if events.TweetCreateEvents == nil {
		return false
	}

	tweets := events.TweetCreateEvents
	if len(tweets) == 0 {
		return false
	}

	tweet := tweets[0]

	if tweet.InReplyToUserID != toUserId {
		return false
	}

	entityMediaList := tweet.Entities.Media
	if entityMediaList == nil || len(entityMediaList) == 0 {
		return false
	}

	if !strings.Contains(tweet.Text, keyword) {
		return false
	}

	return true
}
