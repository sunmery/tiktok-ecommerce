package data

import (
	"context"
	"errors"
	"fmt"

	paymentv1 "backend/api/payment/v1"

	productv1 "backend/api/product/v1"

	"github.com/google/uuid"

	cartv1 "backend/api/cart/v1"
	v1 "backend/api/order/v1"
	userv1 "backend/api/user/v1"
	"backend/application/checkout/internal/biz"
	"backend/constants"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type checkoutRepo struct {
	data *Data
	log  *log.Helper
}

func (c checkoutRepo) Checkout(ctx context.Context, req *biz.CheckoutRequest) (*biz.CheckoutReply, error) {
	// 传递用户ID到购物车微服务
	ctx = metadata.AppendToClientContext(ctx, constants.UserId, req.UserId.String())
	carts, err := c.data.cartv1.GetCart(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}

	if len(carts.Items) == 0 {
		fmt.Println("No cart items found")
		return nil, status.Error(codes.InvalidArgument, "购物车为空")
	}

	var productIds []string
	var merchantIds []string
	var Pictures []string
	var Names []string
	for _, c := range carts.Items {
		productIds = append(productIds, c.ProductId)
		merchantIds = append(merchantIds, c.MerchantId)
		Pictures = append(Pictures, c.Picture)
		Names = append(Names, c.Name)
	}

	// 从商品微服务获取商品信息, 例如价格
	products, perr := c.data.productv1.GetProductsBatch(ctx, &productv1.GetProductsBatchRequest{
		ProductIds:  productIds,
		MerchantIds: merchantIds,
	})
	if perr != nil {
		return nil, perr
	}

	var cartItems []cartv1.CartItem
	var orderItems []*v1.OrderItem
	var amount float64
	for _, cart := range carts.Items {
		for _, p := range products.Items {
			productId, err := uuid.Parse(p.Id)
			if err != nil {
				return nil, err
			}
			merchantId, err := uuid.Parse(p.MerchantId)
			if err != nil {
				return nil, err
			}

			if merchantId.String() == cart.MerchantId && productId.String() == cart.ProductId {
				var picture string
				for _, image := range p.Images {
					if image.IsPrimary {
						picture = image.Url
					}
				}
				cartItems = append(cartItems, cartv1.CartItem{
					MerchantId: cart.MerchantId,
					ProductId:  cart.ProductId,
					Quantity:   cart.Quantity,
					Name:       p.Name,
					Picture:    picture,
				})

				orderItems = append(orderItems, &v1.OrderItem{
					Item: &cartv1.CartItem{
						MerchantId: cart.MerchantId,
						ProductId:  cart.ProductId,
						Quantity:   cart.Quantity,
						Name:       cart.Name,
						Picture:    cart.Picture,
					},
					Cost: p.Price,
				})

				amount += p.Price * float64(cart.Quantity)
				fmt.Printf("amount: %v\n", amount)
			}
		}
	}

	// 获取用户地址
	address, err := c.data.userv1.GetAddress(ctx, &userv1.GetAddressRequest{
		AddressId: req.AddressId,
		UserId:    req.UserId.String(),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取用户地址失败: %v", err))
	}

	log.Debugf("cartItems: %+v", cartItems)

	// 获取用户支付卡信息
	creditCard, err := c.data.userv1.GetCreditCard(ctx, &userv1.GetCreditCardRequest{
		Id:     req.CreditCardId,
		UserId: req.UserId.String(),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取用户支付卡信息失败: %v", err))
	}

	// 调用订单微服务创建订单
	order, orderErr := c.data.orderv1.PlaceOrder(ctx, &v1.PlaceOrderReq{
		Currency: creditCard.Currency,
		Address: &userv1.Address{
			Id:            address.Id,
			UserId:        address.UserId,
			City:          address.City,
			State:         address.State,
			Country:       address.Country,
			ZipCode:       address.ZipCode,
			StreetAddress: address.StreetAddress,
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
	payment, paymentErr := c.data.paymentv1.CreatePayment(ctx, &paymentv1.CreatePaymentRequest{
		OrderId:   order.Order.OrderId,
		UserId:    req.UserId.String(),
		Amount:    fmt.Sprintf("%.2f", amount),
		Currency:  creditCard.Currency,
		Subject:   "支付测试",
		ReturnUrl: "",
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
		PaymentURL: payment.PayUrl,
	}, err
}

func NewCheckoutRepo(data *Data, logger log.Logger) biz.CheckoutRepo {
	return &checkoutRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
