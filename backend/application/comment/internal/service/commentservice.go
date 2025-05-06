package service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	"backend/pkg"

	"github.com/go-kratos/kratos/v2/log"

	"backend/application/comment/internal/biz"

	pb "backend/api/comment/v1"
)

type CommentService struct {
	pb.UnimplementedCommentServer

	uc  *biz.CommentUsecase
	log *log.Helper
}

func NewCommentService(uc *biz.CommentUsecase, logger log.Logger) *CommentService {
	return &CommentService{uc: uc, log: log.NewHelper(logger)}
}

func (s *CommentService) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentReply, error) {
	var userId uuid.UUID
	var err error
	userId, err = pkg.GetMetadataUesrID(ctx)
	if req.UserId != "" {
		userId, err = uuid.Parse(req.UserId)
	}

	if err != nil {
		return nil, fmt.Errorf("get user id from metadata failed: %v", err)
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("parse product id failed: %v", err)
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, fmt.Errorf("parse merchant id failed: %v", err)
	}

	reply, err := s.uc.CreateComment(ctx, &biz.CreateCommentRequest{
		ProductId:  productId,
		MerchantId: merchantId,
		UserId:     userId,
		Score:      req.Score,
		Content:    req.Content,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateCommentReply{
		IsSensitive: reply.IsSensitive,
	}, nil
}

func (s *CommentService) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("parse product id failed: %v", err)
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, fmt.Errorf("parse merchant id failed: %v", err)
	}

	result, err := s.uc.GetComments(ctx, &biz.GetCommentsRequest{
		ProductId:  productId,
		MerchantId: merchantId,
		Page:       int64(req.Page),
		PageSize:   int64(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	comments := make([]*pb.CommentType, 0, len(result.Comments))
	for _, comment := range result.Comments {
		comments = append(comments, &pb.CommentType{
			Id:         comment.Id,
			ProductId:  comment.ProductId.String(),
			MerchantId: comment.MerchantId.String(),
			UserId:     comment.UserId.String(),
			Score:      comment.Score,
			Content:    comment.Content,
			CreatedAt:  timestamppb.New(comment.CreatedAt),
			UpdatedAt:  timestamppb.New(comment.UpdatedAt),
		})
	}

	return &pb.GetCommentsResponse{
		Comments: comments,
		Total:    uint32(result.Total),
	}, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.CommentType, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user id from metadata failed: %v", err)
	}

	result, err := s.uc.UpdateComment(ctx, &biz.UpdateCommentRequest{
		Id:      req.CommentId,
		UserId:  userId,
		Content: req.Content,
		Score:   uint32(req.Score),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CommentType{
		Id:         result.Id,
		ProductId:  result.ProductId.String(),
		MerchantId: result.MerchantId.String(),
		UserId:     result.UserId.String(),
		Score:      result.Score,
		Content:    result.Content,
		CreatedAt:  timestamppb.New(result.CreatedAt),
		UpdatedAt:  timestamppb.New(result.UpdatedAt),
	}, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user id from metadata failed: %v", err)
	}

	result, err := s.uc.DeleteComment(ctx, &biz.DeleteCommentRequest{
		Id:     req.CommentId,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteCommentResponse{
		Success: result.Id > 0,
	}, nil
}
