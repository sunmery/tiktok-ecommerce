GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
# 修改：直接获取当前微服务目录名（如 application/user 下的 user）
APP_RELATIVE_PATH=$(shell basename $$PWD)
# 修改：API 路径调整为从 application/微服务目录 到根目录 api 的相对路径
API_PROTO_FILES=$(shell cd ../../api/$(APP_RELATIVE_PATH) && find . -name *.proto)

INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
KRATOS_VERSION=$(shell go mod graph |grep go-kratos/kratos/v2 |head -n 1 |awk -F '@' '{print $$2}')
KRATOS=$(GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@$(KRATOS_VERSION)
# 修改：APP_NAME 直接使用微服务名（无需路径转换）
APP_NAME=$(APP_RELATIVE_PATH)
DOCKER_IMAGE=go-kratos/beer-$(APP_NAME):0.1.0

.PHONY: init
init:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/google/wire/cmd/wire
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

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

# ... 其他目标保持不变 ...

.PHONY: docker
docker:
	cd ../.. && docker build -f deploy/build/Dockerfile --build-arg APP_RELATIVE_PATH=$(APP_RELATIVE_PATH) -t $(DOCKER_IMAGE) .

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
