package web

import (
	"github.com/LEILEI0628/GoWeb/internal/web/router"
	"github.com/gin-gonic/gin"
)

type RegisterRouters struct {
	routers []router.Routers
}

func NewRegisterRouters(userRouters *router.UserRouters, articleRouters *router.ArticleRouters) *RegisterRouters { // 后期需要不断添加routers
	routers := []router.Routers{userRouters, articleRouters}
	return &RegisterRouters{routers: routers}
}

func (rr *RegisterRouters) RegisterAll(server *gin.Engine) {
	for _, v := range rr.routers {
		v.RegisterRouters(server)
	}
}
