package service

import (
	pb "backend/api/order/v1"
	"backend/application/order/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/metadata"
)

type OrderServiceService struct {
	pb.UnimplementedOrderServiceServer

	uc *biz.OrderUsecase
}

func NewOrderServiceService(uc *biz.OrderService) *OrderServiceService {
	return &OrderServiceService{uc: uc}
}

func (s *OrderServiceService) PlaceOrder(ctx context.Context, req *pb.PlaceOrderReq) (*pb.PlaceOrderResp, error) {
	// 从网关获取 userId 元信息
	var userId string
	if md, ok := metadata.FromServerContext(ctx); ok {
		userId = md.Get("x-md-global-userId")
	}

	// 类型转换
	var orderItems = make([]*biz.OrderItem, 0)
	for _, item := range req.OrderItems {
		orderItems = append(orderItems, &biz.OrderItem{
			Item: biz.CartItem{
				ProductId: item.Item.ProductId,
				Quantity:  item.Item.Quantity,
			},
			Cost: item.Cost,
		})
	}

	// DTO -> DO
	order, err := s.uc.PlaceOrder(ctx, &biz.PlaceOrderReq{
		UserId:   userId,
		Currency: req.Currency,
		Address: biz.Address{
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

	return &pb.PlaceOrderResp{
		Order: &pb.OrderResult{OrderId: order.OrderId.String()},
	}, nil
}
func (s *OrderServiceService) ListOrder(ctx context.Context, req *pb.ListOrderReq) (*pb.ListOrderResp, error) {
	return &pb.ListOrderResp{}, nil
}
func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *pb.MarkOrderPaidReq) (*pb.MarkOrderPaidResp, error) {
	return &pb.MarkOrderPaidResp{}, nil
}
