package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewCartUsecase)

// var (
// 	// ErrCartNotFound is Cart not found.
// 	ErrCartNotFound = errors.NotFound(v1.ErrorReason_Cart_NOT_FOUND.String(), "Cart not found")
// )

type CartRepo interface {
	AddItem(ctx context.Context, req *AddItemReq) (*AddItemResp, error)
	GetCart(ctx context.Context, req *GetCartReq) (*GetCartResp, error)
	EmptyCart(ctx context.Context, req *EmptyCartReq) (*EmptyCartResp, error)
	UpdateItem(ctx context.Context, req *UpdateItemReq) (*UpdateItemResp, error)
	RemoveItem(ctx context.Context, req *RemoveItemReq) (*RemoveItemResp, error)
}

type CartUsecase struct {
	repo CartRepo
	log  *log.Helper
}

func NewCartUsecase(repo CartRepo, logger log.Logger) *CartUsecase {
	return &CartUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (cc *CartUsecase) AddItem(ctx context.Context, req *AddItemReq) (*AddItemResp, error) {
	cc.log.WithContext(ctx).Infof("AddItem request: %+v", req)
	return cc.repo.AddItem(ctx, req)
}

func (cc *CartUsecase) GetCart(ctx context.Context, req *GetCartReq) (*GetCartResp, error) {
	cc.log.WithContext(ctx).Infof("GetCart request: %+v", req)
	return cc.repo.GetCart(ctx, req)
}

func (cc *CartUsecase) EmptyCart(ctx context.Context, req *EmptyCartReq) (*EmptyCartResp, error) {
	cc.log.WithContext(ctx).Infof("EmptyCart request: %+v", req)
	return cc.repo.EmptyCart(ctx, req)
}

func (cc *CartUsecase) UpdateItem(ctx context.Context, req *UpdateItemReq) (*UpdateItemResp, error) {
	cc.log.WithContext(ctx).Infof("UpdateItem request: %+v", req)
	return cc.repo.UpdateItem(ctx, req)
}

func (cc *CartUsecase) RemoveItem(ctx context.Context, req *RemoveItemReq) (*RemoveItemResp, error) {
	cc.log.WithContext(ctx).Infof("RemoveItem request: %+v", req)
	return cc.repo.RemoveItem(ctx, req)
}
