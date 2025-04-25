package biz

import (
	"context"
	"time"

	"backend/constants"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
)

type (
	CancelFreezeRequest struct {
		FreezeId        string // 冻结记录ID (BIGINT as string)
		Reason          string // 取消原因
		IdempotencyKey  string // 幂等键
		ExpectedVersion int32  // 期望的用户余额版本号
	}
	CancelFreezeResponse struct {
		Success    bool
		NewVersion int32
	}
)

type (
	ConfirmTransferRequest struct {
		FreezeId string
		// merchant_id 可以从 freeze_id 关联的 order_id 推出，或者在这里显式传入
		// string merchant_id = 2;
		IdempotencyKey          string
		ExpectedUserVersion     int32
		ExpectedMerchantVersion int32
	}
	ConfirmTransferResponse struct {
		Success            bool
		TransactionId      string
		NewUserVersion     int32
		NewMerchantVersion int32
	}
)

type (
	FreezeBalanceRequest struct {
		UserId          uuid.UUID
		OrderId         string
		Amount          float64
		Currency        constants.Currency
		ExpiresAt       time.Time
		ExpectedVersion int32
		IdempotencyKey  string
	}
	FreezeBalanceResponse struct {
		FreezeId   string
		NewVersion int32
	}
)

type (
	GetMerchantBalanceRequest struct {
		MerchantId uuid.UUID
		Currency   constants.Currency
	}
	BalanceResponse struct {
		Available float64
		Frozen    float64
		Currency  constants.Currency
		Version   int32
	}
)

type (
	GetUserBalanceRequest struct {
		UserId   uuid.UUID
		Currency constants.Currency
	}
)

type (
	RechargeBalanceRequest struct {
		UserId                uuid.UUID
		Amount                float64
		Currency              constants.Currency
		ExternalTransactionId string
		PaymentMethodType     string
		PaymentAccount        string
		IdempotencyKey        string
		ExpectedVersion       int32
	}
	RechargeBalanceResponse struct {
		Success       bool
		TransactionId string
		NewVersion    int32
	}
)

type (
	WithdrawBalanceRequest struct {
		UserId          uuid.UUID
		Amount          float64
		Currency        constants.Currency
		PaymentMethodId string
		IdempotencyKey  string
		ExpectedVersion int32
	}
	WithdrawBalanceResponse struct {
		Success       bool
		TransactionId string
		NewVersion    int32
	}
)

type BalancerRepo interface {
	// CancelFreeze 取消冻结
	CancelFreeze(ctx context.Context, req *CancelFreezeRequest) (*CancelFreezeResponse, error)
	// ConfirmTransfer 确认转账（解冻并转给商家）
	ConfirmTransfer(ctx context.Context, req *ConfirmTransferRequest) (*ConfirmTransferResponse, error)
	// FreezeBalance 冻结用户余额
	FreezeBalance(ctx context.Context, req *FreezeBalanceRequest) (*FreezeBalanceResponse, error)
	// GetMerchantBalance 获取商家余额
	GetMerchantBalance(ctx context.Context, req *GetMerchantBalanceRequest) (*BalanceResponse, error)
	// GetUserBalance 获取用户余额
	GetUserBalance(ctx context.Context, req *GetUserBalanceRequest) (*BalanceResponse, error)
	// RechargeBalance 用户充值
	RechargeBalance(ctx context.Context, req *RechargeBalanceRequest) (*RechargeBalanceResponse, error)
	// WithdrawBalance 用户提现
	WithdrawBalance(ctx context.Context, req *WithdrawBalanceRequest) (*WithdrawBalanceResponse, error)
}

type BalancerUsecase struct {
	repo BalancerRepo
	log  *log.Helper
}

func NewBalancerUsecase(repo BalancerRepo, logger log.Logger) *BalancerUsecase {
	return &BalancerUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (cc *BalancerUsecase) CancelFreeze(ctx context.Context, req *CancelFreezeRequest) (*CancelFreezeResponse, error) {
	cc.log.WithContext(ctx).Debugf("CancelFreeze request: %+v", req)
	return cc.repo.CancelFreeze(ctx, req)
}

func (cc *BalancerUsecase) ConfirmTransfer(ctx context.Context, req *ConfirmTransferRequest) (*ConfirmTransferResponse, error) {
	cc.log.WithContext(ctx).Debugf("ConfirmTransfer request: %+v", req)
	return cc.repo.ConfirmTransfer(ctx, req)
}

func (cc *BalancerUsecase) FreezeBalance(ctx context.Context, req *FreezeBalanceRequest) (*FreezeBalanceResponse, error) {
	cc.log.WithContext(ctx).Debugf("FreezeBalance request: %+v", req)
	return cc.repo.FreezeBalance(ctx, req)
}

func (cc *BalancerUsecase) GetMerchantBalance(ctx context.Context, req *GetMerchantBalanceRequest) (*BalanceResponse, error) {
	cc.log.WithContext(ctx).Debugf("GetMerchantBalance request: %+v", req)
	return cc.repo.GetMerchantBalance(ctx, req)
}

func (cc *BalancerUsecase) GetUserBalance(ctx context.Context, req *GetUserBalanceRequest) (*BalanceResponse, error) {
	cc.log.WithContext(ctx).Debugf("GetUserBalance request: %+v", req)
	return cc.repo.GetUserBalance(ctx, req)
}

func (cc *BalancerUsecase) RechargeBalance(ctx context.Context, req *RechargeBalanceRequest) (*RechargeBalanceResponse, error) {
	cc.log.WithContext(ctx).Debugf("RechargeBalance request: %+v", req)
	return cc.repo.RechargeBalance(ctx, req)
}

func (cc *BalancerUsecase) WithdrawBalance(ctx context.Context, req *WithdrawBalanceRequest) (*WithdrawBalanceResponse, error) {
	cc.log.WithContext(ctx).Debugf("WithdrawBalance request: %+v", req)
	return cc.repo.WithdrawBalance(ctx, req)
}
