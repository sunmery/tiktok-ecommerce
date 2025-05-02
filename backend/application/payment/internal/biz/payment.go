package biz

import (
	"context"
	"net/url"
	"time"

	"backend/constants"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
)

// PaymentStatus 支付状态
type PaymentStatus string

// 通用支付状态
const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusSuccess    PaymentStatus = "SUCCESS"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusClosed     PaymentStatus = "CLOSED"
)

// 支付宝支付状态
const (
	AliPayStatusPending PaymentStatus = "WAIT_BUYER_PAY"
	AliPayStatusClosed  PaymentStatus = "TRADE_CLOSED"
	AliPayStatusSuccess PaymentStatus = "TRADE_SUCCESS"
)

type (
	AliPayCallbackReq struct {
		Charset     string
		OutTradeNo  string
		Method      string
		TotalAmount string
		Sign        string
		TradeNo     int64
		SellerId    string
		AuthAppId   string
		AppId       string
		SignType    string
		Timestamp   string
	}
)

// Payment 支付记录
type Payment struct {
	ID         int64
	OrderID    int64
	ConsumerID uuid.UUID
	Amount     float64
	Currency   string
	Subject    string
	Status     PaymentStatus
	TradeNo    string // 支付宝交易号
	PayURL     string
	NotifyTime time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CreatePaymentReq 创建支付请求
type CreatePaymentReq struct {
	OrderID    int64
	ConsumerID uuid.UUID
	// MerchantID      uuid.UUID
	Amount          string
	Currency        string
	Subject         string
	ReturnURL       string
	FreezeId        int64
	ConsumerVersion int64
	MerchanVersions []int64
}

// CreatePaymentResp 创建支付响应
type CreatePaymentResp struct {
	Payment *Payment
}

type (
	Notification struct {
		AuthAppId           string                `json:"auth_app_id"`
		NotifyTime          string                `json:"notify_time"`
		NotifyType          string                `json:"notify_type"`
		NotifyId            string                `json:"notify_id"`
		AppId               string                `json:"app_id"`
		Charset             string                `json:"charset"`
		Version             string                `json:"version"`
		SignType            string                `json:"sign_type"`
		Sign                string                `json:"sign"`
		TradeNo             string                `json:"trade_no"`
		OutTradeNo          string                `json:"out_trade_no"`
		OutRequestNo        string                `json:"out_request_no"`
		OutBizNo            string                `json:"out_biz_no"`
		BuyerId             string                `json:"buyer_id"`
		BuyerLogonId        string                `json:"buyer_logon_id"`
		BuyerOpenId         string                `json:"buyer_open_id"`
		SellerId            string                `json:"seller_id"`
		SellerEmail         string                `json:"seller_email"`
		TradeStatus         constants.TradeStatus `json:"trade_status"`
		RefundStatus        string                `json:"refund_status"`
		RefundReason        string                `json:"refund_reason"`
		RefundAmount        string                `json:"refund_amount"`
		TotalAmount         string                `json:"total_amount"`
		ReceiptAmount       string                `json:"receipt_amount"`
		InvoiceAmount       string                `json:"invoice_amount"`
		BuyerPayAmount      string                `json:"buyer_pay_amount"`
		PointAmount         string                `json:"point_amount"`
		RefundFee           string                `json:"refund_fee"`
		Subject             string                `json:"subject"`
		Body                string                `json:"body"`
		GmtCreate           string                `json:"gmt_create"`
		GmtPayment          string                `json:"gmt_payment"`
		GmtRefund           string                `json:"gmt_refund"`
		GmtClose            string                `json:"gmt_close"`
		FundBillList        string                `json:"fund_bill_list"`
		PassbackParams      string                `json:"passback_params"`
		VoucherDetailList   string                `json:"voucher_detail_list"`
		AgreementNo         string                `json:"agreement_no"`
		ExternalAgreementNo string                `json:"external_agreement_no"`
		DBackStatus         string                `json:"dback_status"`
		DBackAmount         string                `json:"dback_amount"`
		BankAckTime         string                `json:"bank_ack_time"`
	}
)

// PaymentNotifyReq 支付通知请求
type PaymentNotifyReq struct {
	AppID       string
	AuthAppId   string
	TradeNo     string
	Charset     string
	Method      string
	Sign        string
	SignType    string
	OutTradeNo  string
	TotalAmount string
	SellerId    string
	Params      map[string][]string
}

// PaymentNotifyResp 支付通知响应
type PaymentNotifyResp struct {
	Success bool
	Message string
}

// PaymentCallbackReq 支付回调请求
type PaymentCallbackReq struct {
	Params      map[string]string
	OutTradeNo  string
	TradeNo     string
	TotalAmount string
	Subject     string
	TradeStatus string
}

// PaymentCallbackResp 支付回调响应
type PaymentCallbackResp struct {
	Success bool
	Message string
}

// GetPaymentStatusReq 查询支付状态请求
type GetPaymentStatusReq struct {
	PaymentID int64
}

// GetPaymentStatusResp 查询支付状态响应
type GetPaymentStatusResp struct {
	Payment *Payment
}

type GetPaymentByOrderIDRequest struct {
	OrderID     int64
	TotalAmount string
}

type (
	UpdatePaymentStatusRequest struct {
		PaymentId int64
		OrderId   int64
		TradeNo   string
		Status    PaymentStatus
	}
	UpdatePaymentStatusResponse struct{}
)

// PaymentRepo 支付仓储接口
type PaymentRepo interface {
	// CreatePayment 创建支付记录
	CreatePayment(ctx context.Context, req *CreatePaymentReq) (*CreatePaymentResp, error)
	// GetPaymentStatus 查询支付状态
	GetPaymentStatus(ctx context.Context, req *GetPaymentStatusReq) (*GetPaymentStatusResp, error)
	// HandlePaymentNotify 处理支付通知
	HandlePaymentNotify(ctx context.Context, req url.Values) (*PaymentNotifyResp, error)
	// HandlePaymentCallback 处理支付回调
	HandlePaymentCallback(ctx context.Context, req *PaymentCallbackReq) (*PaymentCallbackResp, error)
	// GetPaymentByOrderID 根据订单ID查询支付记录
	GetPaymentByOrderID(ctx context.Context, req *GetPaymentByOrderIDRequest) (*Payment, error)
}

// PaymentUsecase 支付用例
type PaymentUsecase struct {
	repo PaymentRepo
	log  *log.Helper
	// ordersvc OrderService
}

// // OrderService 订单服务接口
// type OrderService interface {
// 	// 标记订单为已支付
// 	MarkOrderPaid(ctx context.Context, orderID string, userID string) error
// }

// NewPaymentUsecase 创建支付用例
func NewPaymentUsecase(repo PaymentRepo, logger log.Logger) *PaymentUsecase {
	return &PaymentUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// CreatePayment 创建支付
func (uc *PaymentUsecase) CreatePayment(ctx context.Context, req *CreatePaymentReq) (*CreatePaymentResp, error) {
	uc.log.WithContext(ctx).Debugf("Creating payment for order %d", req.OrderID)
	return uc.repo.CreatePayment(ctx, req)
}

// GetPaymentStatus 查询支付状态
func (uc *PaymentUsecase) GetPaymentStatus(ctx context.Context, req *GetPaymentStatusReq) (*GetPaymentStatusResp, error) {
	uc.log.WithContext(ctx).Debugf("Getting payment status for payment %d", req.PaymentID)
	return uc.repo.GetPaymentStatus(ctx, req)
}

// HandlePaymentNotify 处理支付通知
func (uc *PaymentUsecase) HandlePaymentNotify(ctx context.Context, req url.Values) (*PaymentNotifyResp, error) {
	uc.log.WithContext(ctx).Debugf("Handling payment notify for order %s", req)
	return uc.repo.HandlePaymentNotify(ctx, req)
	//
	// // 处理支付宝通知
	// resp, err := uc.repo.HandlePaymentNotify(ctx, req)
	// if err != nil {
	// 	uc.log.WithContext(ctx).Errorf("Failed to handle payment notify: %v", err)
	// 	return nil, err
	// }
	//
	// // 检查支付状态是否为成功
	// if req.TradeStatus == "TRADE_SUCCESS" || req.TradeStatus == "TRADE_FINISHED" {
	// 	// 更新支付状态
	// 	err = uc.repo.UpdatePaymentStatus(ctx, "", req.OutTradeNo, req.TradeNo, PaymentStatusSuccess)
	// 	if err != nil {
	// 		uc.log.WithContext(ctx).Errorf("Failed to update payment status: %v", err)
	// 		return nil, fmt.Errorf("更新支付状态失败: %w", err)
	// 	}
	//
	// 	// 查询支付记录
	// 	payment, err := uc.repo.GetPaymentByOrderID(ctx, models.&GetPaymentByOrderIDRequest)
	// 	if err != nil {
	// 		uc.log.WithContext(ctx).Errorf("Failed to get payment by order ID: %v", err)
	// 		return nil, fmt.Errorf("查询支付记录失败: %w", err)
	// 	}
	//
	// 	// 标记订单为已支付
	// 	err = uc.ordersvc.MarkOrderPaid(ctx, req.OutTradeNo, payment.UserID)
	// 	if err != nil {
	// 		uc.log.WithContext(ctx).Errorf("Failed to mark order as paid: %v", err)
	// 		return nil, fmt.Errorf("标记订单为已支付失败: %w", err)
	// 	}
	//
	// 	uc.log.WithContext(ctx).Debugf("Payment for order %s is successful", req.OutTradeNo)
	// }
	//
	// return resp, nil
}

// HandlePaymentCallback 处理支付回调
func (uc *PaymentUsecase) HandlePaymentCallback(ctx context.Context, req *PaymentCallbackReq) (*PaymentCallbackResp, error) {
	uc.log.WithContext(ctx).Debugf("Handling payment callback for order %s", req.OutTradeNo)
	return uc.repo.HandlePaymentCallback(ctx, req)
}

// GetPaymentByOrderID 根据订单ID查询支付记录
func (uc *PaymentUsecase) GetPaymentByOrderID(ctx context.Context, req *GetPaymentByOrderIDRequest) (*Payment, error) {
	uc.log.WithContext(ctx).Debugf("Updating payment status for payment %d", req)
	return uc.repo.GetPaymentByOrderID(ctx, req)
}
