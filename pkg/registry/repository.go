package registry

import (
	"github.com/mpppk/sutaba-server/pkg/application/repository"
	"github.com/mpppk/sutaba-server/pkg/infra/classifier"
)

type Repository interface {
	NewImageClassifierRepository() repository.ImageClassifierRepository
}

type RepositoryConfig struct {
	ClassifierServerHost string
}

type repositoryImpl struct {
	classifierServerHost string
}

func NewRepository(config *RepositoryConfig) Repository {
	return &repositoryImpl{
		config.ClassifierServerHost,
	}
}

func (r *repositoryImpl) NewImageClassifierRepository() repository.ImageClassifierRepository {
	return classifier.NewImageClassifierServerRepository(r.classifierServerHost)
}
