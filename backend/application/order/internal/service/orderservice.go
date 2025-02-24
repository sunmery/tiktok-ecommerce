package service

import (
	"context"
	"fmt"

	"backend/application/order/internal/biz"

	v1 "backend/api/order/v1"
	"github.com/google/uuid"
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

// ListOrders TODO
func (s *OrderServiceService) ListOrders(ctx context.Context, req *v1.ListOrderReq) (*v1.ListOrderResp, error) {
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

	listReq := &biz.ListOrderReq{
		UserID:   UserMock,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	resp, err := s.uc.ListOrder(ctx, listReq)
	if err != nil {
		return nil, err
	}
	fmt.Printf("resp: %+v\n", resp)

	return &v1.ListOrderResp{Orders: nil}, nil
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
