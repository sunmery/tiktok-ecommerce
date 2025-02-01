package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUsecase)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type UserRepo interface {
	Signin(ctx context.Context, req *SigninRequest) (*SigninReply, error)
	GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoReply, error)
	CreateAddress(ctx context.Context, req *Address) (*Address, error)
	UpdateAddress(ctx context.Context, req *Address) (*Address, error)
	DeleteAddress(ctx context.Context, req *DeleteAddressesRequest) (*DeleteAddressesReply, error)
	GetAddresses(ctx context.Context, req *Request) (*Addresses, error)
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

func (cc *UserUsecase) CreateAddress(ctx context.Context, req *Address) (*Address, error) {
	cc.log.WithContext(ctx).Infof("CreateAddress: %+v", req)
	return cc.repo.CreateAddress(ctx, req)
}

func (cc *UserUsecase) UpdateAddress(ctx context.Context, req *Address) (*Address, error) {
	cc.log.WithContext(ctx).Infof("UpdateAddress: %+v", req)
	return cc.repo.UpdateAddress(ctx, req)
}

func (cc *UserUsecase) DeleteAddress(ctx context.Context, req *DeleteAddressesRequest) (*DeleteAddressesReply, error) {
	cc.log.WithContext(ctx).Infof("DeleteAddress: %+v", req)
	return cc.repo.DeleteAddress(ctx, req)
}

func (cc *UserUsecase) GetAddresses(ctx context.Context, req *Request) (*Addresses, error) {
	cc.log.WithContext(ctx).Infof("GetAddresses: %+v", req)
	return cc.repo.GetAddresses(ctx, req)
}
