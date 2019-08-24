package sutaba

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"

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
		util.LogPrintfInOneLine("twitter event received: %#v\n", events)

		if events.GetEventName() != twitter.TweetCreateEventsEventName {
			return c.NoContent(http.StatusNoContent)
		}

		classifierClient := classifier.NewClassifier(conf.ClassifierServerHost)

		tweets := events.TweetCreateEvents
		for _, tweet := range tweets {
			entityMediaList := tweet.Entities.Media
			if entityMediaList == nil || len(entityMediaList) == 0 {
				util.LogPrintfInOneLine("tweet is ignored because it has no media")
				continue
			}

			if !strings.Contains(tweet.Text, conf.SendUser.TargetKeyword) {
				util.LogPrintfInOneLine("tweet is ignored because it has no keyword")
				continue
			}
			for _, subscribeUser := range conf.SubscribeUsers {
				if tweet.User.Id == conf.SendUser.ID {
					util.LogPrintfInOneLine("tweet is ignored because it is sent by bot")
					continue
				}

				if tweet.InReplyToUserID != subscribeUser.ID {
					util.LogPrintfInOneLine("tweet is ignored because it is not sent to subscribe user")
					continue
				}

				mediaBytes, err := twitter.DownloadEntityMediaFromTweet(tweet, 3, 1)
				if err != nil {
					util.LogPrintfInOneLine("failed to download media: %v", err)
					continue
				}

				f := func() (*anaconda.Tweet, error) {
					predict, err := classifierClient.Predict(mediaBytes)
					if err != nil {
						util.LogPrintfInOneLine("failed to predict: %v")
						return nil, err
					}

					tweetText, err := PredToText(predict)
					if err != nil {
						util.LogPrintfInOneLine("failed to convert predict result to tweet text: %v")
						return nil, err
					}

					postedTweet, err := conf.SendUser.PostByTweetType(tweetText, tweet)
					if err != nil {
						util.LogPrintfInOneLine("failed to tweet predict result: %v")
						return nil, err
					}
					return &postedTweet, nil
				}
				postedTweet, err := f()
				if err != nil {
					util.LogPrintlnInOneLine("error occurred:", err)
					errTweetText := conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())
					if subscribeUser.IsErrorReporter {
						tweet := events.TweetCreateEvents[0]
						subscribeUser.PostErrorTweet(errTweetText, conf.SorryTweetMessage, tweet.User.ScreenName, tweet.IdStr)
					}
					return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
				}
				util.LogPrintlnInOneLine("posted tweet:", postedTweet)
			}
		}

		return c.NoContent(http.StatusNoContent)
	}
}
