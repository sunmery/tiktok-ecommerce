package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"

	orderv1 "backend/api/order/v1"
	"backend/constants"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	consulAPI "github.com/hashicorp/consul/api"

	"backend/application/payment/internal/data/models"

	"github.com/smartwalle/alipay/v3"

	"backend/application/payment/internal/conf"

	"github.com/exaring/otelpgx"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewCache, NewAlipay, NewDiscovery, NewOrderServiceClient, NewPaymentRepo)

type Data struct {
	db      *models.Queries
	pgx     *pgxpool.Pool
	rdb     *redis.Client
	logger  *log.Helper
	alipay  *alipay.Client
	conf    *conf.Pay
	orderv1 orderv1.OrderServiceClient
}

// 使用标准库的私有类型(包级唯一)避免冲突
type contextTxKey struct{}

// NewData .
func NewData(
	db *pgxpool.Pool,
	rdb *redis.Client,
	logger log.Logger,
	alipay *alipay.Client,
	conf *conf.Pay,
	orderv1 orderv1.OrderServiceClient,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:      models.New(db),        // 数据库
		pgx:     db,                    // 数据库事务
		rdb:     rdb,                   // 缓存
		logger:  log.NewHelper(logger), // 注入日志
		alipay:  alipay,
		conf:    conf,
		orderv1: orderv1,
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

// NewAlipay 支付宝
func NewAlipay(c *conf.Pay) *alipay.Client {
	// log.Debugf("config Pay: %+v", c)
	client, err := alipay.New(c.Alipay.AppId, c.Alipay.PrivateKey, false)
	if err != nil {
		panic(fmt.Errorf("new alipay client failed: %v", err))
	}

	// 证书方式
	// 加载应用公钥证书
	// if err := client.LoadAppCertPublicKeyFromFile("/app/appPublicCert.crt"); err != nil {
	if err := client.LoadAppCertPublicKeyFromFile("./appPublicCert.crt"); err != nil {
		panic(fmt.Errorf("load app public cert failed: %v", err))
	}
	// 加载支付宝根证书
	// if err := client.LoadAliPayRootCertFromFile("/app/alipayRootCert.crt"); err != nil {
	if err := client.LoadAliPayRootCertFromFile("./alipayRootCert.crt"); err != nil {
		panic(fmt.Errorf("load alipay root cert failed: %v", err))
	}
	// 加载支付宝公钥证书
	// if err := client.LoadAlipayCertPublicKeyFromFile("/app/alipayPublicCert.crt"); err != nil {
	if err := client.LoadAlipayCertPublicKeyFromFile("./alipayPublicCert.crt"); err != nil {
		panic(fmt.Errorf("load alipay public cert failed: %v", err))
	}

	// // 加载应用公钥证书
	// if err := client.LoadAppCertPublicKey(c.Alipay.AppPublicCert); err != nil {
	// 	panic(fmt.Errorf("load app public cert failed: %v", err))
	// }
	//
	// // 加载支付宝根证书
	// if err := client.LoadAliPayRootCert(c.Alipay.AlipayRootCert); err != nil {
	// 	panic(fmt.Errorf("load alipay root cert failed: %v", err))
	// }
	//
	// // 加载支付宝公钥证书
	// if err := client.LoadAlipayCertPublicKey(c.Alipay.AliPublicKey); err != nil {
	// 	panic(fmt.Errorf("load alipay public cert failed: %v", err))
	// }

	// 设置加密密钥
	// if err := client.SetEncryptKey(c.Alipay.Secret); err != nil {
	// 	panic("设置加密密钥失败")
	// }

	return client
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

// NewDB 数据库
func NewDB(c *conf.Data) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(c.Database.Source)
	if err != nil {
		panic(fmt.Errorf("parse database config failed: %v", err))
	}

	// 连接池设置
	cfg.MaxConns = c.Database.Pool.MaxConns
	cfg.MinConns = c.Database.Pool.MinConns
	cfg.MaxConnLifetime = c.Database.Pool.MaxConnLifetime.AsDuration()
	cfg.HealthCheckPeriod = c.Database.Pool.HealthCheckPeriod.AsDuration()
	cfg.MaxConnIdleTime = c.Database.Pool.MaxConnIdleTime.AsDuration()

	// 链路追踪配置
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	conn, connErr := pgxpool.NewWithConfig(context.Background(), cfg)
	if connErr != nil {
		panic(fmt.Sprintf("connect to database: %v", connErr))
	}

	if err := otelpgx.RecordStats(conn); err != nil {
		panic(fmt.Errorf("unable to record database stats: %w", err))
	}

	return conn
}

// DB 从上下文中获取事务或返回默认DB
// 通过 data.DB(ctx) 自动获取事务或普通连接
// example: db := p.data.DB(ctx)
func (d *Data) DB(ctx context.Context) *models.Queries {
	if tx, ok := ctx.Value(contextTxKey{}).(pgx.Tx); ok {
		// 如果上下文中有事务，使用事务版 Queries
		return models.New(tx)
	}
	// 无事务时使用普通连接
	return d.db
}

// WithTx 将事务存入上下文
func (d *Data) WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, contextTxKey{}, tx)
}

// ExecTx 支持嵌套事务检测
func (d *Data) ExecTx(ctx context.Context, fn func(context.Context) error) error {
	// 如果已经在事务中，直接执行（实现嵌套事务）
	if _, ok := ctx.Value(contextTxKey{}).(pgx.Tx); ok {
		d.logger.WithContext(ctx).Debug("reuse existing transaction")
		return fn(ctx)
	}

	// 开始新事务
	d.logger.WithContext(ctx).Info("begin transaction")
	tx, err := d.pgx.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		d.logger.WithContext(ctx).Errorf("begin transaction failed: %v", err)
		return fmt.Errorf("begin tx failed: %w", err)
	}

	// 将事务存入上下文
	txCtx := d.WithTx(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	// 执行事务操作
	if err := fn(txCtx); err != nil {
		d.logger.WithContext(ctx).Warnf("rollback transaction, reason: %v", err)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			d.logger.Errorf("rollback failed: %v", rbErr)
			return fmt.Errorf("%w (rollback err: %v)", err, rbErr)
		}
		return err
	}

	// 提交事务
	if err := tx.Commit(ctx); err != nil {
		d.logger.WithContext(ctx).Errorf("commit failed: %v", err)
		return fmt.Errorf("commit failed: %w", err)
	}
	d.logger.WithContext(ctx).Info("transaction committed")
	return nil
}
