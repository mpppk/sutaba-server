package usecase

import (
	"fmt"
	"time"

	domain "github.com/mpppk/sutaba-server/pkg/domain/service"

	"github.com/mpppk/sutaba-server/pkg/application/ipresenter"

	"github.com/mpppk/sutaba-server/pkg/domain/model"

	"github.com/mpppk/sutaba-server/pkg/util"

	"golang.org/x/xerrors"
)

type PredictTweetMediaUseCase interface {
	Handle(forUserIDStr string, message *model.Message) (string, error)
}

type PredictTweetMediaInteractorConfig struct {
	BotUser           model.User
	TargetKeyword     string
	ErrorTweetMessage string
	SorryTweetMessage string
	MessagePresenter  ipresenter.MessagePresenter
	ClassifierService domain.ClassifierService
}

type PredictTweetMediaInteractor struct {
	conf              *PredictTweetMediaInteractorConfig
	messagePresenter  ipresenter.MessagePresenter
	classifierService domain.ClassifierService
}

func NewPredictTweetMediaInteractor(conf *PredictTweetMediaInteractorConfig) *PredictTweetMediaInteractor {
	return &PredictTweetMediaInteractor{
		conf:              conf,
		messagePresenter:  conf.MessagePresenter,
		classifierService: conf.ClassifierService,
	}
}

func (p *PredictTweetMediaInteractor) Handle(forUserIDStr string, message *model.Message) (string, error) {
	if forUserIDStr != p.conf.BotUser.GetIDStr() { // FIXME: this is business logic
		return "message is ignored because event is not for bot", nil
	}

	if reason := domain.IsTargetMessage(&p.conf.BotUser, message, p.conf.TargetKeyword); reason == "" {
		f := func() error {
			// ignore error because it is ensured that message has one ore more medias by IsTargetMessage
			media, _ := message.GetFirstMedia()
			classifyResult, err := p.classifierService.Classify(media)
			if err != nil {
				return xerrors.Errorf("failed to classifyResult: %v", err)
			}
			return p.messagePresenter.ReplyResultToMessageWithReference(
				message.User,
				message,
				message,
				classifyResult,
			)
		}
		err := f()
		if err != nil {
			p.notifyError(err)
			return "", xerrors.Errorf("error occurred in Handle: %w", err)
		}
		return "", nil
	} else if !message.HasMessageReference() {
		return reason, nil
	}

	// Check quote message
	if reason := domain.IsTargetMessage(&p.conf.BotUser, message.ReferencedMessage, p.conf.TargetKeyword); reason != "" {
		return "quoted tweet: " + reason, nil
	}

	f := func() error {
		// ignore error because it is ensured that message has one ore more medias by IsTargetMessage
		media, _ := message.ReferencedMessage.GetFirstMedia()
		classifyResult, err := p.classifierService.Classify(media)
		if err != nil {
			return xerrors.Errorf("failed to classifyResult: %v", err)
		}

		if err := p.messagePresenter.ReplyResultToMessageWithReference(
			message.User,
			message,
			message.ReferencedMessage,
			classifyResult,
		); err != nil {
			return xerrors.Errorf("failed to post message: %v", err)
		}

		return nil
	}
	err := f()
	if err != nil {
		p.notifyError(err)
		return "", xerrors.Errorf("error occurred in JudgeAndPostPredictTweetUseCase: %w", err)
	}
	return "", nil
}

func (p *PredictTweetMediaInteractor) notifyError(err error) {
	errTweetText := p.conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())
	if err := p.messagePresenter.PostText(p.conf.BotUser, errTweetText); err != nil {
		util.LogPrintlnInOneLine("failed to message error notify message", err)
	}

	if err := p.messagePresenter.PostText(p.conf.BotUser, p.conf.SorryTweetMessage); err != nil {
		util.LogPrintlnInOneLine("failed to message error notify message", err)
	}
}
