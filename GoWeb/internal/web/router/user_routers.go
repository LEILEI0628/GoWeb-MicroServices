package router

import (
	"github.com/LEILEI0628/GoWeb/internal/web/handler"
	"github.com/gin-gonic/gin"
)

// UserRouters User相关路由
var _ Routers = (*UserRouters)(nil) // 确保UserRouters实现了Routers接口
type UserRouters struct {
	handler *handler.UserHandler
}

func NewUserRouters(handler *handler.UserHandler) *UserRouters {
	return &UserRouters{handler: handler}
}

func (ur UserRouters) RegisterRouters(server *gin.Engine) {
	// 分组路由
	userGroup := server.Group("/users")
	ur.signUpRouter(userGroup)
	ur.signInRouter(userGroup)
	ur.signOutRouter(userGroup)
	ur.profileRouter(userGroup)
	ur.editRouter(userGroup)

}

func (ur UserRouters) signUpRouter(userGroup *gin.RouterGroup) {
	userGroup.POST("/signup", ur.handler.SignUp)
}

func (ur UserRouters) signInRouter(userGroup *gin.RouterGroup) {
	userGroup.POST("/login", ur.handler.SignInByJWT)
}

func (ur UserRouters) signOutRouter(userGroup *gin.RouterGroup) {
	userGroup.GET("/logout", ur.handler.SignOutByJWT)

}

func (ur UserRouters) editRouter(userGroup *gin.RouterGroup) {
	userGroup.POST("/edit", ur.handler.EditByJWT)

}

func (ur UserRouters) profileRouter(userGroup *gin.RouterGroup) {
	userGroup.GET("/profile", ur.handler.ProfileByJWT)

}
