package data

import (
	"backend/application/cart/internal/biz"
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
	return &biz.EmptyCartResp{
		Success: true,
	}, nil
	panic("unimplemented")
}

// GetCart implements biz.CartRepo.
func (c *cartRepo) GetCart(ctx context.Context, req *biz.GetCartReq) (*biz.GetCartResp, error) {

	return &biz.GetCartResp{
		Cart: biz.Cart{
			UserId: 1,
			Items: []biz.CartItem{
				{
					ProductId: 1,
					Quantity:  1,
				},
			},
		},
	}, nil
}

// RemoveCartItem implements biz.CartRepo.
func (c *cartRepo) RemoveCartItem(ctx context.Context, req *biz.RemoveCartItemReq) (*biz.RemoveCartItemResp, error) {
	panic("unimplemented")
}

// UpsertItem implements biz.CartRepo.
func (c *cartRepo) UpsertItem(ctx context.Context, req *biz.UpsertItemReq) (*biz.UpsertItemResp, error) {
	panic("unimplemented")
}
