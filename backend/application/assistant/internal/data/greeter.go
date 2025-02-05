package data

import (
	"context"

	"backend/application/assistant/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type assistantRepo struct {
	data *Data
	log  *log.Helper
}

func (g assistantRepo) Query(ctx context.Context, question *biz.QueryQuestion) (*biz.QueryReply, error) {
	// TODO implement me
	panic("implement me")
}

func NewAssistantRepo(data *Data, logger log.Logger) biz.AssistantRepo {
	return &assistantRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
