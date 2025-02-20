package service

import (
	"backend/application/payment/internal/biz"
	"context"

	pb "backend/api/payment/v1"
)

type PaymentServiceService struct {
	pb.UnimplementedPaymentServiceServer
	oc *biz.PaymentUsecase
}

func NewPaymentServiceService(oc *biz.PaymentUsecase) *PaymentServiceService {
	return &PaymentServiceService{
		oc: oc,
	}
}

func (s *PaymentServiceService) Charge(ctx context.Context, req *pb.ChargeReq) (*pb.ChargeResp, error) {
	return &pb.ChargeResp{}, nil
}
