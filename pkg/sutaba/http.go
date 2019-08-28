package sutaba

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mpppk/sutaba-server/pkg/classifier"

	"github.com/mpppk/sutaba-server/pkg/util"

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

func GeneratePredictHandler(conf *PredictHandlerConfig) func(c echo.Context) error {
	return func(c echo.Context) error {
		events := new(twitter.AccountActivityEvent)
		if err := c.Bind(events); err != nil {
			return err
		}

		if eventJson, err := json.Marshal(events); err == nil {
			util.LogPrintfInOneLine("twitter event received: %s\n", string(eventJson))
		} else {
			util.LogPrintfInOneLine("twitter event received: %#v\n", events)
		}

		if events.GetEventName() != twitter.TweetCreateEventsEventName {
			return c.NoContent(http.StatusNoContent)
		}

		usecase := NewPostPredictTweetUsecase(&PostPredictTweetUseCaseConfig{
			SendUser:          conf.SendUser,
			ClassifierClient:  classifier.NewClassifier(conf.ClassifierServerHost),
			ErrorTweetMessage: conf.ErrorTweetMessage,
			SorryTweetMessage: conf.SorryTweetMessage,
		})
		tweets := events.TweetCreateEvents
		for _, tweet := range tweets {
			postedTweets, ignoreReasons, err := usecase.ReplyToUsers(tweet, conf.SubscribeUsers)
			if err != nil {
				util.LogPrintfInOneLine("error occurred: %v", err)
				return c.String(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
			}

			if len(ignoreReasons) != 0 {
				util.LogPrintfInOneLine("len(ignoreReasons) tweets are ignored. reasons: %v", ignoreReasons)
			}

			if len(postedTweets) > 0 {
				util.LogPrintfInOneLine("posted tweets: %v", postedTweets)
			}
		}

		return c.NoContent(http.StatusNoContent)
	}
}
