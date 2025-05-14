package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	balancev1 "backend/api/balancer/v1"
	consumerOrderv1 "backend/api/consumer/order/v1"

	"backend/application/payment/internal/pkg/id"

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

// processBalanceTransfer 处理余额支付的资金转移
func (r *paymentRepo) processBalanceTransfer(ctx context.Context, payment models.PaymentsPayments) {
	// 查询订单的所有子订单信息
	orderInfo, err := r.data.consumerOrderv1.GetConsumerOrdersWithSuborders(ctx, &consumerOrderv1.GetConsumerOrdersWithSubordersRequest{
		UserId:  payment.ConsumerID.String(),
		OrderId: payment.OrderID,
	})
	if err != nil {
		r.log.Errorf("获取订单信息失败: %v", err)
		return
	}

	// 记录已处理的商家ID，避免重复处理
	processedMerchants := make(map[string]bool)

	for _, subOrder := range orderInfo.Orders {
		merchantId := subOrder.MerchantId

		// 检查该商家是否已经处理过
		if processedMerchants[merchantId] {
			r.log.Infof("商家ID: %s 已经处理过，跳过重复处理", merchantId)
			continue
		}

		// 根据商家ID和子订单金额，调用余额服务进行转账
		for _, v := range payment.MerchantVersions {
			// 调用余额服务确认转账
			params := &balancev1.ConfirmTransferRequest{
				FreezeId:                payment.FreezeID,
				MerchantId:              merchantId,
				IdempotencyKey:          strconv.FormatInt(payment.OrderID, 10),
				ExpectedUserVersion:     int32(payment.ConsumerVersion),
				ExpectedMerchantVersion: int32(v),
				PaymentAccount:          "", // 余额支付不需要支付账号
			}
			r.log.Debugf("确认余额转账参数: %+v", params)

			_, err = r.data.balancev1.ConfirmTransfer(ctx, params)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					r.log.Debugf("ConfirmTransfer error: %+v", pgErr)
				}
				r.log.Infof("确认转账失败，可能已经处理过，商家ID: %s, 错误: %v",
					merchantId, err)
				// 继续处理其他商家
				continue
			}

			r.log.Infof("确认余额转账成功，商家ID: %s", merchantId)
			// 标记该商家已处理
			processedMerchants[merchantId] = true
			break
		}
	}
}

