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

type PaymentCallbackReq struct {
	PaymentId       string
	Status          string
	GatewayResponse string
	ProcessedAt     time.Time
}

type PaymentCallbackResp struct{}

type CreatePaymentReq struct {
	OrderId       uuid.UUID
	Currency      string
	Amount        string
	PaymentMethod string
	Method        string
	Status        string
	GatewayTxID   *string
	Metadata      map[string]string
}

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
type UpdateStatusReq struct{}

type (
	UpdateStatusRes  struct{}
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
type PaymentRepo interface {
	CreatePayment(ctx context.Context, req *CreatePaymentReq) (*CreatePaymentRes, error)
	GetPayment(ctx context.Context, id uuid.UUID) (*PaymentResp, error)
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
