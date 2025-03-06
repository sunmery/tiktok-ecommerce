package service

import (
	pb "backend/api/cart/v1"
	"backend/application/cart/internal/biz"
	"backend/pkg"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type CartServiceService struct {
	pb.UnimplementedCartServiceServer
	cc *biz.CartUsecase
}

func NewCartServiceService(cc *biz.CartUsecase) *CartServiceService {
	return &CartServiceService{cc: cc}
}

func (s *CartServiceService) UpsertItem(ctx context.Context, req *pb.UpsertItemReq) (*pb.UpsertItemResp, error) {
	// 解析 UUID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse userId UUID: %v", err))
	}
	productId, err := uuid.Parse(req.Item.ProductId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse productId UUID: %v", err))
	}
	merchantId, err := uuid.Parse(req.Item.MerchantId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse productId UUID: %v", err))
	}
	resp, err := s.cc.UpsertItem(ctx, &biz.UpsertItemReq{
		UserId: userId,
		Item: biz.CartItem{
			MerchantId: merchantId,
			ProductId:  productId,
			Quantity:   req.Item.Quantity,
		},
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to upsert item: %v", err))
	}
	return &pb.UpsertItemResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) GetCart(ctx context.Context, _ *pb.GetCartReq) (*pb.GetCartResp, error) {
	// 从网关获取用户ID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := s.cc.GetCart(ctx, &biz.GetCartReq{
		UserId:userId,
	})
	if err != nil {
		return nil, pb.ErrorCartitemNotFound("failed to get cart: %v", err)
	}
	items := make([]*pb.CartItem, len(cart.Cart.Items))
	for i, item := range cart.Cart.Items {
		items[i] = &pb.CartItem{
			MerchantId: item.MerchantId.String(),
			ProductId:  item.ProductId.String(),
			Quantity:   item.Quantity,
		}
	}
	return &pb.GetCartResp{
		Cart: &pb.Cart{
			Items:  items,
		},
	}, nil
}

func (s *CartServiceService) EmptyCart(ctx context.Context, _ *pb.EmptyCartReq) (*pb.EmptyCartResp, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse userId UUID: %v", err))
	}

	resp, err := s.cc.EmptyCart(ctx, &biz.EmptyCartReq{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to empty cart: %v", err))
	}
	if !resp.Success {
		return nil, errors.New("failed to empty cart")
	}
	return &pb.EmptyCartResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) RemoveCartItem(ctx context.Context, req *pb.RemoveCartItemReq) (*pb.RemoveCartItemResp, error) {
	// 解析 UUID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse userId UUID: %v", err))
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse productId UUID: %v", err))
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse merchantId UUID: %v", err))
	}

	resp, err := s.cc.RemoveCartItem(ctx, &biz.RemoveCartItemReq{
		UserId:     userId,
		MerchantId: merchantId,
		ProductId:  productId,
	})
	if err != nil {
		return nil, pb.ErrorCartitemNotFound("failed to remove cart item: %v", err)
	}

	return &pb.RemoveCartItemResp{
		Success: resp.Success,
	}, nil
}
