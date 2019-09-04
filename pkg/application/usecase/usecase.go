package usecase

import (
	"fmt"
	"time"

	domain "github.com/mpppk/sutaba-server/pkg/domain/service"

	"github.com/mpppk/sutaba-server/pkg/application/service"

	"github.com/mpppk/sutaba-server/pkg/application/ipresenter"

	"github.com/mpppk/sutaba-server/pkg/application/repository"

	"github.com/mpppk/sutaba-server/pkg/domain/model"

	"github.com/mpppk/sutaba-server/pkg/util"

	"golang.org/x/xerrors"
)

type PredictTweetMediaUseCase interface {
	Handle(forUserIDStr string, message *model.Message) (string, error)
}

type PredictTweetMediaInteractorConfig struct {
	BotUser              model.User
	TargetKeyword        string
	ErrorTweetMessage    string
	SorryTweetMessage    string
	MessagePresenter     ipresenter.MessagePresenter
	ClassifierRepository repository.ImageClassifierRepository
}

type PredictTweetMediaInteractor struct {
	conf             *PredictTweetMediaInteractorConfig
	messagePresenter ipresenter.MessagePresenter
}

func NewPredictTweetMediaInteractor(conf *PredictTweetMediaInteractorConfig) *PredictTweetMediaInteractor {
	return &PredictTweetMediaInteractor{
		conf:             conf,
		messagePresenter: conf.MessagePresenter,
	}
}

func (p *PredictTweetMediaInteractor) Handle(forUserIDStr string, message *model.Message) (string, error) {
	if forUserIDStr != p.conf.BotUser.GetIDStr() { // FIXME: this is business logic
		return "anacondaTweet is ignored because event is not for bot", nil
	}

	if reason := domain.IsTargetMessage(&p.conf.BotUser, message, p.conf.TargetKeyword); reason == "" {
		f := func() error {
			messageText, err := p.predictMessageMedia(message)
			if err != nil {
				return err
			}
			return p.messagePresenter.ReplyWithQuote(
				message.User,
				message.GetIDStr(),
				message.GetIDStr(),
				message.User.Name,
				messageText,
			)
		}
		err := f()
		if err != nil {
			p.notifyError(err)
			return "", xerrors.Errorf("error occurred in Handle: %w", err)
		}
		return "", nil
	} else if !message.HasQuoteTweet() {
		return reason, nil
	}

	// Check quote message
	if reason := domain.IsTargetMessage(&p.conf.BotUser, message.QuoteMessage, p.conf.TargetKeyword); reason != "" {
		return "quoted tweet: " + reason, nil
	}

	f := func() error {
		messageText, err := p.predictMessageMedia(message.QuoteMessage)
		if err != nil {
			return err
		}

		err = p.messagePresenter.ReplyWithQuote(
			message.User,
			message.GetIDStr(),
			message.QuoteMessage.GetIDStr(),
			message.QuoteMessage.User.Name,
			messageText,
		)
		if err != nil {
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

func (p *PredictTweetMediaInteractor) predictMessageMedia(message *model.Message) (*repository.ClassifyResult, error) {
	mediaBytes, err := service.DownloadMediaFromTweet(message, 3, 1)
	if err != nil {
		return nil, err
	}

	classifyResult, err := p.conf.ClassifierRepository.Do(mediaBytes)
	if err != nil {
		return nil, xerrors.Errorf("failed to classifyResult: %v", err)
	}

	return classifyResult, err
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
