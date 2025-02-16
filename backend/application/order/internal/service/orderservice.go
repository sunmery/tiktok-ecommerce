package service

import (
	"backend/application/order/internal/biz"
	"context"

	pb "backend/api/order/v1"
)

type OrderServiceService struct {
	pb.UnimplementedOrderServiceServer
	oc *biz.OrderUseCase
}

func NewOrderServiceService(oc *biz.OrderUseCase) *OrderServiceService {
	return &OrderServiceService{
		oc: oc,
	}
}

func (s *OrderServiceService) PlaceOrder(ctx context.Context, req *pb.PlaceOrderReq) (*pb.PlaceOrderResp, error) {
	return &pb.PlaceOrderResp{}, nil
}
func (s *OrderServiceService) ListOrder(ctx context.Context, req *pb.ListOrderReq) (*pb.ListOrderResp, error) {
	return &pb.ListOrderResp{}, nil
}
func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *pb.MarkOrderPaidReq) (*pb.MarkOrderPaidResp, error) {
	return &pb.MarkOrderPaidResp{}, nil
}
