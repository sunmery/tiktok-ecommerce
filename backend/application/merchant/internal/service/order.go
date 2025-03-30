package service

import (
	"context"

	"backend/pkg"

	v1 "backend/api/order/v1"

	cartv1 "backend/api/cart/v1"
	userv1 "backend/api/user/v1"

	"google.golang.org/protobuf/types/known/timestamppb"

	orderv1 "backend/api/merchant/order/v1"
	"backend/application/merchant/internal/biz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *OrderServiceService) GetMerchantOrders(ctx context.Context, req *orderv1.GetMerchantOrdersReq) (*v1.Orders, error) {
	// 从网关获取用户ID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to get user ID")
	}

	// 构建业务层请求
	listReq := &biz.GetMerchantOrdersReq{
		UserID:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 调用业务层获取订单列表
	resp, err := s.oc.GetMerchantOrders(ctx, listReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	// 转换订单列表为API响应格式
	var orders []*v1.Order
	for _, order := range resp.Orders {
		// 转换地址信息
		address := &userv1.Address{
			StreetAddress: order.Address.StreetAddress,
			City:          order.Address.City,
			State:         order.Address.State,
			Country:       order.Address.Country,
			ZipCode:       order.Address.ZipCode,
		}

		// 转换订单项
		var orderItems []*v1.OrderItem
		for _, subOrder := range order.SubOrders {
			for _, item := range subOrder.Items {
				orderItems = append(orderItems, &v1.OrderItem{
					Item: &cartv1.CartItem{
						MerchantId: item.Item.MerchantId.String(),
						ProductId:  item.Item.ProductId.String(),
						Quantity:   item.Item.Quantity,
					},
					Cost: item.Cost,
				})
			}
		}

		// 转换时间戳
		createdAt := timestamppb.New(order.CreatedAt)

		// 添加订单到响应列表
		orders = append(orders, &v1.Order{
			Items:         orderItems,
			OrderId:       order.OrderID,
			UserId:        order.UserID.String(),
			Currency:      order.Currency,
			Address:       address,
			Email:         order.Email,
			CreatedAt:     createdAt,
			PaymentStatus: mapBizStatusToProto(order.PaymentStatus),
		})
	}

	return &v1.Orders{Orders: orders}, nil
}

// 转换业务层枚举到 Proto int
func mapBizStatusToProto(status biz.PaymentStatus) v1.PaymentStatus {
	switch status {
	case biz.PaymentPending:
		return v1.PaymentStatus_NOT_PAID
	case biz.PaymentProcessing:
		return v1.PaymentStatus_PROCESSING
	case biz.PaymentPaid:
		return v1.PaymentStatus_PAID
	case biz.PaymentFailed:
		return v1.PaymentStatus_FAILED
	case biz.PaymentCancelled:
		return v1.PaymentStatus_CANCELLED
	default:
		return v1.PaymentStatus_NOT_PAID
	}
}
