---
title: tiktok_e-commence
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# tiktok_e-commence

Base URLs:

# Authentication

- HTTP Authentication, scheme: bearer

# Balance

<a id="opIdBalance_FreezeBalance"></a>

## POST 冻结用户余额

POST /v1/balances/freeze

> Body 请求参数

```json
{
  "userId": "string",
  "orderId": "string",
  "amount": 0.1,
  "currency": "string",
  "idempotencyKey": "string",
  "expectedVersion": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|[v1FreezeBalanceRequest](#schemav1freezebalancerequest)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "freezeId": "string",
  "newVersion": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1FreezeBalanceReply](#schemav1freezebalancereply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<a id="opIdBalance_CancelFreeze"></a>

## POST 取消冻结

POST /v1/balances/freezes/{freezeId}/cancel

> Body 请求参数

```json
{
  "reason": "string",
  "idempotencyKey": "string",
  "expectedVersion": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|freezeId|path|string(int64)| 是 |冻结记录ID|
|Authorization|header|string| 否 |none|
|body|body|[BalanceCancelFreezeBody](#schemabalancecancelfreezebody)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "success": true,
  "newVersion": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CancelFreezeReply](#schemav1cancelfreezereply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<a id="opIdBalance_ConfirmTransfer"></a>

## POST 确认转账（解冻并转给商家）

POST /v1/balances/freezes/{freezeId}/confirm

> Body 请求参数

```json
{
  "merchantId": "string",
  "idempotencyKey": "string",
  "expectedUserVersion": 0,
  "expectedMerchantVersion": 0,
  "paymentAccount": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|freezeId|path|string(int64)| 是 |冻结记录ID|
|Authorization|header|string| 否 |none|
|body|body|[BalanceConfirmTransferBody](#schemabalanceconfirmtransferbody)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "success": true,
  "transactionId": "string",
  "newUserVersion": 0,
  "newMerchantVersion": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ConfirmTransferReply](#schemav1confirmtransferreply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<a id="opIdBalance_GetMerchantBalance"></a>

## GET 获取商家余额

GET /v1/balances/merchants/{merchantId}/balance

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|merchantId|path|string| 是 |UUID as string|
|currency|query|string| 否 |指定币种|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "available": 0.1,
  "frozen": 0.1,
  "currency": "string",
  "version": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1BalanceReply](#schemav1balancereply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<a id="opIdBalance_CreateMerchantBalance"></a>

## PUT 创建商家余额
为用户创建指定币种的初始余额记录 (通常在用户注册或首次涉及该币种时调用)

PUT /v1/balances/merchants/{merchantId}/balance

> Body 请求参数

```json
{
  "currency": "string",
  "initialBalance": 0.1,
  "balancerType": "string",
  "isDefault": true,
  "accountDetails": {}
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|merchantId|path|string| 是 |UUID as string|
|Authorization|header|string| 否 |none|
|body|body|[BalanceCreateMerchantBalanceBody](#schemabalancecreatemerchantbalancebody)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "userId": "string",
  "currency": "string",
  "available": 0.1
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateMerchantBalanceReply](#schemav1createmerchantbalancereply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<a id="opIdBalance_GetUserBalance"></a>

## GET 获取用户余额

GET /v1/balances/users/{userId}/balance

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|userId|path|string| 是 |UUID as string|
|currency|query|string| 否 |指定币种|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "available": 0.1,
  "frozen": 0.1,
  "currency": "string",
  "version": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1BalanceReply](#schemav1balancereply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<a id="opIdBalance_CreateConsumerBalance"></a>

## PUT 创建消费者余额
为用户创建指定币种的初始余额记录 (通常在用户注册或首次涉及该币种时调用)

PUT /v1/balances/users/{userId}/balance

> Body 请求参数

```json
{
  "currency": "string",
  "initialBalance": 0.1,
  "balancerType": "string",
  "isDefault": true,
  "accountDetails": {}
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|userId|path|string| 是 |UUID as string|
|Authorization|header|string| 否 |none|
|body|body|[BalanceCreateConsumerBalanceBody](#schemabalancecreateconsumerbalancebody)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "userId": "string",
  "currency": "string",
  "available": 0.1
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateConsumerBalanceReply](#schemav1createconsumerbalancereply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<a id="opIdBalance_RechargeBalance"></a>

## POST 用户充值

POST /v1/balances/users/{userId}/recharge

> Body 请求参数

```json
{
  "amount": 0.1,
  "currency": "string",
  "externalTransactionId": "string",
  "paymentMethodType": "string",
  "paymentAccount": "string",
  "idempotencyKey": "string",
  "expectedVersion": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|userId|path|string| 是 |用户 UUID (string)|
|Authorization|header|string| 否 |none|
|body|body|[BalanceRechargeBalanceBody](#schemabalancerechargebalancebody)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "success": true,
  "transactionId": "string",
  "newVersion": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RechargeBalanceReply](#schemav1rechargebalancereply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<a id="opIdBalance_WithdrawBalance"></a>

## POST 用户提现

POST /v1/balances/users/{userId}/withdraw

> Body 请求参数

```json
{
  "amount": 0.1,
  "currency": "string",
  "paymentMethodId": "string",
  "idempotencyKey": "string",
  "expectedVersion": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|userId|path|string| 是 |用户 UUID (string)|
|Authorization|header|string| 否 |none|
|body|body|[BalanceWithdrawBalanceBody](#schemabalancewithdrawbalancebody)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "success": true,
  "transactionId": "string",
  "newVersion": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1WithdrawBalanceReply](#schemav1withdrawbalancereply)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

# 数据模型

<h2 id="tocS_protobufAny">protobufAny</h2>

<a id="schemaprotobufany"></a>
<a id="schema_protobufAny"></a>
<a id="tocSprotobufany"></a>
<a id="tocsprotobufany"></a>

```json
{
  "@type": "string",
  "property1": "string",
  "property2": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|**additionalProperties**|string|false|none||none|
|@type|string|false|none||none|

<h2 id="tocS_rpcStatus">rpcStatus</h2>

<a id="schemarpcstatus"></a>
<a id="schema_rpcStatus"></a>
<a id="tocSrpcstatus"></a>
<a id="tocsrpcstatus"></a>

```json
{
  "code": 0,
  "message": "string",
  "details": [
    {
      "@type": "string",
      "property1": "string",
      "property2": "string"
    }
  ]
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|code|integer(int32)|false|none||none|
|message|string|false|none||none|
|details|[[protobufAny](#schemaprotobufany)]|false|none||none|

<h2 id="tocS_BalanceCancelFreezeBody">BalanceCancelFreezeBody</h2>

<a id="schemabalancecancelfreezebody"></a>
<a id="schema_BalanceCancelFreezeBody"></a>
<a id="tocSbalancecancelfreezebody"></a>
<a id="tocsbalancecancelfreezebody"></a>

```json
{
  "reason": "string",
  "idempotencyKey": "string",
  "expectedVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|reason|string|false|none|取消原因|none|
|idempotencyKey|string|false|none|幂等键|none|
|expectedVersion|integer(int32)|false|none|期望的用户余额版本号|none|

<h2 id="tocS_BalanceConfirmTransferBody">BalanceConfirmTransferBody</h2>

<a id="schemabalanceconfirmtransferbody"></a>
<a id="schema_BalanceConfirmTransferBody"></a>
<a id="tocSbalanceconfirmtransferbody"></a>
<a id="tocsbalanceconfirmtransferbody"></a>

```json
{
  "merchantId": "string",
  "idempotencyKey": "string",
  "expectedUserVersion": 0,
  "expectedMerchantVersion": 0,
  "paymentAccount": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|merchantId|string|false|none|merchant_id 可以从 freeze_id 关联的 order_id 推出，或者在这里显式传入|none|
|idempotencyKey|string|false|none|幂等键|none|
|expectedUserVersion|integer(int32)|false|none|期望的用户余额版本号|none|
|expectedMerchantVersion|integer(int32)|false|none|期望的商家余额版本号|none|
|paymentAccount|string|false|none|支付账号|none|

<h2 id="tocS_BalanceRechargeBalanceBody">BalanceRechargeBalanceBody</h2>

<a id="schemabalancerechargebalancebody"></a>
<a id="schema_BalanceRechargeBalanceBody"></a>
<a id="tocSbalancerechargebalancebody"></a>
<a id="tocsbalancerechargebalancebody"></a>

```json
{
  "amount": 0.1,
  "currency": "string",
  "externalTransactionId": "string",
  "paymentMethodType": "string",
  "paymentAccount": "string",
  "idempotencyKey": "string",
  "expectedVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|amount|number(double)|false|none|充值金额|none|
|currency|string|false|none|充值币种|none|
|externalTransactionId|string(int64)|false|none|外部支付流水号 (如支付宝/微信订单号)|none|
|paymentMethodType|string|false|none|支付方式类型 (e.g., "ALIPAY", "WECHAT")|none|
|paymentAccount|string|false|none|支付账号快照|none|
|idempotencyKey|string|false|none|幂等键|none|
|expectedVersion|integer(int32)|false|none|期望的用户余额版本号|none|

<h2 id="tocS_BalanceWithdrawBalanceBody">BalanceWithdrawBalanceBody</h2>

<a id="schemabalancewithdrawbalancebody"></a>
<a id="schema_BalanceWithdrawBalanceBody"></a>
<a id="tocSbalancewithdrawbalancebody"></a>
<a id="tocsbalancewithdrawbalancebody"></a>

```json
{
  "amount": 0.1,
  "currency": "string",
  "paymentMethodId": "string",
  "idempotencyKey": "string",
  "expectedVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|amount|number(double)|false|none|提现金额|none|
|currency|string|false|none|提现币种|none|
|paymentMethodId|string|false|none|用户选择的提现方式ID (BIGINT as string from user_payment_methods)|none|
|idempotencyKey|string|false|none|幂等键|none|
|expectedVersion|integer(int32)|false|none|期望的用户余额版本号|none|

<h2 id="tocS_v1FreezeBalanceRequest">v1FreezeBalanceRequest</h2>

<a id="schemav1freezebalancerequest"></a>
<a id="schema_v1FreezeBalanceRequest"></a>
<a id="tocSv1freezebalancerequest"></a>
<a id="tocsv1freezebalancerequest"></a>

```json
{
  "userId": "string",
  "orderId": "string",
  "amount": 0.1,
  "currency": "string",
  "idempotencyKey": "string",
  "expectedVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|userId|string|false|none|用户 UUID (string)|none|
|orderId|string(int64)|false|none|订单id, 用于关联|none|
|amount|number(double)|false|none|冻结金额|none|
|currency|string|false|none|冻结币种|none|
|idempotencyKey|string|false|none|google.protobuf.Timestamp expires_at = 5; // 冻结过期时间|幂等键 (例如使用 order_id 或单独生成)|
|expectedVersion|integer(int32)|false|none|期望的用户余额版本号 (用于乐观锁)|none|

<h2 id="tocS_v1BalanceReply">v1BalanceReply</h2>

<a id="schemav1balancereply"></a>
<a id="schema_v1BalanceReply"></a>
<a id="tocSv1balancereply"></a>
<a id="tocsv1balancereply"></a>

```json
{
  "available": 0.1,
  "frozen": 0.1,
  "currency": "string",
  "version": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|available|number(double)|false|none|可用余额|none|
|frozen|number(double)|false|none|冻结余额  - 对商家可能总为 0|none|
|currency|string|false|none|返回币种|none|
|version|integer(int32)|false|none|当前版本号 (用于乐观锁)|none|

<h2 id="tocS_v1CancelFreezeReply">v1CancelFreezeReply</h2>

<a id="schemav1cancelfreezereply"></a>
<a id="schema_v1CancelFreezeReply"></a>
<a id="tocSv1cancelfreezereply"></a>
<a id="tocsv1cancelfreezereply"></a>

```json
{
  "success": true,
  "newVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|success|boolean|false|none||none|
|newVersion|integer(int32)|false|none|用户余额新版本号|none|

<h2 id="tocS_v1ConfirmTransferReply">v1ConfirmTransferReply</h2>

<a id="schemav1confirmtransferreply"></a>
<a id="schema_v1ConfirmTransferReply"></a>
<a id="tocSv1confirmtransferreply"></a>
<a id="tocsv1confirmtransferreply"></a>

```json
{
  "success": true,
  "transactionId": "string",
  "newUserVersion": 0,
  "newMerchantVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|success|boolean|false|none||none|
|transactionId|string(int64)|false|none|交易流水ID|none|
|newUserVersion|integer(int32)|false|none|用户余额新版本号|none|
|newMerchantVersion|integer(int32)|false|none|商家余额新版本号|none|

<h2 id="tocS_v1FreezeBalanceReply">v1FreezeBalanceReply</h2>

<a id="schemav1freezebalancereply"></a>
<a id="schema_v1FreezeBalanceReply"></a>
<a id="tocSv1freezebalancereply"></a>
<a id="tocsv1freezebalancereply"></a>

```json
{
  "freezeId": "string",
  "newVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|freezeId|string(int64)|false|none|冻结记录ID|none|
|newVersion|integer(int32)|false|none|操作后用户余额的新版本号|none|

<h2 id="tocS_v1RechargeBalanceReply">v1RechargeBalanceReply</h2>

<a id="schemav1rechargebalancereply"></a>
<a id="schema_v1RechargeBalanceReply"></a>
<a id="tocSv1rechargebalancereply"></a>
<a id="tocsv1rechargebalancereply"></a>

```json
{
  "success": true,
  "transactionId": "string",
  "newVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|success|boolean|false|none||none|
|transactionId|string(int64)|false|none|内部交易流水ID|none|
|newVersion|integer(int32)|false|none|用户余额新版本号|none|

<h2 id="tocS_v1WithdrawBalanceReply">v1WithdrawBalanceReply</h2>

<a id="schemav1withdrawbalancereply"></a>
<a id="schema_v1WithdrawBalanceReply"></a>
<a id="tocSv1withdrawbalancereply"></a>
<a id="tocsv1withdrawbalancereply"></a>

```json
{
  "success": true,
  "transactionId": "string",
  "newVersion": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|success|boolean|false|none||none|
|transactionId|string(int64)|false|none|内部交易流水ID  - 初始状态可能是 PENDING|none|
|newVersion|integer(int32)|false|none|用户余额新版本号|none|

<h2 id="tocS_BalanceCreateConsumerBalanceBody">BalanceCreateConsumerBalanceBody</h2>

<a id="schemabalancecreateconsumerbalancebody"></a>
<a id="schema_BalanceCreateConsumerBalanceBody"></a>
<a id="tocSbalancecreateconsumerbalancebody"></a>
<a id="tocsbalancecreateconsumerbalancebody"></a>

```json
{
  "currency": "string",
  "initialBalance": 0.1,
  "balancerType": "string",
  "isDefault": true,
  "accountDetails": {}
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|currency|string|false|none|指定币种|none|
|initialBalance|number(double)|false|none|初始余额|none|
|balancerType|string|false|none|余额类型|none|
|isDefault|boolean|false|none|是否默认余额|none|
|accountDetails|object|false|none|账户详情 (JSON 格式)|none|

<h2 id="tocS_BalanceCreateMerchantBalanceBody">BalanceCreateMerchantBalanceBody</h2>

<a id="schemabalancecreatemerchantbalancebody"></a>
<a id="schema_BalanceCreateMerchantBalanceBody"></a>
<a id="tocSbalancecreatemerchantbalancebody"></a>
<a id="tocsbalancecreatemerchantbalancebody"></a>

```json
{
  "currency": "string",
  "initialBalance": 0.1,
  "balancerType": "string",
  "isDefault": true,
  "accountDetails": {}
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|currency|string|false|none|指定币种|none|
|initialBalance|number(double)|false|none|初始余额|none|
|balancerType|string|false|none|余额类型|none|
|isDefault|boolean|false|none|是否默认余额|none|
|accountDetails|object|false|none|账户详情 (JSON 格式)|none|

<h2 id="tocS_v1CreateConsumerBalanceReply">v1CreateConsumerBalanceReply</h2>

<a id="schemav1createconsumerbalancereply"></a>
<a id="schema_v1CreateConsumerBalanceReply"></a>
<a id="tocSv1createconsumerbalancereply"></a>
<a id="tocsv1createconsumerbalancereply"></a>

```json
{
  "userId": "string",
  "currency": "string",
  "available": 0.1
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|userId|string|false|none||none|
|currency|string|false|none||none|
|available|number(double)|false|none||none|

<h2 id="tocS_v1CreateMerchantBalanceReply">v1CreateMerchantBalanceReply</h2>

<a id="schemav1createmerchantbalancereply"></a>
<a id="schema_v1CreateMerchantBalanceReply"></a>
<a id="tocSv1createmerchantbalancereply"></a>
<a id="tocsv1createmerchantbalancereply"></a>

```json
{
  "userId": "string",
  "currency": "string",
  "available": 0.1
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|userId|string|false|none||none|
|currency|string|false|none||none|
|available|number(double)|false|none||none|

