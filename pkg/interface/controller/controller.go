package controller

import (
	"github.com/mpppk/sutaba-server/pkg/domain/model"

	"github.com/mpppk/sutaba-server/pkg/registry"

	"github.com/mpppk/sutaba-server/pkg/application/usecase"

	"github.com/mpppk/sutaba-server/pkg/util"
)

type TweetCreateEvent struct {
}

type TweetClassificationController struct {
	botUser      *model.TwitterUser
	presenter    registry.Presenter
	converter    registry.Converter
	repository   registry.Repository
	errorMessage string // FIXME move to presenter
	sorryMessage string // FIXME move to presenter
}

type TweetClassificationControllerConfig struct {
	BotUser      *model.TwitterUser
	Presenter    registry.Presenter
	Converter    registry.Converter
	Repository   registry.Repository
	ErrorMessage string
	SorryMessage string
}

func NewTweetClassificationController(config *TweetClassificationControllerConfig) *TweetClassificationController {
	return &TweetClassificationController{
		botUser:    config.BotUser,
		presenter:  config.Presenter,
		converter:  config.Converter,
		repository: config.Repository,
	}
}

func (t *TweetClassificationController) Handle(forUserIDStr string, tweets []*model.Tweet) error {
	if forUserIDStr != t.botUser.GetIDStr() {
		util.LogPrintfInOneLine("anacondaTweet is ignored because event is not for bot")
		return nil
	}

	uc := usecase.NewPostPredictTweetUsecase(&usecase.PostPredictTweetUseCaseConfig{
		TwitterPresenter:     t.presenter.NewMessagePresenter(),
		MessageConverter:     t.converter.NewMessageConverter(),
		SendUser:             *t.botUser,
		ClassifierRepository: t.repository.NewImageClassifierRepository(),
		ErrorTweetMessage:    t.errorMessage,
		SorryTweetMessage:    t.sorryMessage,
	})

	for _, tweet := range tweets {
		util.LogPrintlnInOneLine("user id:", tweet.InReplyToUserID, t.botUser.ID)
		if tweet.InReplyToUserID != t.botUser.ID {
			util.LogPrintfInOneLine("anacondaTweet is ignored because it is not sent to subscribe user")
			continue
		}
		ignoreReason, err := uc.ReplyToUser(tweet)
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
