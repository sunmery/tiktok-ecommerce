package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUsecase)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type UserRepo interface {
	GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error)

	CreateAddress(ctx context.Context, req *Address) (*Address, error)
	UpdateAddress(ctx context.Context, req *Address) (*Address, error)
	DeleteAddress(ctx context.Context, req *DeleteAddressesRequest) (*DeleteAddressesReply, error)
	GetAddresses(ctx context.Context, req *Request) (*Addresses, error)

	CreateCreditCard(ctx context.Context, req *CreditCard) (*emptypb.Empty, error)
	DeleteCreditCard(ctx context.Context, req *DeleteCreditCardRequest) (*emptypb.Empty, error)
	ListCreditCards(ctx context.Context, req *ListCreditCardsRequest) (*CreditCards, error)
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
