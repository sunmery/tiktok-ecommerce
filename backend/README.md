Kratos 使用文档
https://i0jecrneytu.feishu.cn/docx/WoStdOA3zo7KF1xnOM5caU8cnUg?from=from_copylink

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