// CreatePayment 创建支付记录
func (r *paymentRepo) CreatePayment(ctx context.Context, req *biz.CreatePaymentReq) (*biz.CreatePaymentResp, error) {
	// 生成支付ID
	tradeNo := fmt.Sprintf("%d", xid.Next())

	// 转换金额为数值
	numeric, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, errors.New("转换字符串为数值失败")
	}
	amount, err := types.Float64ToNumeric(numeric)
	if err != nil {
		return nil, errors.New("转换数值为数据库类型失败")
	}

	// 确定支付方式和状态
	paymentMethod := string(constants.PaymentMethodAlipay)
	paymentStatus := string(biz.PaymentStatusPending)
	var payUrl *url.URL

	// 根据支付方式处理
	if req.Subject == fmt.Sprintf("Balance Payment for Order %d", req.OrderID) {
		// 余额支付
		paymentMethod = string(constants.PaymentMethodBalance)
		paymentStatus = string(biz.PaymentStatusSuccess) // 余额支付直接标记为成功

		// 余额支付不需要支付链接
		payUrl = &url.URL{}
	} else {
		// 支付宝支付
		pay := alipay.TradePagePay{ // 电脑网站支付
			Trade: alipay.Trade{
				NotifyURL:   r.conf.Alipay.NotifyUrl,  // 异步通知地址
				ReturnURL:   r.conf.Alipay.ReturnUrl,  // 回调地址
				Subject:     req.Subject,              // 订单主题
				OutTradeNo:  tradeNo,                  // 商户订单号
				TotalAmount: req.Amount,               // 订单金额
				ProductCode: "FAST_INSTANT_TRADE_PAY", // 电脑网站支付，产品码为固定值
			},
		}
		payUrl, err = r.data.alipay.TradePagePay(pay) // 生成支付链接
		if err != nil {
			return nil, fmt.Errorf("failed to create payment: %v", err)
		}
	}

	// 创建支付记录
	payments, paymentsErr := r.data.DB(ctx).CreatePaymentQuery(ctx, models.CreatePaymentQueryParams{
		ID:               id.SnowflakeID(),
		OrderID:          req.OrderID,
		ConsumerID:       req.ConsumerID,
		Amount:           amount,
		Currency:         req.Currency,
		Method:           paymentMethod,
		Status:           paymentStatus,
		Subject:          req.Subject,
		TradeNo:          tradeNo,
		FreezeID:         req.FreezeId,
		ConsumerVersion:  req.ConsumerVersion,
		MerchantVersions: req.MerchanVersions,
	})
	if paymentsErr != nil {
		return nil, fmt.Errorf("failed to create payment: %v", paymentsErr)
	}
	r.log.Infof("Created payment: %+v", payments)

	amountResp, err := types.NumericToFloat(payments.Amount)
	if err != nil {
		return nil, fmt.Errorf("转换金额失败: %w", err)
	}

	// 如果是余额支付，直接调用订单服务标记订单为已支付
	if paymentMethod == string(constants.PaymentMethodBalance) {
		// 传递用户 ID
		metadataCtx := metadata.AppendToClientContext(ctx, constants.UserId, req.ConsumerID.String())
		_, markOrderErr := r.data.orderv1.MarkOrderPaid(metadataCtx, &orderv1.MarkOrderPaidReq{
			OrderId: req.OrderID,
		})
		if markOrderErr != nil {
			r.log.Errorf("余额支付标记订单失败: %v", markOrderErr)
			// 不返回错误，继续处理
		}

		// 处理余额转账
		r.processBalanceTransfer(ctx, payments)
	}

	return &biz.CreatePaymentResp{
		Payment: &biz.Payment{
			ID:         payments.ID,
			OrderID:    payments.OrderID,
			ConsumerID: payments.ConsumerID,
			Amount:     amountResp,
			Currency:   payments.Currency,
			Subject:    req.Subject,
			Status:     biz.PaymentStatus(payments.Status),
			TradeNo:    tradeNo,
			PayURL:     payUrl.String(),
			NotifyTime: time.Time{},
			CreatedAt:  payments.CreatedAt,
			UpdatedAt:  payments.UpdatedAt,
		},
	}, nil
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
			ID:         payment.ID,
			OrderID:    payment.OrderID,
			ConsumerID: payment.ConsumerID,
			// MerchantID: payment.MerchantID,
			Amount:     amount,
			Currency:   payment.Currency,
			Subject:    payment.Subject,
			Status:     biz.PaymentStatus(payment.Status),
			TradeNo:    payment.TradeNo,
			PayURL:     "",
			NotifyTime: time.Time{},
			CreatedAt:  payment.CreatedAt,
			UpdatedAt:  payment.UpdatedAt,
		},
	}, nil
}

