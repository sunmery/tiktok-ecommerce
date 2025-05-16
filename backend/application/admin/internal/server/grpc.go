package server

import (
	"context"

	commentv1 "backend/api/admin/comment/v1"

	"backend/application/admin/internal/service"

	orderv1 "backend/api/admin/order/v1"
	"backend/application/admin/internal/conf"
	"backend/constants"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	obs *conf.Observability,
	logger log.Logger,
	orderService *service.AdminOrderService,
	commentService *service.AdminCommentService,
) *grpc.Server {
	// trace start
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// examples:
			// attribute.String("exporter", "otlptracehttp"),
			// attribute.String("environment", "dev"),
			// attribute.Float64("float", 312.23),

			// The service name used to display traces in backends serviceName
			semconv.ServiceNameKey.String(constants.AdminServiceV1),
		),
	)
	if err != nil {
		log.Warnf("There was a problem creating the resource: %v", err)
	}

	_, err2 := initGrpcTracerProvider(ctx, res, obs.Trace.Grpc.Endpoint)
	if err2 != nil {
		log.Errorf("There was a problem initializing the tracer: %v", err)
	}
	// trace end

	opts := []grpc.ServerOption{
		grpc.Middleware(
			metadata.Server(),    // 元信息
			validate.Validator(), // 参数校验
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					// do someting
					return nil
				})),
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

	orderv1.RegisterAdminOrderServer(srv, orderService)
	commentv1.RegisterAdminCommentServer(srv, commentService)

	return srv
}
