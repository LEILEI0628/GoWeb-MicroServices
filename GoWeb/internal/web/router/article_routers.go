package router

import (
	"github.com/LEILEI0628/GoWeb/internal/web/handler"
	"github.com/gin-gonic/gin"
)

// ArticleRouters Article相关路由
var _ Routers = (*ArticleRouters)(nil) // 确保ArticleRouters实现了Routers接口
type ArticleRouters struct {
	handler *handler.ArticleHandler
}

func NewArticleRouters(handler *handler.ArticleHandler) *ArticleRouters {
	return &ArticleRouters{handler: handler}
}

func (ar ArticleRouters) RegisterRouters(server *gin.Engine) {
	group := server.Group("/articles")
	ar.edit(group)
	ar.publish(group)
}

func (ar ArticleRouters) edit(group *gin.RouterGroup) {
	group.POST("edit", ar.handler.Edit)
}

func (ar ArticleRouters) publish(group *gin.RouterGroup) {
	group.POST("publish", ar.handler.Publish)
}
