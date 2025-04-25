package service

import (
	"context"
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

func (s *BalanceService) GetUserBalance(ctx context.Context, req *pb.GetUserBalanceRequest) (*pb.BalanceResponse, error) {
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
	return &pb.BalanceResponse{
		Available: balance.Available,
		Frozen:    balance.Frozen,
		// Currency:  pb.Currency(balance.Currency),
		Version: balance.Version,
	}, nil
}

func (s *BalanceService) FreezeBalance(ctx context.Context, req *pb.FreezeBalanceRequest) (*pb.FreezeBalanceResponse, error) {
	userId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now()
	if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt.AsTime()
	} else {
		// 默认24小时后过期
		expiresAt = time.Now().Add(24 * time.Hour)
	}

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

	return &pb.FreezeBalanceResponse{
		FreezeId:   result.FreezeId,
		NewVersion: result.NewVersion,
	}, nil
}

func (s *BalanceService) ConfirmTransfer(ctx context.Context, req *pb.ConfirmTransferRequest) (*pb.ConfirmTransferResponse, error) {
	result, err := s.uc.ConfirmTransfer(ctx, &biz.ConfirmTransferRequest{
		FreezeId:                req.FreezeId,
		IdempotencyKey:          req.IdempotencyKey,
		ExpectedUserVersion:     req.ExpectedUserVersion,
		ExpectedMerchantVersion: req.ExpectedMerchantVersion,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ConfirmTransferResponse{
		Success:            result.Success,
		TransactionId:      result.TransactionId,
		NewUserVersion:     result.NewUserVersion,
		NewMerchantVersion: result.NewMerchantVersion,
	}, nil
}

func (s *BalanceService) CancelFreeze(ctx context.Context, req *pb.CancelFreezeRequest) (*pb.CancelFreezeResponse, error) {
	result, err := s.uc.CancelFreeze(ctx, &biz.CancelFreezeRequest{
		FreezeId:        req.FreezeId,
		Reason:          req.Reason,
		IdempotencyKey:  req.IdempotencyKey,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CancelFreezeResponse{
		Success:    result.Success,
		NewVersion: result.NewVersion,
	}, nil
}

func (s *BalanceService) GetMerchantBalance(ctx context.Context, req *pb.GetMerchantBalanceRequest) (*pb.BalanceResponse, error) {
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

	return &pb.BalanceResponse{
		Available: balance.Available,
		Frozen:    balance.Frozen,
		Currency:  string(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (s *BalanceService) RechargeBalance(ctx context.Context, req *pb.RechargeBalanceRequest) (*pb.RechargeBalanceResponse, error) {
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

	return &pb.RechargeBalanceResponse{
		Success:       result.Success,
		TransactionId: result.TransactionId,
		NewVersion:    result.NewVersion,
	}, nil
}

func (s *BalanceService) WithdrawBalance(ctx context.Context, req *pb.WithdrawBalanceRequest) (*pb.WithdrawBalanceResponse, error) {
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

	return &pb.WithdrawBalanceResponse{
		Success:       result.Success,
		TransactionId: result.TransactionId,
		NewVersion:    result.NewVersion,
	}, nil
}
