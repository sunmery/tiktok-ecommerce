package data

import (
	"context"
	"fmt"

	cartv1 "backend/api/cart/v1"
	orderv1 "backend/api/order/v1"
	paymentv1 "backend/api/payment/v1"
	productv1 "backend/api/product/v1"
	userv1 "backend/api/user/v1"
	"backend/application/checkout/internal/conf"
	"backend/constants"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewCache, NewCheckoutRepo, NewDiscovery, NewUserServiceClient, NewProductServiceClient, NewCartServiceClient, NewOrderServiceClient, NewPaymentServiceClient)

type Data struct {
	rdb       *redis.Client
	logger    *log.Helper
	cartv1    cartv1.CartServiceClient
	orderv1   orderv1.OrderServiceClient
	userv1    userv1.UserServiceClient
	paymentv1 paymentv1.PaymentServiceClient
	productv1 productv1.ProductServiceClient
}

func NewData(
	rdb *redis.Client,
	logger log.Logger,
	cartv1 cartv1.CartServiceClient,
	userv1 userv1.UserServiceClient,
	orderv1 orderv1.OrderServiceClient,
	paymentv1 paymentv1.PaymentServiceClient,
	productv1 productv1.ProductServiceClient,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		rdb:       rdb,                   // 缓存
		logger:    log.NewHelper(logger), // 注入日志
		userv1:    userv1,                // 用户服务
		cartv1:    cartv1,                // 购物车服务
		orderv1:   orderv1,               // 订单服务
		paymentv1: paymentv1,             // 支付服务
		productv1: productv1,             // 商品服务
	}, cleanup, nil
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

// NewProductServiceClient 商品微服务
func NewProductServiceClient(d registry.Discovery, logger log.Logger) (productv1.ProductServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", constants.ProductServiceV1)),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return productv1.NewProductServiceClient(conn), nil
}

// NewCartServiceClient 购物车服务
func NewCartServiceClient(d registry.Discovery, logger log.Logger) (cartv1.CartServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", constants.CartServiceV1)),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return cartv1.NewCartServiceClient(conn), nil
}

// NewUserServiceClient 用户微服务
func NewUserServiceClient(d registry.Discovery, logger log.Logger) (userv1.UserServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", constants.UserServiceV1)),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return userv1.NewUserServiceClient(conn), nil
}

// NewOrderServiceClient 订单微服务
func NewOrderServiceClient(d registry.Discovery, logger log.Logger) (orderv1.OrderServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", constants.OrderServiceV1)),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return orderv1.NewOrderServiceClient(conn), nil
}

// NewPaymentServiceClient 支付微服务
func NewPaymentServiceClient(d registry.Discovery, logger log.Logger) (paymentv1.PaymentServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", constants.PaymentServiceV1)),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return paymentv1.NewPaymentServiceClient(conn), nil
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
