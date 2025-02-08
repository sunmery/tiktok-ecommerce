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

func (u *userRepo) GetProfile(ctx context.Context, req *biz.GetProfileRequest) (*biz.GetProfileReply, error) {
	authHeader := req.Authorization
	if authHeader == "" {
		return nil, fmt.Errorf("authorization: (%v) header is empty", authHeader)
	}

	token := strings.Split(authHeader, "Bearer ")
	if len(token) < 2 {
		return nil, fmt.Errorf("token is not valid Bearer token : %s", authHeader)
	}

	fmt.Println("token:", token[1])
	profile, err := u.data.authClient.GetUserInfo(ctx, &authV1.GetUserInfoRequest{
		Authorization: token[1],
	})
	if err != nil {
		return nil, err
	}

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
		// Data:  claims.User,
		Data: resp,
	}, nil
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
