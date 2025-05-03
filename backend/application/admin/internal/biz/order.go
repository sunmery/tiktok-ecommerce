package biz

import (
	"context"
	"time"

	userv1 "backend/api/user/v1"

	"github.com/google/uuid"

	"backend/constants"

	cartv1 "backend/api/cart/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type SubOrder struct {
	ID             int64
	MerchantID     uuid.UUID
	TotalAmount    float64
	Currency       string
	PaymentStatus  constants.PaymentStatus
	ShippingStatus constants.ShippingStatus
	Items          []*OrderItem
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
type (
	GetAllOrdersReq struct {
		Page     uint32 // 分页页码，从1开始
		PageSize uint32 // 每页数量
	}
	OrderItem struct {
		Cost            float64
		Item            *cartv1.CartItem
		Email           string
		ConsumerAddress userv1.ConsumerAddress
		UserID          uuid.UUID
		SubOrderID      int64
		TotalAmount     float64
		Currency        string
		PaymentStatus   constants.PaymentStatus
		ShippingStatus  constants.ShippingStatus
		CreatedAt       time.Time
		UpdatedAt       time.Time
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

func (cc *AdminOrderUsecase) GetAllOrders(ctx context.Context, req *GetAllOrdersReq) (*GetAllOrdersReply, error) {
	cc.log.WithContext(ctx).Debugf("GetAllOrders request: %+v", req)
	return cc.repo.GetAllOrders(ctx, req)
}
