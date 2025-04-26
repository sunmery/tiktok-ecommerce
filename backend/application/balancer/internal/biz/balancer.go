package biz

import (
	"context"
	"encoding/json"
	"time"

	"backend/constants"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
)

type (
	CreateConsumerBalanceRequest struct {
		UserId         uuid.UUID
		Currency       constants.Currency
		InitialBalance float64
		BalancerType   string
		IsDefault      bool
		AccountDetails json.RawMessage
	}
	CreateConsumerBalanceReply struct {
		UserId    uuid.UUID
		Currency  constants.Currency
		Available float64
	}
)

type (
	CreateMerchantBalanceRequest struct {
		MerchantId     uuid.UUID
		Currency       constants.Currency
		InitialBalance float64
		BalancerType   string
		IsDefault      bool
		AccountDetails json.RawMessage
	}
	CreateMerchantBalanceReply struct {
		UserId    uuid.UUID
		Currency  constants.Currency
		Available float64
	}
)

type (
	CancelFreezeRequest struct {
		FreezeId        int64  // 冻结记录ID
		Reason          string // 取消原因
		IdempotencyKey  string // 幂等键
		ExpectedVersion int32  // 期望的用户余额版本号
	}
	CancelFreezeReply struct {
		Success    bool
		NewVersion int32
	}
)

type (
	ConfirmTransferRequest struct {
		FreezeId int64
		// merchant_id 可以从 freeze_id 关联的 order_id 推出，或者在这里显式传入
		// string merchant_id = 2;
		PaymentAccount          string
		IdempotencyKey          string
		ExpectedUserVersion     int32
		ExpectedMerchantVersion int32
		MerchantId              uuid.UUID
	}
	ConfirmTransferReply struct {
		Success            bool
		TransactionId      int64
		NewUserVersion     int32
		NewMerchantVersion int32
	}
)

type (
	FreezeBalanceRequest struct {
		ID              int64
		UserId          uuid.UUID
		OrderId         int64
		Amount          float64
		Currency        constants.Currency
		ExpiresAt       time.Time
		ExpectedVersion int32
		IdempotencyKey  string
	}
	FreezeBalanceReply struct {
		FreezeId   int64
		NewVersion int32
	}
)

type (
	GetMerchantBalanceRequest struct {
		MerchantId uuid.UUID
		Currency   constants.Currency
	}
	BalanceReply struct {
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
		MerchantId            uuid.UUID
		Amount                float64
		Currency              constants.Currency
		ExternalTransactionId int64
		PaymentMethodType     string
		PaymentAccount        string
		IdempotencyKey        string
		ExpectedVersion       int32
	}
	RechargeBalanceReply struct {
		Success       bool
		TransactionId int64
		NewVersion    int32
	}
)

type (
	WithdrawBalanceRequest struct {
		UserId          uuid.UUID
		MerchantId      uuid.UUID
		Amount          float64
		Currency        constants.Currency
		PaymentMethodId string
		IdempotencyKey  string
		ExpectedVersion int32
	}
	WithdrawBalanceReply struct {
		Success       bool
		TransactionId int64
		NewVersion    int32
	}
)

type (
	Transactions struct {
		Id                int64
		Type              constants.TransactionType
		Amount            float64
		Currency          string
		FromUserId        uuid.UUID
		ToMerchantId      uuid.UUID
		PaymentMethodType constants.PaymentMethod
		PaymentAccount    string
		PaymentExtra      json.RawMessage
		Status            constants.PaymentStatus
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}
	GetTransactionsRequest struct {
		UserId        uuid.UUID
		Currency      string
		Page          int64
		PageSize      int64
		PaymentStatus constants.PaymentStatus
	}
	GetTransactionsReply struct {
		Transactions []*Transactions
	}
)

type (
	CreateTransactionRequest struct {
		Type              constants.TransactionType
		Amount            float64
		Currency          constants.Currency
		FromUserId        uuid.UUID
		ToMerchantId      uuid.UUID
		PaymentMethodType constants.PaymentMethod
		PaymentAccount    string
		PaymentExtra      json.RawMessage
		Status            constants.PaymentStatus
		IdempotencyKey    string
		FreezeId          int64
		ConsumerVersion   int64
		MerchantVersion   int64
	}
	CreateTransactionReply struct {
		Id int64
	}
)

