package data

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"backend/pkg/types"

	"backend/application/payment/internal/conf"

	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"

	"github.com/google/uuid"

	"backend/application/payment/internal/data/models"

	"backend/application/payment/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

func NewPaymentRepo(data *Data, logger log.Logger, pay *conf.Pay) biz.PaymentRepo {
	return &paymentRepo{
		data: data,
		log:  log.NewHelper(logger),
		pay:  pay,
	}
}

type paymentRepo struct {
	data *Data
	log  *log.Helper
	pay  *conf.Pay
}

func (r *paymentRepo) PaymentNotify(ctx context.Context, req *biz.PaymentNotifyReq) (*biz.PaymentNotifyResp, error) {
	requestForm := map[string][]string{}
	for k, v := range req.Values {
		requestForm[k] = v
	}
	fmt.Printf("requestForm: %v\n", requestForm)
	return &biz.PaymentNotifyResp{
		Code: 200,
		Msg:  "success",
	}, nil
}

func (r *paymentRepo) GetPayment(ctx context.Context, id uuid.UUID) (*biz.PaymentResp, error) {
	panic("TODO")
}

// ProcessPaymentCallback TODO
// 验证支付结果，更新交易状态。
// 调用订单服务的MarkOrderPaid标记订单为已支付。
// 若支付超时（定时取消），触发订单服务的订单取消逻辑。
func (r *paymentRepo) ProcessPaymentCallback(ctx context.Context, req *biz.PaymentCallbackReq) (*biz.PaymentCallbackResp, error) {
	err := r.data.alipay.VerifySign(req.RequestForm)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to verify sign: %v", err))
	}
	requestForm := map[string][]string{}
	for k, v := range req.RequestForm {
		requestForm[k] = v
	}
	fmt.Printf("requestForm: %v\n", requestForm)
	return &biz.PaymentCallbackResp{}, nil
}

func (r *paymentRepo) CreatePayment(ctx context.Context, req *biz.CreatePaymentReq) (*biz.CreatePaymentRes, error) {
	// types.Float64ToNumeric(req.Amount)

	// fmt.Printf("req: %v\n", r.pay.Alipay)
	fmt.Printf("req: %+v\n", req)
	tradeNo := fmt.Sprintf("%d", xid.Next())
	pay := alipay.TradePagePay{ // 电脑网站支付
		Trade: alipay.Trade{
			NotifyURL: r.pay.Alipay.NotifyUrl, // 异步通知地址
			ReturnURL: r.pay.Alipay.ReturnUrl, // 回调地址
			Subject:   "支付测试:" + req.Subject,  // 订单主题
			// OutTradeNo:  fmt.Sprintf("%d", time.Now().Unix()), // 商户订单号，必须唯一
			// OutTradeNo: req.OrderId.String(), // 商户订单号，由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
			OutTradeNo:  tradeNo,                  // 商户订单号，由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
			TotalAmount: req.Amount,               // 订单金额
			ProductCode: "FAST_INSTANT_TRADE_PAY", // 电脑网站支付，产品码为固定值
		},
	}
	url, err := r.data.alipay.TradePagePay(pay) // 生成支付链接
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

	payments, paymentsErr := r.data.DB(ctx).CreatePaymentQuery(ctx, models.CreatePaymentQueryParams{
		PaymentID: uuid.New(),
		OrderID:   req.OrderId,
		Amount:    amount,
		Currency:  req.Currency,
		Method:    req.Method,
		Status:    string(biz.PaymentPending),
		Metadata:  nil,
	})
	if paymentsErr != nil {
		return nil, fmt.Errorf("failed to create payment: %v", paymentsErr)
	}
	fmt.Printf("payments: %+v\n", payments)

	return &biz.CreatePaymentRes{
		PaymentId:  payments.PaymentID.String(),
		Status:     payments.Status,
		PaymentUrl: url.String(),
		CreatedAt:  time.Now(),
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
