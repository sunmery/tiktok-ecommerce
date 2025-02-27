package service

import (
	"backend/application/assistant/internal/biz"
	"context"
	"fmt"

	pb "backend/api/assistant/v1"
)

type AssistantService struct {
	pb.UnimplementedAssistantServer

	uc *biz.AssistantUsecase
}

func NewAssistantService(uc *biz.AssistantUsecase) *AssistantService {
	return &AssistantService{uc: uc}
}

func (s *AssistantService) ProcessQuery(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	process, err := s.uc.Process(ctx, req.Question)
	if err != nil {
		return nil, err
	}
	fmt.Println("process:",process)
	return &pb.ProcessResponse{
		Result: nil,
	}, nil
}
