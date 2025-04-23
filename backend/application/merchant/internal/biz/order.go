package biz

import (
	"context"
	"encoding/json"
	"time"

	"backend/constants"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/google/uuid"
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

type SubOrder struct {
	OrderID        int64
	SubOrderID     int64
	MerchantID     uuid.UUID
	TotalAmount    float64
	Currency       string
	Status         constants.PaymentStatus
	ShippingStatus constants.ShippingStatus
	Items          []*OrderItem
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// 发货
type (
	ShipOrderReq struct {
		Id              int64                    // 物流 ID
		MerchantID      uuid.UUID                // 商家 ID
		SubOrderId      int64                    // 子订单 ID
		TrackingNumber  string                   // 物流单号
		Carrier         string                   // 物流承运商
		ShippingStatus  constants.ShippingStatus // 物流状态
		Delivery        time.Time                // 送达时间
		ShippingAddress json.RawMessage          // 发货地址
		ReceiverAddress []byte                   // 收货地址
		ShippingFee     float64                  // 运费
		CreatedAt       time.Time                // 发货时间
		UpdatedAt       time.Time                // 更新时间
	}
	ShipOrderResp struct {
		Id        int64     // 物流 ID
		CreatedAt time.Time // 发货时间
	}
)

type (
	GetMerchantOrdersReq struct {
		UserID   uuid.UUID
		Page     uint32 // 分页页码，从1开始
		PageSize uint32 // 每页数量
	}

	GetMerchantOrdersReply struct {
		Orders []*SubOrder
	}
)

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

// OrderRepo 订单域方法
type OrderRepo interface {
	GetMerchantOrders(ctx context.Context, req *GetMerchantOrdersReq) (*GetMerchantOrdersReply, error)
	ShipOrder(ctx context.Context, req *ShipOrderReq) (*ShipOrderResp, error)
}

func (oc *OrderUsecase) GetMerchantOrders(ctx context.Context, req *GetMerchantOrdersReq) (*GetMerchantOrdersReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetMerchantOrders:%+v", req)
	return oc.repo.GetMerchantOrders(ctx, req)
}

func (oc *OrderUsecase) ShipOrder(ctx context.Context, req *ShipOrderReq) (*ShipOrderResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order ShipOrder req:%+v", req)
	return oc.repo.ShipOrder(ctx, req)
}
