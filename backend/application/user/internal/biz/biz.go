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
	GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error)

	// 地址接口
	CreateAddress(ctx context.Context, req *Address) (*Address, error)
	UpdateAddress(ctx context.Context, req *Address) (*Address, error)
	DeleteAddress(ctx context.Context, req *DeleteAddressesRequest) (*DeleteAddressesReply, error)
	GetAddresses(ctx context.Context, req *Request) (*Addresses, error)

	// 银行卡接口
	CreateCreditCard(ctx context.Context, req *CreditCards) (*CreditCardsReply, error)
	UpdateCreditCard(ctx context.Context, req *CreditCards) (*CreditCardsReply, error)
	DeleteCreditCard(ctx context.Context, req *DeleteCreditCardsRequest) (*CreditCardsReply, error)
	GetCreditCard(ctx context.Context, req *GetCreditCardsRequest) (*CreditCards, error)
	SearchCreditCards(ctx context.Context, req *GetCreditCardsRequest) ([]*CreditCards, error)
	ListCreditCards(ctx context.Context, req *CreditCardsRequest) ([]*CreditCards, error)
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

func (cc *UserUsecase) GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error) {
	// cc.log.WithContext(ctx).Infof("GetProfile request: %+v", req)
	return cc.repo.GetProfile(ctx, req)
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
