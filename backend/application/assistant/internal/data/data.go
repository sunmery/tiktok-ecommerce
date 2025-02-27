package data

import (
	orderV1 "backend/api/order/v1"
	productv1 "backend/api/product/v1"
	"backend/application/assistant/internal/conf"
	"backend/constants"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewCache, NewDiscovery, NewOrderServiceClient,NewProductServiceClient, NewAssistantRepo)

type Data struct {
	rdb    *redis.Client
	logger *log.Helper
	orderv1 orderV1.OrderServiceClient
	productv1 productv1.ProductServiceClient
}

// NewData .
func NewData(
	rdb *redis.Client,
	logger log.Logger,
	orderv1 orderV1.OrderServiceClient,
	productv1 productv1.ProductServiceClient,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		orderv1:orderv1,
		productv1:productv1,
		rdb:    rdb,                   // 缓存
		logger: log.NewHelper(logger), // 注入日志
	}, cleanup, nil
}

// NewCache 缓存
func NewCache(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Protocol: 3,
		Addr:     c.Cache.Addr,
		// Username:     c.Cache.Username,
		// Password:     c.Cache.Password,
		DialTimeout:  c.Cache.DialTimeout.AsDuration(),
		ReadTimeout:  c.Cache.ReadTimeout.AsDuration(),
		WriteTimeout: c.Cache.WriteTimeout.AsDuration(),
	})

	return rdb
}

// NewDiscovery 配置服务发现功能
func NewDiscovery(conf *conf.Consul) (registry.Discovery, error) {
	c := consulAPI.DefaultConfig()
	c.Address = conf.RegistryCenter.Address
	c.Scheme = conf.RegistryCenter.Scheme
	c.Token = conf.RegistryCenter.AclToken
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		return nil, err
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r, nil
}

// NewOrderServiceClient 订单微服务
func NewOrderServiceClient(d registry.Discovery, logger log.Logger) (orderV1.OrderServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", constants.OrderServiceV1)),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return orderV1.NewOrderServiceClient(conn), nil
}

// NewProductServiceClient 认证微服务
func NewProductServiceClient(d registry.Discovery, logger log.Logger) (productv1.ProductServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", constants.ProductServiceV1)),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return productv1.NewProductServiceClient(conn), nil
}
