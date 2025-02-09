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

// AddItem implements biz.CartRepo.
func (c *cartRepo) AddItem(ctx context.Context, req *biz.AddItemReq) (*biz.AddItemResp, error) {
	panic("unimplemented")
}

// EmptyCart implements biz.CartRepo.
func (c *cartRepo) EmptyCart(ctx context.Context, req *biz.EmptyCartReq) (*biz.EmptyCartResp, error) {
	panic("unimplemented")
}

// GetCart implements biz.CartRepo.
func (c *cartRepo) GetCart(ctx context.Context, req *biz.GetCartReq) (*biz.GetCartResp, error) {
	panic("unimplemented")
}

// RemoveItem implements biz.CartRepo.
func (c *cartRepo) RemoveItem(ctx context.Context, req *biz.RemoveItemReq) (*biz.RemoveItemResp, error) {
	panic("unimplemented")
}

// UpdateItem implements biz.CartRepo.
func (c *cartRepo) UpdateItem(ctx context.Context, req *biz.UpdateItemReq) (*biz.UpdateItemResp, error) {
	panic("unimplemented")
}
