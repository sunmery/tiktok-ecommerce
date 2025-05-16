package biz

import (
	"context"
	"time"

	"github.com/google/uuid"

	"backend/constants"

	"github.com/go-kratos/kratos/v2/log"
)

type ConsumerAddress struct {
	City          string
	State         string
	Country       string
	ZipCode       string
	StreetAddress string
}
type SubOrder struct {
	OrderID         int64
	SubOrderID      int64
	TotalAmount     float64
	ConsumerId      uuid.UUID
	ConsumerAddress *ConsumerAddress
	ConsumerEmail   string
	Currency        constants.Currency
	PaymentStatus   constants.PaymentStatus
	ShippingStatus  constants.ShippingStatus
	SubOrderItems   []*SubOrderItem
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
type CartItem struct {
	MerchantId uuid.UUID
	ProductId  uuid.UUID
	Quantity   uint32
	Name       string
	Picture    string
}

type (
	GetAllOrdersReq struct {
		Page     uint32 // 分页页码，从1开始
		PageSize uint32 // 每页数量
	}
	SubOrderItem struct {
		Cost float64
		Item *CartItem
	}

	GetAllOrdersReply struct {
		Orders []*SubOrder
	}
)

type AdminOrderRepo interface {
	GetAllOrders(ctx context.Context, req *GetAllOrdersReq) (*GetAllOrdersReply, error)
}

type AdminOrderUsecase struct {
	repo AdminOrderRepo
	log  *log.Helper
}

func NewAdminOrderUsecase(repo AdminOrderRepo, logger log.Logger) *AdminOrderUsecase {
	return &AdminOrderUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (oc *AdminOrderUsecase) GetAllOrders(ctx context.Context, req *GetAllOrdersReq) (*GetAllOrdersReply, error) {
	oc.log.WithContext(ctx).Debugf("GetAllOrders request: %+v", req)
	return oc.repo.GetAllOrders(ctx, req)
}
