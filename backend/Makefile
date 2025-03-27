# 默认值
VERSION ?= dev
GOIMAGE ?= golang:1.24.0-alpine3.21
GOOS ?= linux
GOARCH ?= amd64
CGOENABLED ?= 0

# 动态变量
DOCKER_IMAGE=e-commerce/$(APP_NAME):$(VERSION)
SERVICE = $(APP_RELATIVE_PATH)
REPOSITORY = sumery/$(SERVICE)
REGISTER = ccr.ccs.tencentyun.com
ARM64=linux/arm64/v8
AMD64=linux/amd64

GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
APP_RELATIVE_PATH=$(shell basename $$PWD)
API_PROTO_FILES=$(shell cd ../../api/$(APP_RELATIVE_PATH) && find . -name *.proto)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
APP_NAME=$(APP_RELATIVE_PATH)

.PHONY: run stop

# 启动所有服务
.PHONY: run stop

# 启动所有服务
run:
	@echo "Starting all services..."
	@rm -f .pids  # 清理旧的 PID 文件
	@cd application/auth && config_center=$(CONFIG_CENTER) kratos run & echo $$! >> .pids && sleep 1 && \
	cd application/user && config_center=$(CONFIG_CENTER) kratos run & echo $$! >> .pids && sleep 1 && \
	cd application/category && config_center=$(CONFIG_CENTER) kratos run & echo $$! >> .pids && sleep 1 && \
	cd application/product && config_center=$(CONFIG_CENTER) kratos run & echo $$! >> .pids && sleep 1 && \
	cd application/cart && config_center=$(CONFIG_CENTER) kratos run & echo $$! >> .pids && sleep 1 && \
	cd application/checkout && config_center=$(CONFIG_CENTER) kratos run & echo $$! >> .pids && sleep 1 && \
	cd application/order && config_center=$(CONFIG_CENTER) kratos run & echo $$! >> .pids && sleep 1 && \
	cd application/payment && config_center=$(CONFIG_CENTER) kratos run & echo $$! >> .pids
	@echo "All services started. Use 'make stop' to stop them."

# 停止所有服务
stop:
	@echo "Stopping all services..."
	@if [ -f .pids ]; then \
		while read pid; do \
			if ps -p $$pid > /dev/null; then \
				kill $$pid || kill -15 $$pid; \
				echo "Stopped process $$pid"; \
			else \
				echo "Process $$pid not found"; \
			fi; \
		done < .pids; \
		rm -f .pids; \
		echo "All services stopped."; \
	else \
		echo "No services running."; \
	fi

.PHONY: api
# 生成api目录的proto的代码
api:
	@cd ../../api/$(APP_RELATIVE_PATH) && \
	protoc --proto_path=../../.. \
	--proto_path=. \
	--proto_path=../../third_party \
	--go_out=paths=source_relative:. \
	--go-http_out=paths=source_relative:. \
	--go-grpc_out=paths=source_relative:. \
	--go-errors_out=paths=source_relative:. \
	--validate_out=paths=source_relative,lang=go:. \
	--openapiv2_out=. \
	$(API_PROTO_FILES) && echo "ok, no err"

.PHONY: infra
# 启动基础设施
infra:
	docker compose -f infrastructure/consul/compose.yaml up -d
	docker compose -f infrastructure/citus/compose.yaml up -d
	docker compose -f infrastructure/dragonflydb/compose.yaml up -d
	docker compose -f infrastructure/minio/compose.yaml up -d

.PHONY: config
# 重新生成proto的配置
config:
	@protoc --proto_path=. \
    --proto_path=../../third_party \
    --go_out=paths=source_relative:. \
    $(INTERNAL_PROTO_FILES)

