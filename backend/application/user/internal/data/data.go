package data

import (
	"backend/application/user/internal/biz"
	"backend/application/user/internal/conf"
	"backend/application/user/internal/data/models"
	"context"
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewCache, NewUserRepo, NewCasdoor)

type Data struct {
	db     *models.Queries
	rdb    *redis.Client
	cs     *casdoorsdk.Client
	logger *log.Helper
}

// NewData .
func NewData(
	pgx *pgxpool.Pool,
	rdb *redis.Client,
	cs *casdoorsdk.Client,
	logger log.Logger,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:  models.New(pgx),
		rdb: rdb,
		cs:  cs,
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

func NewCasdoor(cc *conf.Auth) *casdoorsdk.Client {
	client := casdoorsdk.NewClient(
		cc.Casdoor.Server.Endpoint,
		cc.Casdoor.Server.ClientId,
		cc.Casdoor.Server.ClientSecret,
		cc.Jwt.Certificate,
		cc.Casdoor.Server.Organization,
		cc.Casdoor.Server.Application,
	)

	return client
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
