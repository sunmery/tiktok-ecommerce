package data

import (
	authV1 "backend/api/auth/v1"
	"backend/application/user/internal/biz"
	"backend/application/user/internal/conf"
	"backend/application/user/internal/data/models"
	"backend/constants"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewCache, NewUserRepo, NewDiscovery, NewAuthServiceClient)

type Data struct {
	db         *models.Queries
	rdb        *redis.Client
	authClient authV1.AuthServiceClient
}

// NewData .
func NewData(
	pgx *pgxpool.Pool,
	rdb *redis.Client,
	authClient authV1.AuthServiceClient,
	logger log.Logger,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:         models.New(pgx),
		rdb:        rdb,
		authClient: authClient,
	}, cleanup, nil
}

func NewDB(c *conf.Data) *pgxpool.Pool {
	fmt.Printf("connecting to the database: %s\n", c.Database)
	conn, err := pgxpool.New(context.Background(), c.Database.Source)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}

	return conn
}

func NewCache(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Protocol:     3,
		Addr:         c.Cache.Addr,
		Username:     c.Cache.Username,
		Password:     c.Cache.Password,
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

// NewAuthServiceClient 认证微服务
func NewAuthServiceClient(d registry.Discovery, logger log.Logger) (authV1.AuthServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", constants.AuthServiceV1)),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return authV1.NewAuthServiceClient(conn), nil
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

func (u *userRepo) CreateCreditCard(ctx context.Context, req *biz.CreditCards) (*biz.CreditCardsReply, error) {
	// TODO implement me
	panic("implement me")
}

func (u *userRepo) UpdateCreditCard(ctx context.Context, req *biz.CreditCards) (*biz.CreditCardsReply, error) {
	// TODO implement me
	panic("implement me")
}

func (u *userRepo) DeleteCreditCard(ctx context.Context, req *biz.DeleteCreditCardsRequest) (*biz.CreditCardsReply, error) {
	// TODO implement me
	panic("implement me")
}

func (u *userRepo) GetCreditCard(ctx context.Context, req *biz.GetCreditCardsRequest) (*biz.CreditCards, error) {
	// TODO implement me
	panic("implement me")
}

func (u *userRepo) SearchCreditCards(ctx context.Context, req *biz.GetCreditCardsRequest) ([]*biz.CreditCards, error) {
	// TODO implement me
	panic("implement me")
}

func (u *userRepo) ListCreditCards(ctx context.Context, req *biz.CreditCardsRequest) ([]*biz.CreditCards, error) {
	// TODO implement me
	panic("implement me")
}
