package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAuthUsecase)

type AuthRepo interface {
	Signin(ctx context.Context, req *SigninRequest) (*SigninReply, error)
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
