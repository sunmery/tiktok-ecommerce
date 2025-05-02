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
		BalanceType    string
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
		BalanceType    string
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
	Transaction struct {
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
		UserId               uuid.UUID
		Currency             string
		Page                 int64
		PageSize             int64
		PaymentStatus        constants.PaymentStatus
		TransactionsUserType constants.TransactionsUserType
	}
	GetTransactionsReply struct {
		Transactions []*Transaction
	}
)

type (
	GetMerchantVersionRequest struct {
		MerchantIds []uuid.UUID
	}
	GetMerchantVersionReply struct {
		Versions    []int64
		MerchantIds []uuid.UUID
	}
)

type (
	GetMerchantVersionByIDRequest struct {
		Version int64
	}
	GetMerchantVersionByIDReply struct {
		MerchantId uuid.UUID
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

type (
	RechargeMerchantBalanceRequest struct {
		MerchantId        uuid.UUID
		Amount            float64
		Currency          constants.Currency
		PaymentMethodType string
		PaymentMethod     constants.PaymentMethod
		PaymentAccount    string
		PaymentExtra      json.RawMessage
		IdempotencyKey    string
		ExpectedVersion   int32
	}
	RechargeMerchantBalanceReply struct {
		TransactionId int64
		NewVersion    int32
	}
)

type BalanceRepo interface {
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
	// RechargeMerchantBalance 商家余额充值
	RechargeMerchantBalance(ctx context.Context, req *RechargeMerchantBalanceRequest) (*RechargeMerchantBalanceReply, error)
	// WithdrawBalance 用户提现
	WithdrawBalance(ctx context.Context, req *WithdrawBalanceRequest) (*WithdrawBalanceReply, error)
	// CreateTransaction 创建交易记录
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionReply, error)
	// GetMerchantVersion 获取商家版本号
	GetMerchantVersion(ctx context.Context, req *GetMerchantVersionRequest) (*GetMerchantVersionReply, error)
}

type BalanceUsecase struct {
	repo BalanceRepo
	log  *log.Helper
}

func NewBalanceUsecase(repo BalanceRepo, logger log.Logger) *BalanceUsecase {
	return &BalanceUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (cc *BalanceUsecase) CreateConsumerBalance(ctx context.Context, req *CreateConsumerBalanceRequest) (*CreateConsumerBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("CreateConsumerBalance request: %+v", req)
	return cc.repo.CreateConsumerBalance(ctx, req)
}

func (cc *BalanceUsecase) CreateMerchantBalance(ctx context.Context, req *CreateMerchantBalanceRequest) (*CreateMerchantBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("CreateMerchantBalance request: %+v", req)
	return cc.repo.CreateMerchantBalance(ctx, req)
}

func (cc *BalanceUsecase) CancelFreeze(ctx context.Context, req *CancelFreezeRequest) (*CancelFreezeReply, error) {
	cc.log.WithContext(ctx).Debugf("CancelFreeze request: %+v", req)
	return cc.repo.CancelFreeze(ctx, req)
}

func (cc *BalanceUsecase) ConfirmTransfer(ctx context.Context, req *ConfirmTransferRequest) (*ConfirmTransferReply, error) {
	cc.log.WithContext(ctx).Debugf("ConfirmTransfer request: %+v", req)
	return cc.repo.ConfirmTransfer(ctx, req)
}

func (cc *BalanceUsecase) FreezeBalance(ctx context.Context, req *FreezeBalanceRequest) (*FreezeBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("FreezeBalance request: %+v", req)
	return cc.repo.FreezeBalance(ctx, req)
}

func (cc *BalanceUsecase) GetMerchantBalance(ctx context.Context, req *GetMerchantBalanceRequest) (*BalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("GetMerchantBalance request: %+v", req)
	return cc.repo.GetMerchantBalance(ctx, req)
}

func (cc *BalanceUsecase) GetUserBalance(ctx context.Context, req *GetUserBalanceRequest) (*BalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("GetUserBalance request: %+v", req)
	return cc.repo.GetUserBalance(ctx, req)
}

func (cc *BalanceUsecase) RechargeBalance(ctx context.Context, req *RechargeBalanceRequest) (*RechargeBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("RechargeBalance request: %+v", req)
	return cc.repo.RechargeBalance(ctx, req)
}

func (cc *BalanceUsecase) RechargeMerchantBalance(ctx context.Context, req *RechargeMerchantBalanceRequest) (*RechargeMerchantBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("RechargeMerchantBalance request: %+v", req)
	return cc.repo.RechargeMerchantBalance(ctx, req)
}

func (cc *BalanceUsecase) WithdrawBalance(ctx context.Context, req *WithdrawBalanceRequest) (*WithdrawBalanceReply, error) {
	cc.log.WithContext(ctx).Debugf("WithdrawBalance request: %+v", req)
	return cc.repo.WithdrawBalance(ctx, req)
}

func (cc *BalanceUsecase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionReply, error) {
	cc.log.WithContext(ctx).Debugf("CreateTransaction request: %+v", req)
	return cc.repo.CreateTransaction(ctx, req)
}

func (cc *BalanceUsecase) GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsReply, error) {
	cc.log.WithContext(ctx).Debugf("GetTransactions request: %+v", req)
	return cc.repo.GetTransactions(ctx, req)
}

func (cc *BalanceUsecase) GetMerchantVersion(ctx context.Context, req *GetMerchantVersionRequest) (*GetMerchantVersionReply, error) {
	cc.log.WithContext(ctx).Debugf("biz/order GetMerchantVersion req:%+v", req)
	return cc.repo.GetMerchantVersion(ctx, req)
}
