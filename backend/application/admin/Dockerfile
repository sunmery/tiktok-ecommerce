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

# 最终镜像
FROM gcr.io/distroless/static-debian12 AS final
ARG SERVICE
COPY --from=compile /app/$SERVICE /app/service
ENTRYPOINT ["/app/service"]

# # 编译阶段
# ARG GO_IMAGE=ghcr.io/ranna-go/gotip:latest
# FROM --platform=$BUILDPLATFORM ${GO_IMAGE} AS compile
#
# ARG TARGETOS=linux
# ARG TARGETARCH
# ARG SERVICE
# ARG VERSION=dev
# ARG GOPROXY=https://goproxy.cn,direct
# ARG CGOENABLED=0
#
# WORKDIR /build
#
# # 挂载依赖缓存
# COPY go.mod go.sum ./
# RUN --mount=type=cache,target=/go/pkg/mod \
#     --mount=type=cache,target=/root/.cache/go-build \
#     go mod download -x
#
# # 编译代码
# # 链接器标志 (-ldflags)
# # -ldflags="<标志列表>": 传递给链接器的标志，用于控制链接过程。这是 优化构建 的关键参数之一。
# # -s: 省略符号表 (symbol table)。这会使可执行文件更小，但会使 GDB 等底层调试器难以使用（但 Go 的 delve 调试器通常不受影响）。
# # -w: 省略 DWARF 调试信息。这会进一步减小可执行文件的大小，并且使得调试更加困难（包括 delve）。
# COPY . .
# RUN --mount=type=cache,target=/go/pkg/mod \
#     --mount=type=cache,target=/root/.cache/go-build \
#     GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=$CGOENABLED \
#     gotip build -ldflags="-s -w -X main.Version=$VERSION" -race -o /app/$SERVICE ./application/${SERVICE}/cmd/${SERVICE}
#
# # 最终镜像
# FROM gcr.io/distroless/static-debian12 AS final
# ARG SERVICE
# COPY --from=compile /app/$SERVICE /app/service
# ENTRYPOINT ["/app/service"]
