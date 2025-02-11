Kratos 使用文档
https://i0jecrneytu.feishu.cn/docx/WoStdOA3zo7KF1xnOM5caU8cnUg?from=from_copylink

## 微服务端口分配
目前每个微服务都有自己的 HTTP 和 gRPC 端口, 根据Kubernetes NodePort Service的端口(默认范围：30000-32767)分配, 从 30000 端口开始, 分配如下:
- auth: 30001,30002
- user: 30003,30004
- product: 30005,30006
- cart: 30007, 30008
- order: 30009, 30010
- checkout: 30011, 300012
- payment: 30013, 30014

## 启动基础设施
```bash
make infra
```

## 初始化数据库

升级
```bash
cd <微服务>
make migrate-up
```

降级:
```bash
cd <微服务>
make migrate-down
```
