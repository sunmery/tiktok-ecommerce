package service

import (
	"context"

	pb "backend/api/admin/comment/v1"
)

type AdminCommentService struct {
	pb.UnimplementedAdminCommentServer
}

func NewAdminCommentService() *AdminCommentService {
	return &AdminCommentService{}
}

func (s *AdminCommentService) SetSensitiveWords(ctx context.Context, req *pb.SetSensitiveWordsReq) (*pb.SetSensitiveWordsReply, error) {
	return &pb.SetSensitiveWordsReply{}, nil
}
func (s *AdminCommentService) GetSensitiveWords(ctx context.Context, req *pb.GetSensitiveWordsReq) (*pb.GetSensitiveWordsReply, error) {
	return &pb.GetSensitiveWordsReply{}, nil
}
