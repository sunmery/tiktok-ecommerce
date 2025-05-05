package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"backend/constants"

	"github.com/google/uuid"
)

type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       string
}

type SubOrder struct {
	OrderId        int64
	SubOrderId     int64
	StreetAddress  string
	City           string
	State          string
	Country        string
	ZipCode        string
	Email          string
	MerchantId     string
	TotalAmount    float64
	PaymentStatus  constants.PaymentStatus
	ShippingStatus constants.ShippingStatus
	Currency       string
	Items          []*OrderItem
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type (
	GetOrdersRequest struct {
		UserId   uuid.UUID
		Page     uint32
		PageSize uint32
	}
	Orders struct {
		Orders []*ConsumerOrder
	}
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

// OrderItem 订单商品, 是以 JSON 存储到数据库中, 需要添加tags
// 最终是给前端展示, 所以使用小驼峰符合前端变量命名规范
type OrderItem struct {
	Item *CartItem `json:"item"`
	Cost float64   `json:"cost"`
}

type OrderResult struct {
	OrderId         int64
	FreezeId        int64
	ConsumerVersion int64
	MerchantVersion []int64
}

type PlaceOrderRequest struct {
	UserId     uuid.UUID
	Currency   string
	Address    *Address
	Email      string
	OrderItems []*OrderItem
}
type PlaceOrderReply struct {
	Order *OrderResult
}

type GetConsumerSubOrderDetailRequest struct {
	UserId     uuid.UUID
	SubOrderId int64
}

// 确认收货
type (
	ConfirmReceivedRequest struct {
		UserId  uuid.UUID
		OrderId int64
	}
	ConfirmReceivedReply struct{}
)

// 查询订单状态
type (
	GetShipOrderStatusRequest struct {
		UserId     uuid.UUID
		SubOrderId int64
	}
	ShippingAddress struct {
		AddressType   string `json:"addressType"`
		City          string `json:"city"`
		ContactPerson string `json:"contactPerson"`
		ContactPhone  string `json:"contactPhone"`
		Country       string `json:"country"`
		State         string `json:"state"`
		StreetAddress string `json:"streetAddress"`
		ZipCode       string `json:"zipCode"`
	}
	ReceiverAddress struct {
		City           string    `json:"city"`
		Country        string    `json:"country"`
		CreatedAt      time.Time `json:"createdAt"`
		Email          string    `json:"email"`
		ID             int64     `json:"id"`
		PaymentStatus  string    `json:"paymentStatus"`
		ShippingStatus string    `json:"shippingStatus"`
		State          string    `json:"state"`
		StreetAddress  string    `json:"streetAddress"`
		UpdatedAt      time.Time `json:"updatedAt"`
		UserID         string    `json:"userId"`
		ZipCode        string    `json:"zipCode"`
	}
	GetShipOrderStatusReply struct {
		Id              int64                    // 物流 ID
		SubOrderId      int64                    // 子订单 ID
		TrackingNumber  string                   // 物流单号
		Carrier         string                   // 物流公司
		ShippingStatus  constants.ShippingStatus // 货运状态
		Delivery        time.Time                // 送达时间
		ShippingFee     float64                  // 运费
		ReceiverAddress ReceiverAddress          // 用户地址
		ShippingAddress ShippingAddress          // 商家地址
		CreatedAt       time.Time
		UpdatedAt       time.Time // 更新时间
	}
)

type (
	GetConsumerOrdersWithSubordersRequest struct {
		UserId  uuid.UUID
		OrderId int64
	}

	GetConsumerOrdersWithSubordersReply struct {
		Orders []*SubOrder
	}
)

type (
	GetConsumerOrdersRequest struct {
		UserId   uuid.UUID
		Page     uint32
		PageSize uint32
	}
	ConsumerOrderItem struct {
		Cost float64
		Item *CartItem
	}
	ConsumerOrder struct {
		OrderId        int64
		SubOrderID     int64
		Items          []*ConsumerOrderItem
		Address        Address
		Currency       string
		PaymentStatus  constants.PaymentStatus
		ShippingStatus constants.ShippingStatus
		Email          string
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
	GetConsumerOrdersReply struct {
		SubOrders []*ConsumerOrder
	}
)

type ConsumerOrderRepo interface {
	PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*PlaceOrderReply, error)
	GetConsumerOrders(ctx context.Context, req *GetConsumerOrdersRequest) (*GetConsumerOrdersReply, error)
	GetConsumerOrdersWithSuborders(ctx context.Context, req *GetConsumerOrdersWithSubordersRequest) (*GetConsumerOrdersWithSubordersReply, error)
	GetConsumerSubOrderDetail(ctx context.Context, req *GetConsumerSubOrderDetailRequest) (*ConsumerOrder, error)
	GetShipOrderStatus(ctx context.Context, req *GetShipOrderStatusRequest) (*GetShipOrderStatusReply, error)
	ConfirmReceived(ctx context.Context, req *ConfirmReceivedRequest) (*ConfirmReceivedReply, error)
}

type ConsumerOrderUsecase struct {
	repo ConsumerOrderRepo
	log  *log.Helper
}

func NewConsumerOrderUsecase(repo ConsumerOrderRepo, logger log.Logger) *ConsumerOrderUsecase {
	return &ConsumerOrderUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (oc *ConsumerOrderUsecase) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*PlaceOrderReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order req:%+v", req)
	return oc.repo.PlaceOrder(ctx, req)
}

func (oc *ConsumerOrderUsecase) GetConsumerOrders(ctx context.Context, req *GetConsumerOrdersRequest) (*GetConsumerOrdersReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetConsumerOrders:%+v", req)
	return oc.repo.GetConsumerOrders(ctx, req)
}

func (oc *ConsumerOrderUsecase) GetConsumerOrdersWithSuborders(ctx context.Context, req *GetConsumerOrdersWithSubordersRequest) (*GetConsumerOrdersWithSubordersReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetConsumerOrdersWithSuborders:%+v", req)
	return oc.repo.GetConsumerOrdersWithSuborders(ctx, req)
}

func (oc *ConsumerOrderUsecase) GetConsumerSubOrderDetail(ctx context.Context, req *GetConsumerSubOrderDetailRequest) (*ConsumerOrder, error) {
	oc.log.WithContext(ctx).Debugf("biz/order getorder req:%+v", req)
	return oc.repo.GetConsumerSubOrderDetail(ctx, req)
}

func (oc *ConsumerOrderUsecase) ConfirmReceived(ctx context.Context, req *ConfirmReceivedRequest) (*ConfirmReceivedReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order ConfirmReceived req:%+v", req)
	return oc.repo.ConfirmReceived(ctx, req)
}

func (oc *ConsumerOrderUsecase) GetShipOrderStatus(ctx context.Context, req *GetShipOrderStatusRequest) (*GetShipOrderStatusReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetOrderStatus req:%+v", req)
	return oc.repo.GetShipOrderStatus(ctx, req)
}
