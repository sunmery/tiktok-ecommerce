GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
APP_RELATIVE_PATH=$(shell basename $$PWD)
API_PROTO_FILES=$(shell cd ../../api/$(APP_RELATIVE_PATH) && find . -name *.proto)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
KRATOS_VERSION=$(shell go mod graph |grep go-kratos/kratos/v2 |head -n 1 |awk -F '@' '{print $$2}')
KRATOS=$(GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@$(KRATOS_VERSION)
APP_NAME=$(APP_RELATIVE_PATH)
DOCKER_IMAGE=e-commence/$(APP_NAME):dev

.PHONY: grpc
grpc:
	cd ../../api/$(APP_RELATIVE_PATH) && protoc --proto_path=. \
           --proto_path=../../third_party \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: http
http:
	cd ../../api/$(APP_RELATIVE_PATH) && protoc --proto_path=. \
           --proto_path=../../third_party \
           --go_out=paths=source_relative:. \
           --go-http_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: errors
errors:
	cd ../../api/$(APP_RELATIVE_PATH) && protoc --proto_path=. \
           --proto_path=../../third_party \
           --go_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: swagger
swagger:
	cd ../../api/$(APP_RELATIVE_PATH) && protoc --proto_path=. \
	        --proto_path=../../third_party \
	        --openapiv2_out . \
	        --openapiv2_opt logtostderr=true \
           $(API_PROTO_FILES)

.PHONY: api
api:
	@make http
	@make grpc
	@make errors
	@make swagger

.PHONY: proto
proto:
	protoc --proto_path=. \
           --proto_path=../../third_party \
           --go_out=paths=source_relative:. \
           $(INTERNAL_PROTO_FILES)

# 默认值
VERSION ?= dev
GOIMAGE ?= golang:1.23.3-alpine3.20
GOOS ?= linux
GOARCH ?= amd64
CGOENABLED ?= 0
HTTP_PORT ?= 30001
GRPC_PORT ?= 30002
# 动态变量
SERVICE = $(APP_RELATIVE_PATH)
REPOSITORY = e-commence/$(SERVICE)

.PHONY: docker-build
docker-build:
	@echo "Building Docker image for service: $(SERVICE)"
	@echo "Docker image: $(REPOSITORY):$(VERSION)"
	REPOSITORY="ecommence/$(APP_RELATIVE_PATH)"
	cd ../.. && docker build . \
	  -f application/$(APP_RELATIVE_PATH)/Dockerfile \
	  --progress=plain \
	  -t $(REPOSITORY):$(VERSION) \
	  --build-arg SERVICE=$(APP_RELATIVE_PATH) \
	  --build-arg CGOENABLED=$(CGOENABLED) \
	  --build-arg GOIMAGE=$(GOIMAGE) \
	  --build-arg GOOS=$(GOOS) \
	  --build-arg GOARCH=$(GOARCH) \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg HTTP_PORT=$(HTTP_PORT) \
	  --build-arg GRPC_PORT=$(GRPC_PORT)

REGISTER="ccr.ccs.tencentyun.com"
# 推送 Docker 镜像
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image for service: $(SERVICE)"
	@echo "Docker image: $(REPOSITORY):$(VERSION)"
	docker tag $(REPOSITORY):$(VERSION) $(REGISTER)/$(REPOSITORY):$(VERSION)
	docker push $(REGISTER)/$(REPOSITORY):$(VERSION)

.PHONY: generate
# generate client code
generate:
	go generate ./...

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: test
# test
test:
	go test -v ./... -cover

.PHONY: run
run:
	cd cmd/server/ && go run .

.PHONY: ent
ent:
	cd internal/data/ && ent generate ./ent/schema

.PHONY: wire
# generate wire
wire:
	cd cmd/server && wire

.PHONY: api
# generate api proto
api: grpc http swagger errors

.PHONY: all
# generate all
all: grpc http proto generate

# 生成迁移文件
# 使用: make new-migrate name=verify_emails
migrate-new:
	migrate create -dir internal/data/migrate -ext sql -seq $(name)

# 升级全部的迁移文件, 先安装https://github.com/golang-migrate/migrate/tree/master
migrate-up:
	#export DB_SOURCE="postgresql://postgres:postgres@localhost:5432/tiktok?sslmode=disable"
	migrate -database "${DB_SOURCE}" -path internal/data/migrate -verbose up

# 向上迁移一个版本, 根据数据库的表schema_migrations的version来决定
migrate-up1:
	migrate -database "${DB_SOURCE}" -path internal/data/migrate -verbose up 1

# 向下全部降级迁移文件, 先安装https://github.com/golang-migrate/migrate/tree/master
migrate-down:
	migrate -database "${DB_SOURCE}" -path internal/data/migrate -verbose down

# 向下降级一个版本, 根据数据库的表schema_migrations的version来决定
migrate-down1:
	migrate -database "${DB_SOURCE}" -path internal/data/migrate -verbose down 1

# Mock DB
# -package 包名
# -destination 包含生成的mock文件位置
# simple_bank/worker TaskDistributor
# 最后参数: 包路径与要生成的interface接口的名称
mock:
	#mockgen -package mockdb -destination mock/store.go example/sqlc Store

# grpc的命令行测试工具
evans:
	evans -r repl --host localhost -p 30002

# 生成sql代码
sqlc:
	sqlc generate

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
