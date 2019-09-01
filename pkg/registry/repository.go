package registry

import (
	"github.com/mpppk/sutaba-server/pkg/domain/twitter"
	itwitter "github.com/mpppk/sutaba-server/pkg/infra/twitter"
)

type Repository interface {
	NewTwitterRepository() twitter.Repository
}

type RepositoryConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type repositoryImpl struct {
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
}

func NewRepository(config RepositoryConfig) Repository {
	return &repositoryImpl{
		config.ConsumerKey,
		config.ConsumerSecret,
		config.AccessToken,
		config.AccessTokenSecret,
	}
}

func (r *repositoryImpl) NewTwitterRepository() twitter.Repository {
	return itwitter.NewRepository(r.consumerKey, r.consumerSecret, r.accessToken, r.accessTokenSecret)
}
