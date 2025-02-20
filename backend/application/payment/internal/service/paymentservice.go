package service

import (
	"backend/application/payment/internal/biz"
	"context"

	pb "backend/api/payment/v1"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	oc *biz.PaymentUsecase
}

func (s *PaymentService) Charge(ctx context.Context, req *pb.ChargeReq) (*pb.ChargeResp, error) {
	return &pb.ChargeResp{}, nil
}
