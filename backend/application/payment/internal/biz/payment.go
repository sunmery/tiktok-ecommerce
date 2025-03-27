package biz

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
)

// PaymentStatus 支付状态
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusSuccess    PaymentStatus = "SUCCESS"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusClosed     PaymentStatus = "CLOSED"
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
	UserID     uuid.UUID
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
	OrderID   int64
	UserID    uuid.UUID
	Amount    string
	Currency  string
	Subject   string
	ReturnURL string
}

// CreatePaymentResp 创建支付响应
type CreatePaymentResp struct {
	Payment *Payment
}

// PaymentNotifyReq 支付通知请求
type PaymentNotifyReq struct {
	AppID       string
	TradeNo     string
	OutTradeNo  string
	TotalAmount string
	Subject     string
	TradeStatus string
	GmtPayment  string
	GmtCreate   string
	Sign        string
	SignType    string
	Params      map[string]string
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
	HandlePaymentNotify(ctx context.Context, req *PaymentNotifyReq) (*PaymentNotifyResp, error)
	// HandlePaymentCallback 处理支付回调
	HandlePaymentCallback(ctx context.Context, req *PaymentCallbackReq) (*PaymentCallbackResp, error)
	// UpdatePaymentStatus 更新支付状态
	UpdatePaymentStatus(ctx context.Context, req *UpdatePaymentStatusRequest) (*UpdatePaymentStatusResponse, error)
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
	uc.log.WithContext(ctx).Infof("Creating payment for order %d", req.OrderID)
	return uc.repo.CreatePayment(ctx, req)
}

// GetPaymentStatus 查询支付状态
func (uc *PaymentUsecase) GetPaymentStatus(ctx context.Context, req *GetPaymentStatusReq) (*GetPaymentStatusResp, error) {
	uc.log.WithContext(ctx).Infof("Getting payment status for payment %d", req.PaymentID)
	return uc.repo.GetPaymentStatus(ctx, req)
}

// HandlePaymentNotify 处理支付通知
func (uc *PaymentUsecase) HandlePaymentNotify(ctx context.Context, req *PaymentNotifyReq) (*PaymentNotifyResp, error) {
	uc.log.WithContext(ctx).Infof("Handling payment notify for order %s", req.OutTradeNo)
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
	// 	uc.log.WithContext(ctx).Infof("Payment for order %s is successful", req.OutTradeNo)
	// }
	//
	// return resp, nil
}

// HandlePaymentCallback 处理支付回调
func (uc *PaymentUsecase) HandlePaymentCallback(ctx context.Context, req *PaymentCallbackReq) (*PaymentCallbackResp, error) {
	uc.log.WithContext(ctx).Infof("Handling payment callback for order %s", req.OutTradeNo)
	return uc.repo.HandlePaymentCallback(ctx, req)
}

// UpdatePaymentStatus 更新支付状态
func (uc *PaymentUsecase) UpdatePaymentStatus(ctx context.Context, req *UpdatePaymentStatusRequest) (*UpdatePaymentStatusResponse, error) {
	uc.log.WithContext(ctx).Infof("Updating payment status for payment %v", req)
	return uc.repo.UpdatePaymentStatus(ctx, req)
}

// GetPaymentByOrderID 根据订单ID查询支付记录
func (uc *PaymentUsecase) GetPaymentByOrderID(ctx context.Context, req *GetPaymentByOrderIDRequest) (*Payment, error) {
	uc.log.WithContext(ctx).Infof("Updating payment status for payment %d", req)
	return uc.repo.GetPaymentByOrderID(ctx, req)
}
