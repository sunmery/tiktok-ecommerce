package server

import (
	v1 "backend/api/cart/v1"
	"backend/application/cart/internal/conf"
	"backend/application/cart/internal/service"
	"backend/constants"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/go-kratos/kratos/v2/errors"
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
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	cart *service.CartServiceService,
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
			// The service name used to display traces in backends
			// serviceName,
			semconv.ServiceNameKey.String(constants.CartServiceV1),
		// attribute.String("exporter", "otlptracehttp"),
		// attribute.String("environment", "dev"),
		// attribute.Float64("float", 312.23),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// shutdownTracerProvider, err := initTracerProvider(ctx, res, obs.Trace.Http.Endpoint)
	_, err2 := initTracerProvider(ctx, res, obs.Trace.Http.Endpoint)
	if err2 != nil {
		log.Fatal(err)
	}
	// trace end
	var opts = []http.ServerOption{
		http.Middleware(
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
		http.Filter(handlers.CORS( // 浏览器跨域
			handlers.AllowedOrigins([]string{"http://localhost:3000", "http://127.0.0.1:3000", "http://127.0.0.1:443", "https://node1.apikv.com"}),
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}),
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
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
	v1.RegisterCartServiceHTTPServer(srv, cart)
	return srv
}

func MultipartFormDataDecoder(r *http.Request, v interface{}) error {
	// 从Request Header的Content-Type中提取出对应的解码器
	codec, ok := http.CodecForRequest(r, "Content-Type")
	// 如果找不到对应的解码器此时会报错
	if !ok {
		//r.Header.Set("Content-Type", "application/json")
		return errors.BadRequest("CODEC", r.Header.Get("Content-Type"))
	}
	// fmt.Printf("method:%s\n", r.Method)
	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.BadRequest("CODEC", err.Error())
		}
		if err = codec.Unmarshal(data, v); err != nil {
			return errors.BadRequest("CODEC", err.Error())
		}
	}

	return nil
}
