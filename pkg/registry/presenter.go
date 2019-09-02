package registry

import (
	"github.com/mpppk/sutaba-server/pkg/application/output"
	"github.com/mpppk/sutaba-server/pkg/interface/presenter"
	"github.com/mpppk/sutaba-server/pkg/interface/view"
)

type Presenter interface {
	NewMessagePresenter() output.MessagePresenter
}

type PresenterConfig struct {
	View view.MessageView
}

type presenterImpl struct {
	messagePresenter *presenter.MessagePresenter
}

func (p *presenterImpl) NewMessagePresenter() output.MessagePresenter {
	return p.messagePresenter
}

func NewPresenter(config *PresenterConfig) Presenter {
	twitterPresenter := presenter.NewPresenter(config.View)
	return &presenterImpl{
		messagePresenter: twitterPresenter,
	}
}
