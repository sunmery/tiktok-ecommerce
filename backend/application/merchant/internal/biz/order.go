package biz

import (
	"context"
	"time"

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

// CartItem 购物车商品, 是以 JSON 存储到数据库中, 需要添加tags
// 最终是给前端展示, 所以使用小驼峰符合前端变量命名规范
type CartItem struct {
	MerchantId uuid.UUID `json:"merchantId"`
	ProductId  uuid.UUID `json:"productId"`
	Quantity   uint32    `json:"quantity"`
}

type OrderItem struct {
	Item *CartItem `json:"item"`
	Cost float64   `json:"cost"`
}
type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       string
}
type Order struct {
	OrderID       int64
	UserID        uuid.UUID
	Currency      string
	Address       *Address
	Email         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	SubOrders     []*SubOrder
	PaymentStatus PaymentStatus // 支付状态
}
type SubOrder struct {
	ID          int64
	MerchantID  uuid.UUID
	TotalAmount float64
	Currency    string
	Status      string
	Items       []*OrderItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type (
	GetMerchantOrdersReq struct {
		UserID   uuid.UUID
		Page     uint32 // 分页页码，从1开始
		PageSize uint32 // 每页数量
	}

	GetMerchantOrdersReply struct {
		Orders []*Order
	}
)

func (oc *OrderUsecase) GetMerchantOrders(ctx context.Context, req *GetMerchantOrdersReq) (*GetMerchantOrdersReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetMerchantOrders:%+v", req)
	return oc.repo.GetMerchantOrders(ctx, req)
}
