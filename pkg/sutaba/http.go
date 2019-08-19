package sutaba

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mpppk/sutaba-server/pkg/util"

	"github.com/ChimeraCoder/anaconda"

	"github.com/labstack/echo/v4"
	"github.com/mpppk/sutaba-server/pkg/twitter"
)

type PredictHandlerConfig struct {
	SendUser             *twitter.User
	SubscribeUsers       []*twitter.User
	ClassifierServerHost string
	ErrorTweetMessage    string
	SorryTweetMessage    string
}

func postPredictTweet(events *twitter.TweetCreateEvents, sendUser, subscribeUser *twitter.User, classifierServerHost string) (*anaconda.Tweet, error) {
	if ok, reason := isTargetTweetCreateEvents(events, sendUser.ID, subscribeUser.ID, subscribeUser.TargetKeyword); !ok {
		util.LogPrintfInOneLine("tweet does not be predicted. reason: %s subscribeUser: %#v\n", reason, subscribeUser)
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

	postedTweet, err := sendUser.PostByTweetType(tweetText, tweet)
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
		util.LogPrintfInOneLine("twitter event received: %#v\n", events)
		for _, subscribeUser := range conf.SubscribeUsers {
			postedTweet, err := postPredictTweet(events, conf.SendUser, subscribeUser, conf.ClassifierServerHost)
			if err != nil {
				util.LogPrintlnInOneLine("error occurred:", err)
				errTweetText := conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())
				if subscribeUser.IsErrorReporter {
					tweet := &events.TweetCreateEvents[0]
					subscribeUser.PostErrorTweet(errTweetText, conf.SorryTweetMessage, tweet.User.ScreenName, tweet.IdStr)
				}
				return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
			}
			util.LogPrintlnInOneLine("posted tweet:", postedTweet)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
