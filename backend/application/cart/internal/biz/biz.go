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
	UpsertItem(ctx context.Context, req *UpsertItemReq) (*UpsertItemResp, error)
	GetCart(ctx context.Context, req *GetCartReq) (*GetCartResp, error)
	EmptyCart(ctx context.Context, req *EmptyCartReq) (*EmptyCartResp, error)
	RemoveCartItem(ctx context.Context, req *RemoveCartItemReq) (*RemoveCartItemResp, error)
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

func (cc *CartUsecase) UpsertItem(ctx context.Context, req *UpsertItemReq) (*UpsertItemResp, error) {
	cc.log.WithContext(ctx).Infof("UpsertItem request: %+v", req)
	return cc.repo.UpsertItem(ctx, req)
}

func (cc *CartUsecase) GetCart(ctx context.Context, req *GetCartReq) (*GetCartResp, error) {
	cc.log.WithContext(ctx).Infof("GetCart request: %+v", req)
	return cc.repo.GetCart(ctx, req)
}

func (cc *CartUsecase) EmptyCart(ctx context.Context, req *EmptyCartReq) (*EmptyCartResp, error) {
	cc.log.WithContext(ctx).Infof("EmptyCart request: %+v", req)
	return cc.repo.EmptyCart(ctx, req)
}

func (cc *CartUsecase) RemoveCartItem(ctx context.Context, req *RemoveCartItemReq) (*RemoveCartItemResp, error) {
	cc.log.WithContext(ctx).Infof("RemoveCartItem request: %+v", req)
	return cc.repo.RemoveCartItem(ctx, req)
}
