package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type (
	SensitiveWord struct {
		ID        int32
		CreatedBy uuid.UUID
		Category  string
		Word      string
		Level     int32
		IsActive  bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	SetSensitiveWordsReq struct {
		SensitiveWords []*SensitiveWord
	}

	SetSensitiveWordsReply struct {
		Rows uint32
	}

	GetSensitiveWordsReq struct {
		CreatedBy *uuid.UUID
		Page      uint32
		PageSize  uint32
	}

	GetSensitiveWordsReply struct {
		Words []*SensitiveWord
	}

	DeleteSensitiveWordReq struct {
		ID        int32
		CreatedBy uuid.UUID
	}
)

type SensitiveWordRepo interface {
	// SetSensitiveWords 设置敏感词
	SetSensitiveWords(ctx context.Context, req *SetSensitiveWordsReq) (*SetSensitiveWordsReply, error)
	// GetSensitiveWords 获取敏感词列表
	GetSensitiveWords(ctx context.Context, req *GetSensitiveWordsReq) (*GetSensitiveWordsReply, error)
	// DeleteSensitiveWord 删除敏感词
	DeleteSensitiveWord(ctx context.Context, req *DeleteSensitiveWordReq) error
}

type SensitiveWordUsecase struct {
	repo SensitiveWordRepo
	log  *log.Helper
}

func NewSensitiveWordUsecase(repo SensitiveWordRepo, logger log.Logger) *SensitiveWordUsecase {
	return &SensitiveWordUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *SensitiveWordUsecase) SetSensitiveWords(ctx context.Context, req *SetSensitiveWordsReq) (*SetSensitiveWordsReply, error) {
	uc.log.WithContext(ctx).Debugf("req:%+v", req)
	return uc.repo.SetSensitiveWords(ctx, req)
}

func (uc *SensitiveWordUsecase) GetSensitiveWords(ctx context.Context, req *GetSensitiveWordsReq) (*GetSensitiveWordsReply, error) {
	uc.log.WithContext(ctx).Debugf("req:%+v", req)
	return uc.repo.GetSensitiveWords(ctx, req)
}

func (uc *SensitiveWordUsecase) DeleteSensitiveWord(ctx context.Context, req *DeleteSensitiveWordReq) error {
	uc.log.WithContext(ctx).Debugf("req:%+v", req)
	return uc.repo.DeleteSensitiveWord(ctx, req)
}
