#!/usr/bin/env bash
# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

save_dir="/home/docker/minio"
mkdir -p "$save_dir/data"
cd "$save_dir"

#docker stop minio || true
#docker rm minio || true
# 不需要修最后一个参数的端口:9001
docker run \
   -itd \
   -p 9000:9000 \
   -p 9001:9001 \
   --name minio \
   -v "$save_dir/data":/data \
   -e "MINIO_ROOT_USER=minio" \
   -e "MINIO_ROOT_PASSWORD=msdnmmi,." \
   quay.io/minio/minio server /data --console-address ":9001"
