package service

import (
	"backend/application/user/internal/biz"
	"backend/pkg/token"
	"context"

	pb "backend/api/user/v1"
)

func (s *UserService) CreateCreditCard(ctx context.Context, req *pb.CreditCards) (*pb.CardsReply, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.uc.CreateCreditCard(ctx, &biz.CreditCards{
		Owner: payload.Owner,
		Name:  payload.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CardsReply{
		Message: "OK",
		Code:    200,
	}, nil
}
func (s *UserService) UpdateCreditCard(ctx context.Context, req *pb.CreditCards) (*pb.CardsReply, error) {
	return &pb.CardsReply{}, nil
}
func (s *UserService) DeleteCreditCard(ctx context.Context, req *pb.DeleteCreditCardsRequest) (*pb.CardsReply, error) {
	return &pb.CardsReply{}, nil
}
func (s *UserService) GetCreditCard(ctx context.Context, req *pb.GetCreditCardsRequest) (*pb.GetCreditCardsReply, error) {
	return &pb.GetCreditCardsReply{}, nil
}
func (s *UserService) ListCreditCards(ctx context.Context, req *pb.ListCreditCardsRequest) (*pb.ListCreditCardsReply, error) {
	return &pb.ListCreditCardsReply{}, nil
}
