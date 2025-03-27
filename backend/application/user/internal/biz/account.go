package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// 用户档案
type (
	GetProfileRequest struct {
		Owner  string
		UserId uuid.UUID
	}
	GetProfileReply struct {
		Id                uuid.UUID // 用户唯一标识
		Role              string
		IsDeleted         bool
		Owner             string // 用户所在的组织
		SignupApplication string // 用户所在的应用
		Name              string // 用户姓名
		Email             string // 用户邮箱
		Avatar            string // 用户头像 string
		CreatedTime       string
		UpdatedTime       string
		DeletedTime       time.Time
		DisplayName       string // 显示的用户名(可作为昵称)
	}
)

func (cc *UserUsecase) GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error) {
	cc.log.WithContext(ctx).Infof("GetProfile: %+v", req)
	return cc.repo.GetProfile(ctx, req)
}
