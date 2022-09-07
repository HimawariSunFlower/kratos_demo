//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"veigit-system/app/system/interface/internal/biz"
	"veigit-system/app/system/interface/internal/conf"
	"veigit-system/app/system/interface/internal/data"
	"veigit-system/app/system/interface/internal/server"
	"veigit-system/app/system/interface/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Auth, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
