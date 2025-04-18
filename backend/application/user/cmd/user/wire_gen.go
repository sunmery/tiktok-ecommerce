// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"backend/application/user/internal/biz"
	"backend/application/user/internal/conf"
	"backend/application/user/internal/data"
	"backend/application/user/internal/server"
	"backend/application/user/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, auth *conf.Auth, consul *conf.Consul, observability *conf.Observability, logger log.Logger) (*kratos.App, func(), error) {
	pool := data.NewDB(confData)
	client := data.NewCache(confData)
	casdoorsdkClient := data.NewCasdoor(auth)
	discovery, err := data.NewDiscovery(consul)
	if err != nil {
		return nil, nil, err
	}
	categoryServiceClient, err := data.NewCategoryClient(discovery, logger)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(pool, client, casdoorsdkClient, categoryServiceClient, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(userRepo, logger)
	userService := service.NewUserService(userUsecase)
	grpcServer := server.NewGRPCServer(confServer, userService, observability, logger)
	httpServer := server.NewHTTPServer(confServer, userService, observability, logger)
	registrar := server.NewRegistrar(consul)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
