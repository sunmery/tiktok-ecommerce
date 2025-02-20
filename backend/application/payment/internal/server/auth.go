package server

import (
	"backend/application/payment/internal/conf"
	"context"
	"crypto/rsa"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwtV5 "github.com/golang-jwt/jwt/v5"
)

// NewWhiteListMatcher 创建 jwt 白名单
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	// 从 api 目录生成的 pb.go 文件找到需要不需要进行鉴权的接口
	// example:
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
		log.Errorf("failed to parse RSA public key: %w", err)
	}
	return publicKey, nil
}

func InitJwtKey(ac *conf.Auth) *rsa.PublicKey {
	publicKey, err := parseRSAPublicKeyFromPEM([]byte(ac.Jwt.Certificate))
	if err != nil {
		log.Errorf("failed to parse RSA public key from PEM: %v", err)
	}
	return publicKey
}
