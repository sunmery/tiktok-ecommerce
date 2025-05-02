package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"backend/constants"

	v1 "backend/api/order/v1"

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
	GetOrdersReq struct {
		UserId   uuid.UUID
		Page     uint32
		PageSize uint32
	}
	Orders struct {
		Orders []*v1.Order
	}
)

// // GetMerchantOrdersReq 获取商家订单
// type (
// 	GetMerchantOrdersReq struct {
// 		MerchantId uuid.UUID
// 		Page       uint32
// 		PageSize   uint32
// 	}
// 	MerchantOrdersSubOrder struct {
// 		OrderID        int64
// 		SubOrderID     int64
// 		MerchantID     uuid.UUID
// 		TotalAmount    float64
// 		Currency       string
// 		PaymentStatus  constants.PaymentStatus
// 		ShippingStatus constants.ShippingStatus
// 		Items          []*OrderItem
// 		CreatedAt      time.Time
// 		UpdatedAt      time.Time
// 	}
// 	GetMerchantOrdersReply struct {
// 		Orders []*MerchantOrdersSubOrder
// 	}
// )

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

type (
	GetAllOrdersReq struct {
		Page     uint32 // 分页页码，从1开始
		PageSize uint32 // 每页数量
	}
	GetAllOrdersReply struct {
		Orders []*SubOrder
	}
)

type MarkOrderPaidResp struct{}

type MarkOrderPaidReq struct {
	UserId  uuid.UUID
	OrderId int64
}

type GetOrderReq struct {
	UserId  uuid.UUID
	OrderId int64
}

// 确认收货
type (
	ConfirmReceivedReq struct {
		UserId  uuid.UUID
		OrderId int64
	}
	ConfirmReceivedResp struct{}
)

// 查询订单状态
type (
	GetShipOrderStatusReq struct {
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
		// EstimatedDelivery string                     // 送达时间
		CreatedAt time.Time
		UpdatedAt time.Time // 更新时间
	}
)

type (
	GetUserOrdersWithSubordersReq struct {
		UserId  uuid.UUID
		OrderId int64
	}
	Suborder struct {
		OrderId        int64
		SubOrderId     int64
		StreetAddress  string
		City           string
		State          string
		Country        string
		ZipCode        string
		Email          string
		MerchantId     string
		PaymentStatus  constants.PaymentStatus
		ShippingStatus constants.ShippingStatus
		TotalAmount    float64
		Currency       string
		Items          []*OrderItem
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
	GetUserOrdersWithSubordersReply struct {
		Orders []*Suborder
	}
)

type OrderRepo interface {
	PlaceOrder(ctx context.Context, req *PlaceOrderReq) (*PlaceOrderResp, error)
	GetOrders(ctx context.Context, req *GetOrdersReq) (*Orders, error)
	GetUserOrdersWithSuborders(ctx context.Context, req *GetUserOrdersWithSubordersReq) (*GetUserOrdersWithSubordersReply, error)
	GetAllOrders(ctx context.Context, req *GetAllOrdersReq) (*GetAllOrdersReply, error)

	MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error)
	GetOrder(ctx context.Context, req *GetOrderReq) (*v1.Order, error)
	GetShipOrderStatus(ctx context.Context, req *GetShipOrderStatusReq) (*GetShipOrderStatusReply, error)
	ConfirmReceived(ctx context.Context, req *ConfirmReceivedReq) (*ConfirmReceivedResp, error)
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

func (oc *OrderUsecase) GetOrders(ctx context.Context, req *GetOrdersReq) (*Orders, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetOrders:%+v", req)
	return oc.repo.GetOrders(ctx, req)
}

func (oc *OrderUsecase) GetUserOrdersWithSuborders(ctx context.Context, req *GetUserOrdersWithSubordersReq) (*GetUserOrdersWithSubordersReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetUserOrdersWithSuborders:%+v", req)
	return oc.repo.GetUserOrdersWithSuborders(ctx, req)
}

func (oc *OrderUsecase) GetAllOrders(ctx context.Context, req *GetAllOrdersReq) (*GetAllOrdersReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetAllOrders:%+v", req)
	return oc.repo.GetAllOrders(ctx, req)
}

func (oc *OrderUsecase) MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order MarkOrderPaid req:%+v", req)
	return oc.repo.MarkOrderPaid(ctx, req)
}

func (oc *OrderUsecase) GetOrder(ctx context.Context, req *GetOrderReq) (*v1.Order, error) {
	oc.log.WithContext(ctx).Debugf("biz/order getorder req:%+v", req)
	return oc.repo.GetOrder(ctx, req)
}

func (oc *OrderUsecase) ConfirmReceived(ctx context.Context, req *ConfirmReceivedReq) (*ConfirmReceivedResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order ConfirmReceived req:%+v", req)
	return oc.repo.ConfirmReceived(ctx, req)
}

func (oc *OrderUsecase) GetShipOrderStatus(ctx context.Context, req *GetShipOrderStatusReq) (*GetShipOrderStatusReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetOrderStatus req:%+v", req)
	return oc.repo.GetShipOrderStatus(ctx, req)
}