type BalancerRepo interface {
	// CreateConsumerBalance 创建用户余额
	CreateConsumerBalance(ctx context.Context, req *CreateConsumerBalanceRequest) (*CreateConsumerBalanceReply, error)
	// CreateMerchantBalance 创建商家余额
	CreateMerchantBalance(ctx context.Context, req *CreateMerchantBalanceRequest) (*CreateMerchantBalanceReply, error)
	// CancelFreeze 取消冻结
	CancelFreeze(ctx context.Context, req *CancelFreezeRequest) (*CancelFreezeReply, error)
	// ConfirmTransfer 确认转账（解冻并转给商家）
	ConfirmTransfer(ctx context.Context, req *ConfirmTransferRequest) (*ConfirmTransferReply, error)
	// FreezeBalance 冻结用户余额
	FreezeBalance(ctx context.Context, req *FreezeBalanceRequest) (*FreezeBalanceReply, error)
	// GetMerchantBalance 获取商家余额
	GetMerchantBalance(ctx context.Context, req *GetMerchantBalanceRequest) (*BalanceReply, error)
	// GetTransactions 获取商家交易记录
	GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsReply, error)
	// GetUserBalance 获取用户余额
	GetUserBalance(ctx context.Context, req *GetUserBalanceRequest) (*BalanceReply, error)
	// RechargeBalance 用户充值
	RechargeBalance(ctx context.Context, req *RechargeBalanceRequest) (*RechargeBalanceReply, error)
	// WithdrawBalance 用户提现
	WithdrawBalance(ctx context.Context, req *WithdrawBalanceRequest) (*WithdrawBalanceReply, error)
	// CreateTransaction 创建交易记录
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionReply, error)
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

func (cc *BalancerUsecase) CreateConsumerBalance(ctx context.Context, req *CreateConsumerBalanceRequest) (*CreateConsumerBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("CreateConsumerBalance request: %+v", req)
	return cc.repo.CreateConsumerBalance(ctx, req)
}

func (cc *BalancerUsecase) CreateMerchantBalance(ctx context.Context, req *CreateMerchantBalanceRequest) (*CreateMerchantBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("CreateMerchantBalance request: %+v", req)
	return cc.repo.CreateMerchantBalance(ctx, req)
}

func (cc *BalancerUsecase) CancelFreeze(ctx context.Context, req *CancelFreezeRequest) (*CancelFreezeReply, error) {
	cc.log.WithContext(ctx).Debugf("CancelFreeze request: %+v", req)
	return cc.repo.CancelFreeze(ctx, req)
}

func (cc *BalancerUsecase) ConfirmTransfer(ctx context.Context, req *ConfirmTransferRequest) (*ConfirmTransferReply, error) {
	cc.log.WithContext(ctx).Debugf("ConfirmTransfer request: %+v", req)
	return cc.repo.ConfirmTransfer(ctx, req)
}

func (cc *BalancerUsecase) FreezeBalance(ctx context.Context, req *FreezeBalanceRequest) (*FreezeBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("FreezeBalance request: %+v", req)
	return cc.repo.FreezeBalance(ctx, req)
}

func (cc *BalancerUsecase) GetMerchantBalance(ctx context.Context, req *GetMerchantBalanceRequest) (*BalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("GetMerchantBalance request: %+v", req)
	return cc.repo.GetMerchantBalance(ctx, req)
}

func (cc *BalancerUsecase) GetUserBalance(ctx context.Context, req *GetUserBalanceRequest) (*BalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("GetUserBalance request: %+v", req)
	return cc.repo.GetUserBalance(ctx, req)
}

func (cc *BalancerUsecase) RechargeBalance(ctx context.Context, req *RechargeBalanceRequest) (*RechargeBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("RechargeBalance request: %+v", req)
	return cc.repo.RechargeBalance(ctx, req)
}

func (cc *BalancerUsecase) WithdrawBalance(ctx context.Context, req *WithdrawBalanceRequest) (*WithdrawBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("WithdrawBalance request: %+v", req)
	return cc.repo.WithdrawBalance(ctx, req)
}

func (cc *BalancerUsecase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionReply, error) {
	cc.log.WithContext(ctx).Debugf("CreateTransaction request: %+v", req)
	return cc.repo.CreateTransaction(ctx, req)
}

func (cc *BalancerUsecase) GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsReply, error) {
	cc.log.WithContext(ctx).Debugf("GetTransactions request: %+v", req)
	return cc.repo.GetTransactions(ctx, req)
}
