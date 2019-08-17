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
	TwitterClient        *anaconda.TwitterApi
	InReplyToUserID      int64
	ClassifierServerHost string
	TweetKeyword         string
	ErrorTweetMessage    string
}

func GeneratePredictHandler(conf *PredictHandlerConfig) func(c echo.Context) error {
	return func(c echo.Context) error {
		events := new(twitter.TweetCreateEvents)
		if err := c.Bind(events); err != nil {
			return err
		}
		log.Printf("twitter event received: %#v\n", events)
		if !IsTargetTweetCreateEvents(events, conf.InReplyToUserID, conf.TweetKeyword) {
			return c.NoContent(http.StatusNoContent)
		}
		errTweetText := conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())

		tweet := &events.TweetCreateEvents[0]
		entityMedia := tweet.Entities.Media[0]
		mediaBytes, err := twitter.DownloadEntityMedia(&entityMedia, 3, 1)
		if err != nil {
			LogAndPostErrorTweet(conf.TwitterClient, errTweetText, err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to download media: %s", err))
		}

		classifier := NewClassifier(conf.ClassifierServerHost)
		predict, err := classifier.Predict(mediaBytes)
		if err != nil {
			LogAndPostErrorTweet(conf.TwitterClient, errTweetText, err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		}

		log.Printf("predict: %#v\n", predict)

		tweetText, err := PredToText(predict)
		if err != nil {
			LogAndPostErrorTweet(conf.TwitterClient, errTweetText, err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		}

		postedTweet, err := twitter.PostQuoteTweet(conf.TwitterClient, tweetText, tweet)
		if err != nil {
			LogAndPostErrorTweet(conf.TwitterClient, errTweetText, err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		}

		log.Println("posted tweet:", postedTweet)
		return c.NoContent(http.StatusNoContent)
	}
}

func LogAndPostErrorTweet(api *anaconda.TwitterApi, text string, err error) {
	log.Println(err)
	if _, err := api.PostTweet(text, nil); err != nil {
		log.Println("failed to tweet error message", err)
	}
}
