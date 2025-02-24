package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"google.golang.org/protobuf/types/known/timestamppb"

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
	orderId, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, err
	}
	result, err := s.uc.CreatePayment(ctx, &biz.CreatePaymentReq{
		OrderId:       orderId,
		Currency:      req.Currency,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Metadata:      req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResp{
		PaymentId:  result.PaymentId,
		Status:     result.Status,
		PaymentUrl: result.PaymentUrl,
		CreatedAt:  timestamppb.New(result.CreatedAt),
	}, nil
}

func (s *PaymentServiceService) GetPayment(ctx context.Context, req *pb.GetPaymentReq) (*pb.PaymentResp, error) {
	paymentId, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, err
	}
	payment, err := s.uc.GetPayment(ctx, paymentId)
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResp{
		PaymentId:  payment.PaymentId,
		Status:     payment.Status,
		PaymentUrl: payment.PaymentUrl,
		CreatedAt:  timestamppb.New(payment.CreatedAt),
	}, nil
}

func (s *PaymentServiceService) ProcessPaymentCallback(ctx context.Context, req *pb.PaymentCallbackReq) (*pb.PaymentCallbackResp, error) {
	callback, err := s.uc.ProcessPaymentCallback(ctx, &biz.PaymentCallbackReq{
		PaymentId:       "",
		Status:          "",
		GatewayResponse: "",
		ProcessedAt:     time.Time{},
	})
	fmt.Printf("callback: %v\n", callback)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
