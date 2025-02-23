package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type UserRepo interface {
	Signin(context.Context, *SigninRequest) (*SigninReply, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (cc *UserUsecase) Signin(ctx context.Context, req *SigninRequest) (*SigninReply, error) {
	cc.log.WithContext(ctx).Debugf("Signin request: %+v", req)
	return cc.repo.Signin(ctx, req)
}
