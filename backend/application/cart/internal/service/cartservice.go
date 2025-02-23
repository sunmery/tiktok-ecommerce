package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

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
		return nil, errors.New("failed to parse user id")
	}
	resp, err2 := s.cc.UpsertItem(ctx, &biz.UpsertItemReq{
		UserId: UserMock,

		Item: biz.CartItem{
			ProductId: int32(req.Item.ProductId),
			Quantity:  req.Item.Quantity,
		},
	})
	if err2 != nil {
		return nil, errors.New(fmt.Sprintf("failed to upsert item: %v", err2))
	}
	return &pb.UpsertItemResp{
		Success: resp.Success,
	}, nil
}
func (s *CartServiceService) GetCart(ctx context.Context, _ *emptypb.Empty) (*pb.GetCartResp, error) {
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
		return nil, errors.New("failed to parse user id")
	}
	cart, err2 := s.cc.GetCart(ctx, &biz.GetCartReq{
		UserId: UserMock,
	})
	if err2 != nil {
		return nil, errors.New("failed to get cart")
	}
	items := make([]*pb.CartItem, len(cart.Cart.Items))
	for i, item := range cart.Cart.Items {
		items[i] = &pb.CartItem{
			ProductId: uint32(item.ProductId),
			Quantity:  item.Quantity,
		}
	}
	return &pb.GetCartResp{
		Cart: &pb.Cart{
			UserId: UserMock.String(),
			Items:  items,
		},
	}, nil
}
func (s *CartServiceService) EmptyCart(ctx context.Context, _ *emptypb.Empty) (*pb.EmptyCartResp, error) {
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
	resp, err := s.cc.EmptyCart(ctx, &biz.EmptyCartReq{
		UserId: UserMock,
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
	resp, err := s.cc.RemoveCartItem(ctx, &biz.RemoveCartItemReq{
		UserId:    UserMock,
		ProductId: req.ProductId,
	})
	if err != nil {
		return nil, errors.New("failed to remove cart item")
	}
	return &pb.RemoveCartItemResp{
		Success: resp.Success,
	}, nil
}
