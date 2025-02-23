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
	CheckCartItem(ctx context.Context, req *CheckCartItemReq) (*CheckCartItemResp, error)
	UncheckCartItem(ctx context.Context, req *UncheckCartItemReq) (*UncheckCartItemResp, error)
	CreateOrder(ctx context.Context, req *CreateOrderReq) (*CreateOrderResp, error)
	CreateCart(ctx context.Context, req *CreateCartReq) (*CreateCartResp, error)
	ListCarts(ctx context.Context, req *ListCartsReq) (*ListCartsResp, error)
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

func (cc *CartUsecase) CheckCartItem(ctx context.Context, req *CheckCartItemReq) (*CheckCartItemResp, error) {
	cc.log.WithContext(ctx).Infof("CheckCartItem request: %+v", req)
	return cc.repo.CheckCartItem(ctx, req)
}

func (cc *CartUsecase) UncheckCartItem(ctx context.Context, req *UncheckCartItemReq) (*UncheckCartItemResp, error) {
	cc.log.WithContext(ctx).Infof("UncheckCartItem request: %+v", req)
	return cc.repo.UncheckCartItem(ctx, req)
}

func (cc *CartUsecase) CreateOrder(ctx context.Context, req *CreateOrderReq) (*CreateOrderResp, error) {
	cc.log.WithContext(ctx).Infof("CreateOrder request: %+v", req)
	return cc.repo.CreateOrder(ctx, req)
}

func (cc *CartUsecase) CreateCart(ctx context.Context, req *CreateCartReq) (*CreateCartResp, error) {
	cc.log.WithContext(ctx).Infof("CreateCart request: %+v", req)
	return cc.repo.CreateCart(ctx, req)
}

func (cc *CartUsecase) ListCarts(ctx context.Context, req *ListCartsReq) (*ListCartsResp, error) {
	cc.log.WithContext(ctx).Infof("ListCarts request: %+v", req)
	return cc.repo.ListCarts(ctx, req)
}
