package domain

import (
	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type ClassifierService interface {
	Classify(media *model.MessageMedia) (*ClassifyResult, error)
}

type ClassifyResult struct {
	Class      string
	Confidence float64
}
