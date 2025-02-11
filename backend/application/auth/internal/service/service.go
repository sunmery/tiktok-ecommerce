package service

import (
	authV1 "backend/api/auth/v1"
	"backend/application/auth/internal/biz"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAuthService)

type AuthService struct {
	authV1.UnimplementedAuthServiceServer

	ac *biz.AuthUsecase
}

// NewAuthService new a Auth service.
func NewAuthService(ac *biz.AuthUsecase) *AuthService {
	return &AuthService{ac: ac}
}
