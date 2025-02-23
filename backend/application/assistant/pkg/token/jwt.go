package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtV5 "github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Type  string `json:"type"`

	jwtV5.Claims `json:"claims"`
}

func ExtractPayload(ctx context.Context) (*Payload, error) {
	user, ok := jwt.FromContext(ctx)
	if !ok {
		return nil, errors.New("invalid token")
	}

	// 检查是否是 Payload 类型
	if payload, ok := user.(*Payload); ok {
		return payload, nil
	}

	// 如果是 MapClaims，尝试转换
	claims, ok := user.(jwtV5.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	payload := &Payload{
		ID:    fmt.Sprintf("%v", claims["id"]),
		Name:  fmt.Sprintf("%v", claims["name"]),
		Owner: fmt.Sprintf("%v", claims["owner"]),
		Type:  fmt.Sprintf("%v", claims["type"]),
	}

	return payload, nil
}
