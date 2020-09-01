package util

import (
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	l, err := zapdriver.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger = l.Sugar()
}
