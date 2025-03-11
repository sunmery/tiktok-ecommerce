package data

import (
	"context"
	"errors"

	"backend/application/auth/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

func (u *authRepo) Signin(ctx context.Context, req *biz.SigninRequest) (*biz.SigninReply, error) {
	code := req.Code
	state := req.State
	token, err := u.data.cs.GetOAuthToken(code, state)
	if err != nil {
		u.log.Errorf("GetOAuthToken() error: %s", err)
		return nil, errors.New("GetOAuthToken() error:" + err.Error())
	}

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
