package data

import (
	"context"
	"fmt"
	"time"

	"backend/application/payment/internal/conf"

	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"

	"github.com/google/uuid"

	"backend/application/payment/internal/data/models"

	"backend/application/payment/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

func NewPaymentRepo(data *Data, logger log.Logger) biz.PaymentRepo {
	return &paymentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

type paymentRepo struct {
	data *Data
	log  *log.Helper
	conf *conf.Pay
}

func (r *paymentRepo) GetPayment(ctx context.Context, id uuid.UUID) (*biz.PaymentResp, error) {
	// TODO implement me
	panic("implement me")
}

func (r *paymentRepo) ProcessPaymentCallback(ctx context.Context, req *biz.PaymentCallbackReq) (*biz.PaymentCallbackResp, error) {
	// TODO implement me
	panic("implement me")
}

func (r *paymentRepo) CreatePayment(ctx context.Context, req *biz.CreatePaymentReq) (*biz.CreatePaymentRes, error) {
	// types.Float64ToNumeric(req.Amount)

	tradeNo := fmt.Sprintf("%d", xid.Next())
	pay := alipay.TradePagePay{ // 电脑网站支付
		Trade: alipay.Trade{
			Subject:     "支付测试:" + tradeNo,                    // 订单主题
			OutTradeNo:  fmt.Sprintf("%d", time.Now().Unix()), // 商户订单号，必须唯一
			TotalAmount: req.Amount,                           // 订单金额
			ProductCode: "FAST_INSTANT_TRADE_PAY",             // 电脑网站支付，产品码为固定值
			NotifyURL:   r.conf.Alipay.NotifyUrl,              // 异步通知地址
			ReturnURL:   r.conf.Alipay.ReturnUrl,              // 回调地址
		},
	}
	url, err := r.data.alipay.TradePagePay(pay) // 生成支付链接
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}
	// amount,_ := types.Float64ToNumeric(req.Amount)
	payments, err := r.data.DB(ctx).CreatePaymentQuery(ctx, models.CreatePaymentQueryParams{
		// PaymentID:   ,
		OrderID: req.OrderId,
		// Amount:      amount,
		Currency:    req.Currency,
		Method:      req.Method,
		Status:      req.Status,
		GatewayTxID: req.GatewayTxID,
		// Metadata:    req.Metadata,
		// CreatedAt:   pgtype.Timestamptz{},
		// UpdatedAt:   pgtype.Timestamptz{},
	})

	fmt.Printf("payments: %v\n", payments)
	if err != nil {
		return nil, err
	}

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
