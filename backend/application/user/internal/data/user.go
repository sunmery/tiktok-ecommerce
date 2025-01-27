package data

import (
	"context"

	"backend/application/user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func (u *userRepo) Register(ctx context.Context, user *biz.RegisterReq) (*biz.RegisterResp, error) {
	// TODO implement me
	panic("implement me")
}

func (u*userRepo) Login(ctx context.Context, user *biz.LoginReq) (*biz.LoginResp, error) {
	// TODO implement me
	panic("implement me")
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
