package registry

import (
	"github.com/mpppk/sutaba-server/pkg/application/ipresenter"
	"github.com/mpppk/sutaba-server/pkg/interface/presenter"
	"github.com/mpppk/sutaba-server/pkg/interface/view"
)

type Presenter interface {
	NewMessagePresenter() ipresenter.MessagePresenter
}

type PresenterConfig struct {
	View view.TwitterView
}

type presenterImpl struct {
	messagePresenter *presenter.MessagePresenter
}

func (p *presenterImpl) NewMessagePresenter() ipresenter.MessagePresenter {
	return p.messagePresenter
}

func NewPresenter(config *PresenterConfig) Presenter {
	twitterPresenter := presenter.NewPresenter(config.View)
	return &presenterImpl{
		messagePresenter: twitterPresenter,
	}
}
