## 项目开发

项目是前端后端分离,仓库地址是 https://github.com/sunmery/tiktok-e-commence 前端以 **submodule** 方式链接到单独的前端项目, 根据需要来决定是否也拉取前端项目仓库代码:

1. 拉取, 包含前端项目. 不需要前端删除`` `-recurse-submodules` `` 即可


```Bash
git clone --recurse-submodules git@github.com:sunmery/tiktok-e-commence.git
```

目录结构:


```
.
├── LICENSE
├── Makefile CLI 生成命令
├── README.md
├── api 应用proto目录
│  ├── cart  服务 proto
│  │  └── v1  proto版本
│  │      └── service.proto
│  ├── checkout
│  │  └── v1
│  │      └── service.proto
│  ├── order
│  │  └── v1
│  │      └── service.proto
│  ├── payment
│  │  └── v1
│  │      └── service.proto
│  ├── product
│  │  └── v1
│  │      └── service.proto
│  └── user
│      └── v1
│          ├── error_reason.pb.go
│          ├── error_reason.proto
│          ├── error_reason.swagger.json
│          ├── service.pb.go
│          ├── service.proto
│          ├── service.swagger.json
│          ├── service_grpc.pb.go
│          └── user
│              └── v1
│                  ├── error_reason.pb.go
│                  ├── service.pb.go
│                  └── service_grpc.pb.go
├── application 微服务目录
│  ├── cart
│  ├── checkout
│   ├── order
│   ├── payment
│  ├── product
│  ├── user
├── third_party 第三方 proto目录
```

## 架构

南北向为 Client 发送请求到 API Gateway， 由 API Gateway 分发流量到对应的微服务

东西向为 微服务之间通过 gRPC 来交互

暂时无法在飞书文档外展示此内容

针对单点和多点 BFF 架构缺陷， 我们使用了 API Gateway， 把安全， 限流， 熔断， 统一认证这些跨横切面的逻辑都上放到了 API Gateway， 由 API Gatewa 来实现， 把 BFF 当成了一个基础设施， 实现关注点， 专门对业务进行兼容， 不对安全， 限流， 熔断， 统一认证这些跨横切面的逻辑这些通用逻辑维护

### 微服务设计

#### 角色服务

Web UI 的前端后端 （BFF） 服务， 有两种角色：

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

- 管理平台用户（普通用户、商家）。

- 审核商品、处理违规行为。

- 配置平台参数（如运费模板、支付方式）。

- 查看平台整体数据（如交易额、用户增长）

- 增删改查用户

- 可以访问所有数据和功能


#### **数据分析服务（Analytics Service）**

- **职责**：

    - 收集和分析用户行为数据。

    - 生成销售报表、用户画像。

    - 提供数据可视化（如 Dashboard）

- **权限**：

    - 商家可以查看自己店铺的数据，管理员可以查看全平台数据。


#### **评价服务（Review Service）**

- **职责**：

    - 管理用户对商品的评价。

    - 提供评价展示、评分统计。

    - 处理评价举报。

- **权限**：

    - 用户可以发布评价，商家可以回复。


#### **通知服务（Notification Service）**

- **职责**：

    - 发送短信、邮件、站内信通知。

    - 管理通知模板。

    - 处理通知队列。

- **权限**：

    - 所有用户都可以接收通知。


#### **搜索服务（Search Service）**

- **职责**：

    - 提供商品、订单、用户的搜索功能。

    - 支持全文检索、模糊搜索。

    - 集成搜索引擎（如 Elasticsearch）。

- **权限**：

    - 用户只能搜索公开数据，管理员可以搜索所有数据。


#### 用户服务

用户微服务包含以下的服务：

- 地址服务， 存储在用户的数据库的地址表中

- 银行卡服务， 存储在用户的数据库的银行卡表中


##### 功能

- 创建用户

    - 对密码进行加密存储到数据库

- 登录

- 用户登出

    - 删除用户 Token

- 删除用户

    - 当用户选择注销时，根据中国《个人信息保护法》及相关法律法规的要求， 我们会对其进行匿名化处理， 用户注销之后我们把用户相关的信息全部匿名化， 把相关的操作， 例如商品记录等一并删除

- 更新用户

- 获取用户身份信息


##### SQL 设计



#### 商品服务

- 创建商品

    - 商户可以创建自己的商品

- 修改商品信息

    - 商家可以修改自己的商品

- 删除商品

    - 商家可以修改自己的商品

- 查询商品信息（单个商品、批量商品）

    - 用户与商户均可以查询商品


#### **购物车服务**

- 创建购物车

- 清空购物车

- 获取购物车信息

- 订单定时取消（高级）

- 修改订单信息（可选）

- 创建订单


#### **订单服务**

- 创建订单

    - 用户只能查看和管理自己的订单，商户只能管理自己店铺的订单。管理订单状态（如待支付、已发货、已完成）

- 修改订单信息

    - 用户只能查看和管理自己的订单，商户只能管理自己店铺的订单。管理订单状态（如待支付、已发货、已完成）

