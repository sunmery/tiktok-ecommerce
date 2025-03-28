// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"backend/application/admin/internal/biz"
	"backend/application/admin/internal/conf"
	"backend/application/admin/internal/data"
	"backend/application/admin/internal/server"
	"backend/application/admin/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, auth *conf.Auth, consul *conf.Consul, trace *conf.Trace, logger log.Logger) (*kratos.App, func(), error) {
	pool := data.NewDB(confData)
	client := data.NewCache(confData)
	casdoorsdkClient := data.NewCasdoor(auth)
	dataData, cleanup, err := data.NewData(pool, client, casdoorsdkClient, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(userRepo, logger)
	userService := service.NewUserService(userUsecase)
	grpcServer := server.NewGRPCServer(confServer, userService, trace, logger)
	httpServer := server.NewHTTPServer(confServer, userService, auth, trace, logger)
	registrar := server.NewRegistrar(consul)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
