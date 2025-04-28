package constants

type (
	RoleType        string // 角色类型
	PaymentStatus   string // 支付状态
	TradeStatus     string // 支付宝交易状态
	ShippingStatus  string // 发货状态
	FreezeStatus    string // 余额冻结状态
	PaymentMethod   string // 支付方式
	TransactionType string // 交易类型
	AddressType     string // 地址类型
)

const (
	WAREHOUSE    AddressType = "WAREHOUSE"    // 仓库地址
	RETURN       AddressType = "RETURN"       // 退货地址
	STORE        AddressType = "STORE"        // 门店地址
	BILLING      AddressType = "BILLING"      // 财务地址
	HEADQUARTERS AddressType = "HEADQUARTERS" // 总部地址
)

const (
	Consumer RoleType = "Consumer" // 消费者
	Merchant RoleType = "Merchant" // 商家
	Admin    RoleType = "Admin"    // 管理员
	Guest    RoleType = "Guest"    // 管理员
)

const (
	TradeStatusWaitBuyerPay TradeStatus = "WAIT_BUYER_PAY" // （交易创建，等待买家付款）
	TradeStatusClosed       TradeStatus = "TRADE_CLOSED"   // （未付款交易超时关闭，或支付完成后全额退款）
	TradeStatusSuccess      TradeStatus = "TRADE_SUCCESS"  // （交易支付成功）
	TradeStatusFinished     TradeStatus = "TRADE_FINISHED" // （交易结束，不可退款）
)

const (
	PaymentPending   PaymentStatus = "PENDING"   // 等待支付
	PaymentPaid      PaymentStatus = "PAID"      // 已支付
	PaymentFailed    PaymentStatus = "FAILED"    // 支付失败
	PaymentCancelled PaymentStatus = "CANCELLED" // 取消支付
)

const (
	ShippingWaitCommand ShippingStatus = "WAIT_COMMAND"     // 等待操作, 支付完成支付和订单分单好之后, 商家获取订单时展示
	ShippingPending     ShippingStatus = "PENDING_SHIPMENT" // 待发货
	ShippingShipped     ShippingStatus = "SHIPPED"          // 已发货
	ShippingInTransit   ShippingStatus = "IN_TRANSIT"       // 运输中
	ShippingDelivered   ShippingStatus = "DELIVERED"        // 已送达
	ShippingConfirmed   ShippingStatus = "CONFIRMED"        // 确认收货
	ShippingCancelled   ShippingStatus = "CANCELLED"        // 已取消发货
)

const (
	PaymentMethodAlipay   PaymentMethod = "ALIPAY"
	PaymentMethodWechat   PaymentMethod = "WECHAT"
	PaymentMethodBalancer PaymentMethod = "BALANCER"
	PaymentMethodBankCard PaymentMethod = "BANK_CARD"
)

const (
	FreezeFrozen    FreezeStatus = "FROZEN"    // 冻结余额
	FreezeConfirmed FreezeStatus = "CONFIRMED" // 确认余额
	FreezeCanceled  FreezeStatus = "CANCELED"  // 取消冻结
)

const (
	TransactionRecharge TransactionType = "RECHARGE" // 充值
	TransactionPayment  TransactionType = "PAYMENT"  // 支付
	TransactionRefund   TransactionType = "REFUND"   // 退款
	TransactionWithdraw TransactionType = "WITHDRAW" // 提现
)
