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

func (s *CartServiceService) CheckCartItem(ctx context.Context, req *pb.CheckCartItemReq) (*pb.CheckCartItemResp, error) {
	resp, err := s.cc.CheckCartItem(ctx, &biz.CheckCartItemReq{
		UserId:     req.UserId,
		MerchantId: req.MerchantId,
		ProductId:  req.ProductId,
	})
	if err != nil {
		return nil, errors.New("failed to check cart item")
	}
	return &pb.CheckCartItemResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) UncheckCartItem(ctx context.Context, req *pb.UncheckCartItemReq) (*pb.UncheckCartItemResp, error) {
	resp, err := s.cc.UncheckCartItem(ctx, &biz.UncheckCartItemReq{
		UserId:     req.UserId,
		MerchantId: req.MerchantId,
		ProductId:  req.ProductId,
	})
	if err != nil {
		return nil, errors.New("failed to uncheck cart item")
	}
	return &pb.UncheckCartItemResp{
		Success: resp.Success,
	}, nil
}

func (s *CartServiceService) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderResp, error) {
	resp, err := s.cc.CreateOrder(ctx, &biz.CreateOrderReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, errors.New("failed to create order")
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
	resp, err := s.cc.CreateCart(ctx, &biz.CreateCartReq{
		UserId:   req.UserId,
		CartName: req.CartName,
	})
	if err != nil {
		return nil, errors.New("failed to create cart")
	}
	return &pb.CreateCartResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

func (s *CartServiceService) ListCarts(ctx context.Context, req *pb.ListCartsReq) (*pb.ListCartsResp, error) {
	carts, err := s.cc.ListCarts(ctx, &biz.ListCartsReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, errors.New("failed to list carts")
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
	resp, err := s.cc.UpsertItem(ctx, &biz.UpsertItemReq{
		UserId: req.UserId,
		Item: biz.CartItem{
			MerchantId: req.Item.MerchantId,
			ProductId:  req.Item.ProductId,
			Quantity:   req.Item.Quantity,
			Selected:   req.Item.Selected,
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
	// var extra string
	// if md, ok := metadata.FromServerContext(ctx); ok {
	// 	extra = md.Get("x-md-global-userid")
	// }
	// fmt.Println(extra)
	cart, err := s.cc.GetCart(ctx, &biz.GetCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, errors.New("failed to get cart")
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
	resp, err := s.cc.EmptyCart(ctx, &biz.EmptyCartReq{
		UserId: req.UserId,
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
		UserId:     req.UserId,
		MerchantId: req.MerchantId,
		ProductId:  req.ProductId,
	})
	if err != nil {
		return &pb.RemoveCartItemResp{
			Success: resp.Success,
		}, errors.New("failed to remove cart item")
	}
	return &pb.RemoveCartItemResp{
		Success: resp.Success,
	}, nil
}
