package registry

import (
	domain "github.com/mpppk/sutaba-server/pkg/domain/service"
	"github.com/mpppk/sutaba-server/pkg/infra/classifier"
	"github.com/mpppk/sutaba-server/pkg/interface/itwitter"
)

type DomainService interface {
	NewClassifierService() domain.ClassifierService
}

type ServiceConfig struct {
	ClassifierServerHost       string
	MediaDownloadRetryNum      int
	MediaDownloadRetryInterval int
	TwitterService             *itwitter.Twitter
}

type domainServiceImpl struct {
	classifierServerHost       string
	mediaDownloadRetryNum      int
	mediaDownloadRetryInterval int
	twitterService             *itwitter.Twitter
}

func NewDomainService(config *ServiceConfig) DomainService {
	return &domainServiceImpl{
		config.ClassifierServerHost,
		config.MediaDownloadRetryNum,
		config.MediaDownloadRetryInterval,
		config.TwitterService,
	}
}

func (r *domainServiceImpl) NewClassifierService() domain.ClassifierService {
	return classifier.NewImageClassifierServerService(
		r.classifierServerHost, r.mediaDownloadRetryNum, r.mediaDownloadRetryInterval, r.twitterService)
}
