package data

import (
	"backend/application/cart/internal/biz"
	"backend/application/cart/internal/data/models"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

func NewCartRepo(data *Data, logger log.Logger) biz.CartRepo {
	return &cartRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// EmptyCart implements biz.CartRepo.
func (c *cartRepo) EmptyCart(ctx context.Context, req *biz.EmptyCartReq) (*biz.EmptyCartResp, error) {
	err := c.data.db.EmptyCart(ctx, int32(req.UserId))
	if err != nil {
		return nil, err
	}
	return &biz.EmptyCartResp{
		Success: true,
	}, nil
}

// GetCart implements biz.CartRepo.
func (c *cartRepo) GetCart(ctx context.Context, req *biz.GetCartReq) (*biz.GetCartResp, error) {
	cart, err := c.data.db.GetCart(ctx, int32(req.UserId))
	if err != nil {
		return nil, err
	}
	var cartItems []biz.CartItem
	for _, item := range cart {
		var cartitem biz.CartItem
		cartitem.ProductId = uint32(item.CartItemID)
		cartitem.Quantity = item.Quantity
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
		UserID:    int32(req.UserId),
		ProductID: int32(req.ProductId),
	})
	if err != nil || dreq == (models.CartSchemaCartItems{}) {
		return nil, err
	}
	c.log.WithContext(ctx).Infof("RemoveCartItem request2 : %+v", dreq)
	return &biz.RemoveCartItemResp{
		Success: true,
	}, nil
}

// UpsertItem implements biz.CartRepo.
func (c *cartRepo) UpsertItem(ctx context.Context, req *biz.UpsertItemReq) (*biz.UpsertItemResp, error) {
	panic("unimplemented")
}
