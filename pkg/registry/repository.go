package registry

import (
	"github.com/mpppk/sutaba-server/pkg/application/repository"
	"github.com/mpppk/sutaba-server/pkg/infra/classifier"
	itwitter "github.com/mpppk/sutaba-server/pkg/infra/twitter"
)

type Repository interface {
	NewTwitterRepository() repository.TwitterRepository
	NewImageClassifierRepository() repository.ImageClassifierRepository
}

type RepositoryConfig struct {
	ConsumerKey          string
	ConsumerSecret       string
	AccessToken          string
	AccessTokenSecret    string
	ClassifierServerHost string
}

type repositoryImpl struct {
	consumerKey          string
	consumerSecret       string
	accessToken          string
	accessTokenSecret    string
	classifierServerHost string
}

func NewRepository(config *RepositoryConfig) Repository {
	return &repositoryImpl{
		config.ConsumerKey,
		config.ConsumerSecret,
		config.AccessToken,
		config.AccessTokenSecret,
		config.ClassifierServerHost,
	}
}

func (r *repositoryImpl) NewTwitterRepository() repository.TwitterRepository {
	return itwitter.NewRepository(r.consumerKey, r.consumerSecret, r.accessToken, r.accessTokenSecret)
}

func (r *repositoryImpl) NewImageClassifierRepository() repository.ImageClassifierRepository {
	return classifier.NewImageClassifierServerRepository(r.classifierServerHost)
}
