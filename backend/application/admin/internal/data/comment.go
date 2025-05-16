package data

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"backend/application/admin/internal/data/models"
	"github.com/jackc/pgx/v5"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	"backend/application/admin/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type AdminCommentRepo struct {
	data *Data
	log  *log.Helper
}

func NewAdminCommentRepo(data *Data, logger log.Logger) biz.AdminCommentRepo {
	return &AdminCommentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (a AdminCommentRepo) SetSensitiveWords(ctx context.Context, req *biz.SetSensitiveWordsReq) (*biz.SetSensitiveWordsReply, error) {
	rows, err := a.data.db.CreateBulkSensitiveWords(ctx, convertToBulkParams(req))
	if err != nil {
		return nil, kerrors.InternalServer("SET_SENSITIVE_WORDS_INTERNAL_SERVER", "set sensitive words failed")
	}

	return &biz.SetSensitiveWordsReply{
		Rows: uint32(rows),
	}, nil
}

func (a AdminCommentRepo) GetSensitiveWords(ctx context.Context, req *biz.GetSensitiveWordsReq) (*biz.GetSensitiveWordsReply, error) {
	page := (req.Page - 1) * req.PageSize

	// 使用准备好的参数调用数据库查询
	sensitiveWords, err := a.data.db.GetSensitiveWords(ctx, models.GetSensitiveWordsParams{
		Page:     int32(page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.log.WithContext(ctx).Infof("no sensitive words found for request: %+v", req)
			return &biz.GetSensitiveWordsReply{
				Words: []*biz.SensitiveWord{},
			}, nil
		}
		a.log.WithContext(ctx).Errorf("a.data.db.GetSensitiveWords failed with params %+v", err)
		return nil, kerrors.InternalServer("GET_SENSITIVE_WORDS_INTERNAL_SERVER", "get sensitive words failed")
	}

	// 保持结果转换逻辑不变
	words := make([]*biz.SensitiveWord, 0, len(sensitiveWords))
	for _, word := range sensitiveWords {
		words = append(words, &biz.SensitiveWord{
			ID:        word.ID,
			CreatedBy: word.CreatedBy,
			Category:  word.Category,
			Word:      word.Word,
			Level:     word.Level,
			IsActive:  word.IsActive,
			CreatedAt: word.CreatedAt,
			UpdatedAt: word.UpdatedAt,
		})
	}

	return &biz.GetSensitiveWordsReply{
		Words: words,
	}, nil
}

func (a AdminCommentRepo) UpdateSensitiveWord(ctx context.Context, req *biz.UpdateSensitiveWordReq) (*biz.UpdateSensitiveWordReply, error) {
	// 使用准备好的参数调用数据库查询
	_, err := a.data.db.UpdateSensitiveWord(ctx, models.UpdateSensitiveWordParams{
		Category:  req.Category,
		CreatedBy: req.CreatedBy,
		Word:      req.Word,
		Level:     req.Level,
		IsActive:  req.IsActive,
		ID:        int32(req.ID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.log.WithContext(ctx).Infof("no sensitive words found for request: %+v", req)
			return nil, nil
		}
		a.log.WithContext(ctx).Errorf("a.data.db.UpdateSensitiveWord failed with params %+v", err)
		return nil, kerrors.InternalServer("UPDATE_SENSITIVE_WORDS_INTERNAL_SERVER", "get sensitive words failed")
	}

	return &biz.UpdateSensitiveWordReply{}, nil
}

func (a AdminCommentRepo) DeleteSensitiveWord(ctx context.Context, req *biz.DeleteSensitiveWordReq) (*biz.DeleteSensitiveWordReply, error) {
	// 使用准备好的参数调用数据库查询
	_, err := a.data.db.DeleteSensitiveWords(ctx, int32(req.ID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.log.WithContext(ctx).Infof("no sensitive words found for request: %+v", req)
			return nil, nil
		}
		a.log.WithContext(ctx).Errorf("a.data.db.UpdateSensitiveWord failed with params %+v", err)
		return nil, kerrors.InternalServer("UPDATE_SENSITIVE_WORDS_INTERNAL_SERVER", "get sensitive words failed")
	}

	return &biz.DeleteSensitiveWordReply{}, nil
}

func convertToBulkParams(req *biz.SetSensitiveWordsReq) models.CreateBulkSensitiveWordsParams {
	params := models.CreateBulkSensitiveWordsParams{
		CreatedBy:  make([]uuid.UUID, 0, len(req.SensitiveWords)),
		Categories: make([]string, 0, len(req.SensitiveWords)),
		Words:      make([]string, 0, len(req.SensitiveWords)),
		Level:      make([]int32, 0, len(req.SensitiveWords)),
		IsActive:   make([]bool, 0, len(req.SensitiveWords)),
	}

	for _, word := range req.SensitiveWords {
		params.Categories = append(params.Categories, word.Category)
		params.Words = append(params.Words, word.Word)
		params.Level = append(params.Level, word.Level)
		params.IsActive = append(params.IsActive, word.IsActive)
		params.CreatedBy = append(params.CreatedBy, word.CreatedBy)
	}

	return params
}
