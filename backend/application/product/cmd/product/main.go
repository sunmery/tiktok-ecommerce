package main

import (
	"backend/application/product/internal/conf"
	"backend/application/product/pkg"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/registry"
	"os"

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
	Name = "ecommence-product-v1"
	// Version 通过环境变量来替换
	Version      string
	flagconf     string
	configCenter string
	configPath   string
	id, _        = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&configCenter, "config_center", "localhost:8500", "config center url, eg: -config_center 127.0.0.1:8500")
	flag.StringVar(&configPath, "config_path", "ecommerce/product/config.yaml", "config center path, eg: -config_center ecommerce/user/account/config.yaml")
	flag.StringVar(&Version, "version", "v0.0.1", "version, eg: -version v0.0.1")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, r registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(fmt.Sprintf("%s-%s", id, Name)),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{
			"configCenter": configCenter,
			"configPath":   configPath,
		}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(r),
	)
}

func main() {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	consulConfig := pkg.ConfigCenter{
		Addr: configCenter,
		Path: configPath,
	}
	cs := pkg.InitConsul(consulConfig)

	c := config.New(
		config.WithSource(cs),
		// config.WithSource(
		// 	file.NewSource(flagconf),
		// ),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		log.Fatal(fmt.Errorf("load config failed:%w", err))
	}

	// 认证和授权
	var ac conf.Auth
	if err := c.Scan(&ac); err != nil {
		log.Fatal(fmt.Errorf("load auth config failed:%w", err))
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		log.Fatal(fmt.Errorf("load bootstrap config failed:%w", err))
	}

	// 注册中心和配置中心
	var cc conf.Consul
	if err := c.Scan(&cc); err != nil {
		log.Fatal(fmt.Errorf("load consul config failed:%w", err))
	}

	// 可观测性
	var obs conf.Observability
	if err := c.Scan(&obs); err != nil {
		log.Fatal(fmt.Errorf("load observability config failed:%w", err))
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, &ac, &cc, &obs, logger)
	if err != nil {
		log.Fatal(fmt.Errorf("load config failed:%w", err))
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Fatal(fmt.Errorf("app run failed:%w", err))
	}
}
