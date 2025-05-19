package ioc

import (
	"github.com/LEILEI0628/GoWeb/internal/web"
	"github.com/gin-gonic/gin"
)

func InitGin(middleware []gin.HandlerFunc, routers *web.RegisterRouters) *gin.Engine {
	server := gin.Default()
	server.Use(middleware...)
	routers.RegisterAll(server)
	return server
}
