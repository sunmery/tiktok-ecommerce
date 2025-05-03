package service

import (
	"context"

	"backend/constants"

	"github.com/go-kratos/kratos/v2/log"

	adminv1 "backend/api/admin/order/v1"
	cartv1 "backend/api/cart/v1"
	orderv1 "backend/api/order/v1"
	userv1 "backend/api/user/v1"
	"backend/application/admin/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "backend/api/admin/order/v1"
)

type AdminOrderService struct {
	pb.UnimplementedAdminOrderServer
	ac *biz.AdminOrderUsecase
}

func NewAdminOrderService(ac *biz.AdminOrderUsecase) *AdminOrderService {
	return &AdminOrderService{ac: ac}
}

func (s *AdminOrderService) GetAllOrders(ctx context.Context, req *adminv1.GetAllOrdersReq) (*adminv1.AdminOrderReply, error) {
	// 调用业务层获取订单列表
	resp, err := s.ac.GetAllOrders(ctx, &biz.GetAllOrdersReq{
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

	// 将业务层返回的订单数据转换为proto消息格式
	merchantOrders := make([]*adminv1.AdminOrderInterface, 0, len(resp.Orders))
	for _, order := range resp.Orders {
		// 创建订单项
		orderItems := make([]*adminv1.OrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			// 创建购物车商品
			cartItem := &cartv1.CartItem{
				MerchantId: item.Item.MerchantId,
				ProductId:  item.Item.ProductId,
				Quantity:   item.Item.Quantity,
			}

			// 创建订单项
			orderItem := &adminv1.OrderItem{
				SubOrderId:     item.SubOrderID,
				Item:           cartItem,
				Cost:           item.Cost,
				Email:          item.Email,
				UserId:         item.UserID.String(),
				Currency:       item.Currency,
				PaymentStatus:  convertToPaymentStatus(item.PaymentStatus),
				ShippingStatus: convertToShippingStatus(item.ShippingStatus),
				CreatedAt:      timestamppb.New(item.CreatedAt),
				UpdatedAt:      timestamppb.New(item.UpdatedAt),
			}

			// 添加地址信息
			if item.ConsumerAddress.StreetAddress != "" {
				orderItem.Address = &userv1.ConsumerAddress{
					StreetAddress: item.ConsumerAddress.StreetAddress,
					City:          item.ConsumerAddress.City,
					State:         item.ConsumerAddress.State,
					Country:       item.ConsumerAddress.Country,
					ZipCode:       item.ConsumerAddress.ZipCode,
				}
			}

			orderItems = append(orderItems, orderItem)
		}

		// 创建商家订单
		merchantOrder := &adminv1.AdminOrderInterface{
			Items:     orderItems,
			CreatedAt: timestamppb.New(order.CreatedAt),
		}

		merchantOrders = append(merchantOrders, merchantOrder)
	}

	return &adminv1.AdminOrderReply{
		Orders: merchantOrders,
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
