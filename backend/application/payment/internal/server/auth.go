package server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwtV5 "github.com/golang-jwt/jwt/v5"
	"backend/application/payment/internal/conf"
)

// NewWhiteListMatcher 创建jwt白名单
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	// example: 从 api 目录生成的 pb.go 文件找到需要不需要进行鉴权的接口
	// whiteList["/admin.v1.AdminService/Login"] = struct{}{}
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
