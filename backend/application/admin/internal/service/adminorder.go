package service

import (
	"context"

	userv1 "backend/api/user/v1"

	"google.golang.org/protobuf/types/known/timestamppb"

	cartv1 "backend/api/cart/v1"

	"backend/constants"

	"github.com/go-kratos/kratos/v2/log"

	adminv1 "backend/api/admin/order/v1"
	orderv1 "backend/api/order/v1"
	"backend/application/admin/internal/biz"

	v1 "backend/api/admin/order/v1"
)

type AdminOrderService struct {
	v1.UnimplementedAdminOrderServer
	oc *biz.AdminOrderUsecase
}

func NewAdminOrderService(oc *biz.AdminOrderUsecase) *AdminOrderService {
	return &AdminOrderService{oc: oc}
}

func (s *AdminOrderService) GetAllOrders(ctx context.Context, req *adminv1.GetAllOrdersReq) (*adminv1.AdminOrderReply, error) {
	// 调用业务层获取订单列表
	resp, err := s.oc.GetAllOrders(ctx, &biz.GetAllOrdersReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// 检查是否有订单
	if len(resp.Orders) == 0 {
		log.Infof("没有订单记录")
		return nil, nil
	}

	orders := make([]*v1.SubOrder, 0, len(resp.Orders))
	for _, o := range resp.Orders {
		subOrderItem := make([]*v1.SubOrderItem, 0, len(o.SubOrderItems))

		for _, item := range o.SubOrderItems {
			subOrderItem = append(subOrderItem, &v1.SubOrderItem{
				Item: &cartv1.CartItem{
					MerchantId: item.Item.MerchantId.String(),
					ProductId:  item.Item.ProductId.String(),
					Quantity:   item.Item.Quantity,
					Name:       item.Item.Name,
					Picture:    item.Item.Picture,
				},
				Cost: item.Cost,
			})
		}

		orders = append(orders, &v1.SubOrder{
			OrderId:     o.OrderID,
			SubOrderId:  o.SubOrderID,
			TotalAmount: o.TotalAmount,
			ConsumerId:  o.ConsumerId.String(),
			Address: &userv1.ConsumerAddress{
				City:          o.ConsumerAddress.City,
				State:         o.ConsumerAddress.State,
				Country:       o.ConsumerAddress.Country,
				ZipCode:       o.ConsumerAddress.ZipCode,
				StreetAddress: o.ConsumerAddress.StreetAddress,
			},
			ConsumerEmail:  o.ConsumerEmail,
			Currency:       string(o.Currency),
			SubOrderItems:  subOrderItem,
			PaymentStatus:  convertToPaymentStatus(o.PaymentStatus),
			ShippingStatus: convertToShippingStatus(o.ShippingStatus),
			CreatedAt:      timestamppb.New(o.CreatedAt),
			UpdatedAt:      timestamppb.New(o.UpdatedAt),
		})
	}

	return &adminv1.AdminOrderReply{
		Orders: orders,
	}, nil
}

func convertToShippingStatus(status constants.ShippingStatus) orderv1.ShippingStatus {
	switch status {
	case constants.ShippingWaitCommand:
		return orderv1.ShippingStatus_WAIT_COMMAND
	case constants.ShippingPending:
		return orderv1.ShippingStatus_PENDING_SHIPMENT
	case constants.ShippingShipped:
		return orderv1.ShippingStatus_SHIPPED
	case constants.ShippingInTransit:
		return orderv1.ShippingStatus_IN_TRANSIT
	case constants.ShippingDelivered:
		return orderv1.ShippingStatus_DELIVERED
	case constants.ShippingConfirmed:
		return orderv1.ShippingStatus_CONFIRMED
	case constants.ShippingCancelled:
		return orderv1.ShippingStatus_CANCELLED_SHIPMENT
	default:
		return orderv1.ShippingStatus_WAIT_COMMAND
	}
}

func convertToPaymentStatus(status constants.PaymentStatus) orderv1.PaymentStatus {
	switch status {
	case constants.PaymentPending:
		return orderv1.PaymentStatus_PENDING
	case constants.PaymentPaid:
		return orderv1.PaymentStatus_PAID
	case constants.PaymentFailed:
		return orderv1.PaymentStatus_FAILED
	case constants.PaymentCancelled:
		return orderv1.PaymentStatus_CANCELLED
	default:
		return orderv1.PaymentStatus_PENDING
	}
}
