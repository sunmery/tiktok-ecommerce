package biz

import (
	"context"

	"github.com/google/uuid"

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

type Favorites struct {
	Items []*Product
}

type GetFavoritesRequest struct {
	UserId   uuid.UUID
	Page     int32
	PageSize int32
}

type (
	UpdateFavoritesRequest struct {
		UserId    uuid.UUID
		ProductId uuid.UUID
	}
	UpdateFavoritesResply struct {
		Message string
		Code    int32
	}
)

type UserRepo interface {
	GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error)
	GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersReply, error)
	DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserReply, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserReply, error)

	GetFavorites(ctx context.Context, req *GetFavoritesRequest) (*Favorites, error)
	DeleteFavorites(ctx context.Context, req *UpdateFavoritesRequest) (*UpdateFavoritesResply, error)
	SetFavorites(ctx context.Context, req *UpdateFavoritesRequest) (*UpdateFavoritesResply, error)

	CreateAddress(ctx context.Context, req *Address) (*Address, error)
	UpdateAddress(ctx context.Context, req *Address) (*Address, error)
	DeleteAddress(ctx context.Context, req *AddressRequest) (*DeleteAddressesReply, error)
	GetAddress(ctx context.Context, req *AddressRequest) (*Address, error)
	GetAddresses(ctx context.Context, req *Request) (*Addresses, error)

	CreateCreditCard(ctx context.Context, req *CreditCard) (*emptypb.Empty, error)
	DeleteCreditCard(ctx context.Context, req *DeleteCreditCardRequest) (*emptypb.Empty, error)
	GetCreditCard(ctx context.Context, req *GetCreditCardRequest) (*CreditCard, error)
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
