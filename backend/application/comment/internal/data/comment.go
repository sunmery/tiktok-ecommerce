package data

import (
	"context"

	"backend/application/comment/internal/pkg/id"

	"backend/pkg/types"

	"backend/application/comment/internal/data/models"

	"backend/application/comment/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type commentRepo struct {
	data *Data
	log  *log.Helper
}

func (c commentRepo) CreateComment(ctx context.Context, req *biz.CreateCommentRequest) (*biz.CreateCommentReply, error) {
	snowflakeId := id.SnowflakeID()
	productId := types.ToPgUUID(req.ProductId)
	merchantId := types.ToPgUUID(req.MerchantId)
	userId := types.ToPgUUID(req.UserId)
	score := int32(req.Score)
	content := req.Content
	isSensitive, err := c.data.db.CreateComment(ctx, models.CreateCommentParams{
		ID:         &snowflakeId,
		ProductID:  productId,
		MerchantID: merchantId,
		UserID:     userId,
		Score:      &score,
		Content:    &content,
	})
	if err != nil {
		return nil, err
	}

	return &biz.CreateCommentReply{
		IsSensitive: isSensitive,
	}, nil
}

func (c commentRepo) GetComments(ctx context.Context, req *biz.GetCommentsRequest) (*biz.GetCommentsResponse, error) {
	// 计算分页参数
	page := (req.Page - 1) * req.PageSize
	if page < 0 {
		page = 0
	}

	// 获取评论列表
	comments, err := c.data.db.GetCommentsByProduct(ctx, models.GetCommentsByProductParams{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
		PageSize:   req.PageSize,
		Page:       page,
	})
	if err != nil {
		return nil, err
	}

	// 获取评论总数
	total, err := c.data.db.GetCommentCount(ctx, models.GetCommentCountParams{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
	})
	if err != nil {
		return nil, err
	}

	// 转换为业务层数据结构
	result := &biz.GetCommentsResponse{
		Comments: make([]*biz.Comment, 0, len(comments)),
		Total:    total,
	}

	for _, comment := range comments {
		// 将int32转换为UUID
		productUUID := uuid.New()
		merchantUUID := uuid.New()
		userUUID := uuid.New()

		result.Comments = append(result.Comments, &biz.Comment{
			Id:         comment.ID,
			ProductId:  productUUID,
			MerchantId: merchantUUID,
			UserId:     userUUID,
			Content:    comment.Content,
			Score:      uint32(comment.Score),
			CreatedAt:  comment.CreatedAt.Time,
			UpdatedAt:  comment.UpdatedAt.Time,
		})
	}

	return result, nil
}

func (c commentRepo) UpdateComment(ctx context.Context, req *biz.UpdateCommentRequest) (*biz.Comment, error) {
	// 更新评论
	comment, err := c.data.db.UpdateComment(ctx, models.UpdateCommentParams{
		ID:      req.Id,
		UserID:  req.UserId,
		Content: req.Content,
		Score:   int32(req.Score),
	})
	if err != nil {
		return nil, err
	}

	// 转换为业务层数据结构
	return &biz.Comment{
		Id:         comment.ID,
		ProductId:  comment.ProductID,
		MerchantId: comment.MerchantID,
		UserId:     comment.UserID,
		Content:    comment.Content,
		Score:      uint32(comment.Score),
		CreatedAt:  comment.CreatedAt.Time,
		UpdatedAt:  comment.UpdatedAt.Time,
	}, nil
}

func (c commentRepo) DeleteComment(ctx context.Context, req *biz.DeleteCommentRequest) (*biz.DeleteCommentResponse, error) {
	// 删除评论
	err := c.data.db.DeleteComment(ctx, models.DeleteCommentParams{
		ID:     req.Id,
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	// 返回删除结果
	return &biz.DeleteCommentResponse{
		Id: req.Id,
	}, nil
}

func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
