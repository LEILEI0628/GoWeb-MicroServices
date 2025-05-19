//go:build wireinject

package main

import (
	repository2 "github.com/LEILEI0628/GoWeb/interactive/repository"
	cache2 "github.com/LEILEI0628/GoWeb/interactive/repository/cache"
	dao2 "github.com/LEILEI0628/GoWeb/interactive/repository/dao"
	service2 "github.com/LEILEI0628/GoWeb/interactive/service"
	"github.com/LEILEI0628/GoWeb/internal/repository"
	articleRepo "github.com/LEILEI0628/GoWeb/internal/repository/article"
	"github.com/LEILEI0628/GoWeb/internal/repository/cache"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao"
	articleDao "github.com/LEILEI0628/GoWeb/internal/repository/dao/article"
	"github.com/LEILEI0628/GoWeb/internal/service"
	"github.com/LEILEI0628/GoWeb/internal/web"
	"github.com/LEILEI0628/GoWeb/internal/web/handler"
	"github.com/LEILEI0628/GoWeb/internal/web/router"
	"github.com/LEILEI0628/GoWeb/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewCachedInteractiveRepository,
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)

var databaseSelect = wire.NewSet(ioc.InitDB)

type DatabaseType string

var MySQL DatabaseType = "mysql"

func InitWebServer() *gin.Engine {
	//database := MySQL
	//if database == MySQL {
	//
	//}
	wire.Build(
		// 初始化最基础的第三方依赖
		//ioc.InitDB,
		databaseSelect,
		ioc.InitRedis,

		// 初始化DAO
		dao.NewUserDAO,
		articleDao.NewGORMArticleDAO,
		// 初始化Cache
		cache.NewUserCache,
		// 初始化Repository
		repository.NewUserRepository,
		articleRepo.NewArticleRepository,
		// 初始化Service
		ioc.InitUserService,
		service.NewArticleService,
		interactiveSvcProvider,
		// 初始化Handler
		handler.NewUserHandler,
		handler.NewArticleHandler,
		// 初始化Routers
		router.NewUserRouters,
		router.NewArticleRouters,

		// 初始化Routers、中间件、server
		web.NewRegisterRouters,
		ioc.InitLimiter,
		ioc.InitGlobalLogger,
		ioc.InitMiddleware,
		ioc.InitGin,
	)

	return new(gin.Engine)
}
