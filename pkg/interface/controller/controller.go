package controller

import (
	"github.com/mpppk/sutaba-server/pkg/domain/model"
	"github.com/mpppk/sutaba-server/pkg/interface/itwitter"

	"github.com/mpppk/sutaba-server/pkg/application/usecase"

	"github.com/mpppk/sutaba-server/pkg/util"
)

type TweetClassificationController struct {
	botUser                  *model.User
	predictTweetMediaUseCase usecase.PredictMessageMediaUseCase
	twitter                  itwitter.Twitter
}

type TweetClassificationControllerConfig struct {
	BotUser                  *model.User
	PredictTweetMediaUseCase usecase.PredictMessageMediaUseCase
	Twitter                  *itwitter.Twitter
}

func NewTweetClassificationController(config *TweetClassificationControllerConfig) *TweetClassificationController {
	return &TweetClassificationController{
		botUser:                  config.BotUser,
		predictTweetMediaUseCase: config.PredictTweetMediaUseCase,
		twitter:                  *config.Twitter,
	}
}

func (t *TweetClassificationController) Handle(forUserIDStr string, tweets []*itwitter.Tweet) error {
	for _, tweet := range tweets {
		if tweet.InReplyToUserID != int64(t.botUser.ID) {
			util.LogPrintfInOneLine("anacondaTweet is ignored because it is not sent to subscribe user")
			continue
		}
		message := t.twitter.NewMessage(tweet)
		ignoreReason, err := t.predictTweetMediaUseCase.Handle(forUserIDStr, message)
		if err != nil {
			util.LogPrintfInOneLine("error occurred: %v", err)
			return err
		}

		if ignoreReason != "" {
			util.LogPrintfInOneLine("anacondaTweet is ignored. reason: %v", ignoreReason)
		}
	}

	return nil
}
