package biz

import (
	"context"
	"time"

	merchantAddressv1 "backend/api/merchant/address/v1"
	userv1 "backend/api/user/v1"

	"backend/constants"

	v1 "backend/api/order/v1"

	"github.com/go-kratos/kratos/v2/log"
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
	GetConsumerOrdersReq struct {
		UserId   uuid.UUID
		Page     uint32
		PageSize uint32
	}
	Orders struct {
		Orders []*v1.Order
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
	OrderId int64
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

// 查询订单状态
type (
	GetOrderStatusReq struct {
		UserId  uuid.UUID
		OrderId int64
	}
	GetOrderStatusReply struct {
		OrderId         int64                              // 主订单 ID
		SubOrderId      int64                              // 子订单 ID
		UserAddress     *userv1.ConsumerAddress            // 用户地址
		MerchantAddress *merchantAddressv1.MerchantAddress // 商家地址
		TrackingNumber  string                             // 物流单号
		Carrier         string                             // 物流公司
		// EstimatedDelivery string                     // 预计送达时间
		PaymentStatus  constants.PaymentStatus  // 支付状态
		ShippingStatus constants.ShippingStatus // 货运状态
	}
)

// 更新订单状态
type (
	UpdateOrderStatusReq struct {
		UserId  uuid.UUID
		OrderId int64
		Status  constants.PaymentStatus
	}
	UpdateOrderStatusResp struct{}
)

// 确认收货
type (
	ConfirmReceivedReq struct {
		UserId  uuid.UUID
		OrderId int64
	}
	ConfirmReceivedResp struct{}
)

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

type OrderRepo interface {
	PlaceOrder(ctx context.Context, req *PlaceOrderReq) (*PlaceOrderResp, error)
	GetConsumerOrders(ctx context.Context, req *GetConsumerOrdersReq) (*Orders, error)
	GetAllOrders(ctx context.Context, req *GetAllOrdersReq) (*GetAllOrdersReply, error)

	MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error)
	GetOrder(ctx context.Context, req *GetOrderReq) (*v1.Order, error)

	GetOrderStatus(ctx context.Context, req *GetOrderStatusReq) (*GetOrderStatusReply, error)
	UpdateOrderStatus(ctx context.Context, req *UpdateOrderStatusReq) (*UpdateOrderStatusResp, error)
	ConfirmReceived(ctx context.Context, req *ConfirmReceivedReq) (*ConfirmReceivedResp, error)
}

func (oc *OrderUsecase) PlaceOrder(ctx context.Context, req *PlaceOrderReq) (*PlaceOrderResp, error) {
	// oc.log.WithContext(ctx).Debugf("biz/order req:%+v", req)
	return oc.repo.PlaceOrder(ctx, req)
}

func (oc *OrderUsecase) GetConsumerOrders(ctx context.Context, req *GetConsumerOrdersReq) (*Orders, error) {
	// oc.log.WithContext(ctx).Debugf("biz/order GetConsumerOrders:%+v", req)
	return oc.repo.GetConsumerOrders(ctx, req)
}

func (oc *OrderUsecase) GetAllOrders(ctx context.Context, req *GetAllOrdersReq) (*GetAllOrdersReply, error) {
	// oc.log.WithContext(ctx).Debugf("biz/order GetAllOrders:%+v", req)
	return oc.repo.GetAllOrders(ctx, req)
}

func (oc *OrderUsecase) MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error) {
	// oc.log.WithContext(ctx).Debugf("biz/order MarkOrderPaid req:%+v", req)
	return oc.repo.MarkOrderPaid(ctx, req)
}

func (oc *OrderUsecase) GetOrder(ctx context.Context, req *GetOrderReq) (*v1.Order, error) {
	// oc.log.WithContext(ctx).Debugf("biz/order getorder req:%+v", req)
	return oc.repo.GetOrder(ctx, req)
}

func (oc *OrderUsecase) GetOrderStatus(ctx context.Context, req *GetOrderStatusReq) (*GetOrderStatusReply, error) {
	oc.log.WithContext(ctx).Debugf("biz/order GetOrderStatus req:%+v", req)
	return oc.repo.GetOrderStatus(ctx, req)
}

func (oc *OrderUsecase) UpdateOrderStatus(ctx context.Context, req *UpdateOrderStatusReq) (*UpdateOrderStatusResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order UpdateOrderStatus req:%+v", req)
	return oc.repo.UpdateOrderStatus(ctx, req)
}

func (oc *OrderUsecase) ConfirmReceived(ctx context.Context, req *ConfirmReceivedReq) (*ConfirmReceivedResp, error) {
	// oc.log.WithContext(ctx).Debugf("biz/order ConfirmReceived req:%+v", req)
	return oc.repo.ConfirmReceived(ctx, req)
}
