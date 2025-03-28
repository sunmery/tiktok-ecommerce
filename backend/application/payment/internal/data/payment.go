package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	orderv1 "backend/api/order/v1"

	"backend/pkg/types"

	"backend/application/payment/internal/conf"

	"backend/application/payment/pkg"

	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"

	"backend/application/payment/internal/data/models"

	"backend/application/payment/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// PaymentRepo 支付仓储实现
type paymentRepo struct {
	data         *Data
	log          *log.Helper
	pay          *conf.Pay
	alipayClient *alipay.Client
}

// NewPaymentRepo 创建支付仓储
func NewPaymentRepo(data *Data, logger log.Logger, pay *conf.Pay) biz.PaymentRepo {
	return &paymentRepo{
		data: data,
		log:  log.NewHelper(logger),
		pay:  pay,
	}
}

func (r *paymentRepo) PaymentNotify(ctx context.Context, req *biz.PaymentNotifyReq) (*biz.PaymentNotifyResp, error) {
	// 将请求参数转换为url.Values格式
	values := make(map[string][]string)
	for k, v := range req.Params {
		values[k] = []string{v}
	}
	fmt.Printf("values: %v\n", values)
	return &biz.PaymentNotifyResp{
		Success: true,
		Message: "OK",
	}, nil
}

// func (r *paymentRepo) GetPayment(ctx context.Context, id uuid.UUID) (*biz.PaymentResp, error) {
// 	panic("TODO")
// }

// ProcessPaymentCallback TODO
// 验证支付结果，更新交易状态。
// 调用订单服务的MarkOrderPaid标记订单为已支付。
// 若支付超时（定时取消），触发订单服务的订单取消逻辑。
func (r *paymentRepo) ProcessPaymentCallback(ctx context.Context, req *biz.AliPayCallbackReq) (*biz.PaymentCallbackResp, error) {
	r.log.Infof("Received payment callback: %+v", req)

	tradeNo := fmt.Sprintf("%v", req.TradeNo)

	requestForm := url.Values{
		"app_id":       []string{req.AppId},
		"auth_app_id":  []string{req.AuthAppId},
		"charset":      []string{req.Charset},
		"method":       []string{req.Method},
		"out_trade_no": []string{req.OutTradeNo},
		"seller_id":    []string{req.SellerId},
		"sign":         []string{req.Sign},
		"sign_type":    []string{req.SignType},
		"total_amount": []string{req.TotalAmount},
		"trade_no":     []string{tradeNo},
		"timestamp":    []string{req.Timestamp},
	}

	// 验证签名
	if err := r.data.alipay.VerifySign(requestForm); err != nil {
		r.log.Errorf("Failed to verify sign: %v", err)
		return nil, fmt.Errorf("signature verification failed: %v", err)
	}

	// TODO: 处理支付成功逻辑
	// 1. 更新支付记录状态
	// 2. 调用订单服务更新订单状态
	// 3. 发送支付成功事件
	_, err := r.data.orderv1.MarkOrderPaid(ctx, &orderv1.MarkOrderPaidReq{
		OrderId: req.TradeNo,
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("订单状态已更新为已支付")

	return &biz.PaymentCallbackResp{
		Success: true,
		Message: "OK",
	}, nil
}

// CreatePayment 创建支付记录
func (r *paymentRepo) CreatePayment(ctx context.Context, req *biz.CreatePaymentReq) (*biz.CreatePaymentResp, error) {
	// 生成支付ID
	tradeNo := fmt.Sprintf("%d", xid.Next())

	// totalAmount := strconv.FormatFloat(req.Amount, 'f', 2, 64)

	pay := alipay.TradePagePay{ // 电脑网站支付
		Trade: alipay.Trade{
			NotifyURL: r.pay.Alipay.NotifyUrl, // 异步通知地址
			ReturnURL: r.pay.Alipay.ReturnUrl, // 回调地址
			Subject:   req.Subject,            // 订单主题
			// OutTradeNo:  fmt.Sprintf("%d", time.Now().Unix()), // 商户订单号，必须唯一
			// OutTradeNo: req.OrderId.String(), // 商户订单号，由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
			OutTradeNo:  tradeNo,                  // 商户订单号，由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
			TotalAmount: req.Amount,               // 直接使用字符串金额               // 订单金额
			ProductCode: "FAST_INSTANT_TRADE_PAY", // 电脑网站支付，产品码为固定值
		},
	}
	payUrl, err := r.data.alipay.TradePagePay(pay) // 生成支付链接
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}
	numeric, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, errors.New("转换字符串为数值失败")
	}
	amount, err := types.Float64ToNumeric(numeric)
	if err != nil {
		return nil, errors.New("转换数值为数据库类型失败")
	}

	// 使用雪花ID
	paymentID := pkg.SnowflakeID()

	// 创建支付记录
	payments, paymentsErr := r.data.DB(ctx).CreatePaymentQuery(ctx, models.CreatePaymentQueryParams{
		ID:       paymentID,
		OrderID:  req.OrderID,
		UserID:   req.UserID,
		Amount:   amount,
		Currency: req.Currency,
		Method:   "alipay",
		Status:   string(biz.PaymentStatusPending),
		// PayUrl:    url.String(),
		Subject:     req.Subject,
		TradeNo:     tradeNo,
		GatewayTxID: tradeNo,
		PayUrl:      payUrl.String(),
		Metadata:    nil,
	})
	if paymentsErr != nil {
		return nil, fmt.Errorf("failed to create payment: %v", paymentsErr)
	}
	fmt.Printf("payments: %+v\n", payments)

	amountResp, err := types.NumericToFloat(payments.Amount)
	if err != nil {
		return nil, fmt.Errorf("转换金额失败: %w", err)
	}

	return &biz.CreatePaymentResp{
		Payment: &biz.Payment{
			ID:        payments.ID,
			OrderID:   payments.OrderID,
			UserID:    payments.UserID,
			Amount:    amountResp,
			Currency:  payments.Currency,
			Subject:   req.Subject,
			Status:    biz.PaymentStatus(payments.Status),
			TradeNo:   tradeNo,
			PayURL:    payments.PayUrl,
			CreatedAt: payments.CreatedAt,
			UpdatedAt: payments.UpdatedAt,
		},
	}, err
}

