package data

import (
	"backend/application/user/internal/biz"
	"backend/application/user/internal/data/models"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (u *userRepo) CreateCreditCard(ctx context.Context, req *biz.CreditCard) (*emptypb.Empty, error) {
	params:=models.InsertCreditCardParams{
		Number:   req.Number,
		Cvv:      req.Cvv,
		ExpYear:  req.ExpYear,
		ExpMonth: req.ExpMonth,
		Owner:    req.Owner,
		Name:     &req.Name,
		Type:     req.Type,
		UserID:   req.UserID,
		Brand:    req.Brand,
		Country:  req.Country,
		Currency: req.Currency,
	}
	fmt.Printf("params: %+v", params)
	err := u.data.db.InsertCreditCard(ctx, params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *userRepo) DeleteCreditCard(ctx context.Context, req *biz.DeleteCreditCardRequest) (*emptypb.Empty, error) {
	err := u.data.db.DeleteCreditCard(ctx, int32(req.Id))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (u *userRepo) ListCreditCards(ctx context.Context, req *biz.ListCreditCardsRequest) (*biz.CreditCards, error) {
	cards, err := u.data.db.ListCreditCards(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	var res []*biz.CreditCard
	for _, card := range cards {
		res = append(res, &biz.CreditCard{
			Number:          card.Number,
			Cvv:             card.Cvv,
			Id:              uint32(card.ID),
			Owner:           card.Owner,
			Name:           *card.Name,
			Brand:           card.Brand,
			Country:         card.Country,
			Type:            card.Type,
			UserID:          card.UserID,
			Currency:        card.Currency,
			ExpYear:  card.ExpYear,
			ExpMonth: card.ExpMonth,
			CreatedAt:       card.CreatedTime,
		})
	}
	return &biz.CreditCards{
		CreditCards: res,
	}, nil
}
