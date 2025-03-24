## 项目开发

目前每个微服务都有自己的 HTTP 和 gRPC 端口, 根据Kubernetes NodePort Service的端口(默认范围：30000-32767)分配, 从 30000 端口开始, 分配如下:
- auth: 30001,30002
- user: 30003,30004
- product: 30005,30006
- cart: 30007,30008
- order: 30009,30010
- checkout: 30011,30012
- payment: 30013,30014
- category: 30015,30016
- assistant: 30017,30018
- merchants: 30019,30020
- admin: 30021, 30022

项目是前端后端分离,仓库地址是 https://github.com/sunmery/tiktok-e-commence 前端以 **submodule** 方式链接到单独的前端项目, 根据需要来决定是否也拉取前端项目仓库代码:

1. 拉取, 包含前端和网关项目. 不需要前端删除`-recurse-submodules`即可
```bash
git submodule update --init --recursive
```

## 架构

南北向为 Client 发送请求到 API Gateway，由 API Gateway 分发流量到对应的微服务。东西向为微服务之间通过 gRPC 来交互。

针对单点和多点 BFF 架构缺陷，我们使用了 API Gateway，把安全，限流，熔断，统一认证这些跨横切面的逻辑都上放到了 API Gateway，由 API Gateway 来实现，把 BFF 当成了一个基础设施，实现关注点，专门对业务进行兼容，不对安全，限流，熔断，统一认证这些跨横切面的逻辑这些通用逻辑维护。

### 微服务设计

#### 角色服务
Web UI 的前端后端（BFF）服务，有两种角色：
- 普通用户（**customer**）
- 商户（**Seller**）
- 管理员（**Admin**）

##### 商户角色
商户角色的权限：
- 处理订单
- 查看销售数据
- 与用户沟通（回复评价）
- 发布商品
  - 品牌
  - 标签
  - 描述
  - 分类
  - 价格
  - 商品名称
  - 库存
- 修改商品
  - 品牌
  - 标签
  - 描述
  - 分类
  - 价格
  - 商品名称
  - 库存
- 删除商品

##### 管理员角色
- 管理平台用户（普通用户、商家）
- 审核商品、处理违规行为
- 配置平台参数（如运费模板、支付方式）
- 查看平台整体数据（如交易额、用户增长）
- 增删改查用户
- 可以访问所有数据和功能

#### 数据分析服务（Analytics Service）
- **职责**：
  - 收集和分析用户行为数据
  - 生成销售报表、用户画像
  - 提供数据可视化（如 Dashboard）
- **权限**：
  - 商家可以查看自己店铺的数据，管理员可以查看全平台数据

#### 评价服务（Review Service）
- **职责**：
  - 管理用户对商品的评价
  - 提供评价展示、评分统计
  - 处理评价举报
- **权限**：
  - 用户可以发布评价，商家可以回复

#### 通知服务（Notification Service）
- **职责**：
  - 发送短信、邮件、站内信通知
  - 管理通知模板
  - 处理通知队列
- **权限**：
  - 所有用户都可以接收通知

#### 搜索服务（Search Service）
- **职责**：
  - 提供商品、订单、用户的搜索功能
  - 支持全文检索、模糊搜索
  - 集成搜索引擎（如 Elasticsearch）
- **权限**：
  - 用户只能搜索公开数据，管理员可以搜索所有数据

#### 用户服务
用户微服务包含以下的服务：
- 地址服务，存储在用户的数据库的地址表中
- 银行卡服务，存储在用户的数据库的银行卡表中

##### 功能
- 创建用户
  - 对密码进行加密存储到数据库
- 登录
- 用户登出
  - 删除用户 Token
- 删除用户
  - 当用户选择注销时，根据中国《个人信息保护法》及相关法律法规的要求，我们会对其进行匿名化处理，用户注销之后我们把用户相关的信息全部匿名化，把相关的操作，例如商品记录等一并删除
- 更新用户
- 获取用户身份信息

#### 商品服务
- 创建商品
  - 商户可以创建自己的商品
- 修改商品信息
  - 商家可以修改自己的商品
- 删除商品
  - 商家可以修改自己的商品
- 查询商品信息（单个商品、批量商品）
  - 用户与商户均可以查询商品

#### 购物车服务
- 创建购物车
- 清空购物车
- 获取购物车信息
- 订单定时取消（高级）
- 修改订单信息（可选）
- 创建订单

#### 订单服务
- 创建订单
  - 用户只能查看和管理自己的订单，商户只能管理自己店铺的订单。管理订单状态（如待支付、已发货、已完成）
- 修改订单信息
  - 用户只能查看和管理自己的订单，商户只能管理自己店铺的订单。管理订单状态（如待支付、已发货、已完成）
- 订单定时取消

