package sutaba

import (
	"strings"

	"github.com/mpppk/sutaba-server/pkg/twitter"
)

type ImagePredictResponse struct {
	Pred       string `json:"pred"`
	Confidence string `json:"confidence"`
}

func isTargetTweetCreateEvents(events *twitter.TweetCreateEvents, ignoreUserId int64, toUserId int64, keyword string) (bool, string) {
	if events.TweetCreateEvents == nil {
		return false, "event is not tweet_create_events"
	}

	tweets := events.TweetCreateEvents
	if len(tweets) == 0 {
		return false, "event has no tweets"
	}

	tweet := tweets[0]

	if tweet.User.Id == ignoreUserId {
		return false, "tweet is sent by ignore user ID"
	}

	if tweet.InReplyToUserID != toUserId {
		return false, "tweet is not target user ID"
	}

	entityMediaList := tweet.Entities.Media
	if entityMediaList == nil || len(entityMediaList) == 0 {
		return false, "tweet has no media"
	}

	if !strings.Contains(tweet.Text, keyword) {
		return false, "tweet has no keyword"
	}

	return true, ""
}
