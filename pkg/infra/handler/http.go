package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mpppk/sutaba-server/pkg/interface/itwitter"

	"github.com/mpppk/sutaba-server/pkg/infra/twitter"

	"github.com/mpppk/sutaba-server/pkg/interface/controller"

	"github.com/mpppk/sutaba-server/pkg/util"

	"github.com/labstack/echo/v4"
)

var tweetIDMap = util.NewIDMap()

func BodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	util.LogPrintfInOneLine("Request Body: %v\n", strings.Replace(string(reqBody), "\n", " ", -1))
	util.LogPrintfInOneLine("Response Body: %v\n", strings.Replace(string(resBody), "\n", " ", -1))
}

type PredictHandlerConfig struct {
	TweetClassificationController *controller.TweetClassificationController
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
			return nil
		}

		var tweets []*itwitter.Tweet
		for _, anacondaTweet := range events.TweetCreateEvents {
			tweet := twitter.ToTweet(anacondaTweet)
			_, loaded := tweetIDMap.LoadOrStore(tweet.ID)
			if loaded {
				util.LogPrintfInOneLine("tweet is ignored because it is already processed. id: %d", tweet.ID)
			} else {
				tweets = append(tweets, tweet)
			}
		}

		if err := conf.TweetClassificationController.Handle(events.ForUserId, tweets); err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		return c.NoContent(http.StatusNoContent)
	}
}
