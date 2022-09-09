// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"veigit-system/app/user/service/internal/biz"
	"veigit-system/app/user/service/internal/conf"
	"veigit-system/app/user/service/internal/data"
	"veigit-system/app/user/service/internal/server"
	"veigit-system/app/user/service/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, jwt *conf.JWT, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData)
	dataData, cleanup, err := data.NewData(confData, logger, db)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(userRepo, logger, jwt)
	userService := service.NewUserService(userUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, jwt, userService, logger)
	app := newApp(logger, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}
