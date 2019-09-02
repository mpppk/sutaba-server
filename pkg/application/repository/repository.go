package repository

type ImageClassifierRepository interface {
	Do(image []byte) (*ClassifyResult, error)
}

type ClassifyResult struct {
	Class      string
	Confidence float64
}
