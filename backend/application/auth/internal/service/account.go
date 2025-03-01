package service

import (
	v1 "backend/api/auth/v1"
	"backend/application/auth/internal/biz"
	"context"
	"fmt"
)

func (s *AuthService) Signin(ctx context.Context, req *v1.SigninRequest) (*v1.SigninReply, error) {
	result, err := s.ac.Signin(ctx, &biz.SigninRequest{
		State: req.State,
		Code:  req.Code,
	})
	fmt.Printf("result:%+v", result)
	if err != nil {
		return nil, err
	}
	return &v1.SigninReply{
		State: result.State,
		Data:  result.Data,
	}, nil
}
