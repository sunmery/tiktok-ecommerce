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
	o, err := s.oc.Create(ctx, &biz.CreateRequest{
		Amount: float64(req.Amount),
		CreditCard: biz.CreditCard{
			Number:          req.CreditCard.CreditCardNumber,
			CVV:             req.CreditCard.CreditCardCvv,
			ExpirationYear:  req.CreditCard.CreditCardExpirationYear,
			ExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
		},
		OrderID: "",
		UserID:  0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ChargeResp{TransactionId: o.TransactionID}, nil
}
