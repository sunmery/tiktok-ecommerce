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

type (
	GetUsersRequest struct {
		AdminId uuid.UUID
	}

	GetUsersReply struct {
		Users []*GetProfileReply
	}
)

type (
	DeleteUserRequest struct {
		Owner  string
		UserId uuid.UUID
		Name   string
	}
	DeleteUserReply struct {
		Code   int32
		Status string
	}
)

type (
	UpdateUserRequest struct {
		Owner             string
		UserId            uuid.UUID
		Name              string
		Email             string
		Avatar            string
		DisplayName       string
		SignupApplication string
	}
	UpdateUserReply struct {
		Code   int32
		Status string
	}
)

func (cc *UserUsecase) GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error) {
	cc.log.WithContext(ctx).Debugf("GetProfile: %+v", req)
	return cc.repo.GetProfile(ctx, req)
}

func (cc *UserUsecase) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersReply, error) {
	cc.log.WithContext(ctx).Debugf("GetUsers: %+v", req)
	return cc.repo.GetUsers(ctx, req)
}

func (cc *UserUsecase) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserReply, error) {
	cc.log.WithContext(ctx).Debugf("DeleteUser: %+v", req)
	return cc.repo.DeleteUser(ctx, req)
}

func (cc *UserUsecase) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserReply, error) {
	cc.log.WithContext(ctx).Debugf("UpdateUser: %+v", req)
	return cc.repo.UpdateUser(ctx, req)
}
