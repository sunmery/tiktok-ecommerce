package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/go-kratos/kratos/v2/log"
	"backend/application/product/internal/biz"
	"strings"
)

func (u *userRepo) Signin(ctx context.Context, req *biz.SigninRequest) (*biz.SigninReply, error) {
	code := req.Code
	state := req.State
	token, err := u.data.cs.GetOAuthToken(code, state)
	if err != nil {
		fmt.Println("GetOAuthToken() error", err)
		return nil, errors.New("GetOAuthToken() error:" + err.Error())
	}

	fmt.Println("GetOAuthToken() token", token)
	return &biz.SigninReply{
		State: "ok",
		Data:  token.AccessToken,
	}, nil
}

func (u *userRepo) GetUserInfo(ctx context.Context, req *biz.GetUserInfoRequest) (*biz.GetUserInfoReply, error) {
	authHeader := req.Authorization
	if authHeader == "" {
		return nil, fmt.Errorf("authorization: (%v) header is empty", authHeader)
	}

	token := strings.Split(authHeader, "Bearer ")
	if len(token) < 2 {
		return nil, fmt.Errorf("token is not valid Bearer token : %s", authHeader)
	}

	claims, err := u.data.cs.ParseJwtToken(token[1])
	if err != nil {
		return nil, fmt.Errorf("ParseJwtToken() error")
	}

	resp := casdoorsdk.User{
		Owner:  claims.Owner,
		Type:   claims.Type,
		Name:   claims.Name,
		Id:     claims.Id,
		Avatar: claims.Avatar,
		Email:  claims.Email,
	}

	return &biz.GetUserInfoReply{
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
