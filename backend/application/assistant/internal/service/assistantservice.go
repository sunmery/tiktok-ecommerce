package service

import (
	pb "backend/api/assistant/v1"
	"backend/application/assistant/internal/biz"
	"context"
)

type AssistantServiceService struct {
	pb.UnimplementedAssistantServiceServer
	ac *biz.AssistantUseCase
}

func NewAssistantServiceService(ac *biz.AssistantUseCase) *AssistantServiceService {
	return &AssistantServiceService{ac: ac}
}

func (s *AssistantServiceService) Query(req *pb.QueryRequest, conn pb.AssistantService_QueryServer) error {
	result, err := s.ac.QueryQuestion(context.Background(), &biz.QueryQuestion{Question: req.GetQuestion()})
	if err != nil {
		return err
	}
	for {
		err := conn.Send(&pb.QueryReply{
			Message: result.Message,
		})
		if err != nil {
			return err
		}
	}
}
