package router

import (
	"github.com/gin-gonic/gin"
)

type Routers interface {
	RegisterRouters(server *gin.Engine)
}
