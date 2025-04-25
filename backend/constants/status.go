package constants

type (
	RoleType        string // 角色类型
	PaymentStatus   string // 支付状态
	ShippingStatus  string // 发货状态
	FreezeStatus    string // 余额冻结状态
	PaymentMethod   string // 支付方式
	TransactionType string // 交易类型
)

const (
	Consumer RoleType = "Consumer" // 消费者
	Merchant RoleType = "Merchant" // 商家
	Admin    RoleType = "Admin"    // 管理员
	Guest    RoleType = "Guest"    // 管理员
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
