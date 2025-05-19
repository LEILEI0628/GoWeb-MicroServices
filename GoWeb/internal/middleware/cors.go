package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func ResolveCORS() gin.HandlerFunc {
	// 需要解决跨域时可以参考前端的preflight请求
	// GO解决跨域的middleware：https://github.com/gin-contrib/cors
	// CORS for Prefix http://localhost and Contains staycool.top origins
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		// AllowMethods:Default value is simple methods (GET, POST, PUT, PATCH, DELETE, HEAD, and OPTIONS)
		//AllowMethods: []string{"GET", "POST"}, // 不写就是支持所有
		// AllowHeaders:允许请求头中携带
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// ExposeHeaders:允许响应标头暴露给浏览器
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowCredentials: true, // 是否允许携带Cookie等
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境：允许域名前缀是"http://localhost"
				return true
			}
			return strings.Contains(origin, "goweb.com") // 公司域名
		},
		MaxAge: 12 * time.Hour,
	})

}
