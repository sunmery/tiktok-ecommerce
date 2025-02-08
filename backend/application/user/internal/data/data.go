package data

import (
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
var ProviderSet = wire.NewSet(NewData, NewDB, NewCache, NewUserRepo)

type Data struct {
	db  *models.Queries
	rdb *redis.Client
	cs  *casdoorsdk.Client
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

func NewCache(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Protocol:     3,
		Addr:         c.Redis.Addr,
		Username:     c.Redis.Username,
		Password:     c.Redis.Password,
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})

	return rdb
}

type userRepo struct {
	data *Data
	log  *log.Helper
}