// HandlePaymentNotify 处理支付通知
func (r *paymentRepo) HandlePaymentNotify(ctx context.Context, req url.Values) (*biz.PaymentNotifyResp, error) {
	// 检查请求参数
	if req == nil {
		r.log.Error("支付通知请求为空")
		return &biz.PaymentNotifyResp{
			Success: false,
			Message: "支付通知请求为空",
		}, nil
	}
	notification, err := r.data.alipay.DecodeNotification(req)
	if err != nil {
		r.log.Errorf("支付宝异步通知解析失败: %v", err)
		// 解析失败不直接返回，继续处理
		// 可能是签名验证失败，但我们仍然需要处理支付状态更新
		r.log.Warnf("签名验证失败，降级处理支付状态更新")
	}

	if r.data.alipay == nil {
		r.log.Error("支付宝客户端未初始化")
		return &biz.PaymentNotifyResp{
			Success: false,
			Message: "支付宝客户端未初始化",
		}, nil
	}

	r.log.Debugf("支付宝通知: %+v", notification)

	if notification.OutTradeNo == "" {
		r.log.Error("订单号为空")
		return &biz.PaymentNotifyResp{
			Success: false,
			Message: "订单号为空",
		}, nil
	}
	// 查询支付记录
	payment, getPayErr := r.data.db.GetPaymentByTradeNo(ctx, notification.OutTradeNo)
	if getPayErr != nil {
		r.log.Errorf("查询支付记录失败: %v", getPayErr)
		return &biz.PaymentNotifyResp{
			Success: false,
			Message: fmt.Sprintf("查询支付记录失败: %v", getPayErr),
		}, nil
	}
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
	_, updatePaymentStatusErr := r.data.db.UpdatePaymentStatus(ctx, models.UpdatePaymentStatusParams{
		ID:      payment.ID,
		OrderID: payment.OrderID,
		Status:  string(status),
	})
	if updatePaymentStatusErr != nil {
		r.log.Errorf("更新支付状态失败: %v", updatePaymentStatusErr)
		return nil, fmt.Errorf("更新支付状态失败: %w", updatePaymentStatusErr)
	}

	// 传递用户 ID
	metadataCtx := metadata.AppendToClientContext(ctx, constants.UserId, payment.ConsumerID.String())
	_, err = r.data.orderv1.MarkOrderPaid(metadataCtx, &orderv1.MarkOrderPaidReq{
		// UserId:  payment.ConsumerID.String(),
		OrderId: payment.OrderID,
	})
	if err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

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

	// 标记订单为已支付
	metadataCtx := metadata.AppendToClientContext(ctx, constants.UserId, payment.ConsumerID.String())
	_, err = r.data.orderv1.MarkOrderPaid(metadataCtx, &orderv1.MarkOrderPaidReq{
		// UserId:  payment.ConsumerID.String(),
		OrderId: payment.OrderID,
	})
	if err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	// 查询订单的所有子订单信息
	orderInfo, err := r.data.consumerOrderv1.GetConsumerOrdersWithSuborders(ctx, &consumerOrderv1.GetConsumerOrdersWithSubordersRequest{
		UserId:  payment.ConsumerID.String(),
		OrderId: payment.OrderID,
	})
	if err != nil {
		r.log.Errorf("获取订单信息失败: %v", err)
		return nil, err
	}

	// 记录已处理的商家ID，避免重复处理
	processedMerchants := make(map[string]bool)

	for _, subOrder := range orderInfo.Orders {
		merchantId := subOrder.MerchantId

		// 检查该商家是否已经处理过
		if processedMerchants[merchantId] {
			r.log.Infof("商家ID: %s 已经处理过，跳过重复处理", merchantId)
			continue
		}

		// 计算子订单金额
		subOrderAmount := subOrder.TotalAmount

		// 根据商家ID和子订单金额，调用余额服务进行转账
		for _, v := range payment.MerchantVersions {
			// 调用余额服务确认转账
			params := &balancev1.ConfirmTransferRequest{
				FreezeId:                payment.FreezeID,
				MerchantId:              merchantId,
				IdempotencyKey:          strconv.FormatInt(payment.OrderID, 10),
				ExpectedUserVersion:     int32(payment.ConsumerVersion),
				ExpectedMerchantVersion: int32(v),
				PaymentAccount:          "", // TODO PaymentAccount
			}
			log.Debugf("params: %+v", params)
			_, err = r.data.balancev1.ConfirmTransfer(ctx, params)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					log.Debugf("ConfirmTransfer error is pgErr: %+v %+v", err, pgErr)
					log.Debugf("ConfirmTransfer Code: %v", pgErr.Code)
					log.Debugf("ConfirmTransfer Message: %v", pgErr.Message)
					log.Debugf("ConfirmTransfer pgErr: %v", pgErr)
				}
				// 如果错误是因为状态已经是CONFIRMED，则认为是成功的
				log.Debugf("ConfirmTransfer error:%+v", err)
				// if errors.Is(err) && (err.Error() == "rpc error: code = InvalidArgument desc = freeze is not in FROZEN status: CONFIRMED") {
				r.log.Infof("冻结状态已经是CONFIRMED，视为转账成功，商家ID: %s, 金额: %f",
					merchantId, subOrderAmount)
				// 标记该商家已处理
				processedMerchants[merchantId] = true
				continue
				// }

				// r.log.Errorf("确认转账失败，商家ID: %s, 金额: %f, 错误: %v",
				// 	merchantId, subOrderAmount, err)
				// // 继续处理其他子订单，不要因为一个子订单失败而中断整个流程
				// continue
			} else {
				r.log.Infof("确认转账成功，商家ID: %s, 金额: %f",
					merchantId, subOrderAmount)
				// 标记该商家已处理
				processedMerchants[merchantId] = true
			}
		}
	}

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
		ID:         payment.ID,
		OrderID:    payment.OrderID,
		ConsumerID: payment.ConsumerID,
		Amount:     amount,
		Currency:   payment.Currency,
		Subject:    payment.Subject,
		Status:     biz.PaymentStatus(payment.Status),
		TradeNo:    payment.TradeNo,
		PayURL:     "",
		NotifyTime: time.Time{},
		CreatedAt:  payment.CreatedAt,
		UpdatedAt:  payment.UpdatedAt,
	}, nil
}
