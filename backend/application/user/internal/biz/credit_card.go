package biz

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type CreditCard struct {
	Number    string    `json:"number,omitempty"`
	Cvv       string    `json:"cvv,omitempty"`
	Id        uint32    `json:"id,omitempty"`
	Owner     string    `json:"owner,omitempty"`
	Name      string    `json:"name,omitempty"`
	Brand     string    `json:"brand,omitempty"`
	Country   string    `json:"country,omitempty"`
	Type      string    `json:"type,omitempty"`
	UserID    uuid.UUID `json:"userId,omitempty"`
	Currency  string    `json:"currency,omitempty"`
	ExpYear   string    `json:"expYear,omitempty"`
	ExpMonth  string    `json:"expMonth,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
type CreditCards struct {
	CreditCards []*CreditCard
}
type ListCreditCardsRequest struct {
	UserID uuid.UUID
}
type CreditCardRequest struct {
	UserID uuid.UUID
}

type DeleteCreditCardRequest struct {
	UserId uuid.UUID
	Id     uint32
}

// CreateCreditCard 创建银行卡
func (cc *UserUsecase) CreateCreditCard(ctx context.Context, req *CreditCard) (*emptypb.Empty, error) {
	cc.log.WithContext(ctx).Infof("CreateCreditCard: %v", req)
	return cc.repo.CreateCreditCard(ctx, req)
}

// DeleteCreditCard 删除银行卡
func (cc *UserUsecase) DeleteCreditCard(ctx context.Context, req *DeleteCreditCardRequest) (*emptypb.Empty, error) {
	return cc.repo.DeleteCreditCard(ctx, req)
}

// ListCreditCards 获取银行卡列表
func (cc *UserUsecase) ListCreditCards(ctx context.Context, req *ListCreditCardsRequest) (*CreditCards, error) {
	return cc.repo.ListCreditCards(ctx, req)
}
