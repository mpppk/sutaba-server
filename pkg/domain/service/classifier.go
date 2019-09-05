package domain

import "github.com/mpppk/sutaba-server/pkg/domain/model"

type ClassifierService interface {
	Classify(message *model.Message) (*ClassifyResult, error)
}

type ClassifyResult struct {
	Class      string
	Confidence float64
}
