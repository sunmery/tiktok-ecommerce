package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       uint32
}

// CartItem 购物车商品
type CartItem struct {
	ProductId uint32 // 商品ID
	Quantity  uint32 // 商品数量
}

// OrderItem 订单项
type OrderItem struct {
	Item CartItem
	Cost float32
}

type PlaceOrderReq struct {
	UserId     string
	Currency   string
	Address    Address
	Email      string
	OrderItems []*OrderItem
}

type PlaceOrderResp struct {
	OrderId uuid.UUID
}

type ListOrderReq struct {
}

type ListOrderResp struct {
}

type MarkOrderPaidReq struct {
}

type MarkOrderPaidResp struct {
}

type OrderRepo interface {
	PlaceOrder(ctx context.Context, req *PlaceOrderReq) (*PlaceOrderResp, error)
	ListOrder(ctx context.Context, req *ListOrderReq) (*ListOrderResp, error)
	MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error)
}

type OrderUsecase struct {
	repo OrderRepo
	log  *log.Helper
}

func NewOrderUsecase(repo OrderRepo, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *OrderUsecase) PlaceOrder(ctx context.Context, req *PlaceOrderReq) (*PlaceOrderResp, error) {
	uc.log.WithContext(ctx).Debugf("req: %+v", req)
	return uc.repo.PlaceOrder(ctx, req)
}
func (uc *OrderUsecase) ListOrder(ctx context.Context, req *ListOrderReq) (*ListOrderResp, error) {
	uc.log.WithContext(ctx).Debugf("req: %+v", req)
	return uc.repo.ListOrder(ctx, req)
}
func (uc *OrderUsecase) MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error) {
	uc.log.WithContext(ctx).Debugf("req: %+v", req)
	return uc.repo.MarkOrderPaid(ctx, req)
}
