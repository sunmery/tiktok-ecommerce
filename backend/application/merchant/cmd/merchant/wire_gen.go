// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"backend/application/merchant/internal/biz"
	"backend/application/merchant/internal/conf"
	"backend/application/merchant/internal/data"
	"backend/application/merchant/internal/server"
	"backend/application/merchant/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, consul *conf.Consul, observability *conf.Observability, logger log.Logger) (*kratos.App, func(), error) {
	pool := data.NewDB(confData)
	client := data.NewCache(confData)
	discovery, err := data.NewDiscovery(consul)
	if err != nil {
		return nil, nil, err
	}
	productServiceClient, err := data.NewProductServiceClient(discovery, logger)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(pool, client, logger, productServiceClient)
	if err != nil {
		return nil, nil, err
	}
	inventoryRepo := data.NewInventoryRepo(dataData, logger)
	inventoryUsecase := biz.NewInventoryUsecase(inventoryRepo, logger)
	inventoryService := service.NewInventoryService(inventoryUsecase)
	productRepo := data.NewProductRepo(dataData, logger)
	productUsecase := biz.NewProductUsecase(productRepo, logger)
	productService := service.NewProductService(productUsecase)
	orderRepo := data.NewOrderRepo(dataData, logger)
	orderUsecase := biz.NewOrderUsecase(orderRepo, logger)
	orderServiceService := service.NewOrderService(orderUsecase)
	grpcServer := server.NewGRPCServer(confServer, observability, logger, inventoryService, productService, orderServiceService)
	httpServer := server.NewHTTPServer(confServer, observability, logger, inventoryService, productService, orderServiceService)
	registrar := server.NewRegistrar(consul)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
