package service

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-kratos/kratos/v2/log"

	structpb "github.com/golang/protobuf/ptypes/struct"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	"backend/constants"

	"backend/application/balancer/internal/biz"
	globalPkg "backend/pkg"

	v1 "backend/api/balancer/v1"
)

type BalanceService struct {
	v1.UnimplementedBalanceServer
	uc *biz.BalancerUsecase
}

func NewBalancerService(uc *biz.BalancerUsecase) *BalanceService {
	return &BalanceService{uc: uc}
}

func (s *BalanceService) CreateConsumerBalance(ctx context.Context, req *v1.CreateConsumerBalanceRequest) (*v1.CreateConsumerBalanceReply, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	accountDetails, err := json.Marshal(req.AccountDetails)
	if err != nil {
		return nil, err
	}

	reply, err := s.uc.CreateConsumerBalance(ctx, &biz.CreateConsumerBalanceRequest{
		UserId:         userId,
		Currency:       constants.Currency(req.Currency),
		InitialBalance: req.InitialBalance,
		BalancerType:   req.BalancerType,
		IsDefault:      req.IsDefault,
		AccountDetails: accountDetails,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateConsumerBalanceReply{
		UserId:    reply.UserId.String(),
		Currency:  string(reply.Currency),
		Available: reply.Available,
	}, nil
}

func (s *BalanceService) CreateMerchantBalance(ctx context.Context, req *v1.CreateMerchantBalanceRequest) (*v1.CreateMerchantBalanceReply, error) {
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, err
	}

	accountDetails, err := json.Marshal(req.AccountDetails)
	if err != nil {
		return nil, err
	}

	reply, err := s.uc.CreateMerchantBalance(ctx, &biz.CreateMerchantBalanceRequest{
		MerchantId:     merchantId,
		Currency:       constants.Currency(req.Currency),
		InitialBalance: req.InitialBalance,
		BalancerType:   req.BalancerType,
		IsDefault:      req.IsDefault,
		AccountDetails: accountDetails,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateMerchantBalanceReply{
		UserId:    reply.UserId.String(),
		Currency:  string(reply.Currency),
		Available: reply.Available,
	}, nil
}

func (s *BalanceService) GetUserBalance(ctx context.Context, req *v1.GetUserBalanceRequest) (*v1.BalanceReply, error) {
	userId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	balance, err := s.uc.GetUserBalance(ctx, &biz.GetUserBalanceRequest{
		UserId:   userId,
		Currency: constants.Currency(req.Currency),
	})
	if err != nil {
		return nil, err
	}
	return &v1.BalanceReply{
		Available: balance.Available,
		Frozen:    balance.Frozen,
		Currency:  string(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (s *BalanceService) FreezeBalance(ctx context.Context, req *v1.FreezeBalanceRequest) (*v1.FreezeBalanceReply, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	// 默认24小时后过期
	expiresAt := time.Now().Add(24 * time.Hour)

	result, err := s.uc.FreezeBalance(ctx, &biz.FreezeBalanceRequest{
		UserId:          userId,
		OrderId:         req.OrderId,
		Amount:          req.Amount,
		Currency:        constants.Currency(req.Currency),
		ExpiresAt:       expiresAt,
		IdempotencyKey:  req.IdempotencyKey,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}

	return &v1.FreezeBalanceReply{
		FreezeId:   result.FreezeId,
		NewVersion: result.NewVersion,
	}, nil
}

func (s *BalanceService) ConfirmTransfer(ctx context.Context, req *v1.ConfirmTransferRequest) (*v1.ConfirmTransferReply, error) {
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, err
	}
	result, err := s.uc.ConfirmTransfer(ctx, &biz.ConfirmTransferRequest{
		FreezeId:                req.FreezeId,
		PaymentAccount:          req.PaymentAccount,
		IdempotencyKey:          req.IdempotencyKey,
		ExpectedUserVersion:     req.ExpectedUserVersion,
		ExpectedMerchantVersion: req.ExpectedMerchantVersion,
		MerchantId:              merchantId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConfirmTransferReply{
		Success:            result.Success,
		TransactionId:      result.TransactionId,
		NewUserVersion:     result.NewUserVersion,
		NewMerchantVersion: result.NewMerchantVersion,
	}, nil
}

func (s *BalanceService) CancelFreeze(ctx context.Context, req *v1.CancelFreezeRequest) (*v1.CancelFreezeReply, error) {
	result, err := s.uc.CancelFreeze(ctx, &biz.CancelFreezeRequest{
		FreezeId:        req.FreezeId,
		Reason:          req.Reason,
		IdempotencyKey:  req.IdempotencyKey,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CancelFreezeReply{
		Success:    result.Success,
		NewVersion: result.NewVersion,
	}, nil
}

func (s *BalanceService) GetMerchantBalance(ctx context.Context, req *v1.GetMerchantBalanceRequest) (*v1.BalanceReply, error) {
	var userId uuid.UUID
	var err error
	if req.MerchantId == "" {
		userId, err = globalPkg.GetMetadataUesrID(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		userId, err = uuid.Parse(req.MerchantId)
		if err != nil {
			return nil, err
		}
	}

	balance, err := s.uc.GetMerchantBalance(ctx, &biz.GetMerchantBalanceRequest{
		MerchantId: userId,
		Currency:   constants.Currency(req.Currency),
	})
	if err != nil {
		return nil, err
	}

	return &v1.BalanceReply{
		Available: balance.Available,
		Frozen:    balance.Frozen,
		Currency:  string(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (s *BalanceService) GetTransactions(ctx context.Context, req *v1.GetTransactionsRequest) (*v1.GetTransactionsReply, error) {
	var err error
	var userId uuid.UUID
	if req.UserId == "" {
		userId, err = globalPkg.GetMetadataUesrID(ctx)
	}
	userId, err = uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	transactions, err := s.uc.GetTransactions(ctx, &biz.GetTransactionsRequest{
		UserId:        userId,
		Currency:      req.Currency,
		Page:          req.Page,
		PageSize:      req.PageSize,
		PaymentStatus: constants.PaymentStatus(req.PaymentStatus),
	})
	if err != nil {
		return nil, err
	}

	var pbTransactions []*v1.Transactions
	for _, t := range transactions.Transactions {
		pbTransaction := &v1.Transactions{
			Id:                t.Id,
			Type:              string(t.Type),
			Amount:            t.Amount,
			Currency:          t.Currency,
			FromUserId:        t.FromUserId.String(),
			ToMerchantId:      t.ToMerchantId.String(),
			PaymentMethodType: string(t.PaymentMethodType),
			PaymentAccount:    t.PaymentAccount,
			Status:            string(t.Status),
		}

		if t.PaymentExtra != nil {
			pbTransaction.PaymentExtra = &structpb.Struct{}
			if err := json.Unmarshal(t.PaymentExtra, &pbTransaction.PaymentExtra.Fields); err != nil {
				return nil, err
			}
		}
		log.Debugf("pbTransaction.PaymentExtra: %+v", pbTransaction.PaymentExtra)

		pbTransaction.CreatedAt = timestamppb.New(t.CreatedAt)
		pbTransaction.UpdatedAt = timestamppb.New(t.UpdatedAt)

		pbTransactions = append(pbTransactions, pbTransaction)
	}

	return &v1.GetTransactionsReply{
		Transactions: pbTransactions,
	}, nil
}

func (s *BalanceService) RechargeBalance(ctx context.Context, req *v1.RechargeBalanceRequest) (*v1.RechargeBalanceReply, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, err
	}

	result, err := s.uc.RechargeBalance(ctx, &biz.RechargeBalanceRequest{
		UserId:                userId,
		MerchantId:            merchantId,
		Amount:                req.Amount,
		Currency:              constants.Currency(req.Currency),
		ExternalTransactionId: req.ExternalTransactionId,
		PaymentMethodType:     req.PaymentMethodType,
		PaymentAccount:        req.PaymentAccount,
		IdempotencyKey:        req.IdempotencyKey,
		ExpectedVersion:       req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RechargeBalanceReply{
		Success:       result.Success,
		TransactionId: result.TransactionId,
		NewVersion:    result.NewVersion,
	}, nil
}

func (s *BalanceService) WithdrawBalance(ctx context.Context, req *v1.WithdrawBalanceRequest) (*v1.WithdrawBalanceReply, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, err
	}

	result, err := s.uc.WithdrawBalance(ctx, &biz.WithdrawBalanceRequest{
		UserId:          userId,
		MerchantId:      merchantId,
		Amount:          req.Amount,
		Currency:        constants.Currency(req.Currency),
		PaymentMethodId: req.PaymentMethodId,
		IdempotencyKey:  req.IdempotencyKey,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}

	return &v1.WithdrawBalanceReply{
		Success:       result.Success,
		TransactionId: result.TransactionId,
		NewVersion:    result.NewVersion,
	}, nil
}

func (s *BalanceService) CreateTransaction(ctx context.Context, req *v1.CreateTransactionRequest) (*v1.CreateTransactionReply, error) {
	consumerId, err := uuid.Parse(req.FromUserId)
	if err != nil {
		return nil, err
	}
	merchantId, err := uuid.Parse(req.ToMerchantId)
	if err != nil {
		return nil, err
	}
	paymentExtra, err := json.Marshal(req.PaymentExtra.AsMap())
	if err != nil {
		return nil, err
	}

	reply, err := s.uc.CreateTransaction(ctx, &biz.CreateTransactionRequest{
		Type:              constants.TransactionType(req.Type),
		Amount:            req.Amount,
		Currency:          constants.Currency(req.Currency),
		FromUserId:        consumerId,
		ToMerchantId:      merchantId,
		PaymentMethodType: constants.PaymentMethod(req.PaymentMethodType),
		PaymentAccount:    req.PaymentAccount,
		PaymentExtra:      paymentExtra,
		Status:            constants.PaymentStatus(req.Status),
		IdempotencyKey:    req.IdempotencyKey,
		FreezeId:          req.FreezeId,
		ConsumerVersion:   req.ConsumerVersion,
		MerchantVersion:   req.MerchantVersion,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateTransactionReply{
		Id: reply.Id,
	}, nil
}

// GetMerchantVersion 获取商家版本号
func (s *BalanceService) GetMerchantVersion(ctx context.Context, req *v1.GetMerchantVersionRequest) (*v1.GetMerchantVersionReply, error) {
	merchantIds := make([]uuid.UUID, 0, len(req.MerchantIds))
	for _, id := range req.MerchantIds {
		merchantIds = append(merchantIds, uuid.MustParse(id))
	}

	reply, err := s.uc.GetMerchantVersion(ctx, &biz.GetMerchantVersionRequest{
		MerchantIds: merchantIds,
	})
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, status.Error(codes.NotFound, "merchant version not found")
	}
	versions := make([]int64, 0, len(reply.Versions))
	for _, v := range reply.Versions {
		versions = append(versions, v)
	}
	merchantPbIds := make([]string, 0, len(reply.MerchantIds))
	for _, v := range reply.MerchantIds {
		log.Debugf("MerchantId: %+v", v)
		merchantIds = append(merchantIds, v)
	}

	return &v1.GetMerchantVersionReply{
		MerchantVersion: versions,
		MerchantIds:     merchantPbIds,
	}, nil
}
