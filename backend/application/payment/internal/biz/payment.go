package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"   // 待支付
	PaymentSucceeded PaymentStatus = "SUCCEEDED" // 支付成功
	PaymentFailed    PaymentStatus = "FAILED"    // 支付失败
)

type (
	PaymentCallbackReq struct {
		PaymentId       string
		Status          PaymentStatus
		GatewayResponse string
		ProcessedAt     time.Time
		RequestForm     map[string][]string
	}
	PaymentCallbackResp struct{}
)

type Payment struct {
	PaymentID   uuid.UUID
	OrderID     uuid.UUID
	Amount      float64
	Currency    string
	Method      string
	Status      PaymentStatus
	GatewayTxID string
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type (
	CreatePaymentReq struct {
		Subject       string // 订单主题, 例如Iphone6 16G, 不可使用特殊字符，如 /，=，& 等。
		OrderId       uuid.UUID
		Currency      string
		Amount        string
		PaymentMethod string
		Method        string
		Status        PaymentStatus
		Metadata      map[string]string
	}
	CreatePaymentRes struct {
		PaymentId  string
		Status     string
		PaymentUrl string
		CreatedAt  time.Time
	}
)

type PaymentResp struct {
	PaymentId  string
	Status     string
	PaymentUrl string
	CreatedAt  time.Time
}

type (
	PaymentNotifyReq struct {
		Values map[string][]string
	}
	PaymentNotifyResp struct {
		Code int
		Msg  string
	}
)

type PaymentRepo interface {
	CreatePayment(ctx context.Context, req *CreatePaymentReq) (*CreatePaymentRes, error)
	GetPayment(ctx context.Context, id uuid.UUID) (*PaymentResp, error)
	PaymentNotify(ctx context.Context, req *PaymentNotifyReq) (*PaymentNotifyResp, error)
	ProcessPaymentCallback(ctx context.Context, req *PaymentCallbackReq) (*PaymentCallbackResp, error)
}

func (pc *PaymentUsecase) CreatePayment(ctx context.Context, req *CreatePaymentReq) (*CreatePaymentRes, error) {
	return pc.repo.CreatePayment(ctx, req)
}

func (pc *PaymentUsecase) GetPayment(ctx context.Context, id uuid.UUID) (*PaymentResp, error) {
	return pc.repo.GetPayment(ctx, id)
}

func (pc *PaymentUsecase) ProcessPaymentCallback(ctx context.Context, req *PaymentCallbackReq) (*PaymentCallbackResp, error) {
	return pc.repo.ProcessPaymentCallback(ctx, req)
}

func (pc *PaymentUsecase) PaymentNotify(ctx context.Context, req *PaymentNotifyReq) (*PaymentNotifyResp, error) {
	return pc.repo.PaymentNotify(ctx, req)
}

type PaymentUsecase struct {
	repo PaymentRepo
	log  *log.Helper
}

func NewPaymentUsecase(repo PaymentRepo, logger log.Logger) *PaymentUsecase {
	return &PaymentUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
