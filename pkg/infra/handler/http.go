package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	usecase2 "github.com/mpppk/sutaba-server/pkg/application/usecase"

	"github.com/mpppk/sutaba-server/pkg/infra/classifier"

	"github.com/mpppk/sutaba-server/pkg/util"

	"github.com/labstack/echo/v4"
	"github.com/mpppk/sutaba-server/pkg/infra/twitter"
)

type PredictHandlerConfig struct {
	SendUser             *twitter.User
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

		userIDStr := strconv.FormatInt(conf.SendUser.ID, 10)
		if events.ForUserId != userIDStr {
			util.LogPrintfInOneLine("tweet is ignored because event is not for bot")
			return c.NoContent(http.StatusNoContent)
		}

		usecase := usecase2.NewPostPredictTweetUsecase(&usecase2.PostPredictTweetUseCaseConfig{
			SendUser:          conf.SendUser,
			ClassifierClient:  classifier.NewClassifier(conf.ClassifierServerHost),
			ErrorTweetMessage: conf.ErrorTweetMessage,
			SorryTweetMessage: conf.SorryTweetMessage,
		})
		tweets := events.TweetCreateEvents
		for _, tweet := range tweets {
			if tweet.InReplyToUserID != conf.SendUser.ID {
				util.LogPrintfInOneLine("tweet is ignored because it is not sent to subscribe user")
				continue
			}
			postedTweet, ignoreReason, err := usecase.ReplyToUser(tweet)
			if err != nil {
				util.LogPrintfInOneLine("error occurred: %v", err)
				return c.String(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
			}

			if ignoreReason != "" {
				util.LogPrintfInOneLine("tweet is ignored. reason: %v", ignoreReason)
			}

			if postedTweet != nil {
				util.LogPrintfInOneLine("posted tweet: %v", postedTweet)
			}
		}

		return c.NoContent(http.StatusNoContent)
	}
}
