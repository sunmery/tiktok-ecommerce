package service

import (
	"context"
	"fmt"

	pb "backend/api/checkout/v1"
	"backend/application/checkout/internal/biz"
	"backend/pkg"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CheckoutServiceService struct {
	pb.UnimplementedCheckoutServiceServer
	uc *biz.CheckoutUsecase
}

func NewCheckoutServiceService(uc *biz.CheckoutUsecase) *CheckoutServiceService {
	return &CheckoutServiceService{uc: uc}
}

func (s *CheckoutServiceService) Checkout(ctx context.Context, req *pb.CheckoutReq) (*pb.CheckoutResp, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("get user id failed: %v", err))
	}

	reply, err := s.uc.Checkout(ctx, &biz.CheckoutRequest{
		UserId:        userId,
		Firstname:     req.Firstname,
		Lastname:      req.Lastname,
		Email:         req.Email,
		CreditCardId:  req.CreditCardId,
		AddressId:     req.AddressId,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Phone:         req.Phone,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CheckoutResp{
		OrderId:    reply.OrderId,
		PaymentId:  reply.PaymentId,
		PaymentUrl: reply.PaymentURL,
	}, nil
}
