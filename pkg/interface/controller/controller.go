package controller

import (
	"github.com/mpppk/sutaba-server/pkg/domain/model"

	"github.com/mpppk/sutaba-server/pkg/application/usecase"

	"github.com/mpppk/sutaba-server/pkg/util"
)

type TweetClassificationController struct {
	botUser                  *model.TwitterUser
	predictTweetMediaUseCase usecase.PredictTweetMediaUseCase
}

type TweetClassificationControllerConfig struct {
	BotUser                  *model.TwitterUser
	PredictTweetMediaUseCase usecase.PredictTweetMediaUseCase
}

func NewTweetClassificationController(config *TweetClassificationControllerConfig) *TweetClassificationController {
	return &TweetClassificationController{
		botUser:                  config.BotUser,
		predictTweetMediaUseCase: config.PredictTweetMediaUseCase,
	}
}

func (t *TweetClassificationController) Handle(forUserIDStr string, tweets []*model.Tweet) error {
	for _, tweet := range tweets {
		util.LogPrintlnInOneLine("user id:", tweet.InReplyToUserID, t.botUser.ID)
		if tweet.InReplyToUserID != t.botUser.ID {
			util.LogPrintfInOneLine("anacondaTweet is ignored because it is not sent to subscribe user")
			continue
		}
		ignoreReason, err := t.predictTweetMediaUseCase.Handle(forUserIDStr, tweet)
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
