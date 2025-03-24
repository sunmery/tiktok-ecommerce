package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "PENDING"
	PaymentProcessing PaymentStatus = "PROCESSING"
	PaymentPaid       PaymentStatus = "PAID"
	PaymentFailed     PaymentStatus = "FAILED"
	PaymentCancelled  PaymentStatus = "CANCELLED"
)

type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       string
}

type SubOrder struct {
	ID          string
	MerchantID  uuid.UUID
	TotalAmount float64
	Currency    string
	Status      string
	Items       []*OrderItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type Order struct {
	OrderID       uuid.UUID
	UserID        uuid.UUID
	Currency      string
	Address       *Address
	Email         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	SubOrders     []*SubOrder
	PaymentStatus PaymentStatus // 支付状态
}

// CartItem 购物车商品, 是以 JSON 存储到数据库中, 需要添加tags
// 最终是给前端展示, 所以使用小驼峰符合前端变量命名规范
type CartItem struct {
	MerchantId uuid.UUID `json:"merchantId"`
	ProductId  uuid.UUID `json:"productId"`
	Quantity   uint32    `json:"quantity"`
}

// OrderItem 订单商品, 是以 JSON 存储到数据库中, 需要添加tags
// 最终是给前端展示, 所以使用小驼峰符合前端变量命名规范
type OrderItem struct {
	Item *CartItem `json:"item"`
	Cost float64   `json:"cost"`
}

type OrderResult struct {
	OrderId string
}

type PlaceOrderReq struct {
	UserId     uuid.UUID
	Currency   string
	Address    *Address
	Email      string
	OrderItems []*OrderItem
}
type PlaceOrderResp struct {
	Order *OrderResult
}

type ListOrderReq struct {
	UserID   uuid.UUID
	Page     uint32 // 分页页码，从1开始
	PageSize uint32 // 每页数量
}

type ListOrderResp struct {
	Orders []*Order
}
type MarkOrderPaidResp struct{}

type MarkOrderPaidReq struct {
	UserId  uuid.UUID
	OrderId string
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

func NewUserUsecase(repo OrderRepo, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (oc *OrderUsecase) PlaceOrder(ctx context.Context, req *PlaceOrderReq) (*PlaceOrderResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order req:%+v", req)
	return oc.repo.PlaceOrder(ctx, req)
}

func (oc *OrderUsecase) ListOrder(ctx context.Context, req *ListOrderReq) (*ListOrderResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order req:%+v", req)
	return oc.repo.ListOrder(ctx, req)
}

func (oc *OrderUsecase) MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order req:%+v", req)
	return oc.repo.MarkOrderPaid(ctx, req)
}
