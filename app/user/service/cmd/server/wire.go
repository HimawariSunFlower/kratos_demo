//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"veigit-system/app/user/service/internal/biz"
	"veigit-system/app/user/service/internal/conf"
	"veigit-system/app/user/service/internal/data"
	"veigit-system/app/user/service/internal/server"
	"veigit-system/app/user/service/internal/service"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, *conf.JWT, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
