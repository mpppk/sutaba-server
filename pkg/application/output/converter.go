package output

import "github.com/mpppk/sutaba-server/pkg/application/repository"

type MessageConverter interface {
	GenerateResultMessage(result *repository.ClassifyResult) (string, error)
}
