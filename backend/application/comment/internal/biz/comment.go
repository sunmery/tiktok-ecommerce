package biz

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
)

type (
	Comment struct {
		Id         int64
		UserId     uuid.UUID
		MerchantId uuid.UUID
		ProductId  uuid.UUID
		Content    string
		Score      uint32 // 评分
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	CreateCommentRequest struct {
		ProductId  uuid.UUID
		MerchantId uuid.UUID
		UserId     uuid.UUID
		Score      uint32
		Content    string
	}
	UpdateCommentRequest struct {
		Id      int64
		UserId  uuid.UUID
		Content string
		Score   uint32
	}
)

type (
	GetCommentsRequest struct {
		ProductId  uuid.UUID
		MerchantId uuid.UUID
		Page       int64
		PageSize   int64
	}

	GetCommentsResponse struct {
		Comments []*Comment
		Total    int64
	}
)

type (
	DeleteCommentRequest struct {
		Id     int64
		UserId uuid.UUID
	}
	DeleteCommentResponse struct {
		Id int64
	}
)

type CommentRepo interface {
	// CreateComment 创建评论
	CreateComment(ctx context.Context, req *CreateCommentRequest) (*Comment, error)
	// GetComments 获取评论
	GetComments(ctx context.Context, req *GetCommentsRequest) (*GetCommentsResponse, error)
	// UpdateComment 更新评论
	UpdateComment(ctx context.Context, req *UpdateCommentRequest) (*Comment, error)
	// DeleteComment 删除评论
	DeleteComment(ctx context.Context, req *DeleteCommentRequest) (*DeleteCommentResponse, error)
}

type CommentUsecase struct {
	repo CommentRepo
	log  *log.Helper
}

func NewCommentUsecase(repo CommentRepo, logger log.Logger) *CommentUsecase {
	return &CommentUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *CommentUsecase) CreateComment(ctx context.Context, req *CreateCommentRequest) (*Comment, error) {
	uc.log.WithContext(ctx).Debugf("req:%+v", req)
	return uc.repo.CreateComment(ctx, req)
}

func (uc *CommentUsecase) GetComments(ctx context.Context, req *GetCommentsRequest) (*GetCommentsResponse, error) {
	uc.log.WithContext(ctx).Debugf("req:%+v", req)
	return uc.repo.GetComments(ctx, req)
}

func (uc *CommentUsecase) UpdateComment(ctx context.Context, req *UpdateCommentRequest) (*Comment, error) {
	uc.log.WithContext(ctx).Debugf("req:%+v", req)
	return uc.repo.UpdateComment(ctx, req)
}

func (uc *CommentUsecase) DeleteComment(ctx context.Context, req *DeleteCommentRequest) (*DeleteCommentResponse, error) {
	uc.log.WithContext(ctx).Debugf("req:%+v", req)
	return uc.repo.DeleteComment(ctx, req)
}
