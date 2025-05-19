//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/conf"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/ioc"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/repository"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/repository/cache"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/repository/dao"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/server"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// thirdPartySet 基础依赖
var thirdPartySet = wire.NewSet(ioc.InitDB,
	ioc.InitGlobalLogger,
	//ioc.InitKratosLogger,
	ioc.InitKafka,
	ioc.InitRedis)

// interactiveSvcProvider 业务依赖
var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

// providerSet is server providers.
var providerSet = wire.NewSet(server.NewGRPCServer, server.NewHTTPServer)

// wireApp init kratos application.
func wireApp(log.Logger, *conf.Server, *conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(
		thirdPartySet,
		interactiveSvcProvider,
		providerSet,
		newApp))
}
