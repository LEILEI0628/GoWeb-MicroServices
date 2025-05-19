package middleware

import (
	jwtx "github.com/LEILEI0628/GinPro/middleware/jwt"
	"github.com/gin-gonic/gin"
	"time"
)

func JWT() gin.HandlerFunc {
	return jwtx.NewBuilder( // 校验JWT
		jwtx.WithVerificationKey("7x9FpL2QaZ8rT4wY6vBcN1mK3jH5gD7s"),
		jwtx.WithExpiresTime(time.Hour*12),
		jwtx.WithLeftTime(time.Minute*10)).
		IgnorePaths("/users/login"). // 链式调用，不同的server可定制（扩展性）
		IgnorePaths("/users/signup").
		IgnorePaths("/hello").Build() // Builder模式为了解决复杂结构构建问题
}
