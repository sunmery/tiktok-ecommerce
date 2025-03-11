package server

import (
	"context"

	v1 "backend/api/auth/v1"
	"backend/application/auth/internal/conf"
	"backend/application/auth/internal/service"
	"backend/constants"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	auth *service.AuthService,
	obs *conf.Observability,
	logger log.Logger,
) *http.Server {
	// InitSentry()
	// publicKey := InitJwtKey(ac)

	// trace start
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			// serviceName,
			semconv.ServiceNameKey.String(constants.AuthServiceV1),
			// attribute.String("exporter", "otlptracehttp"),
			// attribute.String("environment", "dev"),
			// attribute.Float64("float", 312.23),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// shutdownTracerProvider, err := initTracerProvider(ctx, res, obs.Trace.Http.Endpoint)
	_, err2 := initHttpTracerProvider(ctx, res, obs.Trace.Http.Endpoint)
	if err2 != nil {
		log.Fatal(err)
	}
	// trace end
	opts := []http.ServerOption{
		http.Middleware(
			metadata.Server(),
			validate.Validator(), // 参数校验
			tracing.Server(),
			// sentrykratos.Server(), // must after Recovery middleware, because of the exiting order will be reversed
			recovery.Recovery(
				// recovery.WithLogger(log.DefaultLogger),
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					// do someting
					return nil
				}),
			),
			logging.Server(logger), // 在 http.ServerOption 中引入 logging.Server(), 则会在每次收到 gRPC 请求的时候打印详细请求信息
		),
		http.RequestDecoder(MultipartFormDataDecoder),
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
	v1.RegisterAuthServiceHTTPServer(srv, auth)
	return srv
}

func MultipartFormDataDecoder(r *http.Request, v interface{}) error {
	// 从Request Header的Content-Type中提取出对应的解码器
	_, ok := http.CodecForRequest(r, "Content-Type")
	// 如果找不到对应的解码器此时会报错
	if !ok {
		r.Header.Set("Content-Type", "application/json")
		// return errors.BadRequest("CODEC", r.Header.Get("Content-Type"))
	}
	return nil
}
