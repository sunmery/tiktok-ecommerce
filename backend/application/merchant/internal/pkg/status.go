package pkg

import (
	v1 "backend/api/order/v1"
	"backend/constants"
)

// MapPaymentStatusToProto 将字符串状态转换为Proto枚举
func MapPaymentStatusToProto(status string) v1.PaymentStatus {
	switch status {
	case string(constants.PaymentPending):
		return v1.PaymentStatus_NOT_PAID
	case string(constants.PaymentPaid):
		return v1.PaymentStatus_PAID
	case string(constants.PaymentFailed):
		return v1.PaymentStatus_FAILED
	case string(constants.PaymentCancelled):
		return v1.PaymentStatus_CANCELLED_PAID
	default:
		return v1.PaymentStatus_NOT_PAID
	}
}

// MapShippingStatusToProto 将字符串运输状态转换为Proto枚举
func MapShippingStatusToProto(status constants.ShippingStatus) v1.ShippingStatus {
	switch status {
	case constants.ShippingWaitCommand:
		return v1.ShippingStatus_WAIT_COMMAND
	case constants.ShippingPending:
		return v1.ShippingStatus_PENDING_SHIPMENT
	case constants.ShippingShipped:
		return v1.ShippingStatus_SHIPPED
	case constants.ShippingInTransit:
		return v1.ShippingStatus_IN_TRANSIT
	case constants.ShippingDelivered:
		return v1.ShippingStatus_DELIVERED
	case constants.ShippingConfirmed:
		return v1.ShippingStatus_CONFIRMED
	case constants.ShippingCancelled:
		return v1.ShippingStatus_CANCELLED_SHIPMENT
	default:
		return v1.ShippingStatus_WAIT_COMMAND
	}
}
