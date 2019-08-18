package sutaba

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ChimeraCoder/anaconda"

	"github.com/labstack/echo/v4"
	"github.com/mpppk/sutaba-server/pkg/twitter"
)

type PredictHandlerConfig struct {
	Users                []*twitter.User
	ClassifierServerHost string
	ErrorTweetMessage    string
}

func postPredictTweet(events *twitter.TweetCreateEvents, user *twitter.User, classifierServerHost string) (*anaconda.Tweet, error) {
	if ok, reason := isTargetTweetCreateEvents(events, user.ID, user.TargetKeyword); !ok {
		log.Printf("tweet does not be predicted. reason: %s user: %#v\n", reason, user)
		return nil, nil
	}
	tweet := &events.TweetCreateEvents[0]
	entityMedia := tweet.Entities.Media[0]
	mediaBytes, err := twitter.DownloadEntityMedia(&entityMedia, 3, 1)
	if err != nil {
		return nil, err
	}

	classifier := NewClassifier(classifierServerHost)
	predict, err := classifier.Predict(mediaBytes)
	if err != nil {
		return nil, err
	}

	tweetText, err := PredToText(predict)
	if err != nil {
		return nil, err
	}

	postedTweet, err := twitter.PostQuoteTweet(user.Client, tweetText, tweet)
	if err != nil {
		return nil, err
	}
	return &postedTweet, nil
}

func GeneratePredictHandler(conf *PredictHandlerConfig) func(c echo.Context) error {
	return func(c echo.Context) error {
		events := new(twitter.TweetCreateEvents)
		if err := c.Bind(events); err != nil {
			return err
		}
		log.Printf("twitter event received: %#v\n", events)
		for _, user := range conf.Users {
			postedTweet, err := postPredictTweet(events, user, conf.ClassifierServerHost)
			if err != nil {
				errTweetText := conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())
				if user.IsErrorReporter {
					user.LogAndPostErrorTweet(errTweetText, err)
				}
				return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
			}
			log.Println("posted tweet:", postedTweet)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
