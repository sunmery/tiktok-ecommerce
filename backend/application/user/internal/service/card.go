package service

import (
	pb "backend/api/user/v1"
	"backend/application/user/internal/biz"
	"backend/pkg"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserService) CreateCreditCard(ctx context.Context, req *pb.CreditCard) (*emptypb.Empty, error) {
	userID, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %v", err)
	}
	fmt.Printf("req: %+v\n", req)
	_, err = s.uc.CreateCreditCard(ctx, &biz.CreditCard{
		Number:          req.Number,
		Cvv:             req.Cvv,
		Owner:           req.Owner,
		Name:            req.Name,
		Brand:           req.Brand,
		Country:         req.Country,
		Type:            req.Type,
		UserID:          userID,
		Currency:        req.Currency,
		ExpYear:  req.ExpYear,
		ExpMonth: req.ExpMonth,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *UserService) DeleteCreditCard(ctx context.Context, req *pb.DeleteCreditCardsRequest) (*emptypb.Empty, error) {
	userID, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %v", err)
	}
	_, err = s.uc.DeleteCreditCard(ctx, &biz.DeleteCreditCardRequest{
		Id:     req.Id,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *UserService) ListCreditCards(ctx context.Context, _ *emptypb.Empty) (*pb.CreditCards, error) {
	userID, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user id failed: %v", err)
	}

	cards, err := s.uc.ListCreditCards(ctx, &biz.ListCreditCardsRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	var res []*pb.CreditCard
	for _, card := range cards.CreditCards {
		res = append(res, &pb.CreditCard{
			Number:    card.Number,
			Cvv:       card.Cvv,
			Id:        card.Id,
			Owner:     card.Owner,
			Name:      card.Name,
			Brand:     card.Brand,
			Country:   card.Country,
			Type:      card.Type,
			Currency:  card.Currency,
			ExpYear:   card.ExpYear,
			ExpMonth:  card.ExpMonth,
			CreatedAt: timestamppb.New(card.CreatedAt),
		})
	}
	return &pb.CreditCards{
		CreditCards: res,
	}, nil
}
