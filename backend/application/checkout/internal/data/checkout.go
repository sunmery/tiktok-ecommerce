package data

import (
	"context"
	"errors"
	"fmt"

	banancev1 "backend/api/balancer/v1"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	paymentv1 "backend/api/payment/v1"

	productv1 "backend/api/product/v1"

	"github.com/google/uuid"

	cartv1 "backend/api/cart/v1"
	consumerv1 "backend/api/consumer/order/v1"
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
	var orderItems []*consumerv1.OrderItem
	var amount float64
	for _, cart := range carts.Items {
		for _, p := range products.Items {
			productId, err := uuid.Parse(p.Id)
			if err != nil {
				return nil, fmt.Errorf("invalid product ID format: %v", err)
			}
			merchantId, err := uuid.Parse(p.MerchantId)
			if err != nil {
				return nil, fmt.Errorf("invalid merchant ID format: %v", err)
			}

			if merchantId.String() == cart.MerchantId && productId.String() == cart.ProductId {
				itemPrice := p.Price
				itemQuantity := cart.Quantity
				itemTotal := itemPrice * float64(itemQuantity)

				cartItems = append(cartItems, cartv1.CartItem{
					MerchantId: cart.MerchantId,
					ProductId:  cart.ProductId,
					Quantity:   itemQuantity,
					Name:       p.Name,
					Picture:    p.Images[0].Url,
				})

				orderItems = append(orderItems, &consumerv1.OrderItem{
					Item: &cartv1.CartItem{
						MerchantId: cart.MerchantId,
						ProductId:  cart.ProductId,
						Quantity:   itemQuantity,
						Name:       p.Name,
						Picture:    p.Images[0].Url,
					},
					Cost: itemTotal,
				})

				amount += itemTotal
			}
		}
	}

	// 获取用户地址
	address, err := c.data.userv1.GetConsumerAddress(ctx, &userv1.GetConsumerAddressRequest{
		AddressId: req.AddressId,
		UserId:    req.UserId.String(),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取用户地址失败: %v", err))
	}

	log.Debugf("cartItems: %+v", cartItems)

	// Conditional logic for payment method specific preparations
	// var creditCardInfo *userv1.CreditCard
	paymentCurrency := string(constants.CNY) // Default currency, e.g., for balance or Alipay

	switch req.PaymentMethod {
	case string(constants.PaymentMethodBalance):
		log.Debugf("PaymentMethodGetUserBalanceBalance: %+v", req.UserId.String())
		banance, bananceErr := c.data.banancev1.GetUserBalance(ctx, &banancev1.GetUserBalanceRequest{
			UserId:   req.UserId.String(),
			Currency: paymentCurrency,
		})
		if bananceErr != nil {
			return nil, kerrors.NotFound("GET_BALANCE_FAILED", fmt.Sprintf("获取用户余额信息失败: %v", bananceErr))
		}
		paymentCurrency = banance.Currency
	case string(constants.PaymentMethodAlipay):

	default:
		return nil, kerrors.BadRequest("INVALID_PAYMENT_METHOD", fmt.Sprintf("Unsupported payment method: %s", req.PaymentMethod))
	}

	// 调用订单微服务创建订单
	order, orderErr := c.data.consumerOrderv1.PlaceOrder(ctx, &consumerv1.PlaceOrderRequest{
		Currency: paymentCurrency, // Use determined paymentCurrency
		Address: &userv1.ConsumerAddress{
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
	if orderErr != nil {
		return nil, fmt.Errorf("创建订单失败: %v", orderErr)
	}
	if order == nil || order.Order == nil {
		return nil, status.Error(codes.Internal, "订单服务返回无效数据")
	}

	// 调用支付微服务生成支付URL based on payment method
	log.Debugf("Processing payment for method: %s, UserID: %s, Amount: %.2f", req.PaymentMethod, req.UserId.String(), amount)

	var paymentReply *paymentv1.CreatePaymentResponse
	var createPaymentErr error

	switch req.PaymentMethod {
	case string(constants.PaymentMethodBalance):
		// 1. Check balance - IMPORTANT: Assumes c.data.userv1.GetUserBalance method exists and is correctly implemented.
		// If this method is not available or behaves differently, this section will need adjustment.
		log.Debugf("paymentCurrency: %s", paymentCurrency)
		balanceResp, balanceCheckErr := c.data.banancev1.GetUserBalance(ctx, &banancev1.GetUserBalanceRequest{
			UserId:   req.UserId.String(),
			Currency: paymentCurrency,
		})
		if balanceCheckErr != nil {
			return nil, kerrors.InternalServer("GET_BALANCE_FAILED", fmt.Sprintf("Failed to get user balance: %v", balanceCheckErr))
		}
		// Assuming balanceResp.Available is a float64 field representing the available balance.
		if balanceResp.Available < amount {
			return nil, kerrors.BadRequest("INSUFFICIENT_BALANCE", "Insufficient balance for this transaction")
		}

		paymentReply, createPaymentErr = c.data.paymentv1.CreatePayment(ctx, &paymentv1.CreatePaymentRequest{
			OrderId:          order.Order.OrderId,
			ConsumerId:       req.UserId.String(),
			Amount:           fmt.Sprintf("%.2f", amount),
			Currency:         paymentCurrency,
			Subject:          fmt.Sprintf("Balance Payment for Order %d", order.Order.OrderId),
			ReturnUrl:        "", // May not be applicable or could be a success confirmation page
			FreezeId:         order.Order.FreezeId,
			ConsumerVersion:  order.Order.ConsumerVersion,
			MerchantVersions: order.Order.MerchantVersion,
		})
	case string(constants.PaymentMethodAlipay):
		paymentReply, createPaymentErr = c.data.paymentv1.CreatePayment(ctx, &paymentv1.CreatePaymentRequest{
			OrderId:          order.Order.OrderId,
			ConsumerId:       req.UserId.String(),
			Amount:           fmt.Sprintf("%.2f", amount),
			Currency:         paymentCurrency, // Typically CNY for Alipay
			Subject:          "Alipay Payment",
			ReturnUrl:        "", // Alipay handles its own redirects
			FreezeId:         order.Order.FreezeId,
			ConsumerVersion:  order.Order.ConsumerVersion,
			MerchantVersions: order.Order.MerchantVersion,
		})
	default:
		return nil, kerrors.BadRequest("UNHANDLED_PAYMENT_METHOD", fmt.Sprintf("Unhandled payment method: %s", req.PaymentMethod))
	}

	if createPaymentErr != nil {
		// Consider order rollback logic here if payment creation fails
		return nil, kerrors.New(500, "CREATE_PAYMENT_ERR", fmt.Sprintf("创建支付失败 (%s): %v", req.PaymentMethod, createPaymentErr))
	}

	var finalPaymentId int64
	var finalPaymentURL string

	if paymentReply != nil {
		finalPaymentId = paymentReply.PaymentId
		finalPaymentURL = paymentReply.PayUrl
	} else {
		// This case should ideally not be reached if createPaymentErr is handled properly.
		// However, as a safeguard:
		return nil, kerrors.New(500, "PAYMENT_REPLY_NIL", fmt.Sprintf("Payment reply is nil after attempting payment (%s)", req.PaymentMethod))
	}

	if finalPaymentId == 0 {
		if !(req.PaymentMethod == "balance") { // Allow PaymentId to be non-zero for balance even if URL is empty
			// For non-balance payments, or if balance payment also results in zero PaymentId, it's an issue.
			return nil, kerrors.New(500, "PAYMENT_PROCESSING_FAILED", fmt.Sprintf("支付处理失败或返回无效支付ID (%s): %+v", req.PaymentMethod, paymentReply))
		}
	}

	return &biz.CheckoutReply{
		OrderId:    order.Order.OrderId,
		PaymentId:  finalPaymentId,
		PaymentURL: finalPaymentURL,
	}, nil
}

func NewCheckoutRepo(data *Data, logger log.Logger) biz.CheckoutRepo {
	return &checkoutRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
