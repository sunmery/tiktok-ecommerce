package service

import (
	"context"
	"errors"
	"fmt"

	pb "backend/api/cart/v1"
	"backend/application/cart/internal/biz"
	"backend/pkg"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartServiceService struct {
	pb.UnimplementedCartServiceServer
	cc *biz.CartUsecase
}

func NewCartServiceService(cc *biz.CartUsecase) *CartServiceService {
	return &CartServiceService{cc: cc}
}

func (s *CartServiceService) UpsertItem(ctx context.Context, req *pb.UpsertItemReq) (*pb.UpsertItemResp, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse userId UUID: %v", err))
	}
	productId, perr := uuid.Parse(req.ProductId)
	if perr != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse productId UUID: %v", perr))
	}
	merchantId, merr := uuid.Parse(req.MerchantId)
	if merr != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse merchantId UUID: %v", merr))
	}

	resp, uerr := s.cc.UpsertItem(ctx, &biz.UpsertItemReq{
		UserId:     userId,
		MerchantId: merchantId,
		ProductId:  productId,
		Quantity:   req.Quantity,
	})
	if uerr != nil {
		return nil, errors.New(fmt.Sprintf("failed to upsert item: %v", uerr))
	}

	return &pb.UpsertItemResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) GetCart(ctx context.Context, _ *pb.GetCartReq) (*pb.Cart, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse userId UUID: %v", err))
	}

	carts, gerr := s.cc.GetCart(ctx, &biz.GetCartReq{
		UserId: userId,
	})
	if gerr != nil {
		return nil, fmt.Errorf("failed to get cart: %v", gerr)
	}

	var items []*pb.CartItem
	for _, item := range carts.Cart.Items {
		items = append(items, &pb.CartItem{
			MerchantId: item.MerchantId.String(),
			ProductId:  item.ProductId.String(),
			// Price:      item.Price,
			Quantity: item.Quantity,
			// TotalPrice: item.Price * float64(item.Quantity),
			Name:    item.Name,
			Picture: item.Picture,
		})
	}
	return &pb.Cart{
		Items: items,
	}, nil
}

func (s *CartServiceService) EmptyCart(ctx context.Context, _ *pb.EmptyCartReq) (*pb.EmptyCartResp, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse userId UUID: %v", err))
	}

	resp, rerr := s.cc.EmptyCart(ctx, &biz.EmptyCartReq{
		UserId: userId,
	})
	if rerr != nil {
		return nil, errors.New(fmt.Sprintf("failed to empty cart: %v", rerr))
	}

	return &pb.EmptyCartResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) RemoveCartItem(ctx context.Context, req *pb.RemoveCartItemReq) (*pb.RemoveCartItemResp, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}
	productId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchant ID")
	}

	resp, err := s.cc.RemoveCartItem(ctx, &biz.RemoveCartItemReq{
		UserId:     userId,
		MerchantId: merchantId,
		ProductId:  productId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove cart item: %v", err)
	}

	return &pb.RemoveCartItemResp{
		Success: resp.Success,
	}, nil
}
