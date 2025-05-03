package biz

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
)

type SensitiveWord struct {
	ID        int32
	CreatedBy uuid.UUID
	Category  string
	Word      string
	Level     int32
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
type (
	SetSensitiveWordsReq struct {
		SensitiveWords []*SensitiveWord
	}
	SetSensitiveWordsReply struct {
		Rows uint32
	}
)

type (
	GetSensitiveWordsReq struct {
		Page      uint32
		PageSize  uint32
		CreatedBy *uuid.UUID
		Level     *int32
		IsActive  *bool
		Category  *string
	}
	GetSensitiveWordsReply struct {
		Words []*SensitiveWord
	}
)

type AdminCommentRepo interface {
	SetSensitiveWords(ctx context.Context, req *SetSensitiveWordsReq) (*SetSensitiveWordsReply, error)
	GetSensitiveWords(ctx context.Context, req *GetSensitiveWordsReq) (*GetSensitiveWordsReply, error)
}

type AdminCommentUsecase struct {
	repo AdminCommentRepo
	log  *log.Helper
}

func NewAdminCommentUsecase(repo AdminCommentRepo, logger log.Logger) *AdminCommentUsecase {
	return &AdminCommentUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (cc *AdminCommentUsecase) SetSensitiveWords(ctx context.Context, req *SetSensitiveWordsReq) (*SetSensitiveWordsReply, error) {
	cc.log.WithContext(ctx).Debugf("SetSensitiveWords request: %+v", req)
	return cc.repo.SetSensitiveWords(ctx, req)
}

func (cc *AdminCommentUsecase) GetSensitiveWords(ctx context.Context, req *GetSensitiveWordsReq) (*GetSensitiveWordsReply, error) {
	cc.log.WithContext(ctx).Debugf("GetSensitiveWords request: %+v", req)
	return cc.repo.GetSensitiveWords(ctx, req)
}