- 订单定时取消


#### **结算**

- 订单结算


**支付**

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


我们在选型 `微服务框架` 时， 我们当时是三个 Kit: B 站开源的 `Kratos` 和 `Go-Zero` 和字节的 `CloudWeGo`

我们在快速体验它们之后的感受：

- 学习成本： 字节的 `CloudWeGo`文档比较好， `Go-Zero`其次， Kratos 会难一点上手， `CloudWeGo`最容易上手

- CLI 工具： 它们都各自提供了对应的脚手架， kratos 的体验最好， 生成快

- 相关工具链： `Go-Zero` 魔改 go struct 的 API 定义， 徒增学习成本和框架绑定成本， 放弃。




中间件： 参考了 CNCF 优秀的开源项目， 选型没有太大的争议。



我们团队经过多次讨论最终选择了字节的 `CloudWeGo` 技术栈作为我们的基础 Kit。

基于 **Kitex + Hertz + Protobuf + Citus + Dragonfly + RocketMQ** 技术实现的抖音电商微服务项目



项目整体是整洁模型和`DDD`的微服务架构设计思想， 通过`依赖注入`减少全局变量污染， 具有高内聚、低耦合、`关注点分离`等特点， 方便团队各个成员进行良好的团队协作， 规范化模块开发， 联调测试， **自动化部署**微服务在 **线上 Kubernetes 集群**中



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

我们支持两种方式对容器进行编排

- Kubernetes（**自动部署， 水平扩缩， 回滚**， **服务发现和负载均衡， 自我修复，密钥与配置管理**）

- ArgoCD: 灰度发布

- Dockerfile + docker-compose （构建应用， 单机测试， 一键部署）


### Gateway

- Cilium Gateway(L7, Gateway API, Network policy)


### 流传输

- RocketMQ（延迟消息队列）

- Dragonfly（简单任务队列）


### 可观测性相关

- Opentelemetry（中间件， 采集 Merits/Logs/Traces 并转发到对应的后端）

- Loki（日志存储）

- Grafana（可视化）

- Jaeger（链路追踪， 可视化）

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

- 关系型和非结构化数据库： Citus

- 时序数据库： VictoriaMetrics


### 日志

- Opentelemetry logs

- Kitex logs


### KV

- Dragonfly


## 规范

### Git 规范

#### 推荐的默认配置

##### 全局配置

- 默认的主分支名称为 main， 统一项目的主分支名


```Bash
git config --global init.defaultBranch main 
```

- 提交时转换为 LF，检出时不转换


```Bash
git config --global core.autocrlf input
```

- 避免不必要的合并提交


当你执行 `git pull` 时，Git 默认会执行 `fetch` 和 `merge` 操作，这可能会导致本地分支的合并冲突或其他变化

```Bash
git config pull.rebase true
```

验证设置

```Bash
git config pull.rebase
```

##### 项目配置

1. 把当前仓库的 git 配置文件移动到。git 目录的 config 中， 在当前项目中生效


```Bash
cp .gitconfig .git/config
```

#### 提交规范

2. 提交规范


语法：

<type>: (<emoji>) subject

- type: 本次提交内容的类型， 例如： fix: 修复某个错误， feat: 添加某个接口

