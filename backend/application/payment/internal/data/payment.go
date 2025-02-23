package data

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"backend/application/payment/internal/data/models"

	"backend/application/payment/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type paymentRepo struct {
	data *Data
	log  *log.Helper
}

func (r *paymentRepo) CreatePayment(ctx context.Context, req *biz.Payment) error {
	payments, err := r.data.DB(ctx).CreatePaymentQuery(ctx, models.CreatePaymentQueryParams{
		// PaymentID:   req.PaymentID,
		// OrderID:     uuid.UUID{},
		// Amount:      pgtype.Numeric{},
		// Currency:    "",
		// Method:      "",
		// Status:      "",
		// GatewayTxID: nil,
		// Metadata:    nil,
		// CreatedAt:   pgtype.Timestamptz{},
		// UpdatedAt:   pgtype.Timestamptz{},
	})

	fmt.Printf("payments: %v\n", payments)
	if err != nil {
		return err
	}

	return err
}

func (r *paymentRepo) UpdateStatus(ctx context.Context, paymentID uuid.UUID, status biz.PaymentStatus, gatewayTxID string) error {
	query, err := r.data.DB(ctx).UpdateStatusQuery(ctx, models.UpdateStatusQueryParams{
		// PaymentID:   uuid.UUID{},
		// Status:      "",
		// GatewayTxID: nil,
		// UpdatedAt:   pgtype.Timestamptz{},
	})
	fmt.Printf("query: %v\n", query)
	if err != nil {
		return err
	}
	return err
}

func (r *paymentRepo) GetByID(ctx context.Context, paymentID uuid.UUID) (*biz.Payment, error) {
	query, err := r.data.DB(ctx).GetByIDQuery(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	return &biz.Payment{
		PaymentID: query.PaymentID,
		OrderID:   query.OrderID,
		// Amount:      query.Amount,
		Currency: query.Currency,
		Method:   query.Method,
		Status:   biz.PaymentStatus(query.Status),
		// GatewayTxID: query.GatewayTxID,
		// Metadata:    query.Metadata,
		CreatedAt: query.CreatedAt,
		UpdatedAt: query.UpdatedAt,
	}, nil
}

func (r *paymentRepo) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*biz.Payment, error) {
	query, err := r.data.DB(ctx).GetByOrderIDQuery(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return &biz.Payment{
		PaymentID: query.PaymentID,
		OrderID:   query.OrderID,
		// Amount:      query.Amount,
		Currency: query.Currency,
		Method:   query.Method,
		Status:   biz.PaymentStatus(query.Status),
		// GatewayTxID: query.GatewayTxID,
		// Metadata:    query.Metadata,
		CreatedAt: query.CreatedAt,
		UpdatedAt: query.UpdatedAt,
	}, nil
}

func NewPaymentRepo(data *Data,
	log *log.Helper,
) biz.PaymentRepo {
	return &paymentRepo{
		data: data,
		log:  log,
	}
}
