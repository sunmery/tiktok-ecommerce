package service

import (
	"context"

	pb "api/cart/v1"
)

type CartServiceService struct {
	pb.UnimplementedCartServiceServer
}

func NewCartServiceService() *CartServiceService {
	return &CartServiceService{}
}

func (s *CartServiceService) AddItem(ctx context.Context, req *pb.AddItemReq) (*pb.AddItemResp, error) {
	return &pb.AddItemResp{}, nil
}
func (s *CartServiceService) GetCart(ctx context.Context, req *pb.GetCartReq) (*pb.GetCartResp, error) {
	return &pb.GetCartResp{}, nil
}
func (s *CartServiceService) EmptyCart(ctx context.Context, req *pb.EmptyCartReq) (*pb.EmptyCartResp, error) {
	return &pb.EmptyCartResp{}, nil
}
func (s *CartServiceService) UpdateItem(ctx context.Context, req *pb.UpdateItemReq) (*pb.UpdateItemResp, error) {
	return &pb.UpdateItemResp{}, nil
}
func (s *CartServiceService) RemoveItem(ctx context.Context, req *pb.RemoveItemReq) (*pb.RemoveItemResp, error) {
	return &pb.RemoveItemResp{}, nil
}
