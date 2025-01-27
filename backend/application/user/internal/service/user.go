package service

import (
	v1 "backend/api/user/v1"
	"backend/application/user/internal/biz"
	"context"
)

type UserService struct {
	v1.UnimplementedUserServiceServer

	uc *biz.UserUsecase
}

// NewUserService new a User service.
func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) Register(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterResp, error) {
	res, err := s.uc.Register(ctx, &biz.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &v1.RegisterResp{
		UserId: res.UserId,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginResp, error) {
	res, err := s.uc.Login(ctx, &biz.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &v1.LoginResp{
		UserId:res.UserId,
	}, nil
}
