ARG GO_IMAGE=golang:1.23.3-alpine3.20
FROM ${GO_IMAGE} AS builder

ARG SERVICE
ARG VERSION=dev
ARG CGO_ENABLED=0
ARG GOPROXY=https://proxy.golang.com.cn,direct

WORKDIR /build

# 仅复制根模块文件
COPY go.mod go.sum ./

# 下载全仓依赖
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

# 复制全部代码
COPY . .

# 编译指定服务（注意路径调整）
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=${CGO_ENABLED} GOOS=linux GOARCH=amd64 \
    go build -ldflags="-X main.Version=${VERSION}" -o /app/${SERVICE} ./application/${SERVICE}/cmd/${SERVICE}

# 使用最小化运行时镜像
FROM alpine:3.20

ARG SERVICE
ENV TZ=Asia/Shanghai

RUN apk add --no-cache tzdata ca-certificates && \
    ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo ${TZ} > /etc/timezone

COPY --from=builder /app/${SERVICE} /app/service

EXPOSE 30003 30004

ENTRYPOINT ["/app/service"]
