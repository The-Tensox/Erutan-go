package log

import (
	"go.uber.org/zap"
)

var (
	Zap *zap.Logger
)

func Init() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Zap = l
}

