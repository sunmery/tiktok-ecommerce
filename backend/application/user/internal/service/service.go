package service

import (
	userV1 "backend/api/user/v1"
	"backend/application/user/internal/biz"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService)

type UserService struct {
	userV1.UnimplementedUserServiceServer

	uc *biz.UserUsecase
	ac *biz.AddressesUsecase
}
