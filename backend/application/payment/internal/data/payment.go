package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"backend/constants"

	"github.com/go-kratos/kratos/v2/metadata"

	orderv1 "backend/api/order/v1"

	"github.com/go-kratos/kratos/v2/transport/http"

	"backend/pkg/types"

	"backend/application/payment/internal/conf"

	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"

	"backend/application/payment/internal/data/models"

	"backend/application/payment/internal/biz"
	pkg2 "backend/application/payment/internal/pkg"

	"github.com/go-kratos/kratos/v2/log"
)

// PaymentRepo 支付仓储实现
type paymentRepo struct {
	data *Data
	log  *log.Helper
	conf *conf.Pay
}

// NewPaymentRepo 创建支付仓储
func NewPaymentRepo(data *Data, logger log.Logger, conf *conf.Pay) biz.PaymentRepo {
	return &paymentRepo{
		data: data,
		log:  log.NewHelper(logger),
		conf: conf,
	}
}

// CreatePayment 创建支付记录
func (r *paymentRepo) CreatePayment(ctx context.Context, req *biz.CreatePaymentReq) (*biz.CreatePaymentResp, error) {
	// 生成支付ID
	tradeNo := fmt.Sprintf("%d", xid.Next())

	// totalAmount := strconv.FormatFloat(req.Amount, 'f', 2, 64)

	pay := alipay.TradePagePay{ // 电脑网站支付
		Trade: alipay.Trade{
			NotifyURL: r.conf.Alipay.NotifyUrl, // 异步通知地址
			ReturnURL: r.conf.Alipay.ReturnUrl, // 回调地址
			Subject:   req.Subject,             // 订单主题
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
	paymentID := pkg2.SnowflakeID()

	// 创建支付记录
	payments, paymentsErr := r.data.DB(ctx).CreatePaymentQuery(ctx, models.CreatePaymentQueryParams{
		ID:       paymentID,
		OrderID:  req.OrderID,
		UserID:   req.UserID,
		Amount:   amount,
		Currency: req.Currency,
		Method:   "alipay",
		Status:   string(biz.PaymentStatusPending),
		Subject:  req.Subject,
		TradeNo:  tradeNo,
		Metadata: nil,
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
			ID:       payments.ID,
			OrderID:  payments.OrderID,
			UserID:   payments.UserID,
			Amount:   amountResp,
			Currency: payments.Currency,
			Subject:  req.Subject,
			Status:   biz.PaymentStatus(payments.Status),
			TradeNo:  tradeNo,
			// PayURL:    payments.PayURL,
			PayURL:    payUrl.String(),
			CreatedAt: payments.CreatedAt,
			UpdatedAt: payments.UpdatedAt,
		},
	}, err
}

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
			CreatedAt: payment.CreatedAt,
			UpdatedAt: payment.UpdatedAt,
		},
	}, nil
}

