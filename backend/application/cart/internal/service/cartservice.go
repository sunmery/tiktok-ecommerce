package service

import (
	"context"
	"errors"

	pb "backend/api/cart/v1"
	"backend/application/cart/internal/biz"
)

type CartServiceService struct {
	pb.UnimplementedCartServiceServer
	cc *biz.CartUsecase
}

func NewCartServiceService(cc *biz.CartUsecase) *CartServiceService {
	return &CartServiceService{cc: cc}
}

func (s *CartServiceService) UpsertItem(ctx context.Context, req *pb.UpsertItemReq) (*pb.UpsertItemResp, error) {
	resp, err := s.cc.UpsertItem(ctx, &biz.UpsertItemReq{
		Owner: req.Owner,
		Name:  req.Name,
		Item: biz.CartItem{
			ProductId: req.Item.ProductId,
			Quantity:  req.Item.Quantity,
		},
	})
	if err != nil {
		return nil, errors.New("failed to upsert item")
	}
	return &pb.UpsertItemResp{
		Success: resp.Success,
	}, nil
}
func (s *CartServiceService) GetCart(ctx context.Context, req *pb.GetCartReq) (*pb.GetCartResp, error) {
	cart, err := s.cc.GetCart(ctx, &biz.GetCartReq{
		Owner: req.Owner,
		Name:  req.Name,
	})
	if err != nil {
		return nil, errors.New("failed to get cart")
	}
	items := make([]*pb.CartItem, len(cart.Cart.Items))
	for i, item := range cart.Cart.Items {
		items[i] = &pb.CartItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}
	return &pb.GetCartResp{
		Cart: &pb.Cart{
			Owner: cart.Cart.Owner,
			Name:  cart.Cart.Name,
			Items: items,
		},
	}, nil
}
func (s *CartServiceService) EmptyCart(ctx context.Context, req *pb.EmptyCartReq) (*pb.EmptyCartResp, error) {
	resp, err := s.cc.EmptyCart(ctx, &biz.EmptyCartReq{
		Owner: req.Owner,
		Name:  req.Name,
	})
	if err != nil {
		return nil, errors.New("failed to empty cart")
	}
	if !resp.Success {
		return nil, errors.New("failed to empty cart")
	}
	return &pb.EmptyCartResp{
		Success: resp.Success,
	}, nil
}
func (s *CartServiceService) RemoveCartItem(ctx context.Context, req *pb.RemoveCartItemReq) (*pb.RemoveCartItemResp, error) {
	resp, err := s.cc.RemoveCartItem(ctx, &biz.RemoveCartItemReq{
		Owner:     req.Owner,
		Name:      req.Name,
		ProductId: req.ProductId,
	})
	if err != nil {
		return nil, errors.New("failed to remove cart item")
	}
	return &pb.RemoveCartItemResp{
		Success: resp.Success,
	}, nil
}
