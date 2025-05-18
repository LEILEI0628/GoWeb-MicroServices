//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/LEILEI0628/GoWeb-MicroServices/app/user/internal/biz"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/user/internal/conf"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/user/internal/data"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/user/internal/server"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/user/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
