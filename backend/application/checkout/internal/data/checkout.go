package data

import (
	cartv1 "backend/api/cart/v1"
	v1 "backend/api/order/v1"
	paymentv1 "backend/api/payment/v1"
	userv1 "backend/api/user/v1"
	"backend/application/checkout/internal/biz"
	"backend/constants"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type checkoutRepo struct {
	data *Data
	log  *log.Helper
}

func (c checkoutRepo) Checkout(ctx context.Context, req *biz.CheckoutRequest) (*biz.CheckoutReply, error) {
	// 传递用户ID到购物车微服务
	ctx = metadata.AppendToClientContext(ctx, constants.UserId, req.UserId.String())
	cartItems, err := c.data.cartv1.GetCart(ctx, &cartv1.GetCartReq{
		// UserId: req.UserId.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("get cart failed: %w", err)
	}

	fmt.Printf("cartItems: %+v\n", cartItems)

	// 转成订单商品类型
	var orderItems []*v1.OrderItem
	var amount = 0.0
	for _, item := range cartItems.Cart.Items {
		amount += float64(item.Quantity) * item.Price
		orderItems = append(orderItems, &v1.OrderItem{
			Item: &cartv1.CartItem{
				MerchantId: item.MerchantId,
				ProductId:  item.ProductId,
				Quantity:   item.Quantity,
				Price:      item.Price,
			},
			Cost: float64(item.Quantity) * item.Price,
		})
	}

	if len(orderItems) == 0 {
		return nil, status.Error(codes.InvalidArgument, "购物车为空")
	}

	// 调用订单微服务创建订单
	order, orderErr := c.data.orderv1.PlaceOrder(ctx, &v1.PlaceOrderReq{
		Currency: req.Currency,
		Address: &userv1.Address{
			Id:            req.Address.Id,
			UserId:        req.UserId.String(),
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
			StreetAddress: req.Address.StreetAddress,
		},
		Email:      req.Email,
		OrderItems: orderItems,
	})
	if orderErr != nil || order == nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}
	if order.Order == nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("订单服务返回无效数据: %v", orderErr))
	}

	// 调用支付微服务生成支付URL
	payment, paymentErr := c.data.paymentv1.CreatePayment(ctx, &paymentv1.CreatePaymentReq{
		OrderId:       order.Order.OrderId,
		Currency:      req.Currency,
		Amount:        strconv.Itoa(int(amount)),
		PaymentMethod: req.PaymentMethod,
	})
	if paymentErr != nil || payment == nil {
		// 订单回滚
		// _, _ = c.data.orderv1.CancelOrder(ctx, &v1.CancelOrderReq{
		// 	OrderId: orderResp.Order.OrderId,
		// })
		return nil, status.Error(codes.Internal, fmt.Sprintf("创建支付失败: %v", paymentErr))
	}

	return &biz.CheckoutReply{
		OrderId:    order.Order.OrderId,
		PaymentId:  payment.PaymentId,
		PaymentURL: order.Url,
	}, err
}

func NewCheckoutRepo(data *Data, logger log.Logger) biz.CheckoutRepo {
	return &checkoutRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
