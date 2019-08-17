package sutaba

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mpppk/sutaba-server/internal/option"

	"github.com/ChimeraCoder/anaconda"
	"github.com/labstack/echo/v4"
	"github.com/mpppk/sutaba-server/pkg/twitter"
)

func GeneratePredictHandler(conf *option.ServerCmdConfig) func(c echo.Context) error {
	return func(c echo.Context) error {
		events := new(twitter.TweetCreateEvents)
		if err := c.Bind(events); err != nil {
			return err
		}
		fmt.Printf("tweet_create_events received: %#v\n", events)
		if !IsTargetTweetCreateEvents(events, 1354555700, conf.TweetKeyword) {
			return c.NoContent(http.StatusNoContent)
		}
		api := anaconda.NewTwitterApiWithCredentials(
			conf.TwitterAccessToken,
			conf.TwitterAccessTokenSecret,
			conf.TwitterConsumerKey,
			conf.TwitterConsumerSecret,
		)
		errTweetText := conf.ErrorMessage + fmt.Sprintf(" %v", time.Now())

		tweet := &events.TweetCreateEvents[0]
		entityMedia := tweet.Entities.Media[0]
		mediaBytes, err := twitter.DownloadEntityMedia(&entityMedia, 3, 1)
		if err != nil {
			LogAndPostErrorTweet(api, errTweetText, err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to download media: %s", err))
		}

		classifier := NewClassifier("https://sutaba-lkui2qyzba-an.a.run.app")
		predict, err := classifier.Predict(mediaBytes)
		if err != nil {
			LogAndPostErrorTweet(api, errTweetText, err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		}

		log.Printf("predict: %#v\n", predict)

		tweetText, err := PredToText(predict)
		if err != nil {
			LogAndPostErrorTweet(api, errTweetText, err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		}

		postedTweet, err := twitter.PostQuoteTweet(api, tweetText, tweet)
		if err != nil {
			LogAndPostErrorTweet(api, errTweetText, err)
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
