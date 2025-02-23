package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"backend/application/payment/internal/biz"

	pb "backend/api/payment/v1"
)

type PaymentServiceService struct {
	pb.UnimplementedPaymentServiceServer

	uc *biz.PaymentUsecase
}

func NewPaymentServiceService(uc *biz.PaymentUsecase) *PaymentServiceService {
	return &PaymentServiceService{uc: uc}
}

func (s *PaymentServiceService) CreatePayment(ctx context.Context, req *pb.CreatePaymentReq) (*pb.PaymentResp, error) {
	err := s.uc.CreatePayment(ctx, &biz.Payment{
		PaymentID:   uuid.UUID{},
		OrderID:     uuid.UUID{},
		Amount:      0,
		Currency:    "",
		Method:      "",
		Status:      "",
		GatewayTxID: "",
		Metadata:    nil,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	})
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResp{
		PaymentId:  "",
		Status:     "",
		PaymentUrl: "",
		CreatedAt:  nil,
	}, nil
}

func (s *PaymentServiceService) GetPayment(ctx context.Context, req *pb.GetPaymentReq) (*pb.PaymentResp, error) {
	payment, err := s.uc.GetByID(ctx, uuid.UUID{})
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResp{
		PaymentId:  payment.PaymentID.String(),
		Status:     string(payment.Status),
		PaymentUrl: "",
		CreatedAt:  nil,
	}, nil
}

func (s *PaymentServiceService) ProcessPaymentCallback(ctx context.Context, req *pb.PaymentCallbackReq) (*pb.PaymentCallbackResp, error) {
	s.uc.
	return nil, nil
}
