package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/mpppk/sutaba-server/pkg/sutaba"

	"github.com/mpppk/sutaba-server/pkg/twitter"

	"github.com/labstack/echo/v4/middleware"

	"github.com/mpppk/sutaba-server/internal/option"

	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	log.Printf("Request Body: %v\n", string(reqBody))
	log.Printf("Response Body: %v\n", string(resBody))
}

func newStartCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewStartCmdConfigFromViper()
			if err != nil {
				return err
			}

			predictHandlerConfig := &sutaba.PredictHandlerConfig{
				Users: []*twitter.User{
					twitter.NewUser(
						conf.OwnerTwitterAccessToken,
						conf.OwnerTwitterAccessTokenSecret,
						conf.TwitterConsumerKey,
						conf.TwitterConsumerSecret,
						conf.OwnerTwitterUserID,
						conf.TweetKeyword,
						true,
						twitter.Reply,
					),
					twitter.NewUser(
						conf.BotTwitterAccessToken,
						conf.BotTwitterAccessTokenSecret,
						conf.TwitterConsumerKey,
						conf.TwitterConsumerSecret,
						conf.BotTwitterUserID,
						conf.TweetKeyword,
						true,
						twitter.QuoteTweet,
					),
				},
				ClassifierServerHost: conf.ClassifierServerHost,
				ErrorTweetMessage:    conf.ErrorTweetMessage,
			}

			e := echo.New()
			e.Use(middleware.BodyDump(bodyDumpHandler))

			endpoint := "/twitter/aaa"
			e.GET(endpoint, twitter.GenerateCRCTestHandler(conf.TwitterConsumerSecret))
			e.POST(endpoint, sutaba.GeneratePredictHandler(predictHandlerConfig))

			port := "1323"
			if conf.Port != "" {
				port = conf.Port
			}
			e.Logger.Fatal(e.Start(":" + port))
			return nil
		},
	}

	errorMessageFlag := &option.StringFlag{
		Flag: &option.Flag{
			Name:      "error-message",
			Usage:     "text of tweet for error notification",
			ViperName: "ErrorMessage",
		},
	}
	if err := option.RegisterStringFlag(cmd, errorMessageFlag); err != nil {
		return nil, err
	}

	tweetKeywordFlag := &option.StringFlag{
		Flag: &option.Flag{
			Name:      "keyword",
			Usage:     "process only tweets which contain this value",
			ViperName: "TweetKeyword",
		},
	}
	if err := option.RegisterStringFlag(cmd, tweetKeywordFlag); err != nil {
		return nil, err
	}

	inReplyToUserIDFlag := &option.Int64Flag{
		Flag: &option.Flag{
			Name:      "reply-id",
			Usage:     "process only tweets which reply to this user id",
			ViperName: "OwnerTwitterUserID",
		},
	}
	if err := option.RegisterInt64Flag(cmd, inReplyToUserIDFlag); err != nil {
		return nil, err
	}

	classifierServerHostFlag := &option.StringFlag{
		Flag: &option.Flag{
			Name:      "classifier-server",
			Usage:     "classifier server url",
			ViperName: "ClassifierServerHost",
		},
	}
	if err := option.RegisterStringFlag(cmd, classifierServerHostFlag); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("PORT"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("TWITTER_CONSUMER_KEY"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("TWITTER_CONSUMER_SECRET"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("OWNER_TWITTER_ACCESS_TOKEN"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("OWNER_TWITTER_ACCESS_TOKEN_SECRET"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("BOT_TWITTER_ACCESS_TOKEN"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("BOT_TWITTER_ACCESS_TOKEN_SECRET"); err != nil {
		fmt.Println(err)
		os.Exit(1)
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

	cmdGenerators = append(cmdGenerators, newStartCmd)
}