- emoji: 可选，emoji， github 内置了一系列的 [emoji](https://gitmoji.dev/) ，使用可以参考 [lingxd/gitemoji](https://github.com/lingxd/gitemoji)

- subject: 从动词开始（比如“fix”），每行 50 个字符


1. 提交的`类型`格式为下列类型：

    1. feat: 特性

    2. fix: 修复

    3. perf: 优化

    4. style: 样式

    5. docs: 文档

    6. test: 测试

    7. refactor: 重构代码

    8. ci: 工作流

    9. chore: 其它

    10. revert: 回滚

    11. types: 类型

    12. release: 新版本


例如：

- feat: 新特性

- fix: 修复错误


2. 提交的`内容`， 例如： 修复系统过于流畅的 Bug


示例：

- fix: 提交的类型

- 修复系统过于流畅的 Bug: 提交的内容


```Bash
git commit -m "fix: 修复系统过于流畅的Bug"
```

#### Git 最佳实践

不要直接把代码推送到主分支， 先创建一个新的分支， 在经过 review 之后， 确认无误之后在合并代码到主分支

```Bash
git checkout -b dev

git add .
git commit -m "feat: feat"
git push main dev
```

### 数据库规范

#### **开发与测试**

- 在开发和测试环境中使用与生产环境相同的数据库配置。

- 使用迁移工具 go-migrate 管理数据库 schema 变更


**主键**：

- 每个表必须有一个主键（`PRIMARY KEY`）。

- 主键使用 `SERIAL` 自增类型


**外键**：不使用

**数据类型**：

- 选择合适的数据类型（例如 `INTEGER`、`VARCHAR`、`TEXT`、`TIMESTAMP` 等）。

- 避免使用过大的数据类型（例如 `TEXT` 代替 `VARCHAR(255)` 除非确实需要）。


**约束**：

- 使用 `NOT NULL` 约束确保字段不为空。

- 使用 `UNIQUE` 约束确保字段值唯一。

- 使用 `CHECK` 约束确保字段值符合特定条件。


#### 索引

- **索引类型**：

    - 主键字段自动创建索引。

    - 对经常用于查询条件的字段创建索引（例如 `WHERE`、`JOIN`、`ORDER BY` 中的字段）。

    - 对频繁查询的组合条件创建复合索引。

- **避免过度索引**：

    - 索引会增加写操作的开销，因此不要为所有字段都创建索引。

- **索引维护**：

    - 定期使用 `REINDEX` 或 `VACUUM` 维护索引。


#### **SQL 编写**

- **SQL 格式化**：

    - 使用一致的缩进和换行，使 SQL 语句易于阅读。

    - 关键字使用大写（例如 `SELECT`、`FROM`、`WHERE`）。

- **避免** `SELECT *`：

    - 明确指定需要查询的字段，避免不必要的数据传输。

- **使用参数化查询**：

    - 防止 SQL 注入攻击，同时提高查询性能。

- **事务管理**：

    - 对写操作（`INSERT`、`UPDATE`、`DELETE`）使用事务（`BEGIN`、`COMMIT`、`ROLLBACK`）。

    - 避免长时间运行的事务，以减少锁争用。


---

#### **性能优化**

##### 查询优化

- 使用 `EXPLAIN` 和 `EXPLAIN ANALYZE` 分析查询性能。

- 避免在 `WHERE` 子句中对字段进行函数操作（例如 `WHERE LOWER(name) = 'john'`），这会导致索引失效。

- 使用 `LIMIT` 和 `OFFSET` 分页时，确保查询效率。


##### 连接池

- 使用连接池（如 `pgbouncer`）来管理数据库连接，减少连接开销。


##### 分区表

- 对于大表，使用分区表（`PARTITION BY`）来提高查询性能。


#### 维护

- 定期运行 `VACUUM` 和 `ANALYZE` 来清理和优化数据库。

- 使用 `REINDEX` 重建索引。

- 定期检查数据库的健康状态（如锁争用、连接数、磁盘使用率等）。


---

#### **安全性**

##### 访问控制

- 使用角色（`ROLE`）和权限（`GRANT`、`REVOKE`）管理数据库访问。

- 遵循最小权限原则，只授予用户必要的权限。

- 避免使用超级用户（`postgres`）进行日常操作。


##### 数据加密

- 使用 SSL/TLS 加密数据库连接。

- 对敏感数据（如密码）使用加密存储（由应用层面实现）。


#### **扩展**

##### 扩展性

- 使用读写分离（通过逻辑复制或流复制）来提高读性能。

- 使用分片（Sharding）来水平扩展数据库。


#### **监控与日志**

##### 日志配置

- 启用慢查询日志（`log_min_duration_statement`）以捕获性能问题。

- 配置日志轮转（`log_rotation_age` 和 `log_rotation_size`）以避免日志文件过大。


##### 监控

- 定期检查数据库的健康状态（如锁争用、连接数、磁盘使用率等）。


### Docker Compose

1. 在编写 Docker Compose 文件时， 文件名称不再需要是`docker-compose.yaml`， 而是使用`compose.yaml`， Docker 会识别到

2. 在编写`compose.yaml`时， 如果有多个 service， 使用空行来区分， 例如：


```YAML
services:

  proxy:
    image: nginx
    volumes:
      - type: bind
        source: ./proxy/nginx.conf
        target: /etc/nginx/conf.d/default.conf
        read_only: true
    ports:
      - 80:80
    depends_on:
      - backend

  backend:
    build:
      context: backend
      target: builder
```

3. 编写`compose.yaml`时， 不需要添加`version`字段， Compose 不用于选择确切的架构来验证 Compose 文件，而是首选实现时的最新架构

4. 不应该在 Dockerfile 中写固定的端口， 最佳实践是使用 `ARG` 来定义一个默认的值， 在 `build` 构建时可以替换为其它值， 在 build 时即使不添加该变量也可以使用默认的值

5. 在选择 restart 的时候不要选择`always`， 容器出错的情况很大概率是配置出错了， 不要让它一直重启， 浪费资源。 应该使用`restart: on-failure:5`， 让它重试 5 次


### 微服务编写规范

#### Protobuf 编写规范

##### 标准文件格式

- 保持每行长度不超过 **80 个字符**。

- 使用 **2 个空格** 作为缩进。

- 优先使用 **双引号** 表示字符串。


##### 使用合理的参数校验

更多的示例可以在 [proto-gen-validate](https://github.com/envoyproxy/protoc-gen-validate) 文档中查看

```ProtoBuf
syntax="proto3";

package user.service.v1;

import "validate/validate.proto"; // 导入参数校验包

option go_package="/user";

message RegisterReq {
  string email = 1;
  string password = 2 [(validate.rules).string.min_len = 6]; // 密码不少于 6 个字符
  string confirm_password = 3 [(validate.rules).string.min_len =  6]; // 重复密码不少于 6 个字符
}
```

##### 文件结构

1. 文件命名


- 文件名应使用 **小写蛇形命名法**（`lower_snake_case`），并以 `.proto` 结尾。 例如：`song_service.proto`。


2. 文件内容顺序


文件内容应按以下顺序排列：

1. **许可证头**（如果适用）

2. **文件概述**

3. **语法声明**（`syntax`）

4. **包声明**（`package`）

5. **导入声明**（`import`，按字母顺序排序）

6. **文件选项**（`option`）

7. **其他内容**（消息、枚举、服务等）


---

##### 包命名

- 包名应使用 **小写字母**。

- 包名应基于项目名称，并可能基于文件路径以确保唯一性。 例如：`package music_service;`


---

##### 消息和字段命名

8. 消息命名


- 使用 **帕斯卡命名法**（`PascalCase`），首字母大写。 例如：`SongServerRequest`。

- 缩写应作为单个单词大写，例如：`GetDnsRequest` 而不是 `GetDNSRequest`。


9. 字段命名


- 使用 **小写蛇形命名法**（`lower_snake_case`）。 例如：`song_name`。

- 如果字段名包含数字，数字应出现在字母之后，而不是下划线之后。 例如：使用 `song_name1` 而不是 `song_name_1`。


示例：

```ProtoBuf
message SongServerRequest {
  optional string song_name = 1;
}
```

---

##### 重复字段

- 对于重复字段，使用 **复数形式** 的字段名。


示例：

```ProtoBuf
repeated string keys = 1;
repeated MyMessage accounts = 17;
```

---

##### 枚举

10. 枚举类型命名


- 使用 **帕斯卡命名法**（`PascalCase`），首字母大写。


11. 枚举值命名


- 使用 **大写蛇形命名法**（`CAPITALS_WITH_UNDERSCORES`）。

- 每个枚举值应以 **分号** 结尾，而不是逗号。


12. 未指定值


- 枚举的零值应使用 `UNSPECIFIED` 作为后缀。 例如：`FOO_BAR_UNSPECIFIED = 0;`


示例：

```ProtoBuf
enum FooBar {
  FOO_BAR_UNSPECIFIED = 0;
  FOO_BAR_FIRST_VALUE = 1;
  FOO_BAR_SECOND_VALUE = 2;
}
```

13. 避免命名冲突


- 如果枚举是顶级枚举，建议为每个枚举值添加枚举名前缀，或将枚举嵌套在消息中。


---

##### Service

14. 服务命名


- 使用 **帕斯卡命名法**（`PascalCase`），首字母大写。

- 在 Service 起名时在末尾添加`Service`， 使用`` `FooService` `` 而不是`` `Foo` `` 的 Service 名称


15. RPC 方法命名


- 使用 **帕斯卡命名法**（`PascalCase`），首字母大写。


示例：

```ProtoBuf
service FooService {
  rpc GetSomething(GetSomethingRequest) returns (GetSomethingResponse);
  rpc ListSomething(ListSomethingRequest) returns (ListSomethingResponse);
}
```

---

##### 其他最佳实践

16. 为每个方法创建唯一的 Proto


- 每个 RPC 方法应有独立的请求和响应消息，避免复用。


17. 不要在顶层请求或响应消息中包含原始类型


- 请求和响应消息应使用自定义消息类型，而不是直接使用原始类型（如 `string`、`int32` 等）。


18. 将消息定义放在单独的文件中


- 每个消息类型应定义在独立的 `.proto` 文件中，以提高模块化和可维护性。


一个完整的示例：

```ProtoBuf
// 许可证头
// 文件概述：定义音乐服务的 Protocol Buffers。

syntax = "proto3";

package music_service;

import "google/protobuf/empty.proto";

option java_multiple_files = true;
option java_package = "com.example.musicservice";

// 请求消息
message SongServerRequest {
  optional string song_name = 1;
}

// 响应消息
message SongServerResponse {
  repeated string song_names = 1;
}

// 枚举
enum SongType {
  SONG_TYPE_UNSPECIFIED = 0;
  SONG_TYPE_POP = 1;
  SONG_TYPE_ROCK = 2;
}

// 服务定义
service MusicService {
  rpc GetSong(SongServerRequest) returns (SongServerResponse);
  rpc ListSongs(google.protobuf.Empty) returns (SongServerResponse);
}
```

#### 格式化

使用更严格的格式化工具 `gofumpt` 来对项目进行代码格式化

```Bash
go install mvdan.cc/gofumpt@latest
```

##### Goland

![](https://i0jecrneytu.feishu.cn/space/api/box/stream/download/asynccode/?code=YTgwYjBkZTk4ODdhZWU1Mjk5OWNiMTUwNzY4MjdhNzRfQ1dIcU1USkg4VDJtMnJacDBvdVFBTnZBRWY1Qzh1b1pfVG9rZW46VGZGNmJDWGVMb2Z5N3l4eVJsRWNwV1VibmllXzE3Mzc5ODEyNTA6MTczNzk4NDg1MF9WNA)

##### VSCode

`settings.json`

```JSON
{
    "go.useLanguageServer": true,
    "gopls": {
        "formatting.gofumpt": true,
    },
}
```

#### 注册微服务

##### **结合 Google 和行业最佳实践的推荐命名方式**

根据 Google 的 API 设计指南和行业实践，推荐以下命名结构：

- `[项目名]-[功能或模块]-[服务类型]-[版本]`

    - **示例**：`ecommerce-order-service-v1`

    - **优点**：

        - 清晰描述服务的功能（如 `order`）和类型（如 `service`），便于理解和管理。

        - 包含版本号（如 `v1`），支持语义化版本控制，便于后续迭代和兼容性管理18。

        - 适合大多数场景，兼顾简洁性和可扩展性。


例如：

1. `ecommerce-user-service-v1`

2. `ecommerce-user-address-v1`


```Go
package main

var (
        Name = "ecommerce-order-service-v1"
        Version = "1.0.0"
)
```

### Dockerfile

```Dockerfile
# syntax=docker/dockerfile:1
# https://docs.docker.com/go/dockerfile-reference/

# 定义基础镜像的 Golang 版本
ARG GOIMAGE=golang:1.23.3-alpine3.20

FROM --platform=$BUILDPLATFORM ${GOIMAGE} AS build
WORKDIR /src

# 版本号
ARG VERSION=latest

# Go的环境变量, 例如alpine镜像不内置gcc,则关闭CGO很有效
ARG GOOS=linux
ARG GOARCH=amd64
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
    go build -o /bin ./...
   # 带版本的形式: go build -ldflags="-X main.Version=${VERSION}" -o /bin/main .
   # 多个服务的形式: go build -o /bin/ ./...
COPY ./configs /bin/configs
FROM alpine:latest AS final
# 从构建阶段复制编译好的 Go 应用程序到运行阶段
COPY --from=build /bin/product /bin/
COPY --from=build /bin/configs /bin/configs

# 用户进程ID
ARG UID=10001

# 后端程序的HTTP/gRPC端口
ARG HTTP_PORT=30001
ARG GRPC_PORT=30002

# 修改镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 安装应用运行必需的系统证书和时区数据包
# RUN --mount=type=cache,target=/var/cache/apk \
#    apk --update add ca-certificates tzdata && update-ca-certificates

# RUN chmod 1777 /tmp
# # 创建一个非特权用户来运行应用，增强容器安全性
# RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "${UID}" appuser

# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# USER appuser

# 指定容器对外暴露的端口号
EXPOSE $HTTP_PORT $GRPC_PORT

VOLUME /data/conf

# 设置容器启动时执行的命令
CMD ["/bin/product", "-conf", "/data/conf"]
```

#### 构建

##### 单平台架构

构建 Docker 所属的当前平台与架构的二进制文件， 进到当前的 backend 目录 VERSION=dev REPOSITORY="team/backend" Docker 容器在 Linux 内核上运行，即便是在 macOS 或 Windows 环境中。 如果使用 docker 构建时传递 GOOS=darwin 会导致构建的二进制文件不兼容于 Linux 环境，从而出现 exec format error 所以在使用 docker 构建时的目标平台的 GOOS 应该为 linux，而非 darwin。

```Bash
docker build . \
  --progress=plain \
  -t dev/app:dev \
  --build-arg CGOENABLED=0 \
  --build-arg GOIMAGE=golang:1.23.3-alpine3.20 \
  --build-arg GOOS=linux \
  --build-arg GOARCH=arm64 \
  --build-arg VERSION=$VERSION \
  --build-arg HTTP_PORT=30001 \
  --build-arg GRPC_PORT=30002
```

##### 构建多平台架构

构建多架构的二进制文件， 需要在 Docker Desktop 启用 containerd 映像存储

```Bash
https://docs.docker.com/desktop/containerd/#enable-the-containerd-image-store
VERSION=v1.0.0
REPOSITORY="tiktok/products"
GOOS=linux
GOARCH=amd64
HTTP_PORT=30001
GRPC_PORT=30002
PLATFORM_1=linux/amd64
docker buildx build . \
--progress=plain \
-t $REPOSITORY:$VERSION \
--build-arg CGOENABLED=0 \
--build-arg GOIMAGE=golang:1.23.3-alpine3.20 \
--build-arg VERSION=$VERSION \
--build-arg HTTP_PORT=$PORT \
--build-arg GRPC_PORT=$PORT \
--build-arg GOOS=$GOOS \
--build-arg GOARCH=$GOARCH \
--platform $PLATFORM_1 \
--load
```

#### 推送

```Bash
register="ccr.ccs.tencentyun.com"
docker tag $REPOSITORY:$VERSION $register/$REPOSITORY:$VERSION
docker push $register/$REPOSITORY:$VERSION
```

#### 拉取

```Bash
GOOS=linux
GOARCH=amd64
docker pull $register/$REPOSITORY:$VERSION --platform $GOOS/GOARCH
```

#### 运行

##### compose

编写使用如下的 compose 文件， 只需要一条命令即可启动

```YAML
services:

  user:
    pull_policy: build
    build:
      context: user
      dockerfile: Dockerfile
      args:
        HTTP_PORT: 30001
        GRPC_PORT: 30002
        VERSION: v1.0.0
    platform: linux/amd64
    ports:
      - "30001:30001"
      - "30002:30002"
    container_name: user
    restart: always
    networks:
      - app-network
    environment:
      - DB_SOURCE=postgresql://postgres:postgres@postgres17:5432/tiktok?sslmode=disable
      - REDIS_ADDRESS=redis7:6379
    command:
      - "/bin/user"
      - "-conf"
      - "/bin/configs"
```

选项：

- up: 启动

- -d: 以后台的方式运行

- restart: 重启容器

- down: 关闭容器


```Bash
docker compose up -d
```

##### docker run

这种方式不推荐使用， docker run 的所有配置（如环境变量、端口映射、卷挂载等）都需要通过命令行参数传递，命令会变得非常冗长且难以维护。 使用 compose 所有配置都可以在 docker-compose.yml 文件中定义，便于版本控制和共享，同时也更易于阅读和维护。

```Bash
docker run \
--rm \
-p 30001:30001 \
-p 30002:30002 \
-e GIN_MODE=release \
-e DB_SOURCE="postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" \
$REPOSITORY:$VERSION
```

### 可观测性规范

#### 日志

##### 日志等级

从最高到最低分表是：

- Error: 程序错误导致应用运行异常

- Warning: 当调用 a 函数时遇到错误但可以使用 b 函数来作为降级使用， 例如上传文件到线上的 S3 存储时失败， 可以临时降级调用 b 函数保存到本机

- Info: 输出


#### 指标

采集如下：

- 系统各个指标， 例如 CPU，内存， 网络 I/O

- 用户购物车商品数量


#### 链路追踪

框架默认集成， 暂时没有特别的需求

## 基础设施清单

### 开发阶段

`compose.yaml`示例：

```YAML
services:

  user:
    pull_policy: build
    build:
      context: <微服务目录>
      dockerfile: Dockerfile
      args:
        HTTP_PORT: <HTTP 端口>
        GRPC_PORT: <GRPC_PORT 端口>
        VERSION: dev
    platform: linux/amd64
    ports:
      - "<外部HTTP 端口>:<容器HTTP 端口>"
      - "<外部GRPC 端口>:<容器GRPC 端口>"
    container_name: <微服务名称>
    restart: restart: on-failure:5
    environment:
      - DB_SOURCE=postgresql://postgres:postgres@postgres17:5432/tiktok?sslmode=disable
      - REDIS_ADDRESS=redis7:6379
    command:
      - "/bin/user"
      - "-conf"
      - "/bin/configs"
```

### 预发布阶段

预先创建一个 `Docker Network`

```Bash
docker network create app-network
```

配置说明：

假设 user 微服务与 cart 微服务需要进行预发布， 需要模拟在生产环境中测试， 需要加入数据库和缓存来测试接口，

推荐的做法：

1. 将这些微服务与数据库缓存都放在一个 `Docker Network` 中，通过`service name` 来引用它们


```YAML
services:

  # 微服务 1
  user:
    pull_policy: build
    build:
      context: user
      dockerfile: Dockerfile
      args:
        HTTP_PORT: 30001
        GRPC_PORT: 30002
        VERSION: v1.0.0
    platform: linux/amd64
    ports:
      - "30001:30001"
      - "30002:30002"
    container_name: user
    restart: always
    networks:
      - app-network
    environment:
      - DB_SOURCE=postgresql://postgres:mypass@citus:5432/tiktok?sslmode=disable
      - REDIS_ADDRESS=dragonfly:6379
      - REDIS_PASSWORD=mypass
    command:
      - "/bin/user"
      - "-conf"
      - "/bin/configs"
  
  # 微服务 2
  cart:
    build:
      context: cart
      dockerfile: Dockerfile
      args:
        HTTP_PORT: 30003
        GRPC_PORT: 30004
        VERSION: v1.0.0
    platform: linux/amd64
    ports:
      - "30003:30003"
      - "30004:30004"
    container_name: cart
    restart: always
    networks:
      - app-network
    environment:
      - DB_SOURCE=postgresql://postgres:postgres@citus:5432/tiktok?sslmode=disable
      - REDIS_ADDRESS=dragonfly:6379
      - REDIS_PASSWORD=mypass
    command:
      - "/bin/cart"
      - "-conf"
      - "/bin/configs"
  
  # 数据库
  citus:
    image: citusdata/citus:postgres_15
    container_name: citus
    restart: on-failure:5
    ports:
      - 5432:5432
    environment:
      - POSTGRES_PASSWORD: mypass
  
  # 缓存
  dragonfly:
    image: 'docker.dragonflydb.io/dragonflydb/dragonfly'
    ulimits:
      memlock: -1
    ports:
      - "6379:6379"
    environment:
      - DRAGONFLYDB_USERNAME=default # 设置用户名
      - DRAGONFLYDB_PASSWORD=mypass # 设置密码

networks:
  app-network:
    external: true
```

### 生产环境

单独创建一个清单仓库， 单独维护微服务的清单列表

```Bash
git init
mkdir -pv manifests/kustomize/base
mkdir -pv manifests/kustomize/overlays
```

#### Kustomize 清单

```Bash
cat > manifests/kustomize/base/kustomization.yaml <<EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - deployment.yaml
  - service.yaml
  # - ingress-http.yaml
  # - ingress-grpc.yaml

generatorOptions:
  annotations:
    argocd.argoproj.io/compare-options: IgnoreExtraneous
EOF
```

#### Deployment 清单

示例： 假设我需要创建一个`product`微服务， 其它的微服务套用时只需要把`service`的值`product` 替换为微服务名称和 HTTP/gRPC 端口， 如果不需要 gRPC 或 HTTP， 删掉对应的数组元素即可

```Bash
service=product
http_port=30013
grpc_port=30014
cat > manifests/kustomize/base/deployment.yaml <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  # https://argo-cd.readthedocs.io/en/stable/$service-guide/compare-options/
  # 在某些情况下，您可能希望从应用程序的整体同步状态中排除资源。例如，如果它们是由工具生成的。这可以通过在您想要排除的资源上添加此注释来完成
  annotations:
    argocd.argoproj.io/compare-options: IgnoreExtraneous
  name: e-commence-$service
  labels:
    app: e-commence-$service
spec:
  replicas: 1
  selector:
    matchLabels:
      app: e-commence-$service
  template:
    metadata:
      name: e-commence-$service
      labels:
        app: e-commence-$service
    spec:
      restartPolicy: Always
      containers:
        - name: e-commence-$service
          image: example
          imagePullPolicy: Always
          args:
            - /bin/$service
            - -conf
            - /bin/configs
          ports:
            - containerPort: $http_port
              protocol: TCP
              name: http-server
            - containerPort: $grpc_port
              protocol: TCP
              name: grpc-server
EOF
```

#### Service 清单

Service 同理， 这是一个`product`微服务， 其它的微服务套用时只需要把`product` 这个名称替换为微服务名称和 HTTP/gRPC`containerPort`端口， 如果不需要 gRPC 或 HTTP， 删掉对应的数组元素即可

```Bash
service=product
http_port=30013
grpc_port=30014
cat > manifests/kustomize/base/service.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  name: e-commence-$service-service
spec:
  selector:
    app: e-commence-$service
  ports:
    - name: http
      port: $http_port
      protocol: TCP
      targetPort: $http_port
      nodePort: $http_port
    - name: grpc
      port: $grpc_port
      protocol: TCP
      targetPort: $grpc_port
      nodePort: $grpc_port
  type: NodePort
EOF
```

#### Namespace 清单

直接复制粘贴即可， 不需要修改

```Bash
cat > manifests/kustomize/base/namespace.yaml <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: tiktok
EOF
```

### 数据库

#### 介绍

使用 citus（Postgres 插件）从单个 Postgress 节点到分布式数据库集群， 把表或架构分布在多个节点之间，并行化查询和事务。

Citus 支持基于 Schema 的分片，这允许在许多机器上分发常规数据库 Schema。这种分片方法非常适合典型的微服务架构，其中存储完全由服务拥有，因此不能与其他租户共享相同的架构定义。

基于 Schema 的分片是一种更容易采用的模型，只需创建一个新 Schema 并在您的微服务中设置`search_path`，就可以开始了。

使用 Citus 进行微服务的优势：

- 允许在服务之间水平扩展状态，解决微服务[的主要问题](https://stackoverflow.blog/2020/11/23/the-macro-problem-with-microservices/)之一

- 将业务数据从微服务摄取到常见的分布式表中进行分析

- 通过平衡多台计算机上的服务来高效使用硬件

- 将嘈杂的服务隔离到自己的节点

- 易于理解的分片模型

- 快速采用


#### 清单

此清单的 用户名是： postgres 密码是： tiktok 按需修改

```YAML
services:

  citus:
    image: citusdata/citus:postgres_15
    container_name: citus
    platform: linux/amd64
    restart: on-failure:5
    ports:
      - 5432:5432
    environment:
      POSTGRES_PASSWORD: tiktok
```

启动命令。 清单文件同一个目录的情况下：

```Bash
docker compose up -d
```

如果不是同一个目录的情况下：

```Bash
docker compose -f <compose文件的路径> up -d
```

### 缓存

使用 `Dragonfly` 优势：

1. Redis 和 Memcached API 完全兼容。

2. Dragonfly 在多线程、无共享架构之上实现了新颖的算法和数据结构。

3. 因此，与 Redis 相比，Dragonfly 的性能达到了 25 倍，并且在单个实例上支持数百万 QPS。

4. Dragonfly 与 Redis 生态系统完全兼容，无需更改代码即可实现。

5. 可以使用它来实现分布式锁， 限流，消息队列， 延时队列， 分布式 Session， 统计， 排行等


#### 清单

此清单的 用户名是： default 密码是： tiktok 按需修改

```YAML
services:

  dragonfly:
    image: 'docker.dragonflydb.io/dragonflydb/dragonfly'
    ulimits:
      memlock: -1
    ports:
      - "6379:6379"
    environment:
      - DRAGONFLYDB_USERNAME=default # 设置用户名
      - DRAGONFLYDB_PASSWORD=tiktok # 设置密码
    # For better performance, consider `host` mode instead `port` to avoid docker NAT.
    # `host` mode is NOT currently supported in Swarm Mode.
    # https://docs.docker.com/compose/compose-file/compose-file-v3/#network_mode
    # network_mode: "host"
#     volumes:
#       - dragonflydata:/data
# volumes:
#   dragonflydata:
```

## CI/CD

1. 自动生成 CHANGELOG、创建 GitHub 版本、 和项目的版本升级， 查看文档 https://github.com/googleapis/release-please


下载：

```Shell
npm i release-please -g
```

然后生成 CHANGELOG:

```Shell
release-please --help
```

### 优化过程

Docker build 构建未优化前：

![](https://i0jecrneytu.feishu.cn/space/api/box/stream/download/asynccode/?code=OGZkM2YwNDYxNzBlN2JhNDYwNmY5Nzg2YmJhNThjYWJfb0p4Vk5LTTFaMUp0WW5xSWJqNDRmeUl3MEoxcFJ6eWZfVG9rZW46RXpiamJ2UUMzb3R5Z2F4NERFc2NUY3JsbjFkXzE3Mzc5ODEyNTA6MTczNzk4NDg1MF9WNA)

之前在测试使用这种方法： 在 `matrix` 中定义了三个数组：

- `service`（9 个值）

- `http_port`（9 个值）

- `grpc_port`（9 个值）


这将导致生成的任务数量为： `9 (service) * 9 (http_port) * 9 (grpc_port) = 729`，这个数量已经超出了 GitHub Actions 的默认限制 `256`。

```YAML
service: [ addresses, balances, cart, checkout, credit_cards, order, payment, product, user ]  # 定义要并行执行的服务
http_port: [30015, 30017,  30003, 30005 ,30007, 30009, 30011, 30013, 30001]
grpc_port: [30016, 30018,  30004, 30006 ,30008, 30010, 30012, 30014, 30002]
```

当意识到该问题后， 将代码优化一个数组包含三个值， 优化之后变成 9 个任务数量

```YAML
- { service: addresses, http_port: 30015, grpc_port: 30016 }
- { service: balances, http_port: 30017, grpc_port: 30018 }
- { service: cart, http_port: 30003, grpc_port: 30004 }
- { service: checkout, http_port: 30005, grpc_port: 30006 }
- { service: credit_cards, http_port: 30007, grpc_port: 30008 }
- { service: order, http_port: 30009, grpc_port: 30010 }
- { service: payment, http_port: 30011, grpc_port: 30012 }
- { service: product, http_port: 30013, grpc_port: 30014 }
- { service: user, http_port: 30001, grpc_port: 30002 }
```

![](https://i0jecrneytu.feishu.cn/space/api/box/stream/download/asynccode/?code=N2I5OGIyNTliOTM0ODQzNjExNThlMmRiNDEyNDcwYzNfRm5qeHNldTRlcUZ3TjNFcnZsT1BEb2RsQjhzRkNlWnVfVG9rZW46VXhIcmJ6Z0RJb0M2Ymh4YTlSUGM4U1JsbnpoXzE3Mzc5ODEyNTA6MTczNzk4NDg1MF9WNA)

符合实际的需求，

在一个 ci 文件使用了 matrix 就可以解决多个微服务的部署了， 这很方便， 看起来没有问题。

但在我们部署了多个微服务之后， 我们发现：

1. 每个任务都得等前面的任务全部执行完才继续执行下一个任务

2. 当一个 job 中的任务失败了， 该 job 标记被失败


这显然不是我们需要的

![](https://i0jecrneytu.feishu.cn/space/api/box/stream/download/asynccode/?code=NmU5ZDhiMjY5YTQyN2E0Y2M0M2UxZWQyZTExNTcwNDFfaHBtTDdZcGdhQmdrVUNZNlJPSEVLcUtlS3VUQmZNWUxfVG9rZW46S0NvQ2JqQmg0b0U0V214NFdUYmN5T2Z2bkFoXzE3Mzc5ODEyNTA6MTczNzk4NDg1MF9WNA)

我们重新梳理了需求：

~~例如： user 这个 test 完成了就进入到 user 的 build， 而不是 build 都要等 test 的其它的任务这个 job 全部执行完才继续~~
