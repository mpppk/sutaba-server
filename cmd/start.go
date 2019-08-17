package cmd

import (
	"log"
	"time"

	"github.com/ChimeraCoder/anaconda"

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
				TwitterClient: anaconda.NewTwitterApiWithCredentials(
					conf.TwitterAccessToken,
					conf.TwitterAccessTokenSecret,
					conf.TwitterConsumerKey,
					conf.TwitterConsumerSecret,
				),
				InReplyToUserID:      conf.InReplyToUserID,
				ClassifierServerHost: conf.ClassifierServerHost,
				TweetKeyword:         conf.TweetKeyword,
				ErrorTweetMessage:    conf.ErrorMessage,
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
			ViperName: "InReplyToUserID",
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
