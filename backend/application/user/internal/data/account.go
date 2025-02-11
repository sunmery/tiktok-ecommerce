package data

import (
	authV1 "backend/api/auth/v1"
	"backend/application/user/internal/biz"
	"context"
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"

	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

// GetProfile 获取用户档案
func (u *userRepo) GetProfile(ctx context.Context, req *biz.GetProfileRequest) (*biz.GetProfileReply, error) {
	// 获取Authorization头
	authHeader := req.Authorization
	if authHeader == "" {
		return nil, fmt.Errorf("authorization: (%v) header is empty", authHeader)
	}

	// 获取Authorization的值
	token := strings.Split(authHeader, "Bearer ")
	if len(token) < 2 {
		return nil, fmt.Errorf("token is not valid Bearer token : %s", authHeader)
	}

	// 调用 Auth 认证微服务的 GetUserInfo 方法
	profile, err := u.data.authClient.GetUserInfo(ctx, &authV1.GetUserInfoRequest{
		Authorization: token[1],
	})
	if err != nil {
		return nil, err
	}

	// 只返回需要的值
	resp := casdoorsdk.User{
		Owner:  profile.Data.Owner,
		Type:   profile.Data.Type,
		Name:   profile.Data.Name,
		Id:     profile.Data.Id,
		Avatar: profile.Data.Avatar,
		Email:  profile.Data.Email,
	}

	return &biz.GetProfileReply{
		State: "ok",
		Data:  resp,
	}, nil
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
