package service

import (
	pb "backend/api/checkout/v1"
	"backend/application/checkout/internal/biz"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type CheckoutServiceService struct {
	pb.UnimplementedCheckoutServiceServer
	uc *biz.CheckoutUsecase
}

func NewCheckoutServiceService(uc *biz.CheckoutUsecase) *CheckoutServiceService {
	return &CheckoutServiceService{uc: uc}
}

func (s *CheckoutServiceService) Checkout(ctx context.Context, req *pb.CheckoutReq) (*pb.CheckoutResp, error) {
	// userID, err := pkg.GetMetadataUesrID(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("get user id from metadata error: %v", err)
	// }
	userID := "77d08975-972c-4a06-8aa4-d2d23f374bb1"
	userId,err := uuid.Parse(userID)
	if err!= nil {
		return nil, fmt.Errorf("invalid user id")
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
		OrderId:       reply.OrderId,
		PaymentId:     reply.PaymentId,
		PaymentUrl:    reply.PaymentURL,
	}, nil
}
