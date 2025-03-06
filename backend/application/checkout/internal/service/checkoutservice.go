package service

import (
	"backend/application/checkout/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/metadata"

	pb "backend/api/checkout/v1"
)

type CheckoutServiceService struct {
	pb.UnimplementedCheckoutServiceServer
	uc *biz.CheckoutUsecase
}

func NewCheckoutServiceService(uc *biz.CheckoutUsecase) *CheckoutServiceService {
	return &CheckoutServiceService{uc: uc}
}

func (s *CheckoutServiceService) Checkout(ctx context.Context, req *pb.CheckoutReq) (*pb.CheckoutResp, error) {
	var userMdKey string
	if md, ok := metadata.FromServerContext(ctx); ok{
		userMdKey = md.Get("x-md-global-user-id")
	}
	reply, err := s.uc.Checkout(ctx, &biz.CheckoutRequest{
		UserId:     userMdKey,
		Firstname:  req.Firstname,
		Lastname:   req.Lastname,
		Email:      req.Email,
		CreditCard: nil,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CheckoutResp{
		OrderId:       reply.OrderId,
		TransactionId: reply.TransactionId,
	}, nil
}
