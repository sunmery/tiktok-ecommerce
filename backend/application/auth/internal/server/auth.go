package server

import (
	"backend/application/auth/internal/conf"
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwtV5 "github.com/golang-jwt/jwt/v5"
)

// NewWhiteListMatcher 创建jwt白名单
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	// whiteList["/admin.v1.AdminService/Login"] = struct{}{}
	whiteList["/api.auth.v1.AuthService/Signin"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// parseRSAPublicKeyFromPEM 解析 RSA 公钥
func parseRSAPublicKeyFromPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	publicKey, err := jwtV5.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}
	return publicKey, nil
}

func InitJwtKey(ac *conf.Auth) *rsa.PublicKey {
	publicKey, err := parseRSAPublicKeyFromPEM([]byte(ac.Jwt.Certificate))
	if err != nil {
		panic("failed to parse public key")
	}
	return publicKey
}
