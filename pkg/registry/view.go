package registry

import (
	"github.com/mpppk/sutaba-server/pkg/infra/twitter"
	"github.com/mpppk/sutaba-server/pkg/interface/view"
)

type View interface {
	NewMessageView() view.TwitterView
}

type ViewConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type viewImpl struct {
	twitterView *twitter.View
}

func (p *viewImpl) NewMessageView() view.TwitterView {
	return p.twitterView
}

func NewView(config *ViewConfig) View {
	twitterView := twitter.NewView(config.AccessToken, config.AccessTokenSecret, config.ConsumerKey, config.ConsumerSecret)
	return &viewImpl{
		twitterView: twitterView,
	}
}