// HandlePaymentNotify 处理支付通知
func (r *paymentRepo) HandlePaymentNotify(ctx context.Context, req *biz.PaymentNotifyReq) (*biz.PaymentNotifyResp, error) {
	// 检查请求参数
	if req == nil {
		r.log.Error("支付通知请求为空")
		return &biz.PaymentNotifyResp{
			Success: false,
			Message: "支付通知请求为空",
		}, nil
	}

	// 将请求参数转换为url.Values格式
	values := make(url.Values)

	r.log.Debugf("308 HandlePaymentNotify req.Params:%+v", req.Params)
	r.log.Debugf("309 HandlePaymentNotify values:%+v", values)

	// AppID
	if req.AppID != "" {
		values.Set("app_id", req.AppID)
	}
	// AuthAppId
	if req.AuthAppId != "" {
		values.Set("auth_app_id", req.AuthAppId)
	}
	// Charset
	if req.Charset != "" {
		values.Set("charset", req.Charset)
	}
	// 交易号
	if req.TradeNo != "" {
		values.Set("trade_no", req.TradeNo)
	}
	// 支付方式, 沙箱
	if req.Method != "" {
		values.Set("method", req.Method)
	}
	// 签名
	if req.Sign != "" {
		values.Set("sign", req.Sign)
	}
	// 签名类型
	if req.SignType != "" {
		values.Set("sign_type", req.SignType)
	}
	// 商家订单号
	if req.OutTradeNo != "" {
		values.Set("out_trade_no", req.OutTradeNo)
	}
	// 订单金额
	if req.TotalAmount != "" {
		values.Set("total_amount", req.TotalAmount)
	}
	// 商家 ID
	if req.SellerId != "" {
		values.Set("seller_id", req.SellerId)
	}

	r.log.Debugf("343 TradeNo:%v OutTradeNo:%v TotalAmount:%v", req.TradeNo, req.OutTradeNo, req.TotalAmount)

	if r.data.alipay == nil {
		r.log.Error("支付宝客户端未初始化")
		return &biz.PaymentNotifyResp{
			Success: false,
			Message: "支付宝客户端未初始化",
		}, nil
	}

	notification, err := r.data.alipay.DecodeNotification(values)
	if err != nil {
		r.log.Errorf("支付宝异步通知解析失败: %v", err)
		// 解析失败不直接返回，继续处理
		// 可能是签名验证失败，但我们仍然需要处理支付状态更新
		r.log.Warnf("签名验证失败，降级处理支付状态更新")

		// 使用请求参数中的订单号继续处理
		if req.OutTradeNo == "" {
			r.log.Error("订单号为空，无法继续处理")
			return &biz.PaymentNotifyResp{
				Success: false,
				Message: "订单号为空，无法继续处理",
			}, nil
		}
	}

	r.log.Debugf("支付宝通知: %+v", notification)

	if notification.OutTradeNo == "" {
		r.log.Error("订单号为空")
		return &biz.PaymentNotifyResp{
			Success: false,
			Message: "订单号为空",
		}, nil
	}

	r.log.Debugf("HandlePaymentNotify notification.OutTradeNo: %+v", notification.OutTradeNo)

	// 查询支付记录
	payment, getPayErr := r.data.db.GetPaymentByTradeNo(ctx, notification.OutTradeNo)
	if getPayErr != nil {
		r.log.Errorf("查询支付记录失败: %v", getPayErr)
		return &biz.PaymentNotifyResp{
			Success: false,
			Message: fmt.Sprintf("查询支付记录失败: %v", getPayErr),
		}, nil
	}

	r.log.Debugf("data notification.TradeStatus: %+v", notification.TradeStatus)

	// 更新支付记录
	var status biz.PaymentStatus
	switch notification.TradeStatus {
	case alipay.TradeStatus(biz.AliPayStatusSuccess):
		status = biz.PaymentStatusSuccess
	case alipay.TradeStatus(biz.AliPayStatusClosed):
		status = biz.PaymentStatusClosed
	case alipay.TradeStatus(biz.AliPayStatusPending):
		status = biz.PaymentStatusPending
	default:
		status = biz.PaymentStatusProcessing
	}

	// 更新支付状态
	updatePaymentStatusResult, err := r.data.db.UpdatePaymentStatus(ctx, models.UpdatePaymentStatusParams{
		ID:      payment.ID,
		OrderID: payment.OrderID,
		Status:  string(status),
	})
	if err != nil {
		r.log.Errorf("更新支付状态失败: %v", err)
		return nil, fmt.Errorf("更新支付状态失败: %w", err)
	}

	log.Debugf("updatePaymentStatusResult: %+v", updatePaymentStatusResult)

	// 传递用户 ID
	r.log.Debugf("userId:%+v", payment.UserID.String())
	metadataCtx := metadata.AppendToClientContext(ctx, constants.UserId, payment.UserID.String())
	result, err := r.data.orderv1.MarkOrderPaid(metadataCtx, &orderv1.MarkOrderPaidReq{
		// UserId:  payment.UserID.String(),
		OrderId: payment.OrderID,
	})
	if err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	r.log.Debugf("result: %+v", result)

	// 如果通知消息没有问题，我们需要确认收到通知消息，不然支付宝后续会继续推送相同的消息
	// 安全地获取HTTP上下文
	httpCtx, ok := ctx.(http.Context)
	if ok {
		r.data.alipay.ACKNotification(httpCtx.Response())
	} else {
		r.log.Warn("无法获取HTTP上下文，无法确认收到通知消息")
	}

	return &biz.PaymentNotifyResp{
		Success: true,
		Message: "success",
	}, nil
}

