package biz

import (
	"context"

	// v1 "backend/api/helloworld/v1"

	// "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type QueryQuestion struct {
	Question string
}
type QueryReply struct {
	Message string
}

type AssistantRepo interface {
	Query(context.Context, *QueryQuestion) (*QueryReply, error)
}

// AssistantUseCase is a Assistant usecase.
type AssistantUseCase struct {
	repo AssistantRepo
	log  *log.Helper
}

func NewAssistantUseCase(repo AssistantRepo, logger log.Logger) *AssistantUseCase {
	return &AssistantUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *AssistantUseCase) QueryQuestion(ctx context.Context, req *QueryQuestion) (*QueryReply, error) {
	uc.log.WithContext(ctx).Infof("Query: %v", req.Question)
	return uc.repo.Query(ctx, req)
}
