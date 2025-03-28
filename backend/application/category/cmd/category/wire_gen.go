// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"backend/application/category/internal/biz"
	"backend/application/category/internal/conf"
	"backend/application/category/internal/data"
	"backend/application/category/internal/server"
	"backend/application/category/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, consul *conf.Consul, observability *conf.Observability, logger log.Logger) (*kratos.App, func(), error) {
	pool := data.NewDB(confData)
	client := data.NewCache(confData)
	dataData, cleanup, err := data.NewData(pool, client, logger)
	if err != nil {
		return nil, nil, err
	}
	categoryRepo := data.NewCategoryRepo(dataData, logger)
	categoryUsecase := biz.NewCategoryUsecase(categoryRepo, logger)
	categoryServiceService := service.NewCategoryServiceService(categoryUsecase)
	grpcServer := server.NewGRPCServer(categoryServiceService, confServer, observability, logger)
	httpServer := server.NewHTTPServer(confServer, categoryServiceService, observability, logger)
	registrar := server.NewRegistrar(consul)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
