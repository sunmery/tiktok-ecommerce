package biz

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
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

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"   // 待支付
	PaymentSucceeded PaymentStatus = "SUCCEEDED" // 支付成功
	PaymentFailed    PaymentStatus = "FAILED"    // 支付失败
)

type PaymentRepo interface {
	CreatePayment(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, paymentID uuid.UUID) (*Payment, error)
	UpdateStatus(ctx context.Context, paymentID uuid.UUID, status PaymentStatus, gatewayTxID string) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*Payment, error)
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

func (pc *PaymentUsecase) CreatePayment(ctx context.Context, payment *Payment) error {
	return pc.repo.CreatePayment(ctx, payment)
}

func (pc *PaymentUsecase) GetByID(ctx context.Context, paymentID uuid.UUID) (*Payment, error) {
	return pc.repo.GetByID(ctx, paymentID)
}

func (pc *PaymentUsecase) UpdateStatus(ctx context.Context, paymentID uuid.UUID, status PaymentStatus, gatewayTxID string) error {
	return pc.repo.UpdateStatus(ctx, paymentID, status, gatewayTxID)
}

func (pc *PaymentUsecase) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*Payment, error) {
	return pc.repo.GetByOrderID(ctx, orderID)
}
