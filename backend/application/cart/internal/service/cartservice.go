package service

import (
	"context"

	apiErrors "backend/api/cart/v1" // 导入错误包
	pb "backend/api/cart/v1"
	"backend/application/cart/internal/biz"

	"github.com/google/uuid"
)

type CartServiceService struct {
	pb.UnimplementedCartServiceServer
	cc *biz.CartUsecase
}

func NewCartServiceService(cc *biz.CartUsecase) *CartServiceService {
	return &CartServiceService{cc: cc}
}

func (s *CartServiceService) CheckCartItem(ctx context.Context, req *pb.CheckCartItemReq) (*pb.CheckCartItemResp, error) {
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
	userid, err := uuid.Parse("77d08975-972c-4a06-8aa4-d2d23f374bb1")
	if err != nil {
		return nil, apiErrors.ErrorInvalidStatus("failed to parse UUID: %v", err)
	}
	resp, err := s.cc.CheckCartItem(ctx, &biz.CheckCartItemReq{
		UserId: userid.String(),
		//UserId:     req.UserId,
		MerchantId: req.MerchantId,
		ProductId:  req.ProductId,
	})
	if err != nil {
		return nil, apiErrors.ErrorCartitemNotFound("failed to check cart item: %v", err)
	}
	return &pb.CheckCartItemResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) UncheckCartItem(ctx context.Context, req *pb.UncheckCartItemReq) (*pb.UncheckCartItemResp, error) {
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
	resp, err := s.cc.UncheckCartItem(ctx, &biz.UncheckCartItemReq{
		UserId:     req.UserId,
		MerchantId: req.MerchantId,
		ProductId:  req.ProductId,
	})
	if err != nil {
		return nil, apiErrors.ErrorCartitemNotFound("failed to uncheck cart item: %v", err)
	}
	return &pb.UncheckCartItemResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderResp, error) {
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

	resp, err := s.cc.CreateOrder(ctx, &biz.CreateOrderReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, apiErrors.ErrorInvalidStatus("failed to create order: %v", err)
	}
	items := make([]*pb.CartItem, len(resp.Items))
	for i, item := range resp.Items {
		items[i] = &pb.CartItem{
			MerchantId: item.MerchantId,
			ProductId:  item.ProductId,
			Quantity:   item.Quantity,
			Selected:   item.Selected,
		}
	}
	return &pb.CreateOrderResp{
		Success: resp.Success,
		Items:   items,
	}, nil
}

func (s *CartServiceService) CreateCart(ctx context.Context, req *pb.CreateCartReq) (*pb.CreateCartResp, error) {
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

	resp, err := s.cc.CreateCart(ctx, &biz.CreateCartReq{
		UserId:   req.UserId,
		CartName: req.CartName,
	})
	if err != nil {
		return nil, apiErrors.ErrorInvalidAuditAction("failed to create cart: %v", err)
	}
	return &pb.CreateCartResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

func (s *CartServiceService) ListCarts(ctx context.Context, req *pb.ListCartsReq) (*pb.ListCartsResp, error) {
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

	carts, err := s.cc.ListCarts(ctx, &biz.ListCartsReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, apiErrors.ErrorInvalidAuditAction("failed to list carts: %v", err)
	}
	cartList := make([]*pb.CartSummary, len(carts.Carts))
	for i, cart := range carts.Carts {
		cartList[i] = &pb.CartSummary{
			CartId:   cart.CartId,
			CartName: cart.CartName,
		}
	}
	return &pb.ListCartsResp{
		Carts: cartList,
	}, nil
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
		return nil, apiErrors.ErrorInvalidStatus("failed to parse UUID: %v", err)
	}

	resp, err := s.cc.UpsertItem(ctx, &biz.UpsertItemReq{
		UserId: UserMock.String(),
		//UserId: req.UserId,
		Item: biz.CartItem{
			MerchantId: req.Item.MerchantId,
			ProductId:  req.Item.ProductId,
			Quantity:   req.Item.Quantity,
			Selected:   req.Item.Selected,
		},
	})
	if err != nil {
		return nil, apiErrors.ErrorInvalidAuditAction("failed to upsert item: %v", err)
	}
	return &pb.UpsertItemResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) GetCart(ctx context.Context, req *pb.GetCartReq) (*pb.GetCartResp, error) {
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

	cart, err := s.cc.GetCart(ctx, &biz.GetCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, apiErrors.ErrorCartitemNotFound("failed to get cart: %v", err)
	}
	items := make([]*pb.CartItem, len(cart.Cart.Items))
	for i, item := range cart.Cart.Items {
		items[i] = &pb.CartItem{
			MerchantId: item.MerchantId,
			ProductId:  item.ProductId,
			Quantity:   item.Quantity,
			Selected:   item.Selected,
		}
	}
	return &pb.GetCartResp{
		Cart: &pb.Cart{
			UserId: cart.Cart.UserId,
			Items:  items,
		},
	}, nil
}

func (s *CartServiceService) EmptyCart(ctx context.Context, req *pb.EmptyCartReq) (*pb.EmptyCartResp, error) {
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

	resp, err := s.cc.EmptyCart(ctx, &biz.EmptyCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, apiErrors.ErrorInvalidStatus("failed to empty cart: %v", err)
	}
	if !resp.Success {
		return nil, apiErrors.ErrorInvalidStatus("failed to empty cart")
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

	resp, err := s.cc.RemoveCartItem(ctx, &biz.RemoveCartItemReq{
		UserId:     req.UserId,
		MerchantId: req.MerchantId,
		ProductId:  req.ProductId,
	})
	if err != nil {
		return &pb.RemoveCartItemResp{
			Success: resp.Success,
		}, apiErrors.ErrorCartitemNotFound("failed to remove cart item: %v", err)
	}
	return &pb.RemoveCartItemResp{
		Success: resp.Success,
	}, nil
}
