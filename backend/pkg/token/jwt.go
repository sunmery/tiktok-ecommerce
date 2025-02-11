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
		fmt.Println("❌ 无法从 Context 提取 JWT，可能请求没有带 Authorization 头")
		return nil, errors.New("invalid token")
	}

	fmt.Printf("✅ 提取到的 Token: %+v\n", user)

	// 检查是否是 Payload 类型
	if payload, ok := user.(*Payload); ok {
		return payload, nil
	}

	// 如果是 MapClaims，尝试转换
	claims, ok := user.(jwtV5.MapClaims)
	if !ok {
		fmt.Println("❌ JWT 解析失败，格式不匹配")
		return nil, errors.New("invalid claims type")
	}

	payload := &Payload{
		ID:    fmt.Sprintf("%v", claims["id"]),
		Name:  fmt.Sprintf("%v", claims["name"]),
		Owner: fmt.Sprintf("%v", claims["owner"]),
		Type:  fmt.Sprintf("%v", claims["type"]),
	}
	fmt.Printf("✅ Token 解析成功: %+v\n", payload)
	return payload, nil
}
