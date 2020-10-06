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
	Handle(forUserIDStr string, message *model.Message) (string, error)
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

func (p *PredictMessageMediaInteractor) Handle(forUserIDStr string, message *model.Message) (string, error) {
	if forUserIDStr != p.conf.BotUser.GetIDStr() { // FIXME: this is business logic
		msg := fmt.Sprintf("message is ignored because event is not for bot(id: %s) forUserID: %s", p.conf.BotUser.GetIDStr(), forUserIDStr)
		return msg, nil
	}

	var referredMsg *model.Message = nil
	if msgIsTarget, refMsgIsTarget, reason := domain.IsTargetMessage(&p.conf.BotUser, message); msgIsTarget {
		referredMsg = message
	} else if refMsgIsTarget {
		referredMsg = message.ReferencedMessage
	} else {
		return reason, nil
	}

	f := func() error {
		classifyResult, err := p.classifierService.Classify(referredMsg)
		if err != nil {
			return xerrors.Errorf("failed to classifyResult: %v", err)
		}

		if err := p.messagePresenter.ReplyResultToMessageWithReference(
			message,
			referredMsg,
			classifyResult,
			message.IsDebugMode(),
		); err != nil {
			return xerrors.Errorf("failed to post message: %v", err)
		}

		return nil
	}
	err := f()
	if err != nil {
		p.notifyError(message, err)
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
