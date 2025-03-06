package data

import (
	v1 "backend/api/order/v1"
	paymentv1 "backend/api/payment/v1"
	userv1 "backend/api/user/v1"
	"backend/application/checkout/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type checkoutRepo struct {
	data *Data
	log  *log.Helper
}

func (c checkoutRepo) Checkout(ctx context.Context, req *biz.CheckoutRequest) (*biz.CheckoutReply, error) {
	// TODO implement me
	// c.data.DB(ctx).

	// 获取购物车商品
	cartItems, err :=

		// 创建订单
	order, err := c.data.orderv1.PlaceOrder(ctx, &v1.PlaceOrderReq{
		Currency: "CNY",
		Address: &userv1.Address{
			Id:            0,
			UserId:        req.UserId.String(),
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
			StreetAddress: req.Address.StreetAddress,
		},
		Email:      req.Email,
		OrderItems: orderv1.Or,
	})
	if err != nil {
		return nil, err
	}

	payment, err := c.data.paymentv1.CreatePayment(ctx, &paymentv1.CreatePaymentReq{
		OrderId:       order.Order.OrderId,
		Currency:      "CNY",
		Amount:        req.,
		PaymentMethod: "",
	})
	if err != nil {
		return nil, err
	}

	return &biz.CheckoutReply{
		OrderId:       order.Order.OrderId,
		TransactionId: payment.PaymentId,
		PaymentURL:    order.Url,
	}, err
}

func NewCheckoutRepo(data *Data, logger log.Logger) biz.CheckoutRepo {
	return &checkoutRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
