package service

import (
	"context"

	pb "backend/api/user/v1"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceService() *UserServiceService {
	return &UserServiceService{}
}

func (s *UserServiceService) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninReply, error) {
	return &pb.SigninReply{}, nil
}
func (s *UserServiceService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	return &pb.GetUserInfoResponse{}, nil
}
func (s *UserServiceService) CreateAddresses(ctx context.Context, req *pb.Address) (*pb.Address, error) {
	return &pb.Address{}, nil
}
func (s *UserServiceService) UpdateAddresses(ctx context.Context, req *pb.Address) (*pb.Address, error) {
	return &pb.Address{}, nil
}
func (s *UserServiceService) DeleteAddresses(ctx context.Context, req *pb.DeleteAddressesRequest) (*pb.DeleteAddressesReply, error) {
	return &pb.DeleteAddressesReply{}, nil
}
func (s *UserServiceService) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesReply, error) {
	return &pb.GetAddressesReply{}, nil
}
func (s *UserServiceService) CreateCreditCard(ctx context.Context, req *pb.CreditCards) (*pb.CardsReply, error) {
	return &pb.CardsReply{}, nil
}
func (s *UserServiceService) UpdateCreditCard(ctx context.Context, req *pb.CreditCards) (*pb.CardsReply, error) {
	return &pb.CardsReply{}, nil
}
func (s *UserServiceService) DeleteCreditCard(ctx context.Context, req *pb.DeleteCreditCardsRequest) (*pb.CardsReply, error) {
	return &pb.CardsReply{}, nil
}
func (s *UserServiceService) GetCreditCard(ctx context.Context, req *pb.GetCreditCardsRequest) (*pb.GetCreditCardsReply, error) {
	return &pb.GetCreditCardsReply{}, nil
}
func (s *UserServiceService) ListCreditCards(ctx context.Context, req *pb.ListCreditCardsRequest) (*pb.ListCreditCardsReply, error) {
	return &pb.ListCreditCardsReply{}, nil
}
