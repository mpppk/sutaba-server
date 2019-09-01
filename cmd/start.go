package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mpppk/sutaba-server/pkg/registry"

	twitter2 "github.com/mpppk/sutaba-server/pkg/domain/twitter"

	"github.com/mpppk/sutaba-server/pkg/infra/handler"

	"github.com/mpppk/sutaba-server/pkg/util"

	"github.com/spf13/viper"

	"github.com/mpppk/sutaba-server/pkg/infra/twitter"

	"github.com/labstack/echo/v4/middleware"

	"github.com/mpppk/sutaba-server/internal/option"

	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	util.LogPrintfInOneLine("Request Body: %v\n", strings.Replace(string(reqBody), "\n", " ", -1))
	util.LogPrintfInOneLine("Response Body: %v\n", strings.Replace(string(resBody), "\n", " ", -1))
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

			user := twitter2.NewTwitterUser(conf.BotTwitterUserID, "sutaba_police2")

			repositoryConfig := registry.RepositoryConfig{
				ConsumerKey:       conf.TwitterConsumerKey,
				ConsumerSecret:    conf.TwitterConsumerSecret,
				AccessToken:       conf.BotTwitterAccessToken,
				AccessTokenSecret: conf.BotTwitterAccessTokenSecret,
			}
			predictHandlerConfig := &handler.PredictHandlerConfig{
				SendUser:             &user,
				ClassifierServerHost: conf.ClassifierServerHost,
				ErrorTweetMessage:    conf.ErrorTweetMessage,
				SorryTweetMessage:    conf.SorryTweetMessage,
				Repository:           registry.NewRepository(repositoryConfig),
			}

			e := echo.New()
			e.Use(middleware.BodyDump(bodyDumpHandler))

			endpoint := "/twitter/aaa"
			e.GET(endpoint, twitter.GenerateCRCTestHandler(conf.TwitterConsumerSecret))
			e.POST(endpoint, handler.GeneratePredictHandler(predictHandlerConfig))

			port := "1323"
			if conf.Port != "" {
				port = conf.Port
			}
			e.Logger.Fatal(e.Start(":" + port))
			return nil
		},
	}

	stringFlags := []*option.StringFlag{
		{
			Flag: &option.Flag{
				Name:      "error-message",
				Usage:     "text of tweet for error notification",
				ViperName: "ERROR_TWEET_MESSAGE",
			},
		}, {
			Flag: &option.Flag{
				Name:      "sorry-message",
				Usage:     "text of tweet to send to user if process is failed",
				ViperName: "SORRY_TWEET_MESSAGE",
			},
		},
		{
			Flag: &option.Flag{
				Name:      "keyword",
				Usage:     "process only tweets which contain this value",
				ViperName: "TWEET_KEYWORD",
			},
		},
		{
			Flag: &option.Flag{
				Name:      "classifier-server",
				Usage:     "classifier server url",
				ViperName: "CLASSIFIER_SERVER_HOST",
			},
		},
	}
	for _, flag := range stringFlags {
		if err := option.RegisterStringFlag(cmd, flag); err != nil {
			return nil, err
		}
	}

	botId := &option.Int64Flag{
		Flag: &option.Flag{
			Name:      "bot-id",
			Usage:     "bot twitter user id",
			ViperName: "OwnerTwitterUserID",
		},
	}
	if err := option.RegisterInt64Flag(cmd, botId); err != nil {
		return nil, err
	}

	envStrs := []string{
		"PORT", "ERROR_TWEET_MESSAGE", "SORRY_TWEET_MESSAGE",
		"CLASSIFIER_SERVER_HOST",
		"TWEET_KEYWORD", "BOT_TWITTER_USER_ID",
		"TWITTER_CONSUMER_KEY", "TWITTER_CONSUMER_SECRET",
		"BOT_TWITTER_ACCESS_TOKEN", "BOT_TWITTER_ACCESS_TOKEN_SECRET",
	}

	for _, envStr := range envStrs {
		if err := viper.BindEnv(envStr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
