package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtV5 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"backend/application/checkout/constants"
	"backend/application/checkout/internal/conf"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,

	ac *conf.Auth,
	obs *conf.Observability,
	logger log.Logger,
) *http.Server {
	// InitSentry()
	publicKey := InitJwtKey(ac)

	// trace start
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// examples:
			// attribute.String("exporter", "otlptracehttp"),
			// attribute.String("environment", "dev"),
			// attribute.Float64("float", 312.23),

			// The service name used to display traces in backends serviceName
			semconv.ServiceNameKey.String(constants.ServiceNameV1),
		),
	)
	if err != nil {
		log.Warnf("There was a problem creating the resource: %v", err)
	}

	_, err2 := initTracerProvider(ctx, res, obs.Trace.Http.Endpoint)
	if err2 != nil {
		log.Errorf("There was a problem initializing the tracer: %v", err)
	}
	// trace end
	var opts = []http.ServerOption{
		http.Middleware(
			validate.Validator(), // 参数校验
			tracing.Server(),
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					// do someting
					return nil
				}),
			),
			logging.Server(logger), // 在 http.ServerOption 中引入 logging.Server(), 则会在每次收到 HTTP 请求的时候打印详细请求信息
			selector.Server(
				jwt.Server(
					func(token *jwtV5.Token) (interface{}, error) {
						// 检查是否使用了正确的签名方法
						if _, ok := token.Method.(*jwtV5.SigningMethodRSA); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
						}
						return publicKey, nil
					},
					jwt.WithSigningMethod(jwtV5.SigningMethodRS256),
				),
			).
				Match(NewWhiteListMatcher()).Build(),
		),
		// 浏览器跨域
		http.Filter(handlers.CORS(
			// 允许的端点列表:
			handlers.AllowedOrigins([]string{"http://localhost:3000", "http://127.0.0.1:3000", "http://127.0.0.1:443", "https://node1.apikv.com"}),
			// 允许请求的方法:
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "HEAD", "PATCH"}),
			// 允许的 Headers:
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
			// 允许跨域请求能够携带用户的凭据（例如 cookies 或 HTTP 认证信息）
			handlers.AllowCredentials(),
		)),
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
	// v1.RegisterUserServiceHTTPServer(srv, user)
	return srv
}

func MultipartFormDataDecoder(r *http.Request, v interface{}) error {
	// example: 自定义解析
	// fmt.Printf("method:%s\n", r.Method)
	// if r.Method == "POST" {
	// 	data, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		return errors.BadRequest("CODEC", err.Error())
	// 	}
	// 	if err = codec.Unmarshal(data, v); err != nil {
	// 		return errors.BadRequest("CODEC", err.Error())
	// 	}
	// }
	return nil
}
