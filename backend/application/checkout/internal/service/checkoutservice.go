package service

import (
	pb "backend/api/checkout/v1"
	"backend/application/checkout/internal/biz"
	"backend/pkg"
	"context"
	"fmt"
)

type CheckoutServiceService struct {
	pb.UnimplementedCheckoutServiceServer
	uc *biz.CheckoutUsecase
}

func NewCheckoutServiceService(uc *biz.CheckoutUsecase) *CheckoutServiceService {
	return &CheckoutServiceService{uc: uc}
}

func (s *CheckoutServiceService) Checkout(ctx context.Context, req *pb.CheckoutReq) (*pb.CheckoutResp, error) {
	userID, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user id from metadata error: %v", err)
	}
	reply, err := s.uc.Checkout(ctx, &biz.CheckoutRequest{
		UserId:     userID,
		Firstname:  req.Firstname,
		Lastname:   req.Lastname,
		Email:      req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CheckoutResp{
		OrderId:       reply.OrderId,
		TransactionId: reply.TransactionId,
	}, nil
}