.PHONY: validate
# 生成proto参数校验
validate:
	@cd ../../api/$(APP_RELATIVE_PATH) && \
	protoc --proto_path=../../.. \
	--proto_path=. \
	--proto_path=../../third_party \
    --go_out=paths=source_relative:. \
	--validate_out=paths=source_relative,lang=go:. \
	$(API_PROTO_FILES) && echo "ok, no err"

.PHONY: build
# 构建 $(APP_NAME) 应用版本: $(VERSION)
build:
	@echo "构建 $(APP_NAME) 应用版本: $(VERSION)"
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: docker-build
# 使用 docker 构建镜像
docker-build:
	@echo "构建的微服务: $(SERVICE)"
	@echo "系统: $(GOOS) | CPU架构: $(GOARCH)"
	@echo "镜像名: $(REPOSITORY):$(VERSION)"
	docker build . \
	  -f application/$(SERVICE)/Dockerfile \
	  --progress=plain \
	  -t ecommerce/$(SERVICE):$(VERSION) \
	  --build-arg SERVICE=$(SERVICE) \
	  --build-arg CGOENABLED=$(CGOENABLED) \
	  --build-arg GOIMAGE=$(GOIMAGE) \
	  --build-arg GOOS=$(GOOS) \
	  --build-arg GOARCH=$(GOARCH) \
	  --build-arg VERSION=$(VERSION) \
	  --platform $(GOOS)/$(GOARCH)

.PHONY: docker-deployx
# 使用 docker 构建多平台架构镜像
docker-deployx:
	@echo "构建的微服务: $(SERVICE)"
	@echo "平台1: $(ARM64)"
	@echo "平台2: $(AMD64)"
	@echo "镜像名: $(REPOSITORY):$(VERSION)"
	docker buildx build . \
	  -f application/$(SERVICE)/Dockerfile \
	  --progress=plain \
	  -t $(REGISTER)/$(REPOSITORY):$(VERSION) \
	  --build-arg SERVICE=$(SERVICE) \
	  --build-arg CGOENABLED=$(CGOENABLED) \
	  --build-arg GOIMAGE=$(GOIMAGE) \
	  --build-arg VERSION=$(VERSION) \
	  --platform $(ARM64),$(AMD64) \
	  --push

# 使用方式: make docker-push SERVICE=微服务名
.PHONY: docker-push
docker-push:
	@echo "使用方式: make docker-push SERVICE=微服务名"
	@echo "OS: $(GOOS) | ARCH: $(GOARCH)"
	@echo "Docker image: $(REPOSITORY):$(VERSION)"
	docker tag ecommerce/$(SERVICE):$(VERSION) $(REGISTER)/$(REPOSITORY):$(VERSION)
	docker push $(REGISTER)/$(REPOSITORY):$(VERSION)

.PHONY: docker-deploy
docker-deploy:
	@echo "使用方式: make docker-deploy SERVICE=微服务名"
	@echo "SERVICE=$(SERVICE)"
	make docker-build SERVICE=$(SERVICE)
	@echo "SERVICE=$(SERVICE)"
	make docker-push SERVICE=$(SERVICE)

.PHONY: docker-deploy-all
# 遍历所有微服务并执行部署
docker-deploy-all:
	@for service in $(SERVICES); do \
		service_name=$$(basename $$service); \
		echo "正在部署微服务: $$service_name"; \
		make docker-deploy SERVICE=$$service_name; \
	done

.PHONY: generate
# generate client code
generate:
	go generate ./...

.PHONY: wire
# generate wire
wire:
	go generate ./...

.PHONY: test
# test go file
test:
	@echo "测试 $(APP_NAME) 正在输出覆盖率和查找潜在的goruntine竞争"
	go test -v ./... -cover -race

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
	mockgen -package mockdb -destination mock/store.go example/sqlc Store

.PHONY: evans
# grpc的命令行测试工具
evans:
	evans -r repl --host localhost -p $(GRPC_PROT)

.PHONY: sqlc
# 生成sql代码
sqlc:
	rm -rf internal/data/models/*
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
