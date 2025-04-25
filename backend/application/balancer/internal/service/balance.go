package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"backend/constants"

	"backend/application/balancer/internal/biz"
	globalPkg "backend/pkg"

	pb "backend/api/balancer/v1"
)

type BalanceService struct {
	pb.UnimplementedBalanceServer
	uc *biz.BalancerUsecase
}

func NewBalancerService(uc *biz.BalancerUsecase) *BalanceService {
	return &BalanceService{uc: uc}
}

func (s *BalanceService) CreateConsumerBalance(ctx context.Context, req *pb.CreateConsumerBalanceRequest) (*pb.CreateConsumerBalanceReply, error) {
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

	return &pb.CreateConsumerBalanceReply{
		UserId:    reply.UserId.String(),
		Currency:  string(reply.Currency),
		Available: reply.Available,
	}, nil
}

func (s *BalanceService) CreateMerchantBalance(ctx context.Context, req *pb.CreateMerchantBalanceRequest) (*pb.CreateMerchantBalanceReply, error) {
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

	return &pb.CreateMerchantBalanceReply{
		UserId:    reply.UserId.String(),
		Currency:  string(reply.Currency),
		Available: reply.Available,
	}, nil
}

func (s *BalanceService) GetUserBalance(ctx context.Context, req *pb.GetUserBalanceRequest) (*pb.BalanceReply, error) {
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
	return &pb.BalanceReply{
		Available: balance.Available,
		Frozen:    balance.Frozen,
		Currency:  string(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (s *BalanceService) FreezeBalance(ctx context.Context, req *pb.FreezeBalanceRequest) (*pb.FreezeBalanceReply, error) {
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

	return &pb.FreezeBalanceReply{
		FreezeId:   result.FreezeId,
		NewVersion: result.NewVersion,
	}, nil
}

func (s *BalanceService) ConfirmTransfer(ctx context.Context, req *pb.ConfirmTransferRequest) (*pb.ConfirmTransferReply, error) {
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

	return &pb.ConfirmTransferReply{
		Success:            result.Success,
		TransactionId:      result.TransactionId,
		NewUserVersion:     result.NewUserVersion,
		NewMerchantVersion: result.NewMerchantVersion,
	}, nil
}

func (s *BalanceService) CancelFreeze(ctx context.Context, req *pb.CancelFreezeRequest) (*pb.CancelFreezeReply, error) {
	result, err := s.uc.CancelFreeze(ctx, &biz.CancelFreezeRequest{
		FreezeId:        req.FreezeId,
		Reason:          req.Reason,
		IdempotencyKey:  req.IdempotencyKey,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CancelFreezeReply{
		Success:    result.Success,
		NewVersion: result.NewVersion,
	}, nil
}

func (s *BalanceService) GetMerchantBalance(ctx context.Context, req *pb.GetMerchantBalanceRequest) (*pb.BalanceReply, error) {
	var err error
	var userId uuid.UUID
	userId, err = globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
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

	return &pb.BalanceReply{
		Available: balance.Available,
		Frozen:    balance.Frozen,
		Currency:  string(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (s *BalanceService) RechargeBalance(ctx context.Context, req *pb.RechargeBalanceRequest) (*pb.RechargeBalanceReply, error) {
	userId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := s.uc.RechargeBalance(ctx, &biz.RechargeBalanceRequest{
		UserId:                userId,
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

	return &pb.RechargeBalanceReply{
		Success:       result.Success,
		TransactionId: result.TransactionId,
		NewVersion:    result.NewVersion,
	}, nil
}

func (s *BalanceService) WithdrawBalance(ctx context.Context, req *pb.WithdrawBalanceRequest) (*pb.WithdrawBalanceReply, error) {
	userId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := s.uc.WithdrawBalance(ctx, &biz.WithdrawBalanceRequest{
		UserId:          userId,
		Amount:          req.Amount,
		Currency:        constants.Currency(req.Currency),
		PaymentMethodId: req.PaymentMethodId,
		IdempotencyKey:  req.IdempotencyKey,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}

	return &pb.WithdrawBalanceReply{
		Success:       result.Success,
		TransactionId: result.TransactionId,
		NewVersion:    result.NewVersion,
	}, nil
}
