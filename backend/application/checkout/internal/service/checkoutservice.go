package service

import (
	pb "backend/api/checkout/v1"
	"backend/application/checkout/internal/biz"
	"backend/pkg"
	"context"
	"fmt"

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
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}
	fmt.Printf("user id: %v\n", userId)
	reply, err := s.uc.Checkout(ctx, &biz.CheckoutRequest{
		UserId: userId,
		// Currency:   req.Currency,
		Currency:      "CNY",
		Firstname:     req.Firstname,
		Lastname:      req.Lastname,
		Email:         req.Email,
		CreditCard:    req.CreditCard,
		Address:       req.Address,
		PaymentMethod: "alipay",
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
