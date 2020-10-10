package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/mpppk/sutaba-server/pkg/domain/model"

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
			// read req / res body
			// https://github.com/labstack/echo/blob/master/middleware/body_dump.go#L68-L79
			// Request
			reqBody := []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			// Response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			if err := next(c); err != nil {
				c.Error(err)
			}

			req := c.Request()
			echoRes := c.Response()

			var res http.Response
			res.StatusCode = echoRes.Status

			reqRaw := json.RawMessage(reqBody)
			resRaw := json.RawMessage(resBody.Bytes())
			reqField := zap.Any("req", &reqRaw)
			resField := zap.Any("res", &resRaw)
			httpField := zapdriver.HTTP(zapdriver.NewHTTP(req, &res))

			n := res.StatusCode
			switch {
			case n >= 500:
				log.Error("Server error", httpField, reqField, resField)
			case n >= 400:
				log.Warn("Client error", httpField, reqField, resField)
			case n >= 300:
				log.Info("Redirection", httpField, reqField, resField)
			default:
				log.Info("Success", httpField, reqField, resField)
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

		userId, err := strconv.Atoi(events.ForUserId)
		if err != nil {
			return fmt.Errorf("failed to convert user id from string to int: %s: %w", events.ForUserId, err)
		}

		if err := conf.TweetClassificationController.Handle(model.UserID(userId), tweets); err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		return c.NoContent(http.StatusNoContent)
	}
}
