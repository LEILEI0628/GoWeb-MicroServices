//go:build wireinject

package startup

import (
	jwtx "github.com/LEILEI0628/GinPro/middleware/jwt"
	"github.com/LEILEI0628/GoWeb/internal/repository"
	"github.com/LEILEI0628/GoWeb/internal/repository/article"
	"github.com/LEILEI0628/GoWeb/internal/repository/cache"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao"
	"github.com/LEILEI0628/GoWeb/internal/service"
	"github.com/LEILEI0628/GoWeb/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		//articlSvcProvider,
		dao.NewGORMArticleDAO,
		article.NewArticleRepository,
		// service 部分
		service.NewArticleService,
		// handler 部分
		web.NewUserHandler,
		web.NewArticleHandler,
		jwtx.NewRedisJWTHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider,
		dao.NewGORMArticleDAO,
		service.NewArticleService,
		web.NewArticleHandler,
		article.NewArticleRepository,
	)
	return &web.ArticleHandler{}
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitJwtHdl() jwtx.Handler {
	wire.Build(thirdProvider, jwtx.NewRedisJWTHandler)
	return jwtx.NewRedisJWTHandler(nil)
}
