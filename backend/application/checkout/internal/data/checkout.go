package data

import (
	cartv1 "backend/api/cart/v1"
	v1 "backend/api/order/v1"
	paymentv1 "backend/api/payment/v1"
	userv1 "backend/api/user/v1"
	"backend/application/checkout/internal/biz"
	"context"
	"fmt"
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
	cartItems, err := c.data.cartv1.GetCart(ctx, &cartv1.GetCartReq{
		UserId: req.UserId.String(),
	})
	if err!= nil {
		return nil, fmt.Errorf("get cart failed: %w", err)
	}

	fmt.Println(cartItems)

	var orderItems []*v1.OrderItem
	for _, item := range cartItems.Cart.Items {
		orderItems = append(orderItems, &v1.OrderItem{
			Item: &cartv1.CartItem{
				MerchantId: item.MerchantId,
				ProductId:  item.ProductId,
				Quantity:   item.Quantity,
			},
			Cost: 0,
		})
	}

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
		Email: req.Email,
		OrderItems: ,
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