// func RedirectFilter(next ghttp.Handler) ghttp.Handler {
// 	return ghttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path == "/v1/payments" {
// 			ghttp.Redirect(w, r, RedirectURL, ghttp.StatusMovedPermanently)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// GetPaymentStatus 查询支付状态
func (r *paymentRepo) GetPaymentStatus(ctx context.Context, req *biz.GetPaymentStatusReq) (*biz.GetPaymentStatusResp, error) {
	// 查询支付记录
	payment, err := r.data.db.GetPayment(ctx, req.PaymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("支付记录不存在")
		}
		return nil, fmt.Errorf("查询支付记录失败: %w", err)
	}

	// 返回支付状态
	amount, err := types.NumericToFloat(payment.Amount)
	if err != nil {
		return nil, fmt.Errorf("转换金额失败: %w", err)
	}

	return &biz.GetPaymentStatusResp{
		Payment: &biz.Payment{
			ID:        payment.ID,
			OrderID:   payment.OrderID,
			UserID:    payment.UserID,
			Amount:    amount,
			Currency:  payment.Currency,
			Subject:   payment.Subject,
			Status:    biz.PaymentStatus(payment.Status),
			TradeNo:   payment.TradeNo,
			PayURL:    payment.PayUrl,
			CreatedAt: payment.CreatedAt,
			UpdatedAt: payment.UpdatedAt,
		},
	}, nil
}

