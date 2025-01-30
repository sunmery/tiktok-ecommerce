package biz

import (
	"context"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"

	// "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type SigninRequest struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}
type SigninReply struct {
	State string `json:"state,omitempty"`
	Data  string `json:"data,omitempty"`
}
type GetUserInfoRequest struct {
	Authorization string
}

type GetUserInfoReply struct {
	State string          `json:"state,omitempty"`
	Data  casdoorsdk.User `json:"data"`
}

type UserRepo interface {
	Signin(ctx context.Context, req *SigninRequest) (*SigninReply, error)
	GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoReply, error)
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
	cc.log.WithContext(ctx).Infof("Signin request: %+v", req)
	return cc.repo.Signin(ctx, req)
}

func (cc *UserUsecase) GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoReply, error) {
	cc.log.WithContext(ctx).Infof("GetUserInfo request: %+v", req)
	return cc.repo.GetUserInfo(ctx, req)
}
