package service

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	globalPkg "backend/pkg"

	"backend/application/admin/internal/biz"

	pb "backend/api/admin/comment/v1"
)

type AdminCommentService struct {
	pb.UnimplementedAdminCommentServer
	ac *biz.AdminCommentUsecase
}

func NewAdminCommentService(ac *biz.AdminCommentUsecase) *AdminCommentService {
	return &AdminCommentService{ac: ac}
}

func (ac *AdminCommentService) SetSensitiveWords(ctx context.Context, req *pb.SetSensitiveWordsReq) (*pb.SetSensitiveWordsReply, error) {
	adminId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	collection := make([]*biz.SensitiveWord, 0, len(req.SensitiveWords))
	for _, word := range req.SensitiveWords {
		collection = append(collection, &biz.SensitiveWord{
			CreatedBy: adminId,
			Category:  word.Category,
			Word:      word.Word,
			Level:     word.Level,
			IsActive:  word.IsActive,
		})
	}

	reply, err := ac.ac.SetSensitiveWords(ctx, &biz.SetSensitiveWordsReq{
		SensitiveWords: collection,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SetSensitiveWordsReply{
		Rows: reply.Rows,
	}, nil
}

func (ac *AdminCommentService) GetSensitiveWords(ctx context.Context, req *pb.GetSensitiveWordsReq) (*pb.GetSensitiveWordsReply, error) {
	adminId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	reply, getSensitiveWordsErr := ac.ac.GetSensitiveWords(ctx, &biz.GetSensitiveWordsReq{
		CreatedBy: &adminId,
		Page:      req.Page,
		PageSize:  req.PageSize,
		// Level:     req.Level,
		// IsActive:  req.IsActive,
		// Category:  req.Category,
	})
	if getSensitiveWordsErr != nil {
		return nil, getSensitiveWordsErr
	}

	if reply.Words == nil {
		return &pb.GetSensitiveWordsReply{
			Words: []*pb.SensitiveWord{},
		}, nil
	}

	words := make([]*pb.SensitiveWord, 0, len(reply.Words))
	for _, word := range reply.Words {
		words = append(words, &pb.SensitiveWord{
			Id:        &word.ID,
			CreatedBy: word.CreatedBy.String(),
			Category:  word.Category,
			Word:      word.Word,
			Level:     word.Level,
			IsActive:  word.IsActive,
			CreatedAt: timestamppb.New(word.CreatedAt),
			UpdatedAt: timestamppb.New(word.UpdatedAt),
		})
	}

	return &pb.GetSensitiveWordsReply{
		Words: words,
	}, nil
}

func (ac *AdminCommentService) DeleteSensitiveWord(ctx context.Context, req *pb.DeleteSensitiveWordReq) (*pb.DeleteSensitiveWordReply, error) {
	adminId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	_, err = ac.ac.DeleteSensitiveWord(ctx, &biz.DeleteSensitiveWordReq{
		ID:        uint32(req.Id),
		CreatedBy: adminId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteSensitiveWordReply{
		Success: true,
	}, nil
}

func (ac *AdminCommentService) UpdateSensitiveWord(ctx context.Context, req *pb.UpdateSensitiveWordReq) (*pb.UpdateSensitiveWordReply, error) {
	adminId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	_, err = ac.ac.UpdateSensitiveWord(ctx, &biz.UpdateSensitiveWordReq{
		ID:        req.Id,
		CreatedBy: adminId,
		Level:     int32(req.Level),
		IsActive:  req.IsActive,
		Category:  req.Category,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSensitiveWordReply{
		Success: true,
	}, nil
}