// HandlePaymentNotify 处理支付通知
func (r *paymentRepo) HandlePaymentNotify(ctx context.Context, req *biz.PaymentNotifyReq) (*biz.PaymentNotifyResp, error) {
	// 将请求参数转换为url.Values格式
	values := make(map[string][]string)
	for k, v := range req.Params {
		values[k] = []string{v}
	}

	// 验证签名
	if r.alipayClient == nil {
		return nil, fmt.Errorf("支付宝客户端未初始化")
	}

	notification, err := r.alipayClient.DecodeNotification(values)
	if err != nil {
		r.log.Errorf("支付宝异步通知解析失败: %v", err)
		return nil, fmt.Errorf("支付宝异步通知解析失败: %w", err)
	}

	// 获取订单号
	outTradeNo := req.OutTradeNo

	// 查询支付记录
	payment, err := r.data.db.GetPaymentByTradeNo(ctx, outTradeNo)
	if err != nil {
		return nil, fmt.Errorf("查询支付记录失败: %w", err)
	}

	// 更新支付记录
	var status biz.PaymentStatus
	switch req.TradeStatus {
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		status = biz.PaymentStatusSuccess
	case "TRADE_CLOSED":
		status = biz.PaymentStatusClosed
	case "WAIT_BUYER_PAY":
		status = biz.PaymentStatusPending
	default:
		status = biz.PaymentStatusProcessing
	}

	// 更新支付状态
	updatePaymentStatusResult, err := r.data.db.UpdatePaymentStatus(ctx, models.UpdatePaymentStatusParams{
		ID:          payment.ID,
		OrderID:     payment.OrderID,
		GatewayTxID: notification.TradeNo,
		Status:      string(status),
	})
	if err != nil {
		r.log.Errorf("更新支付状态失败: %v", err)
		return nil, fmt.Errorf("更新支付状态失败: %w", err)
	}

	log.Debugf("updatePaymentStatusResult: %+v", updatePaymentStatusResult)

	return &biz.PaymentNotifyResp{
		Success: true,
		Message: "success",
	}, nil
}

// HandlePaymentCallback 处理支付回调
func (r *paymentRepo) HandlePaymentCallback(ctx context.Context, req *biz.PaymentCallbackReq) (*biz.PaymentCallbackResp, error) {
	// 将请求参数转换为url.Values格式
	values := make(map[string][]string)
	for k, v := range req.Params {
		values[k] = []string{v}
	}

	// 验证签名
	if r.alipayClient == nil {
		return nil, fmt.Errorf("支付宝客户端未初始化")
	}

	notification, err := r.alipayClient.DecodeNotification(values)
	if err != nil {
		r.log.Errorf("支付宝异步通知解析失败: %v", err)
		return nil, fmt.Errorf("支付宝异步通知解析失败: %w", err)
	}

	log.Debugf("notification: %+v", notification)

	// 获取订单号
	outTradeNo := req.OutTradeNo

	// 查询支付记录
	payment, err := r.data.db.GetPaymentByTradeNo(ctx, outTradeNo)
	if err != nil {
		return nil, fmt.Errorf("查询支付记录失败: %w", err)
	}

	log.Debugf("payment: %+v", payment)

	// 返回支付结果
	return &biz.PaymentCallbackResp{
		Success: true,
		Message: "支付成功",
	}, nil
}

// UpdatePaymentStatus 更新支付状态
func (r *paymentRepo) UpdatePaymentStatus(ctx context.Context, req *biz.UpdatePaymentStatusRequest) (*biz.UpdatePaymentStatusResponse, error) {
	// 更新支付状态
	updatePaymentStatusResult, err := r.data.db.UpdatePaymentStatus(ctx, models.UpdatePaymentStatusParams{
		ID:          req.PaymentId,
		OrderID:     req.OrderId,
		GatewayTxID: req.TradeNo,
		Status:      string(req.Status),
	})
	if err != nil {
		r.log.Errorf("更新支付状态失败: %v", err)
		return nil, fmt.Errorf("更新支付状态失败: %w", err)
	}

	log.Debugf("updatePaymentStatusResult: %+v", updatePaymentStatusResult)

	return &biz.UpdatePaymentStatusResponse{}, nil
}

// GetPaymentByOrderID 根据订单ID查询支付记录
func (r *paymentRepo) GetPaymentByOrderID(ctx context.Context, req *biz.GetPaymentByOrderIDRequest) (*biz.Payment, error) {
	payment, err := r.data.db.GetPaymentByOrderID(ctx, req.OrderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("支付记录不存在")
		}
		return nil, fmt.Errorf("查询支付记录失败: %w", err)
	}

	amount, err := types.NumericToFloat(payment.Amount)
	if err != nil {
		return nil, fmt.Errorf("转换金额失败: %w", err)
	}

	return &biz.Payment{
		ID:        payment.ID,
		OrderID:   payment.OrderID,
		UserID:    payment.UserID,
		Amount:    amount,
		Currency:  payment.Currency,
		Subject:   payment.Subject,
		Status:    biz.PaymentStatus(payment.Status),
		TradeNo:   payment.TradeNo,
		PayURL:    payment.PayUrl,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}, nil
}
