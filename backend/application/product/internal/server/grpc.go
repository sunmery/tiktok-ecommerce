package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	v1 "backend/api/product/v1"
	"backend/application/product/internal/conf"
	"backend/application/product/internal/service"
	"backend/constants"
	"github.com/go-kratos/kratos/v2/middleware/metadata"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	product *service.ProductService,
	c *conf.Server,
	obs *conf.Observability,
	logger log.Logger,
) *grpc.Server {
	// trace start
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			// serviceName,
			semconv.ServiceNameKey.String(constants.ProductServiceV1),
			// attribute.String("exporter", "otlptracehttp"),
			// attribute.String("environment", "dev"),
			// attribute.Float64("float", 312.23),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// shutdownTracerProvider, err := initTracerProvider(ctx, res, tr.Jaeger.Http.Endpoint)
	_, err2 := initGrpcTracerProvider(ctx, res, obs.Trace.Grpc.Endpoint)
	if err2 != nil {
		log.Fatal(err)
	}
	// trace end

	opts := []grpc.ServerOption{
		grpc.Middleware(
			metadata.Server(), // 元数据
			// validate.Validator(),   // 参数校验
			recovery.Recovery(),    // 异常恢复
			logging.Server(logger), // 在 grpc.ServerOption 中引入 logging.Server(), 则会在每次收到 gRPC 请求的时候打印详细请求信息
			tracing.Server(),       // trace 链路追踪
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
	v1.RegisterProductServiceServer(srv, product)
	return srv
}
