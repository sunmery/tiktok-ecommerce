package data

import (
	"backend/application/user/internal/biz"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

// GetProfile 获取用户档案
func (u *userRepo) GetProfile(ctx context.Context, req *biz.GetProfileRequest) (*biz.GetProfileReply, error) {
	user, err := u.data.cs.GetUser(req.UserId.String())
	if err != nil {
		return nil, err
	}

	// 用户是否被注销
	if user.IsDeleted {
		return nil, fmt.Errorf(fmt.Sprintf("user %s is deleted", user.Name))
	}

	// 组装数据
	return &biz.GetProfileReply{
		Owner:  user.Owner,
		Name:   user.Name,
		Id:     req.UserId,
		Avatar: user.Avatar,
		Email:  user.Email,
		Role:   user.Roles[0].Name,
	}, nil
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