// HandlePaymentCallback 处理支付回调
func (r *paymentRepo) HandlePaymentCallback(ctx context.Context, req *biz.PaymentCallbackReq) (*biz.PaymentCallbackResp, error) {
	// 检查请求参数
	if req == nil {
		r.log.Error("支付回调请求为空")
		return &biz.PaymentCallbackResp{
			Success: false,
			Message: "支付回调请求为空",
		}, nil
	}

	// 将请求参数转换为url.Values格式
	values := make(url.Values)

	// 首先从请求参数中获取数据
	// 使用请求中的参数构建values
	if req.Params != nil {
		for k, v := range req.Params {
			values.Set(k, v)
		}
	}

	// 交易号
	if req.TradeNo != "" {
		values.Set("trade_no", req.TradeNo)
	}
	// 商户订单号
	if req.OutTradeNo != "" {
		values.Set("out_trade_no", req.OutTradeNo)
	}
	// 订单金额
	if req.TotalAmount != "" {
		values.Set("total_amount", req.TotalAmount)
	}
	// 订单标题
	if req.Subject != "" {
		values.Set("subject", req.Subject)
	}
	// 支付状态
	if req.TradeStatus != "" {
		values.Set("trade_status", req.TradeStatus)
	}
	r.log.Debugf("477 TradeNo:%v OutTradeNo:%v TotalAmount:%v Subject:%v TradeStatus:%v", req.TradeNo, req.OutTradeNo, req.TotalAmount, req.Subject, req.TradeStatus)
	r.log.Debugf("478 支付回调参数: %v", values)

	if r.data.alipay == nil {
		r.log.Error("支付宝客户端未初始化")
		return &biz.PaymentCallbackResp{
			Success: false,
			Message: "支付宝客户端未初始化",
		}, nil
	}

	// 检查订单号
	outTradeNo := req.OutTradeNo
	if outTradeNo == "" {
		r.log.Error("订单号为空")
		return &biz.PaymentCallbackResp{
			Success: false,
			Message: "订单号为空",
		}, nil
	}

	// 解析通知
	notification, err := r.data.alipay.DecodeNotification(values)
	if err != nil {
		r.log.Errorf("支付宝异步通知解析失败: %v", err)
		// 解析失败不直接返回，继续处理
		// 可能是签名验证失败，但我们仍然需要处理支付状态更新
		r.log.Info("尽管签名验证失败，仍将继续处理支付状态更新")
	} else {
		r.log.Debugf("支付宝通知解析成功: %+v", notification)
	}

	// 查询支付记录
	payment, err := r.data.db.GetPaymentByTradeNo(ctx, outTradeNo)
	if err != nil {
		r.log.Errorf("查询支付记录失败: %v", err)
		// 检查是否是记录不存在的错误
		if errors.Is(err, sql.ErrNoRows) {
			return &biz.PaymentCallbackResp{
				Success: false,
				Message: "支付记录不存在",
			}, nil
		}
		// 其他数据库错误
		return &biz.PaymentCallbackResp{
			Success: false,
			Message: fmt.Sprintf("查询支付记录失败: %v", err),
		}, nil
	}

	r.log.Debugf("支付记录: %+v", payment)

	// 更新支付状态为成功
	updateParams := models.UpdatePaymentStatusParams{
		ID:      payment.ID,
		OrderID: payment.OrderID,
		Status:  string(biz.PaymentStatusSuccess),
	}

	// 添加防御性检查
	if payment.ID == 0 || payment.OrderID == 0 {
		r.log.Error("支付记录ID或订单ID为空，无法更新状态")
		return &biz.PaymentCallbackResp{
			Success: false,
			Message: "支付记录数据不完整，无法更新状态",
		}, nil
	}

	// 更新支付表的支付状态
	_, updateErr := r.data.DB(ctx).UpdatePaymentStatus(ctx, updateParams)

	if updateErr != nil {
		r.log.Errorf("更新支付状态失败: %v", updateErr)
		return &biz.PaymentCallbackResp{
			Success: false,
			Message: fmt.Sprintf("更新支付状态失败: %v", updateErr),
		}, nil
	}

	// 更新订单表的支付状态
	r.log.Debugf("userId:%+v", payment.UserID.String())
	metadataCtx := metadata.AppendToClientContext(ctx, constants.UserId, payment.UserID.String())
	result, err := r.data.orderv1.MarkOrderPaid(metadataCtx, &orderv1.MarkOrderPaidReq{
		// UserId:  payment.UserID.String(),
		OrderId: payment.OrderID,
	})
	if err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	r.log.Debugf("result: %+v", result)

	// 返回支付结果
	return &biz.PaymentCallbackResp{
		Success: true,
		Message: "支付成功",
	}, nil
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
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}, nil
}
