package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewAuthUsecase)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type AuthRepo interface {
	Signin(ctx context.Context, req *SigninRequest) (*SigninReply, error)
	GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoReply, error)
}

type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
}

func NewAuthUsecase(repo AuthRepo, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (cc *AuthUsecase) Signin(ctx context.Context, req *SigninRequest) (*SigninReply, error) {
	cc.log.WithContext(ctx).Infof("Signin request: %+v", req)
	return cc.repo.Signin(ctx, req)
}

func (cc *AuthUsecase) GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoReply, error) {
	cc.log.WithContext(ctx).Infof("GetUserInfo request: %+v", req)
	return cc.repo.GetUserInfo(ctx, req)
}