#### 结算
- 订单结算

#### 支付
- 取消支付
- 定时取消支付
- 支付
  - 用户只能发起支付，商户和管理员可以查看支付记录

#### AI 大模型
- 订单查询
- 模拟自动下单

## 技术栈

我们在选择技术栈和工具时的主要考量：
1. 是否易于学习
2. 是否使用 Go 作为主要的编程语言编写的
3. 文档是否友好
4. 是否易于部署
5. 是否支持分布式
6. 是否契合微服务架构

我们在选型 `微服务框架` 时，我们当时是三个 Kit: B 站开源的 `Kratos` 和 `Go-Zero` 和字节的 `CloudWeGo`。

我们在快速体验它们之后的感受：
- 学习成本：字节的 `CloudWeGo`文档比较好，`Go-Zero`其次，Kratos 会难一点上手，`CloudWeGo`最容易上手
- CLI 工具：它们都各自提供了对应的脚手架，kratos 的体验最好，生成快
- 相关工具链：`Go-Zero` 魔改 go struct 的 API 定义，徒增学习成本和框架绑定成本，放弃

中间件：参考了 CNCF 优秀的开源项目，选型没有太大的争议。

我们团队经过多次讨论最终选择了字节的 `CloudWeGo` 技术栈作为我们的基础 Kit。

基于 **Kitex + Hertz + Protobuf + Citus + Dragonfly + RocketMQ** 技术实现的抖音电商微服务项目。

项目整体是整洁模型和`DDD`的微服务架构设计思想，通过`依赖注入`减少全局变量污染，具有高内聚、低耦合、`关注点分离`等特点，方便团队各个成员进行良好的团队协作，规范化模块开发，联调测试，**自动化部署**微服务在 **线上 Kubernetes 集群**中。

### 工具下载
- https://github.com/grpc-ecosystem/grpc-gateway/releases

### 线上体验
- [前端](https://node8.apikv.com:30020/)
- 后端 [API 文档](https://app.apifox.com/invite/project?token=3j_9favhOghU57nmpYK8n)
- 注册中心
- 配置中心

## Language
- Golang（主要开发语言）
- Protobuf(IDL)

### 容器编排
我们支持两种方式对容器进行编排：
- Kubernetes（**自动部署，水平扩缩，回滚**，**服务发现和负载均衡，自我修复，密钥与配置管理**）
- ArgoCD: 灰度发布
- Dockerfile + docker-compose（构建应用，单机测试，一键部署）

### Gateway
- Cilium Gateway(L7, Gateway API, Network policy)

### 流传输
- RocketMQ（延迟消息队列）
- Dragonfly（简单任务队列）

### 可观测性相关
- Opentelemetry（中间件，采集 Merits/Logs/Traces 并转发到对应的后端）
- Loki（日志存储）
- Grafana（可视化）
- Jaeger（链路追踪，可视化）
- VictoriaMetrics（时序存储）

### CI/CD
我们提供两种 CI/CD 自动化部署：
- GitHub Actions + Kustomize + ArgoCD + Kubernetes
- GitLab CI

### RPC
- gRPC

### 服务发现/注册
- Consul（健康检查）

### 配置
#### 线上
- Consul Config
- Kubernetes ConfigMap
- Kubernetes Secrets

#### 开发
- Local env
- File env

### 数据库
- 关系型和非结构化数据库：Citus
- 时序数据库：VictoriaMetrics

### 日志
- Opentelemetry logs
- Kitex logs

### KV
- Dragonfly

## 规范

### Git 规范

#### 推荐的默认配置

##### 全局配置
- 默认的主分支名称为 main，统一项目的主分支名
```bash
git config --global init.defaultBranch main 
```

- 提交时转换为 LF，检出时不转换
```bash
git config --global core.autocrlf input
```

- 避免不必要的合并提交
当你执行 `git pull` 时，Git 默认会执行 `fetch` 和 `merge` 操作，这可能会导致本地分支的合并冲突或其他变化
```bash
git config pull.rebase true
```

验证设置
```bash
git config pull.rebase
```

##### 项目配置
1. 把当前仓库的 git 配置文件移动到。git 目录的 config 中，在当前项目中生效
```bash
cp .gitconfig .git/config
```

#### 提交规范
语法：

<type>: (<emoji>) subject

- type: 本次提交内容的类型，例如：fix: 修复某个错误，feat: 添加某个接口
- emoji: 可选，emoji，github 内置了一系列的 [emoji](https://gitmoji.dev/)，使用可以参考 [lingxd/gitemoji](https://github.com/lingxd/gitemoji)
- subject: 从动词开始（比如"fix"），每行 50 个字符

1. 提交的`类型`格式为下列类型：
   1. feat: 特性
   2. fix: 修复
