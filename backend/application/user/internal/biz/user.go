package biz

import (
	"context"

	// "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResp struct {
	UserId int32 `json:"user_id"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	UserId int32 `json:"user_id"`
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Register(context.Context, *RegisterReq) (*RegisterResp, error)
	Login(context.Context, *LoginReq) (*LoginResp, error)
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) Register(ctx context.Context, req *RegisterReq) (*RegisterResp, error) {
	uc.log.WithContext(ctx).Infof("Register: %v", req)
	return uc.repo.Register(ctx, req)
}

func (uc *UserUsecase) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	uc.log.WithContext(ctx).Infof("ctx: %v", req)
	return uc.repo.Login(ctx, req)
}
