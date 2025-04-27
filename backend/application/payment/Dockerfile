# 编译阶段
ARG GO_IMAGE=golang:1.24.0-alpine3.21
FROM --platform=$BUILDPLATFORM ${GO_IMAGE} AS compile

ARG TARGETOS=linux
ARG TARGETARCH
ARG SERVICE
ARG VERSION=dev
ARG GOPROXY=https://goproxy.cn,direct
ARG CGOENABLED=0

WORKDIR /build

# 挂载依赖缓存
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download -x

# 编译代码
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=$CGOENABLED \
    go build -ldflags="-s -w -X main.Version=$VERSION" -o /app/$SERVICE ./application/${SERVICE}/cmd/${SERVICE}

COPY ./alipayPublicCert.crt /app/alipayPublicCert.crt
COPY ./alipayRootCert.crt /app/alipayRootCert.crt
COPY ./appPublicCert.crt /app/appPublicCert.crt

# 最终镜像
FROM gcr.io/distroless/static-debian12 AS final
ARG SERVICE
COPY --from=compile /app/$SERVICE /app/service
ENTRYPOINT ["/app/service"]
