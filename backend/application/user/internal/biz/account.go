package biz

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// 用户档案
type (
	GetProfileRequest struct {
		Owner string
		UserId    uuid.UUID
	}
	GetProfileReply struct {
		Owner  string
		Name   string
		DisplayName string
		Id     uuid.UUID
		Avatar string
		Email  string
		// Roles: []string
		Role string
		IsDeleted bool
		CreatedTime time.Time
		UpdatedTime time.Time
	}
)

func (cc *UserUsecase) GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error) {
	cc.log.WithContext(ctx).Infof("GetProfile: %+v", req)
	return cc.repo.GetProfile(ctx, req)
}
