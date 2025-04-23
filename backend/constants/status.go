package constants

type (
	PaymentStatus  string // 支付状态
	ShippingStatus string // 发货状态
)

const (
	PaymentPending    PaymentStatus = "PENDING"    // 等待支付
	PaymentProcessing PaymentStatus = "PROCESSING" // 支付中
	PaymentPaid       PaymentStatus = "PAID"       // 已支付
	PaymentFailed     PaymentStatus = "FAILED"     // 支付失败
	PaymentCancelled  PaymentStatus = "CANCELLED"  // 取消支付
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
