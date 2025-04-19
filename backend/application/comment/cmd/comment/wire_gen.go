// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"backend/application/comment/internal/biz"
	"backend/application/comment/internal/conf"
	"backend/application/comment/internal/data"
	"backend/application/comment/internal/server"
	"backend/application/comment/internal/service"
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
	dataData, cleanup, err := data.NewData(pool, client, logger)
	if err != nil {
		return nil, nil, err
	}
	commentRepo := data.NewCommentRepo(dataData, logger)
	commentUsecase := biz.NewCommentUsecase(commentRepo, logger)
	commentService := service.NewCommentService(commentUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, commentService, observability, logger)
	httpServer := server.NewHTTPServer(commentService, confServer, observability, logger)
	registrar := server.NewRegistrar(consul)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
