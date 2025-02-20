package service

import (
	"context"

	pb "backend/api/payment/v1"
)

type PaymentServiceService struct {
	pb.UnimplementedPaymentServiceServer
}

func NewPaymentServiceService() *PaymentServiceService {
	return &PaymentServiceService{}
}

func (s *PaymentServiceService) Charge(ctx context.Context, req *pb.ChargeReq) (*pb.ChargeResp, error) {
	return &pb.ChargeResp{}, nil
}
