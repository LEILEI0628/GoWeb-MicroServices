package middleware

import (
	"github.com/LEILEI0628/GinPro/middleware/limiter"
	"github.com/gin-gonic/gin"
)

func GlobalLimiter(limiterParam limiter.Limiter) gin.HandlerFunc {
	// 构造限流中间件
	return limiter.NewBuilder(limiterParam).
		Prefix("ip-limiter").
		KeyType(limiter.IP).
		Build()
}
