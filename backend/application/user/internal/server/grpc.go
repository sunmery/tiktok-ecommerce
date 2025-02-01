package server

import (
	v1 "backend/api/user/v1"
	"backend/application/user/internal/conf"
	"backend/application/user/internal/service"
	"context"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	sentrykratos "github.com/go-kratos/sentry"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.UserService, tr *conf.Trace, logger log.Logger) *grpc.Server {
	// trace start
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			// serviceName,
			semconv.ServiceNameKey.String(tr.Jaeger.ServiceName),
			// attribute.String("exporter", "otlptracehttp"),
			// attribute.String("environment", "dev"),
			// attribute.Float64("float", 312.23),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// shutdownTracerProvider, err := initTracerProvider(ctx, res, tr.Jaeger.Http.Endpoint)
	_, err2 := initTracerProvider(ctx, res, tr.Jaeger.Http.Endpoint)
	if err2 != nil {
		log.Fatal(err)
	}
	// trace end

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			validate.Validator(), // 参数校验
			recovery.Recovery(
				// recovery.WithLogger(log.DefaultLogger),
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					// do someting
					return nil
				})),
			sentrykratos.Server(),  // must after Recovery middleware, because of the exiting order will be reversed
			logging.Server(logger), // 在 grpc.ServerOption 中引入 logging.Server(), 则会在每次收到 gRPC 请求的时候打印详细请求信息
			// tracing.Server(), // trace 链路追踪
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServiceServer(srv, greeter)
	return srv
}
