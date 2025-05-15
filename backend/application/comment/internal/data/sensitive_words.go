package data

import (
	"context"

	"backend/application/comment/internal/biz"
	"backend/application/comment/internal/data/models"

	"github.com/go-kratos/kratos/v2/log"
)

type sensitiveWordRepo struct {
	data *Data
	log  *log.Helper
}

func (r *sensitiveWordRepo) SetSensitiveWords(ctx context.Context, req *biz.SetSensitiveWordsReq) (*biz.SetSensitiveWordsReply, error) {
	var rows int64
	for _, word := range req.SensitiveWords {
		affectedRows, err := r.data.db.SetSensitiveWords(ctx, models.SetSensitiveWordsParams{
			CreatedBy: word.CreatedBy,
			Category:  word.Category,
			Word:      word.Word,
			Level:     word.Level,
			IsActive:  word.IsActive,
		})
		if err != nil {
			return nil, err
		}
		rows += affectedRows
	}

	return &biz.SetSensitiveWordsReply{
		Rows: uint32(rows),
	}, nil
}

func (r *sensitiveWordRepo) GetSensitiveWords(ctx context.Context, req *biz.GetSensitiveWordsReq) (*biz.GetSensitiveWordsReply, error) {
	// 计算分页参数
	page := (req.Page - 1) * req.PageSize
	if page < 0 {
		page = 0
	}

	// 获取敏感词列表
	words, err := r.data.db.GetSensitiveWords(ctx, models.GetSensitiveWordsParams{
		CreatedBy: req.CreatedBy,
		PageSize:  int64(req.PageSize),
		Page:      int64(page),
	})
	if err != nil {
		return nil, err
	}

	// 转换为业务层数据结构
	result := &biz.GetSensitiveWordsReply{
		Words: make([]*biz.SensitiveWord, 0, len(words)),
	}

	for _, word := range words {
		result.Words = append(result.Words, &biz.SensitiveWord{
			ID:        word.ID,
			CreatedBy: word.CreatedBy,
			Category:  word.Category,
			Word:      word.Word,
			Level:     word.Level,
			IsActive:  word.IsActive,
			CreatedAt: word.CreatedAt.Time,
			UpdatedAt: word.UpdatedAt.Time,
		})
	}

	return result, nil
}

func (r *sensitiveWordRepo) DeleteSensitiveWord(ctx context.Context, req *biz.DeleteSensitiveWordReq) error {
	return r.data.db.DeleteSensitiveWord(ctx, models.DeleteSensitiveWordParams{
		ID:        req.ID,
		CreatedBy: req.CreatedBy,
	})
}

func NewSensitiveWordRepo(data *Data, logger log.Logger) biz.SensitiveWordRepo {
	return &sensitiveWordRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
