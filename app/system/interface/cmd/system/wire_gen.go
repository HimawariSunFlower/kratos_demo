// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"veigit-system/app/system/interface/internal/biz"
	"veigit-system/app/system/interface/internal/conf"
	"veigit-system/app/system/interface/internal/data"
	"veigit-system/app/system/interface/internal/server"
	"veigit-system/app/system/interface/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, auth *conf.Auth, logger log.Logger) (*kratos.App, func(), error) {
	userServiceClient := data.NewUserServiceClient(auth)
	dataData, cleanup, err := data.NewData(confData, userServiceClient, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	authUseCase := biz.NewAuthUseCase(auth, userRepo)
	userUseCase := biz.NewUserUseCase(userRepo, logger, authUseCase)
	greeterService := service.NewGreeterService(userUseCase)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}