// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"backend/application/payment/internal/biz"
	"backend/application/payment/internal/conf"
	"backend/application/payment/internal/data"
	"backend/application/payment/internal/server"
	"backend/application/payment/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, consul *conf.Consul, pay *conf.Pay, observability *conf.Observability, logger log.Logger) (*kratos.App, func(), error) {
	pool := data.NewDB(confData)
	client := data.NewCache(confData)
	alipayClient := data.NewAlipay(pay)
	discovery, err := data.NewDiscovery(consul)
	if err != nil {
		return nil, nil, err
	}
	orderServiceClient, err := data.NewOrderServiceClient(discovery, logger)
	if err != nil {
		return nil, nil, err
	}
	balanceClient, err := data.NewBalancerServiceClient(discovery, logger)
	if err != nil {
		return nil, nil, err
	}
	consumerOrderClient, err := data.NewConsumerOrderServiceClient(discovery, logger)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(pool, client, logger, alipayClient, pay, orderServiceClient, balanceClient, consumerOrderClient)
	if err != nil {
		return nil, nil, err
	}
	paymentRepo := data.NewPaymentRepo(dataData, logger, pay)
	paymentUsecase := biz.NewPaymentUsecase(paymentRepo, logger)
	paymentService := service.NewPaymentService(paymentUsecase, logger)
	grpcServer := server.NewGRPCServer(paymentService, confServer, observability, logger)
	httpServer := server.NewHTTPServer(confServer, paymentService, observability, logger)
	registrar := server.NewRegistrar(consul)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
