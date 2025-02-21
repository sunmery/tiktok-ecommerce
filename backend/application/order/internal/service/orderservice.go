package service

import (
	"backend/application/order/internal/biz"
	"context"
	"fmt"
	"github.com/google/uuid"

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
		orderItems = append(orderItems, &biz.OrderItem{
			Item: &biz.CartItem{
				ProductId: item.Item.ProductId,
				Quantity:  item.Item.Quantity,
			},
			Cost: 0,
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
	listOrder, err := s.uc.ListOrder(ctx, &biz.ListOrderReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	var orders []*v1.Order
	for i, order := range listOrder.Orders {
		var orderItems []*v1.OrderItem
		for _, item := range order.OrderItems {
			orderItems = append(orderItems, &v1.OrderItem{
				Item: &v1.CartItem{
					ProductId: item.Item.ProductId,
					Quantity:  item.Item.Quantity,
				},
				Cost: item.Cost,
			})
		}
		orders[i] = &v1.Order{
			OrderItems: orderItems,
			OrderId:    order.OrderId,
			UserId:     order.UserId,
			Currency:   order.UserCurrency,
			Address: &v1.Address{
				StreetAddress: order.Address.StreetAddress,
				City:          order.Address.City,
				State:         order.Address.State,
				Country:       order.Address.Country,
				ZipCode:       order.Address.ZipCode,
			},
			Email:         order.Email,
			CreatedAt:     order.CreatedAt,
			PaymentStatus: order.PaymentStatus,
		}
	}
	return &v1.ListOrderResp{
		Orders: orders,
	}, nil
}

func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *v1.MarkOrderPaidReq) (*v1.MarkOrderPaidResp, error) {
	orderPaid, err := s.uc.MarkOrderPaid(ctx, &biz.MarkOrderPaidReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("orderPaid:%+v\n", orderPaid)

	return &v1.MarkOrderPaidResp{}, nil
}
