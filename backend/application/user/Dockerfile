# https://docs.docker.com/go/dockerfile-reference/

# 定义基础镜像的 Golang 版本
ARG GOIMAGE=golang:1.23.3-alpine3.20
ARG GOOS=
ARG GOARCH=amd64

FROM ${GOIMAGE} AS build
WORKDIR /src

# 服务名
ARG SERVICE

# 版本号
ARG VERSION=dev

# Go的环境变量, 例如alpine镜像不内置gcc,则关闭CGO很有效

ARG CGOENABLED=0

# Go的环境变量, 例如alpine镜像不内置gcc,则关闭CGO很有效
ARG GO_PROXY=https://proxy.golang.com.cn,direct

COPY . .

# 设置环境变量
# RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOPROXY=$GO_PROXY

# 利用 Docker 层缓存机制，单独下载依赖项，提高后续构建速度。
# 使用缓存挂载和绑定挂载技术，避免不必要的文件复制到容器中。
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# 获取代码版本号，用于编译时标记二进制文件
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOOS=$GOOS \
    GOARCH=$GOARCH \
    CGOENABLED=$CGOENABLED \
    go build -o /bin/ ./...
#     go build -o /bin/$SERVICE ./application/$SERVICE/cmd/$SERVICE/...
    # 带版本的形式: go build -ldflags="-X main.Version=${VERSION}" -o /bin/main .

# COPY application/$SERVICE/configs /bin/configs

FROM alpine:latest AS final
# 从构建阶段复制编译好的 Go 应用程序到运行阶段
COPY --from=build /bin /bin/

# 用户进程ID
ARG UID=10001

# 后端程序的HTTP/gRPC端口
ARG HTTP_PORT=30001
ARG GRPC_PORT=30002

# 修改镜像源
# RUN echo "https://mirrors.ustc.edu.cn/alpine/latest-stable/main" > /etc/apk/repositories && \
#     echo "https://mirrors.ustc.edu.cn/alpine/latest-stable/community" >> /etc/apk/repositories

# 安装应用运行必需的系统证书和时区数据包
# RUN --mount=type=cache,target=/var/cache/apk \
#    apk --update add ca-certificates tzdata && update-ca-certificates

# RUN chmod 1777 /tmp
# # 创建一个非特权用户来运行应用，增强容器安全性
# RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "${UID}" appuser

# 设置时区为上海
# RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# RUN echo 'Asia/Shanghai' >/etc/timezone

# USER appuser

# 指定容器对外暴露的端口号
EXPOSE $HTTP_PORT
EXPOSE $GRPC_PORT

# RUN mkdir -p /data/conf

# 定义服务名
ARG SERVICE
ENV SERVICE=$SERVICE

# 设置容器启动时执行的命令
ENTRYPOINT ["/bin/sh", "-c", "/bin/$SERVICE -conf /data/conf"]
CMD ["/bin/sh", "-c", "/bin/$SERVICE"]

# 构建Docker所属的当前平台与架构的二进制文件, 进到当前的backend目录
# Docker 容器在 Linux 内核上运行，即便是在 macOS 或 Windows 环境中。
# 如果使用docker构建时传递 GOOS=darwin 会导致构建的二进制文件不兼容于 Linux 环境，从而出现 exec format error
# 所以在使用docker构建时的目标平台的GOOS应该为 linux，而非 darwin。
# SERVICE=user
# REPOSITORY="e-commence/user"
# VERSION=dev
# GOOS=linux
# GOARCH=amd64
# docker buildx build . \
#   -f application/user/Dockerfile \
#   --progress=plain \
#   -t $REPOSITORY:$VERSION \
#   --build-arg SERVICE=$SERVICE \
#   --build-arg CGOENABLED=0 \
#   --build-arg GOIMAGE=golang:1.23.3-alpine3.20 \
#   --build-arg GOOS=$GOOS \

# 构建多架构的二进制文件, 需要在Docker Desktop 启用 containerd 映像存储
# https://docs.docker.com/desktop/containerd/#enable-the-containerd-image-store
# SERVICE=user
# REPOSITORY="e-commence/user"
# VERSION=dev
# GOOS=linux
# GOARCH=amd64
# docker buildx build . \
#   -f application/user/Dockerfile \
#   --progress=plain \
#   -t $REPOSITORY:$VERSION \
#   --build-arg SERVICE=$SERVICE \
#   --build-arg CGOENABLED=0 \
#   --build-arg GOIMAGE=golang:1.23.3-alpine3.20 \
#   --build-arg GOOS=$GOOS \
#   --build-arg GOARCH=$GOARCH \
#   --build-arg VERSION=$VERSION \
#   --build-arg HTTP_PORT=30001 \
#   --build-arg GRPC_PORT=30002 \
#   --platform linux/amd64 \
#  --load

# 推送
# register="ccr.ccs.tencentyun.com"
# docker tag $REPOSITORY:$VERSION $register/$REPOSITORY:$VERSION
# docker push $register/$REPOSITORY:$VERSION

# 拉取
# GOOS=linux
# GOARCH=amd64
# docker pull $register/$REPOSITORY:$VERSION --platform $GOOS/GOARCH

# 运行
# docker run \
# --rm \
# -p 30001:30001 \
# -p 30002:30002 \
# -e GIN_MODE=release \
# -e DB_SOURCE="postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" \
# $REPOSITORY:$VERSION /bin/user -conf /bin/configs
