package middleware

import (
	"context"
	"github.com/LEILEI0628/GinPro/middleware/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GlobalLogger(l loggerx.Logger) gin.HandlerFunc {
	builder := loggerx.NewBuilder(func(ctx context.Context, al *loggerx.AccessLog) {
		l.Debug("HTTP请求", loggerx.Field{Key: "al", Value: al})
	})
	flag := viper.GetBool("web.log.req.flag")
	maxLen := viper.GetInt64("web.log.req.maxLen")
	builder.AllowReqBody(flag, maxLen)
	viper.OnConfigChange(func(in fsnotify.Event) {
		flag := viper.GetBool("web.log.req.flag")
		maxLen := viper.GetInt64("web.log.req.maxLen")
		builder.AllowReqBody(flag, maxLen)
	})
	return builder.Build()
}
