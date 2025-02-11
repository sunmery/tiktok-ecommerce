package biz

import (
	"context"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

type GetProfileRequest struct {
	Authorization string
}

type GetProfileReply struct {
	State string          `json:"state,omitempty"`
	Data  casdoorsdk.User `json:"data"`
}

func (cc *UserUsecase) GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error) {
	cc.log.WithContext(ctx).Infof("GetProfile: %+v", req)
	return cc.repo.GetProfile(ctx, req)
}
