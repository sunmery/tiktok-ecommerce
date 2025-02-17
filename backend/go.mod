module backend

go 1.22.0

require (
	github.com/casdoor/casdoor-go-sdk v1.3.0
	github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/exaring/otelpgx v0.9.0
	github.com/getsentry/sentry-go v0.31.1
	github.com/gin-gonic/gin v1.4.0
	github.com/go-kratos/kratos/contrib/config/consul/v2 v2.0.0-20250127095200-20970020a5ef
	github.com/go-kratos/kratos/contrib/registry/consul/v2 v2.0.0-20250127095200-20970020a5ef
	github.com/go-kratos/kratos/v2 v2.8.3
	github.com/go-kratos/sentry v0.0.0-20211021071616-de3a2011c4e4
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.6.0
	github.com/gorilla/handlers v1.5.2
	github.com/hashicorp/consul/api v1.31.0
	github.com/jackc/pgx/v5 v5.7.2
	github.com/minio/minio-go/v7 v7.0.84
	github.com/redis/go-redis/v9 v9.7.0
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/viper v1.3.2
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.34.0
	go.opentelemetry.io/otel/sdk v1.34.0
	go.uber.org/automaxprocs v1.6.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250127172529-29210b9bc287
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/magiconair/properties v1.8.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/spf13/afero v1.10.0 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/exp v0.0.0-20230817173708-d852ddb80c63 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250124145028-65684f501c47 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
