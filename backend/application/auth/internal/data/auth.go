package data

import (
	"backend/application/auth/internal/biz"
	"context"
	"errors"
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/go-kratos/kratos/v2/log"
)

func (u *authRepo) Signin(ctx context.Context, req *biz.SigninRequest) (*biz.SigninReply, error) {
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

func (u *authRepo) GetUserInfo(ctx context.Context, req *biz.GetUserInfoRequest) (*biz.GetUserInfoReply, error) {
	claims, err := u.data.cs.ParseJwtToken(req.Authorization)
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

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
