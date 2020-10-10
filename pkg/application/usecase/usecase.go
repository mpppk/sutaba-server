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

type PredictMessageMediaUseCase interface {
	Handle(messageEvent *model.MessageEvent) (string, error)
}

type PredictMessageMediaInteractorConfig struct {
	BotUser           model.User
	ErrorTweetMessage string
	SorryTweetMessage string
	MessagePresenter  ipresenter.MessagePresenter
	ClassifierService domain.ClassifierService
}

type PredictMessageMediaInteractor struct {
	conf              *PredictMessageMediaInteractorConfig
	messagePresenter  ipresenter.MessagePresenter
	classifierService domain.ClassifierService
}

func NewPredictMessageMediaInteractor(conf *PredictMessageMediaInteractorConfig) *PredictMessageMediaInteractor {
	return &PredictMessageMediaInteractor{
		conf:              conf,
		messagePresenter:  conf.MessagePresenter,
		classifierService: conf.ClassifierService,
	}
}

func (p *PredictMessageMediaInteractor) Handle(messageEvent *model.MessageEvent) (string, error) {
	var referredMsg *model.Message = nil
	if msgIsTarget, refMsgIsTarget, reason := domain.IsTargetMessageEvent(&p.conf.BotUser, messageEvent); msgIsTarget {
		referredMsg = messageEvent.Message
	} else if refMsgIsTarget {
		referredMsg = messageEvent.Message.ReferencedMessage
	} else {
		return reason, nil
	}

	f := func() error {
		classifyResult, err := p.classifierService.Classify(referredMsg)
		if err != nil {
			return xerrors.Errorf("failed to classifyResult: %v", err)
		}

		if err := p.messagePresenter.ReplyResultToMessageWithReference(
			messageEvent.Message,
			referredMsg,
			classifyResult,
			messageEvent.Message.IsDebugMode(),
		); err != nil {
			return xerrors.Errorf("failed to post messageEvent: %v", err)
		}

		return nil
	}
	err := f()
	if err != nil {
		p.notifyError(messageEvent.Message, err)
		return "", xerrors.Errorf("error occurred in JudgeAndPostPredictTweetUseCase: %w", err)
	}
	return "", nil
}

func (p *PredictMessageMediaInteractor) notifyError(toMessage *model.Message, err error) {
	errTweetText := p.conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())
	if err := p.messagePresenter.PostText(errTweetText); err != nil {
		util.Logger.Errorw("failed to message error notify message", "error", err)
	}

	if err := p.messagePresenter.ReplyToMessage(toMessage, p.conf.SorryTweetMessage); err != nil {
		util.Logger.Errorw("failed to message error notify message", "error", err)
	}
}
