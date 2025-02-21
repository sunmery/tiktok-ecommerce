GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
APP_RELATIVE_PATH=$(shell basename $$PWD)
API_PROTO_FILES=$(shell cd ../../api/$(APP_RELATIVE_PATH) && find . -name *.proto)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
KRATOS_VERSION=$(shell go mod graph |grep go-kratos/kratos/v2 |head -n 1 |awk -F '@' '{print $$2}')
KRATOS=$(GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@$(KRATOS_VERSION)
APP_NAME=$(APP_RELATIVE_PATH)

# 默认值
VERSION ?= dev
GOIMAGE ?= golang:1.24.0-alpine3.20
GOOS ?= linux
GOARCH ?= amd64
CGOENABLED ?= 0
HTTP_PORT ?= 30001
GRPC_PORT ?= 30002

# 动态变量
DOCKER_IMAGE=e-commerce/$(APP_NAME):dev
SERVICE = $(APP_RELATIVE_PATH)
REPOSITORY = e-commence/$(SERVICE)
REGISTER="ccr.ccs.tencentyun.com"

.PHONY: infra
infra:
	docker compose -f infrastructure/consul/compose.yaml up -d
	docker compose -f infrastructure/citus/compose.yaml up -d
	docker compose -f infrastructure/dragonflydb/compose.yaml up -d
	docker compose -f infrastructure/minio/compose.yaml up -d

.PHONY: api
api:
	@cd ../../api/$(APP_RELATIVE_PATH) && \
	protoc --proto_path=. \
    --proto_path=../../third_party \
    --go_out=paths=source_relative:. \
    --go-http_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    --go-errors_out=paths=source_relative:. \
    --openapiv2_out . \
    --openapiv2_opt logtostderr=true \
    $(API_PROTO_FILES) && echo "ok, no err"

.PHONY: config
config:
	@protoc --proto_path=. \
    --proto_path=../../third_party \
    --go_out=paths=source_relative:. \
    $(INTERNAL_PROTO_FILES)

# 使用 docker 构建镜像, 需要传递该服务的HTTP_PORT, GRPC_PORT 变量
.PHONY: docker-build
docker-build:
	@echo "Building Docker image for service: $(SERVICE)"
	@echo "OS: $(GOOS) | ARCH: $(GOARCH)"
	@echo "Docker image: $(REPOSITORY):dev"
	REPOSITORY="ecommence/$(SERVICE)"
	cd ../.. && docker build . \
	  -f application/$(SERVICE)/Dockerfile \
	  --progress=plain \
	  -t $(REPOSITORY):dev \
	  --build-arg SERVICE=$(APP_RELATIVE_PATH) \
	  --build-arg CGOENABLED=$(CGOENABLED) \
	  --build-arg GOIMAGE=$(GOIMAGE) \
	  --build-arg GOOS=$(GOOS) \
	  --build-arg GOARCH=$(GOARCH) \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg HTTP_PORT=$(HTTP_PORT) \
	  --build-arg GRPC_PORT=$(GRPC_PORT) \
	  --platform $(GOOS)/$(GOARCH)

# 推送 Docker 镜像, 需要传递REPOSITORY变量, 例如
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image for service: $(SERVICE)"
	@echo "OS: $(GOOS) | ARCH: $(GOARCH)"
	@echo "Docker image: $(REPOSITORY):dev"
	docker tag $(REPOSITORY):dev $(REGISTER)/$(REPOSITORY):dev
	docker push $(REGISTER)/$(REPOSITORY):dev

# 运行Docker 镜像, 需要传递该服务的HTTP_PORT, GRPC_PORT 变量
.PHONY: docker-run
docker-run:
	docker run \
	  --rm \
	  -e config_path=ecommerce/user/account/config.yaml \
	  -e config_center=159.75.231.54:8500 \
	  -p $(HTTP_PORT):$(HTTP_PORT) \
	  -p $(GRPC_PORT):$(GRPC_PORT) \
	  $(REGISTER)/$(REPOSITORY):dev

.PHONY: generate
# generate client code
generate:
	go generate ./...

.PHONY: build
build:
	@echo "构建 $(APP_NAME) 应用版本: $(VERSION)"
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

# test
.PHONY: test
test:
	@echo "测试 $(APP_NAME) 正在输出覆盖率和查找潜在的goruntine竞争"
	go test -v ./... -cover -race

.PHONY: run
run:
	cd cmd/server/ && go run .

.PHONY: wire
# generate wire
wire:
	go generate ./...

.PHONY: all
# 生成全部代码
all:
	@make api
	@make config
	@make wire
	@make sqlc
	@echo "ok, no err"

.PHONY: migrate-new
# 生成迁移文件, example: make migrate-new name=verify_emails
migrate-new:
	migrate create -dir internal/data/migrate -ext sql -seq $(name)

.PHONY: migrate-up
# 升级全部的迁移文件, 先安装https://github.com/golang-migrate/migrate/tree/master
migrate-up:
	#export DB_SOURCE="postgresql://postgres:postgres@localhost:5432/ecommerce?sslmode=disable"
	migrate -database "${DB_SOURCE}" -path internal/data/migrate -verbose up

.PHONY: migrate-up1
# 向上迁移一个版本, 根据数据库的表schema_migrations的version来决定
migrate-up1:
	migrate -database "${DB_SOURCE}" -path internal/data/migrate -verbose up 1

.PHONY: migrate-down
# 向下全部降级迁移文件, 先安装https://github.com/golang-migrate/migrate/tree/master
migrate-down:
	migrate -database "${DB_SOURCE}" -path internal/data/migrate -verbose down

.PHONY: migrate-down1
# 向下降级一个版本, 根据数据库的表schema_migrations的version来决定
migrate-down1:
	migrate -database "${DB_SOURCE}" -path internal/data/migrate -verbose down 1

.PHONY: mock
# Mock DB -package 包名 -destination 包含生成的mock文件位置 dir/dir TaskDistributor 最后参数: 包路径与要生成的interface接口的名称
mock:
	#mockgen -package mockdb -destination mock/store.go example/sqlc Store

.PHONY: evans
# grpc的命令行测试工具
evans:
	evans -r repl --host localhost -p 30002

.PHONY: sqlc
# 生成sql代码
sqlc:
	sqlc generate

.PHONY: help
# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
