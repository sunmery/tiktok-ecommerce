package service

import (
	"context"
	"fmt"

	cartv1 "backend/api/cart/v1"

	userv1 "backend/api/user/v1"

	"backend/application/order/internal/biz"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "backend/api/order/v1"
)

type OrderServiceService struct {
	v1.UnimplementedOrderServiceServer

	uc *biz.OrderUsecase
}

func NewOrderServiceService(uc *biz.OrderUsecase) *OrderServiceService {
	return &OrderServiceService{uc: uc}
}

func (s *OrderServiceService) PlaceOrder(ctx context.Context, req *v1.PlaceOrderReq) (*v1.PlaceOrderResp, error) {
	// 从网关获取用户ID
	// var userIdStr string
	// if md, ok := metadata.FromServerContext(ctx); ok {
	// 	userIdStr = md.Get("x-md-global-user-id")
	// }
	// 解析 UUID
	// userId, err := uuid.Parse(userIdStr)
	// if err != nil {
	// 	return nil, err
	// }
	UserMock, err := uuid.Parse("77d08975-972c-4a06-8aa4-d2d23f374bb1")

	// DTO -> DO

	var orderItems []*biz.OrderItem
	for _, item := range req.OrderItems {
		merchantId, err := uuid.Parse(item.Item.MerchantId)
		if err != nil {
			return nil, fmt.Errorf("invalid merchant id: %s", item.Item.MerchantId)
		}

		ProductId := uuid.MustParse(item.Item.ProductId)
		orderItems = append(orderItems, &biz.OrderItem{
			Item: &biz.CartItem{
				MerchantId: merchantId,
				ProductId:  ProductId,
				Quantity:   item.Item.Quantity,
			},
			Cost: item.Cost,
		})
	}

	order, err := s.uc.PlaceOrder(ctx, &biz.PlaceOrderReq{
		// UserId:   userId,
		UserId:   UserMock,
		Currency: req.Currency,
		Address: &biz.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
		},
		Email:      req.Email,
		OrderItems: orderItems,
	})
	if err != nil {
		return nil, err
	}

	return &v1.PlaceOrderResp{
		Order: &v1.OrderResult{
			OrderId: order.Order.OrderId,
		},
	}, nil
}

func (s *OrderServiceService) ListOrder(ctx context.Context, req *v1.ListOrderReq) (*v1.ListOrderResp, error) {
	// 从网关获取用户ID
	// var userIdStr string
	// if md, ok := metadata.FromServerContext(ctx); ok {
	// 	userIdStr = md.Get("x-md-global-user-id")
	// }
	// 解析 UUID
	// userId, err := uuid.Parse(userIdStr)
	// if err != nil {
	// 	return nil, err
	// }
	UserMock, err := uuid.Parse("77d08975-972c-4a06-8aa4-d2d23f374bb1")
	if err != nil {
		return nil, fmt.Errorf("解析用户ID失败: %w", err)
	}

	listOrder, orderErr := s.uc.ListOrder(ctx, &biz.ListOrderReq{
		UserID:        UserMock,
		DateRangeType: req.DateRangeType,
		StartTime:     req.StartTime.AsTime(),
		EndTime:       req.EndTime.AsTime(),
		Page:          int(req.Page),
		PageSize:      int(req.PageSize),
	})
	if orderErr != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", orderErr)
	}

	var orders []*v1.Order

	// 遍历订单列表
	for _, order := range listOrder.Orders {
		// 转换订单项
		var orderItems []*v1.OrderItem
		for _, item := range orderItems {
			orderItems = append(orderItems, &v1.OrderItem{
				Item: &cartv1.CartItem{
					ProductId:  item.Item.ProductId,
					Quantity:   item.Item.Quantity,
					MerchantId: item.Item.MerchantId,
				},
				Cost: item.Cost,
			})
		}

		// 转换地址信息
		address := &userv1.Address{
			StreetAddress: order.Address.StreetAddress,
			City:          order.Address.City,
			State:         order.Address.State,
			Country:       order.Address.Country,
			ZipCode:       order.Address.ZipCode,
		}

		// 确保支付状态正确映射
		paymentStatus := mapPaymentStatus(order.PaymentStatus)

		// 创建 Order 对象
		orders = append(orders, &v1.Order{
			OrderItems:    orderItems,
			OrderId:       order.OrderID,
			UserId:        order.UserID.String(),
			Currency:      order.Currency,
			Address:       address,
			Email:         order.Email,
			CreatedAt:     timestamppb.New(order.CreatedAt),
			PaymentStatus: paymentStatus,
		})
	}

	// 返回结果
	return &v1.ListOrderResp{
		Orders:     orders,
		Stats:      nil,
		Pagination: nil,
	}, nil
}

func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *v1.MarkOrderPaidReq) (*v1.MarkOrderPaidResp, error) {
	// 从网关获取用户ID
	// var userIdStr string
	// if md, ok := metadata.FromServerContext(ctx); ok {
	// 	userIdStr = md.Get("x-md-global-user-id")
	// }
	// 解析 UUID
	// userId, err := uuid.Parse(userIdStr)
	// if err != nil {
	// 	return nil, err
	// }
	UserMock, err := uuid.Parse("77d08975-972c-4a06-8aa4-d2d23f374bb1")

	orderPaid, err := s.uc.MarkOrderPaid(ctx, &biz.MarkOrderPaidReq{
		UserId:  UserMock,
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("orderPaid:%+v\n", orderPaid)

	return &v1.MarkOrderPaidResp{}, nil
}

// 辅助函数：将原始支付状态映射到 v1.PaymentStatus 枚举
func mapPaymentStatus(status int) v1.PaymentStatus {
	switch status {
	case 0:
		return v1.PaymentStatus_NOT_PAID
	case 1:
		return v1.PaymentStatus_PROCESSING
	case 2:
		return v1.PaymentStatus_PAID
	case 3:
		return v1.PaymentStatus_FAILED
	default:
		// 如果状态未知，默认返回 NOT_PAID
		return v1.PaymentStatus_NOT_PAID
	}
}
