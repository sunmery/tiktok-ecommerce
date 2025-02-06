package service

import (
	"backend/api/assistant/v1"
	"backend/application/assistant/internal/biz"
	"context"
)

type AssistantServiceService struct {
	v1.UnimplementedAssistantServiceServer
	ac *biz.AssistantUseCase
}

func NewAssistantServiceService(ac *biz.AssistantUseCase) *AssistantServiceService {
	return &AssistantServiceService{ac: ac}
}

func (s *AssistantServiceService) Query(ctx context.Context, req *v1.QueryRequest) (*v1.QueryReply, error) {
	result, err := s.ac.QueryQuestion(context.Background(), &biz.QueryQuestion{Question: req.GetQuestion()})
	if err != nil {
		return nil, err
	}
	return &v1.QueryReply{
		Message: result.Message,
	}, nil
}
