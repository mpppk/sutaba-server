package handler

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/blendle/zapdriver"

	"github.com/mpppk/anacondaaaa"

	"github.com/mpppk/sutaba-server/pkg/interface/itwitter"

	"github.com/mpppk/sutaba-server/pkg/infra/twitter"

	"github.com/mpppk/sutaba-server/pkg/interface/controller"

	"github.com/mpppk/sutaba-server/pkg/util"

	"github.com/labstack/echo/v4"
)

var tweetIDMap = util.NewIDMap(60*5, 60*10)

// ZapLogger is an example of echo middleware that logs requests using logger "zap"
func ZapLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			echoRes := c.Response()

			var res http.Response
			res.StatusCode = echoRes.Status

			n := res.StatusCode
			switch {
			case n >= 500:
				log.Error("Server error", zapdriver.HTTP(zapdriver.NewHTTP(req, &res)))
			case n >= 400:
				log.Warn("Client error", zapdriver.HTTP(zapdriver.NewHTTP(req, &res)))
			case n >= 300:
				log.Info("Redirection", zapdriver.HTTP(zapdriver.NewHTTP(req, &res)))
			default:
				log.Info("Success", zapdriver.HTTP(zapdriver.NewHTTP(req, &res)))
			}

			return nil
		}
	}
}

type PredictHandlerConfig struct {
	TweetClassificationController *controller.TweetClassificationController
}

func GeneratePredictHandler(conf *PredictHandlerConfig) func(c echo.Context) error {
	tweetIDMap.StartExpirationCheck()
	return func(c echo.Context) error {
		events := new(anacondaaaa.AccountActivityEvent)
		if err := c.Bind(events); err != nil {
			return err
		}

		util.Logger.Infow("twitter event received", "event", events)

		if events.GetEventName() != anacondaaaa.TweetCreateEventsEventName {
			return nil
		}

		var tweets []*itwitter.Tweet
		for _, anacondaTweet := range events.TweetCreateEvents {
			tweet := twitter.ToTweet(anacondaTweet)
			_, loaded := tweetIDMap.LoadOrStore(tweet.ID)
			if loaded {
				util.Logger.Infow("tweet is ignored", "reason", "already processed", "tweetId", tweet.ID)
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
