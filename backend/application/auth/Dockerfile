

ARG GO_IMAGE=golang:1.24.0-alpine3.21
FROM --platform=$BUILDPLATFORM golang:alpine AS compile

ARG TARGETOS=linux
ARG TARGETARCH
ARG SERVICE
ARG VERSION=dev
ARG GOPROXY=https://goproxy.cn,direct
ARG CGOENABLED=0

WORKDIR /build

# 仅复制根模块文件
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download -x

# 复制全部代码
COPY . .

RUN echo "TARGETOS=$CGO_ENABLED"
RUN echo "CGO_ENABLED=$CGO_ENABLED"

# 编译指定服务（注意路径调整）
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    CGOENABLED=$CGOENABLED \
    go build -ldflags="-X main.Version=$VERSION" -o /app/$SERVICE ./application/${SERVICE}/cmd/${SERVICE}

# 使用最小化运行时镜像
FROM alpine:3.21 AS final

ARG SERVICE
ENV TZ=Asia/Shanghai

# 安装依赖并配置时区
RUN apk add --no-cache tzdata ca-certificates && \
    ln -sf "/usr/share/zoneinfo/$TZ" /etc/localtime && \
    echo "$TZ" > /etc/timezone

COPY --from=compile /app/$SERVICE /app/service
EXPOSE 30001 30002

ENTRYPOINT ["/app/service"]
