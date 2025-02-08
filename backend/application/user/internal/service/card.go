package service

import (
	"context"

	pb "backend/api/user/v1"
)

func (s *UserService) CreateCreditCard(ctx context.Context, req *pb.CreditCards) (*pb.CardsReply, error) {
	s.uc.
	return &pb.CardsReply{

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
