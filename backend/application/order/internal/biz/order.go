package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       string
}

//	type Order struct {
//		OrderItems []*OrderItem
//		OrderId    string
//		UserId     uuid.UUID
//		Currency   string
//		Address    *Address
//		Email      string
//		CreatedAt  time.Time
//		// 新增支付状态
//		PaymentStatus int
//	}
type SubOrder struct {
	ID          string
	MerchantID  uuid.UUID
	TotalAmount string // 金额建议用字符串避免精度丢失
	Currency    string
	Status      string
	Items       []OrderItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type Order struct {
	OrderID       string
	UserID        uuid.UUID
	Currency      string
	Address       *Address
	Email         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	SubOrders     []*SubOrder
	PaymentStatus int // 根据业务逻辑后续补充
}

// CartItem 购物车商品
type CartItem struct {
	MerchantId uuid.UUID
	// 商品ID
	ProductId uuid.UUID
	// 商品数量
	Quantity uint32
}

type OrderItem struct {
	Item *CartItem
	Cost float64
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
type Pagination struct {
	Total       uint32 // 总记录数
	CurrentPage uint32 // 当前页码
	PageSize    uint32 // 每页数量
	TotalPages  uint32 // 总页数
}

type ListOrderReq struct {
	UserID        uuid.UUID
	DateRangeType string    // 支持：today/yesterday/last7days/custom
	StartTime     time.Time // 当DateRangeType=custom时必填
	EndTime       time.Time // 当DateRangeType=custom时必填
	Page          int       // 分页页码，从1开始
	PageSize      int       // 每页数量
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
