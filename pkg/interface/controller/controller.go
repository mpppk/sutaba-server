package controller

import (
	"strings"

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

func (t *TweetClassificationController) Handle(forUserID model.UserID, tweets []*itwitter.Tweet) error {
	for _, tweet := range tweets {
		// FIXME: move to IsTargetMessage
		if tweet.InReplyToUserID != int64(t.botUser.ID) && !strings.Contains(tweet.Text, "@"+string(t.botUser.Name)) {
			util.LogPrintfInOneLine("tweet is ignored because it is not sent to subscribe user(%d): receiver(%d)", t.botUser.ID, tweet.InReplyToUserID)
			continue
		}

		messageEvent := t.twitter.NewMessageEvent(forUserID, tweet)
		ignoreReason, err := t.predictTweetMediaUseCase.Handle(messageEvent)
		if err != nil {
			util.LogPrintfInOneLine("error occurred while tweet media predicting: %v", err)
			return err
		}

		if ignoreReason != "" {
			util.LogPrintfInOneLine("tweet is ignored. reason: %v", ignoreReason)
		}
	}

	return nil
}
