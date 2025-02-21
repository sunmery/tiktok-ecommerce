package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/runtime/protoimpl"
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
	ZipCode       uint32
}
type Order struct {
	OrderItems   []*OrderItem
	OrderId      string
	UserId       uint32
	UserCurrency string
	Address      *Address
	Email        string
	CreatedAt    uint32
	// 新增支付状态
	PaymentStatus string
}

// CartItem 购物车商品
type CartItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     uint32                 `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"` // 商品ID
	Quantity      uint32                 `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`                    // 商品数量
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

type OrderItem struct {
	Item *CartItem
	Cost float32
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
	UserId uint32
}

type ListOrderResp struct {
	Orders []*Order
}
type MarkOrderPaidResp struct{}

type MarkOrderPaidReq struct {
	UserId  uint32
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
