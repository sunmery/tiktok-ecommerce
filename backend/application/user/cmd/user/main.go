package main

import (
	"flag"
	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"os"
	"time"

	"backend/application/user/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name = "ecommerce-user-account-v1"
	// Version 通过环境变量来替换
	Version      string
	flagconf     string
	configCenter string
	configPath   string
	id, _        = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&configCenter, "config_center", "159.75.231.54:8500", "config center url, eg: -config_center 127.0.0.1:8500")
	flag.StringVar(&configPath, "config_path", "ecommerce/user/account/config.yaml", "config center path, eg: -config_center ecommerce/user/account/config.yaml")
	flag.StringVar(&Version, "version", "v0.0.1", "version, eg: -version v0.0.1")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, r registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(r),
	)
}

func InitConsul(filePath string) config.Source {
	consulClient, err := api.NewClient(&api.Config{
		Address:  configCenter,
		Scheme:   "http",
		WaitTime: time.Second * 15,
	})
	if err != nil {
		panic(err)
	}
	cs, err := consul.New(consulClient, consul.WithPath(filePath))
	if err != nil {
		panic(err)
	}
	return cs
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	cs := InitConsul(configPath)

	c := config.New(
		config.WithSource(cs),
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	// 认证和授权
	var ac conf.Auth
	if err := c.Scan(&ac); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 注册中心和配置中心
	var cc conf.Consul
	if err := c.Scan(&cc); err != nil {
		panic(err)
	}

	// 链路追踪
	var tc conf.Trace
	if err := c.Scan(&tc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, &ac, &cc, &tc, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
