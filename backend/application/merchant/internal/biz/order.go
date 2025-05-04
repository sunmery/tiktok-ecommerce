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
	Name       string    `json:"name"`
	Picture    string    `json:"picture"`
}

type OrderItem struct {
	Cost           float64
	Item           *CartItem
	Email          string
	Address        Address
	UserID         uuid.UUID
	SubOrderID     int64
	TotalAmount    float64
	Currency       string
	PaymentStatus  constants.PaymentStatus
	ShippingStatus constants.ShippingStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       string
}

type SubOrder struct {
	OrderID   int64
	Items     []*OrderItem
	CreatedAt time.Time
}

// 发货
type (
	CreateOrderShipReq struct {
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
	CreateOrderShipResp struct {
		Id        int64     // 物流 ID
		CreatedAt time.Time // 发货时间
	}
)

type (
	GetMerchantByOrderIdReq struct {
		OrderId int64
	}
	GetMerchantByOrderIdReply struct {
		MerchantId uuid.UUID
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

// 更新订单状态
type (
	UpdateOrderShippingStatusReq struct {
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
	UpdateOrderShippingStatusResply struct {
		ID        int64     // 物流 ID
		UpdatedAt time.Time // 更新时间
	}
)

type OrderUsecase struct {
	repo OrderRepo
	log  *log.Helper
}

// OrderRepo 订单域方法
type OrderRepo interface {
	GetMerchantByOrderId(ctx context.Context, req *GetMerchantByOrderIdReq) (*GetMerchantByOrderIdReply, error)
	GetMerchantOrders(ctx context.Context, req *GetMerchantOrdersReq) (*GetMerchantOrdersReply, error)
	CreateOrderShip(ctx context.Context, req *CreateOrderShipReq) (*CreateOrderShipResp, error)
	UpdateOrderShippingStatus(ctx context.Context, req *UpdateOrderShippingStatusReq) (*UpdateOrderShippingStatusResply, error)
}

func NewOrderUsecase(repo OrderRepo, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (oc *OrderUsecase) GetMerchantByOrderId(ctx context.Context, req *GetMerchantByOrderIdReq) (*GetMerchantByOrderIdReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetMerchantOrder:%+v", req)
	return oc.repo.GetMerchantByOrderId(ctx, req)
}

func (oc *OrderUsecase) GetMerchantOrders(ctx context.Context, req *GetMerchantOrdersReq) (*GetMerchantOrdersReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetMerchantOrders:%+v", req)
	return oc.repo.GetMerchantOrders(ctx, req)
}

func (oc *OrderUsecase) CreateOrderShip(ctx context.Context, req *CreateOrderShipReq) (*CreateOrderShipResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order CreateOrderShip req:%+v", req)
	return oc.repo.CreateOrderShip(ctx, req)
}

func (oc *OrderUsecase) UpdateOrderShippingStatus(ctx context.Context, req *UpdateOrderShippingStatusReq) (*UpdateOrderShippingStatusResply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order UpdateOrderShippingStatus req:%+v", req)
	return oc.repo.UpdateOrderShippingStatus(ctx, req)
}
