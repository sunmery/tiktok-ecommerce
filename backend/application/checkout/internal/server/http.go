package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/metadata"

	v1 "backend/api/checkout/v1"
	"backend/application/checkout/internal/conf"
	"backend/application/checkout/internal/service"
	"backend/constants"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	checkout *service.CheckoutServiceService,
	obs *conf.Observability,
	logger log.Logger,
) *http.Server {
	// trace start
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// examples:
			// attribute.String("exporter", "otlptracehttp"),
			// attribute.String("environment", "dev"),
			// attribute.Float64("float", 312.23),

			// The service name used to display traces in backends serviceName
			semconv.ServiceNameKey.String(constants.CartServiceV1),
		),
	)
	if err != nil {
		log.Warnf("There was a problem creating the resource: %v", err)
	}

	_, err2 := initHttpTracerProvider(ctx, res, obs.Trace.Http.Endpoint)
	if err2 != nil {
		log.Errorf("There was a problem initializing the tracer: %v", err)
	}
	// trace end
	opts := []http.ServerOption{
		http.Middleware(
			metadata.Server(),    // 元数据
			validate.Validator(), // 参数校验
			tracing.Server(),
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err any) error {
					// do someting
					return nil
				}),
			),
			logging.Server(logger), // 在 http.ServerOption 中引入 logging.Server(), 则会在每次收到 HTTP 请求的时候打印详细请求信息
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterCheckoutServiceHTTPServer(srv, checkout)
	return srv
}
