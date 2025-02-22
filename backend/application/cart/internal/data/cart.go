package data

import (
	"backend/application/cart/internal/biz"
	"backend/application/cart/internal/data/models"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type cartRepo struct {
	data *Data
	log  *log.Helper
}

func NewCartRepo(data *Data, logger log.Logger) biz.CartRepo {
	return &cartRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CheckCartItem implements biz.CartRepo.
func (c *cartRepo) CheckCartItem(ctx context.Context, req *biz.CheckCartItemReq) (*biz.CheckCartItemResp, error) {
	c.log.WithContext(ctx).Infof("CheckCartItem request : %+v", req)
	err := c.data.db.CheckCartItem(ctx, models.CheckCartItemParams{
		UserID:     req.UserId,
		MerchantID: req.MerchantId,
		ProductID:  int32(req.ProductId),
		CartName:   "cart",
	})
	if err != nil {
		return &biz.CheckCartItemResp{
			Success: false,
		}, err
	}
	return &biz.CheckCartItemResp{
		Success: true,
	}, nil
}

// CreateCart implements biz.CartRepo.
func (c *cartRepo) CreateCart(ctx context.Context, req *biz.CreateCartReq) (*biz.CreateCartResp, error) {
	resp, err := c.data.db.CreateCart(ctx, models.CreateCartParams{
		UserID:   req.UserId,
		CartName: req.CartName,
	})
	if err != nil {
		return nil, err
	}
	c.log.WithContext(ctx).Infof("CreateCart resp: %+v", resp)
	return &biz.CreateCartResp{
		Success: true,
		Message: "created CartID: " + string(resp.CartID+'0'),
	}, nil

}

// CreateOrder implements biz.CartRepo.
func (c *cartRepo) CreateOrder(ctx context.Context, req *biz.CreateOrderReq) (*biz.CreateOrderResp, error) {
	resp, err := c.data.db.CreateOrder(ctx, models.CreateOrderParams{
		UserID:   req.UserId,
		CartName: "cart",
	})
	if err != nil {
		return nil, err
	}
	var cartItems []biz.CartItem
	for _, item := range resp {
		var cartitem biz.CartItem
		cartitem.MerchantId = item.MerchantID
		cartitem.ProductId = uint32(item.ProductID)
		cartitem.Quantity = item.Quantity
		cartitem.Selected = item.Selected
		cartItems = append(cartItems, cartitem)
	}
	return &biz.CreateOrderResp{
		Success: true,
		Items:   cartItems,
	}, nil
}

// ListCarts implements biz.CartRepo.
func (c *cartRepo) ListCarts(ctx context.Context, req *biz.ListCartsReq) (*biz.ListCartsResp, error) {
	resp, err := c.data.db.ListCarts(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	c.log.WithContext(ctx).Infof("ListCarts resp: %+v", resp)
	var carts []biz.CartSummary
	for _, cart := range resp {
		var cartitem biz.CartSummary
		cartitem.CartId = uint32(cart.CartID)
		cartitem.CartName = cart.CartName
		carts = append(carts, cartitem)
	}
	c.log.WithContext(ctx).Infof("ListCarts resp: %+v", carts)
	return &biz.ListCartsResp{
		Carts: carts,
	}, nil

}

// UncheckCartItem implements biz.CartRepo.
func (c *cartRepo) UncheckCartItem(ctx context.Context, req *biz.UncheckCartItemReq) (*biz.UncheckCartItemResp, error) {
	err := c.data.db.UncheckCartItem(ctx, models.UncheckCartItemParams{
		UserID:     req.UserId,
		MerchantID: req.MerchantId,
		ProductID:  int32(req.ProductId),
		CartName:   "cart",
	})
	if err != nil {
		return nil, err
	}
	return &biz.UncheckCartItemResp{
		Success: true,
	}, nil
}

// EmptyCart implements biz.CartRepo.
func (c *cartRepo) EmptyCart(ctx context.Context, req *biz.EmptyCartReq) (*biz.EmptyCartResp, error) {
	err := c.data.db.EmptyCart(ctx, models.EmptyCartParams{
		UserID:   req.UserId,
		CartName: "cart",
	})
	if err != nil {
		return nil, err
	}
	return &biz.EmptyCartResp{
		Success: true,
	}, nil
}

// GetCart implements biz.CartRepo.
func (c *cartRepo) GetCart(ctx context.Context, req *biz.GetCartReq) (*biz.GetCartResp, error) {
	cart, err := c.data.db.GetCart(ctx, models.GetCartParams{
		UserID:   req.UserId,
		CartName: "cart",
	})
	if err != nil {
		return nil, err
	}
	var cartItems []biz.CartItem
	for _, item := range cart {
		var cartitem biz.CartItem
		cartitem.MerchantId = item.MerchantID
		cartitem.ProductId = uint32(item.ProductID)
		cartitem.Quantity = item.Quantity
		cartitem.Selected = item.Selected
		cartItems = append(cartItems, cartitem)
	}

	return &biz.GetCartResp{
		Cart: biz.Cart{
			UserId: req.UserId,
			Items:  cartItems,
		},
	}, nil
}

// RemoveCartItem implements biz.CartRepo.
func (c *cartRepo) RemoveCartItem(ctx context.Context, req *biz.RemoveCartItemReq) (*biz.RemoveCartItemResp, error) {
	c.log.WithContext(ctx).Infof("RemoveCartItem request1 : %+v", req)
	dreq, err := c.data.db.RemoveCartItem(ctx, models.RemoveCartItemParams{
		UserID:     req.UserId,
		MerchantID: req.MerchantId,
		ProductID:  int32(req.ProductId),
		CartName:   "cart",
	})
	if err != nil || dreq == (models.CartSchemaCartItems{}) {
		return &biz.RemoveCartItemResp{
			Success: false,
		}, err
	}
	c.log.WithContext(ctx).Infof("RemoveCartItem request2 : %+v", dreq)
	return &biz.RemoveCartItemResp{
		Success: true,
	}, nil
}

// UpsertItem implements biz.CartRepo.
func (c *cartRepo) UpsertItem(ctx context.Context, req *biz.UpsertItemReq) (*biz.UpsertItemResp, error) {
	resp, err := c.data.db.UpsertItem(ctx, models.UpsertItemParams{
		UserID:     req.UserId,
		MerchantID: req.Item.MerchantId,
		ProductID:  int32(req.Item.ProductId),
		Quantity:   int32(req.Item.Quantity),
		CartName:   "cart",
	})
	if resp == (models.CartSchemaCartItems{}) {
		return &biz.UpsertItemResp{
			Success: false,
		}, err
	}
	if err != nil {
		return nil, err
	}
	c.log.WithContext(ctx).Infof("UpsertItem request1 : %+v", resp)
	return &biz.UpsertItemResp{
		Success: true,
	}, nil
}
