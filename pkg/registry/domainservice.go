package registry

import (
	domain "github.com/mpppk/sutaba-server/pkg/domain/service"
	"github.com/mpppk/sutaba-server/pkg/infra/classifier"
)

type DomainService interface {
	NewClassifierService() domain.ClassifierService
}

type ServiceConfig struct {
	ClassifierServerHost       string
	MediaDownloadRetryNum      int
	MediaDownloadRetryInterval int
}

type domainServiceImpl struct {
	classifierServerHost       string
	mediaDownloadRetryNum      int
	mediaDownloadRetryInterval int
}

func NewDomainService(config *ServiceConfig) DomainService {
	return &domainServiceImpl{
		config.ClassifierServerHost,
		config.MediaDownloadRetryNum,
		config.MediaDownloadRetryInterval,
	}
}

func (r *domainServiceImpl) NewClassifierService() domain.ClassifierService {
	return classifier.NewImageClassifierServerService(
		r.classifierServerHost, r.mediaDownloadRetryNum, r.mediaDownloadRetryInterval)
}
