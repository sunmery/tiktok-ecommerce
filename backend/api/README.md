# API 接口汇总

## 商品服务 (product/v1)
- `UploadProductFile`: 上传商品文件
- `CreateProduct`: 创建商品
- `SubmitForAudit`: 提交商品审核
- `AuditProduct`: 审核商品
- `ListRandomProducts`: 获取随机商品列表
- `GetCategoryProducts`: 获取分类商品
- `GetCategoryWithChildrenProducts`: 获取分类及其子分类商品
- `GetProductsBatch`: 批量获取商品
- `GetProduct`: 获取单个商品详情
- `SearchProductsByName`: 按名称搜索商品
- `ListProductsByCategory`: 按分类获取商品列表
- `DeleteProduct`: 删除商品

## 购物车服务 (cart/v1)
- `UpsertItem`: 添加或更新购物车商品
- `GetCart`: 获取购物车内容
- `EmptyCart`: 清空购物车
- `RemoveCartItem`: 移除购物车商品

## 订单服务 (order/v1)
- `PlaceOrder`: 下单
- `GetAllOrders`: 获取所有订单
- `MarkOrderPaid`: 标记订单为已支付

## 分类服务 (category/v1)
- `CreateCategory`: 创建分类
- `GetLeafCategories`: 获取叶子分类
- `BatchGetCategories`: 批量获取分类
- `GetCategory`: 获取单个分类
- `UpdateCategory`: 更新分类
- `DeleteCategory`: 删除分类
- `GetSubTree`: 获取分类子树
- `GetDirectSubCategories`: 获取直接子分类
- `GetCategoryPath`: 获取分类路径
- `GetClosureRelations`: 获取分类闭包关系
- `UpdateClosureDepth`: 更新分类闭包深度

## 用户服务 (user/v1)
- `GetUserProfile`: 获取用户资料
- `GetUsers`: 获取用户列表
- `DeleteUser`: 删除用户
- `CreateAddresses`: 创建地址
- `UpdateAddresses`: 更新地址
- `DeleteAddresses`: 删除地址
- `GetAddress`: 获取地址
- `CreateCreditCard`: 创建信用卡
- `GetAddresses`: 获取地址列表
- `ListCreditCards`: 获取信用卡列表
- `UpdateUser`: 更新用户信息
- `GetCreditCard`: 获取信用卡详情
- `DeleteCreditCard`: 删除信用卡

## 库存服务 (merchant/inventory/v1)
- `SetStockAlert`: 设置库存预警
- `GetStockAlerts`: 获取库存预警
- `GetLowStockProducts`: 获取低库存商品
- `RecordStockAdjustment`: 记录库存调整
- `GetStockAdjustmentHistory`: 获取库存调整历史
- `GetProductStock`: 获取商品库存
- `UpdateProductStock`: 更新商品库存

## 商家商品服务 (merchant/product/v1)
- `GetMerchantProducts`: 获取商家商品
- `UpdateProduct`: 更新商品

## 商家订单服务 (merchant/order/v1)
- `GetMerchantOrders`: 获取商家订单

## 助手服务 (assistant/v1)
- `ProcessQuery`: 处理查询

## 结账服务 (checkout/v1)
- `Checkout`: 结账

## 认证服务 (auth/v1)
- `Signin`: 登录