package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CreditCard struct {
	Number    string
	Cvv       string
	Id        uint32
	Owner     string
	Name      string
	Brand     string
	Country   string
	Type      string
	UserID    uuid.UUID
	Currency  string
	ExpYear   string
	ExpMonth  string
	CreatedAt time.Time
}
type CreditCards struct {
	CreditCards []*CreditCard
}

type GetCreditCardRequest struct {
	UserID uuid.UUID
	Id     uint32
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

// GetCreditCard 根据ID获取单个银行卡
func (cc *UserUsecase) GetCreditCard(ctx context.Context, req *GetCreditCardRequest) (*CreditCard, error) {
	return cc.repo.GetCreditCard(ctx, req)
}
