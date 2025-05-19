package ioc

import (
	"github.com/LEILEI0628/GinPro/middleware/logger"
	"go.uber.org/zap"
)

func InitGlobalLogger() loggerx.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return loggerx.NewZapLogger(l)
}
