package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mpppk/sutaba-server/pkg/sutaba"

	"github.com/mpppk/sutaba-server/pkg/twitter"

	"github.com/labstack/echo/v4/middleware"

	"github.com/ChimeraCoder/anaconda"

	"github.com/mpppk/sutaba-server/internal/option"

	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}

func newServerCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewServerCmdConfigFromViper()
			if err != nil {
				return err
			}

			e := echo.New()
			e.Use(middleware.BodyDump(bodyDumpHandler))

			endpoint := "/twitter/aaa"
			e.GET(endpoint, twitter.GenerateCRCTestHandler(conf.TwitterConsumerSecret))

			e.POST(endpoint, func(c echo.Context) error {
				events := new(twitter.TweetCreateEvents)
				if err = c.Bind(events); err != nil {
					return err
				}
				fmt.Printf("tweet_create_events received: %#v\n", events)
				if !sutaba.IsTargetTweetCreateEvents(events, 1354555700, conf.TweetKeyword) {
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
					log.Println(err)
					if _, err := api.PostTweet(errTweetText, nil); err != nil {
						log.Println("failed to tweet error message", err)
					}
					return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to download media: %s", err))
				}
				classifier := sutaba.NewClassifier("https://sutaba-lkui2qyzba-an.a.run.app")
				predict, err := classifier.Predict(mediaBytes)
				if err != nil {
					if _, err := api.PostTweet(errTweetText, nil); err != nil {
						log.Println("failed to tweet error message", err)
					}
					return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
				}
				log.Printf("predict: %#v\n", predict)

				confidence, err := strconv.ParseFloat(predict.Confidence, 32)
				if err != nil {
					log.Println(err)
					if _, err := api.PostTweet(errTweetText, nil); err != nil {
						log.Println("failed to tweet error message", err)
					}
					return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
				}

				predStr := ""
				switch predict.Pred {
				case "sutaba":
					if confidence > 0.8 {
						predStr = "間違いなくスタバ"
					} else if confidence > 0.5 {
						predStr = "スタバ"
					} else {
						predStr = "たぶんスタバ"
					}
				case "ramen":
					if confidence > 0.8 {
						predStr = "どう見てもラーメン"
					} else if confidence > 0.5 {
						predStr = "ラーメン"
					} else {
						predStr = "ラーメン...?"
					}
				case "other":
					if confidence > 0.8 {
						predStr = "スタバではない"
					} else if confidence > 0.5 {
						predStr = "スタバとは言えない"
					} else {
						predStr = "なにこれ...スタバではない気がする"
					}
				}

				userName := tweet.User.ScreenName
				tweetIdStr := tweet.IdStr
				tweetText := fmt.Sprintf("判定:%s\n確信度:%.2f", predStr, confidence*100) + "%"
				tweetText += " " + twitter.BuildTweetUrl(userName, tweetIdStr)
				postedTweet, err := api.PostTweet(tweetText, nil)
				if err != nil {
					log.Println(err)
					if _, err := api.PostTweet(errTweetText, nil); err != nil {
						log.Println("failed to tweet error message", err)
					}
					return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
				}

				log.Println("posted tweet:", postedTweet)
				return c.NoContent(http.StatusNoContent)
			})

			port := "1323"
			envPort := os.Getenv("PORT")
			if envPort != "" {
				port = envPort
			}
			e.Logger.Fatal(e.Start(":" + port))
			return nil
		},
	}
	return cmd, nil
}

func init() {
	location := "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	cmdGenerators = append(cmdGenerators, newServerCmd)
}
