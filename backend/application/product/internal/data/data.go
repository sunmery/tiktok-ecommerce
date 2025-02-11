package data

import (
	"backend/application/product/internal/conf"
	"backend/application/product/internal/data/models"
	"context"
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewCache, NewCasdoor, NewProductRepo)

type Data struct {
	db  *models.Queries
	rdb *redis.Client
	cs  *casdoorsdk.Client
}

func (d *Data) DB(ctx context.Context) any {
	panic("unimplemented")
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

// Remove the unused userRepo type declaration
// type userRepo struct {
//     data *Data
//     log  *log.Helper
// }
