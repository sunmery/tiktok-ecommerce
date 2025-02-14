#!/usr/bin/env bash
# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

# 需要单独安装kustomize: https://github.com/kubernetes-sigs/kustomize

# 0. 编辑通用的模板的部署文件: backend/deploy/application/base 目录下的文件 其它的参数由不同的环境(开发/生产)进行添加和修改

# 部署到生产环境:
# 1. 从项目根目录进入到生产环境所需的部署文件中

# 2. 替换镜像为生产环境的镜像
kustomize edit set image example=registry/user/repo:tag

# 3. 部署
kustomize build . | kubectl apply -f -
