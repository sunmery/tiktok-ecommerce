package data

import (
	"backend/application/auth/internal/biz"
	"context"
	"errors"
	"fmt"

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

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
