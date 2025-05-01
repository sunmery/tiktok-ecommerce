package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	merchantorderv1 "backend/api/merchant/order/v1"

	"github.com/google/uuid"

	"backend/application/balancer/internal/pkg/id"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"

	"backend/constants"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	"backend/pkg/types"

	"backend/application/balancer/internal/data/models"

	"backend/application/balancer/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type balanceRepo struct {
	data *Data
	log  *log.Helper
}

func NewBalanceRepo(data *Data, logger log.Logger) biz.BalanceRepo {
	return &balanceRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (b balanceRepo) CreateTransaction(ctx context.Context, req *biz.CreateTransactionRequest) (*biz.CreateTransactionReply, error) {
	amount, err := types.Float64ToNumeric(req.Amount)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	transactionId, err := b.data.db.CreateTransaction(ctx, models.CreateTransactionParams{
		ID:                id.SnowflakeID(),
		Type:              string(req.Type),
		Amount:            amount,
		Currency:          string(req.Currency),
		FromUserID:        req.FromUserId,
		ToMerchantID:      req.ToMerchantId,
		PaymentMethodType: string(req.PaymentMethodType),
		PaymentAccount:    req.PaymentAccount,
		PaymentExtra:      req.PaymentExtra,
		Status:            string(req.Status),
		FreezeID:          req.FreezeId,
		IdempotencyKey:    req.IdempotencyKey,
		ConsumerVersion:   req.ConsumerVersion,
		MerchantVersion:   req.MerchantVersion,
	})
	if err != nil {
		return nil, err
	}
	return &biz.CreateTransactionReply{
		Id: transactionId,
	}, nil
}

func (b balanceRepo) GetMerchantVersion(ctx context.Context, req *biz.GetMerchantVersionRequest) (*biz.GetMerchantVersionReply, error) {
	result, err := b.data.db.GetMerchantVersions(ctx, req.MerchantIds)
	if err != nil {
		b.log.WithContext(ctx).Errorf("Failed to get merchant version: %v", err)
		return nil, fmt.Errorf("failed to get merchant version: %w", err)
	}

	versions := make([]int64, 0, len(result))
	merchantIds := make([]uuid.UUID, 0, len(result))
	for _, v := range result {
		versions = append(versions, int64(v.Version))
		merchantIds = append(merchantIds, v.MerchantID)
	}

	return &biz.GetMerchantVersionReply{
		Versions:    versions,
		MerchantIds: merchantIds,
	}, nil
}

func (b balanceRepo) GetTransactions(ctx context.Context, req *biz.GetTransactionsRequest) (*biz.GetTransactionsReply, error) {
	page := (req.Page - 1) * req.PageSize
	pageSize := req.PageSize

	var transactions []*biz.Transaction
	switch req.TransactionsUserType {
	case constants.TransactionsUserTypeConsumer:
		return b.getConsumerTransactions(ctx, req, page, pageSize)
	case constants.TransactionsUserTypeMerchant:
		return b.getMerchantTransactions(ctx, req, page, pageSize)
	}

	return &biz.GetTransactionsReply{
		Transactions: transactions,
	}, nil
}

// FreezeBalance 冻结余额
func (b balanceRepo) FreezeBalance(ctx context.Context, req *biz.FreezeBalanceRequest) (*biz.FreezeBalanceReply, error) {
	// 1. 开始事务
	tx := b.data.DB(ctx)

	// 2. 冻结用户余额
	amount, err := types.Float64ToNumeric(req.Amount)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	// 尝试获取已存在的冻结记录
	existingFreeze, err := tx.GetFreezeByOrderForUser(ctx, models.GetFreezeByOrderForUserParams{
		UserID:  req.UserId,
		OrderID: req.OrderId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// 执行冻结操作
			rows, err := tx.FreezeUserBalance(ctx, models.FreezeUserBalanceParams{
				UserID:          req.UserId,
				Currency:        string(req.Currency),
				Amount:          amount,
				ExpectedVersion: req.ExpectedVersion,
			})
			if err != nil {
				return nil, kerrors.New(500, "FREEZE_BALANCE_FAILED", "freeze balance failed")
			}
			if rows == 0 {
				return nil, kerrors.New(409, "OPTIMISTIC_LOCK_FAILED", fmt.Sprintf("balance version mismatch or insufficient funds:%v", err))
			}

			// 3. 创建冻结记录
			freezeId, err := tx.CreateFreeze(ctx, models.CreateFreezeParams{
				ID:       id.SnowflakeID(),
				UserID:   req.UserId,
				OrderID:  req.OrderId,
				Currency: string(req.Currency),
				Amount:   amount,
				ExpiresAt: pgtype.Timestamptz{
					Time:  req.ExpiresAt,
					Valid: true,
				},
			})
			if err != nil {
				return nil, kerrors.New(500, "CREATE_FREEZE_FAILED", fmt.Sprintf("create freeze failed: %v", err))
			}

			return &biz.FreezeBalanceReply{
				FreezeId:   freezeId,
				NewVersion: req.ExpectedVersion + 1, // 版本号+1
			}, nil
		}
		return nil, kerrors.New(500, "GET_FREEZE_FAILED", "failed to check existing freeze")
	}

	// 如果已存在冻结记录，直接返回现有记录ID
	return &biz.FreezeBalanceReply{
		FreezeId:   existingFreeze.ID,
		NewVersion: req.ExpectedVersion, // 保持版本不变，因为没有执行冻结操作
	}, nil
}

// CancelFreeze 取消冻结余额
func (b balanceRepo) CancelFreeze(ctx context.Context, req *biz.CancelFreezeRequest) (*biz.CancelFreezeReply, error) {
	// 1. 开始事务
	tx := b.data.DB(ctx)

	// 2. 获取冻结记录
	freeze, err := tx.GetFreeze(ctx, req.FreezeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "FREEZE_NOT_FOUND", "freeze record not found")
		}
		return nil, kerrors.New(500, "GET_FREEZE_FAILED", "get freeze record failed")
	}

	// 3. 检查冻结状态
	if freeze.Status != string(constants.FreezeFrozen) { // 将 constants.FreezeFrozen 转换为 string
		return nil, kerrors.New(400, "INVALID_FREEZE_STATUS", fmt.Sprintf("freeze is not in FROZEN status: %v not in %v", freeze.Status, constants.FreezeFrozen))
	}

	// 4. 更新冻结记录状态为取消
	rows, err := tx.UpdateFreezeStatus(ctx, models.UpdateFreezeStatusParams{
		Status:        string(constants.FreezeCanceled),
		ID:            req.FreezeId,
		CurrentStatus: string(constants.FreezeFrozen),
	})
	if err != nil {
		return nil, kerrors.New(500, "UPDATE_FREEZE_STATUS_FAILED", "update freeze status failed")
	}
	if rows == 0 {
		return nil, kerrors.New(409, "FREEZE_STATUS_CHANGED", "freeze status has been changed")
	}

	// 5. 解冻用户余额（增加可用余额，减少冻结余额）
	rows, err = tx.UnfreezeUserBalance(ctx, models.UnfreezeUserBalanceParams{
		UserID:          freeze.UserID,
		Currency:        freeze.Currency,
		Amount:          freeze.Amount,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, kerrors.New(500, "UNFREEZE_BALANCE_FAILED", "unfreeze balance failed")
	}
	if rows == 0 {
		return nil, kerrors.New(409, "OPTIMISTIC_LOCK_FAILED", "balance version mismatch")
	}

	return &biz.CancelFreezeReply{
		Success:    true,
		NewVersion: req.ExpectedVersion + 1, // 版本号+1
	}, nil
}

func (b balanceRepo) GetMerchantBalance(ctx context.Context, req *biz.GetMerchantBalanceRequest) (*biz.BalanceReply, error) {
	balance, err := b.data.db.GetMerchantBalance(ctx, models.GetMerchantBalanceParams{
		MerchantID: req.MerchantId,
		Currency:   string(req.Currency),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "BALANCE_NOT_FOUND", fmt.Sprintf("getMerchantBalance balance record not found err:%+v", err))
		}
		return nil, err
	}
	available, err := types.NumericToFloat(balance.Available)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AVAILABLE_FAILED", "convert available to float64 failed")
	}

	return &biz.BalanceReply{
		Available: available,
		Frozen:    0, // 商家账户没有冻结余额
		Currency:  constants.Currency(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (b balanceRepo) CreateConsumerBalance(ctx context.Context, req *biz.CreateConsumerBalanceRequest) (*biz.CreateConsumerBalanceReply, error) {
	tx := b.data.DB(ctx)

	// 首先检查用户余额是否已存在
	_, err := tx.GetUserBalance(ctx, models.GetUserBalanceParams{
		UserID:   req.UserId,
		Currency: string(req.Currency),
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		b.log.WithContext(ctx).Errorf("Failed to check existing consumer balance: %v", err)
		return nil, fmt.Errorf("failed to check existing consumer balance: %w", err)
	}

	// 创建用户余额
	initialBalance, err := types.Float64ToNumeric(req.InitialBalance)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	err = tx.CreateConsumerPaymentMethods(ctx, models.CreateConsumerPaymentMethodsParams{
		ID:             id.SnowflakeID(),
		UserID:         req.UserId,
		Type:           req.BalanceType,
		IsDefault:      req.IsDefault,
		AccountDetails: req.AccountDetails,
	})
	if err != nil {
		b.log.WithContext(ctx).Errorf("Failed to create consumer payment method: %v", err)
		return nil, fmt.Errorf("failed to create consumer payment method: %w", err)
	}

	reply, err := tx.CreateConsumerBalance(ctx, models.CreateConsumerBalanceParams{
		UserID:    req.UserId,
		Currency:  string(req.Currency),
		Available: initialBalance,
	})
	if err != nil {
		b.log.WithContext(ctx).Errorf("Failed to create consumer balance: %v", err)
		return nil, fmt.Errorf("failed to create consumer balance: %w", err)
	}

	available, err := types.NumericToFloat(reply.Available)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AVAILABLE_FAILED", "convert available to float64 failed")
	}

	return &biz.CreateConsumerBalanceReply{
		UserId:    reply.UserID,
		Currency:  constants.Currency(reply.Currency),
		Available: available,
	}, nil
}

func (b balanceRepo) CreateMerchantBalance(ctx context.Context, req *biz.CreateMerchantBalanceRequest) (*biz.CreateMerchantBalanceReply, error) {
	tx := b.data.db

	// 首先检查商家余额是否已存在
	_, err := tx.GetMerchantBalance(ctx, models.GetMerchantBalanceParams{
		MerchantID: req.MerchantId,
		Currency:   string(req.Currency),
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		b.log.WithContext(ctx).Errorf("Failed to check existing merchant balance: %v", err)
		return nil, fmt.Errorf("failed to check existing merchant balance: %w", err)
	}

	// 创建商家余额
	amount, err := types.Float64ToNumeric(req.InitialBalance)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	err = tx.CreateMerchantPaymentMethods(ctx, models.CreateMerchantPaymentMethodsParams{
		ID:             id.SnowflakeID(),
		MerchantID:     req.MerchantId,
		Type:           req.BalanceType,
		IsDefault:      req.IsDefault,
		AccountDetails: req.AccountDetails,
	})
	if err != nil {
		b.log.WithContext(ctx).Errorf("Failed to create merchant payment method: %v", err)
		return nil, fmt.Errorf("failed to create merchant payment method: %w", err)
	}

	merchantBalance, err := tx.CreateMerchantBalance(ctx, models.CreateMerchantBalanceParams{
		MerchantID: req.MerchantId,
		Currency:   string(req.Currency),
		Available:  amount,
	})
	if err != nil {
		b.log.WithContext(ctx).Errorf("Failed to create merchant balance: %v", err)
		return nil, fmt.Errorf("failed to create merchant balance: %w", err)
	}

	available, err := types.NumericToFloat(merchantBalance.Available)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AVAILABLE_FAILED", "convert available to float64 failed")
	}
	return &biz.CreateMerchantBalanceReply{
		UserId:    merchantBalance.MerchantID,
		Currency:  constants.Currency(merchantBalance.Currency),
		Available: available,
	}, nil
}

func (b balanceRepo) GetUserBalance(ctx context.Context, req *biz.GetUserBalanceRequest) (*biz.BalanceReply, error) {
	balance, err := b.data.db.GetUserBalance(ctx, models.GetUserBalanceParams{
		UserID:   req.UserId,
		Currency: string(req.Currency),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "BALANCE_NOT_FOUND", fmt.Sprintf("getUserBalance balance record not found err:%+v", err))
		}
		return nil, err
	}
	available, err := types.NumericToFloat(balance.Available)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AVAILABLE_FAILED", "convert available to float64 failed")
	}
	frozen, err := types.NumericToFloat(balance.Frozen)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "BALANCE_NOT_FOUND", fmt.Sprintf("getUserBalance balance record not found err:%+v", err))
		}
		return nil, kerrors.New(500, "CONVERT_FROZEN_FAILED", "convert available to float64 failed")
	}

	return &biz.BalanceReply{
		Available: available,
		Frozen:    frozen,
		Currency:  constants.Currency(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (b balanceRepo) RechargeBalance(ctx context.Context, req *biz.RechargeBalanceRequest) (*biz.RechargeBalanceReply, error) {
	// 1. 开始事务
	tx := b.data.DB(ctx)

	// 2. 增加用户可用余额
	amount, err := types.Float64ToNumeric(req.Amount)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	rows, err := tx.IncreaseUserAvailableBalance(ctx, models.IncreaseUserAvailableBalanceParams{
		UserID:          req.UserId,
		Currency:        string(req.Currency),
		Amount:          amount,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, kerrors.New(500, "INCREASE_BALANCE_FAILED", "increase balance failed")
	}
	if rows == 0 {
		return nil, kerrors.New(409, "OPTIMISTIC_LOCK_FAILED", "balance version mismatch")
	}

	// 构建额外信息JSON
	paymentExtra := map[string]interface{}{
		"external_transaction_id": req.ExternalTransactionId,
		"idempotency_key":         req.IdempotencyKey,
	}
	paymentExtraJson, err := json.Marshal(paymentExtra)
	if err != nil {
		return nil, kerrors.New(500, "JSON_MARSHAL_FAILED", "marshal payment extra failed")
	}

	// 创建交易记录
	transactionId, err := tx.CreateTransaction(ctx, models.CreateTransactionParams{
		ID:         id.SnowflakeID(),
		Type:       string(constants.TransactionRecharge),
		Amount:     amount,
		Currency:   string(req.Currency),
		FromUserID: req.UserId,
		// ToMerchantID:      merchantId,
		PaymentMethodType: req.PaymentMethodType,
		PaymentAccount:    req.PaymentAccount,
		PaymentExtra:      paymentExtraJson,
		Status:            string(constants.PaymentPaid), // 充值通常是已支付状态
	})
	if err != nil {
		return nil, kerrors.New(500, "CREATE_TRANSACTION_FAILED", fmt.Sprintf("create transaction failed: %v", err))
	}

	return &biz.RechargeBalanceReply{
		Success:       true,
		TransactionId: transactionId,
		NewVersion:    req.ExpectedVersion + 1, // 版本号+1
	}, nil
}

func (b balanceRepo) RechargeMerchantBalance(ctx context.Context, req *biz.RechargeMerchantBalanceRequest) (*biz.RechargeMerchantBalanceReply, error) {
	tx := b.data.db

	// 转换金额为数据库格式
	amount, err := types.Float64ToNumeric(req.Amount)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	// 增加商家余额
	rowsAffected, err := tx.IncreaseMerchantAvailableBalance(ctx, models.IncreaseMerchantAvailableBalanceParams{
		MerchantID:      req.MerchantId,
		Currency:        string(req.Currency),
		Amount:          amount,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		// 提供更详细的错误信息
		return nil, kerrors.New(409, "OPTIMISTIC_LOCK_FAILED", "数据已被修改，请刷新后重试")
	}

	// 创建交易记录
	paymentExtra, err := json.Marshal(req.PaymentExtra)
	if err != nil {
		return nil, kerrors.New(500, "MARSHAL_PAYMENT_EXTRA_FAILED", "marshal payment extra failed")
	}

	transactionId, err := tx.CreateTransaction(ctx, models.CreateTransactionParams{
		ID:                id.SnowflakeID(),
		Type:              string(constants.TransactionRecharge),
		Amount:            amount,
		Currency:          string(req.Currency),
		FromUserID:        uuid.Nil, // 充值没有来源用户
		ToMerchantID:      req.MerchantId,
		PaymentMethodType: string(req.PaymentMethod),
		PaymentAccount:    req.PaymentAccount,
		PaymentExtra:      paymentExtra,
		Status:            string(constants.PaymentPaid), // 修改为正确的状态值
		FreezeID:          0,
		IdempotencyKey:    req.IdempotencyKey,
		ConsumerVersion:   0,
		MerchantVersion:   int64(req.ExpectedVersion + 1),
	})
	if err != nil {
		return nil, err
	}

	return &biz.RechargeMerchantBalanceReply{
		TransactionId: transactionId,
		NewVersion:    req.ExpectedVersion + 1,
	}, nil
}

// WithdrawBalance 用户提现
func (b balanceRepo) WithdrawBalance(ctx context.Context, req *biz.WithdrawBalanceRequest) (*biz.WithdrawBalanceReply, error) {
	// 1. 开始事务
	tx := b.data.DB(ctx)

	// 2. 获取用户支付方式详情
	paymentMethodId, err := strconv.ParseInt(req.PaymentMethodId, 10, 64)
	if err != nil {
		return nil, kerrors.New(400, "INVALID_PAYMENT_METHOD_ID", "invalid payment method id format")
	}

	paymentMethod, err := tx.GetUserPaymentMethod(ctx, models.GetUserPaymentMethodParams{
		ID:     paymentMethodId,
		UserID: req.UserId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "PAYMENT_METHOD_NOT_FOUND", "payment method not found")
		}
		return nil, kerrors.New(500, "GET_PAYMENT_METHOD_FAILED", "get payment method failed")
	}

	// 3. 减少用户可用余额
	amount, err := types.Float64ToNumeric(req.Amount)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	rows, err := tx.DecreaseUserAvailableBalance(ctx, models.DecreaseUserAvailableBalanceParams{
		UserID:          req.UserId,
		Currency:        string(req.Currency),
		Amount:          amount,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, kerrors.New(500, "DECREASE_BALANCE_FAILED", "decrease balance failed")
	}
	if rows == 0 {
		return nil, kerrors.New(409, "OPTIMISTIC_LOCK_FAILED", fmt.Sprintf("balance version mismatch or insufficient funds:%v", err))
	}

	// 4. 创建交易记录
	merchantId := req.MerchantId

	// 从支付方式中提取账号信息
	var accountInfo string
	// if account, ok := paymentMethod.AccountDetails["account"].(string); ok {
	// 	accountInfo = account
	// } else {
	// 	accountInfo = "unknown"
	// }

	// 构建额外信息JSON
	paymentExtra := map[string]interface{}{
		"idempotency_key":   req.IdempotencyKey,
		"payment_method_id": req.PaymentMethodId,
	}
	paymentExtraJson, err := json.Marshal(paymentExtra)
	if err != nil {
		return nil, kerrors.New(500, "JSON_MARSHAL_FAILED", "marshal payment extra failed")
	}

	// 创建交易记录
	transactionId, err := tx.CreateTransaction(ctx, models.CreateTransactionParams{
		ID:                id.SnowflakeID(),
		Type:              string(constants.TransactionWithdraw),
		Amount:            amount,
		Currency:          string(req.Currency),
		FromUserID:        req.UserId,
		ToMerchantID:      merchantId,
		PaymentMethodType: paymentMethod.Type,
		PaymentAccount:    accountInfo,
		PaymentExtra:      paymentExtraJson,
		Status:            "PENDING", // 提现通常是待处理状态
	})
	if err != nil {
		return nil, kerrors.New(500, "CREATE_TRANSACTION_FAILED", fmt.Sprintf("create transaction failed: %v", err))
	}

	return &biz.WithdrawBalanceReply{
		Success:       true,
		TransactionId: transactionId,
		NewVersion:    req.ExpectedVersion + 1, // 版本号+1
	}, nil
}

// ConfirmTransfer 确认转账
func (b balanceRepo) ConfirmTransfer(ctx context.Context, req *biz.ConfirmTransferRequest) (*biz.ConfirmTransferReply, error) {
	// 1. 开始事务
	tx := b.data.DB(ctx)

	// 2. 获取冻结记录
	freeze, err := tx.GetFreeze(ctx, req.FreezeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "FREEZE_NOT_FOUND", "freeze record not found")
		}
		return nil, kerrors.New(500, "GET_FREEZE_FAILED", "get freeze record failed")
	}

	// 3. 检查冻结状态
	if freeze.Status != string(constants.FreezeFrozen) { // 将 constants.FreezeFrozen 转换为 string
		return nil, kerrors.New(400, "INVALID_FREEZE_STATUS", fmt.Sprintf("freeze is not in FROZEN status: %s", freeze.Status))
	}

	// 4. 检查冻结是否过期
	if time.Now().After(freeze.ExpiresAt) {
		// 如果冻结已过期，返回错误，阻止后续操作
		// 可选：在这里可以添加逻辑将冻结状态更新为 EXPIRED 并解冻余额，但这通常由单独的清理任务完成
		return nil, kerrors.New(400, "FREEZE_EXPIRED", fmt.Sprintf("freeze record %d has expired", req.FreezeId))
	}

	// 5. 更新冻结记录状态为已使用
	rowsAffected, UpdateFreezeStatusErr := tx.UpdateFreezeStatus(ctx, models.UpdateFreezeStatusParams{
		Status:        string(constants.FreezeConfirmed),
		ID:            req.FreezeId,
		CurrentStatus: string(constants.FreezeFrozen),
	})
	if UpdateFreezeStatusErr != nil {
		return nil, kerrors.New(500, "UPDATE_FREEZE_STATUS_FAILED", fmt.Sprintf("update freeze status failed:%v", UpdateFreezeStatusErr.Error()))
	}
	if rowsAffected == 0 {
		return nil, kerrors.New(409, "FREEZE_STATUS_CHANGED", fmt.Sprintf("freeze status has been changed:%v", UpdateFreezeStatusErr))
	}

	// 5. 确认用户冻结（减少冻结余额）
	ConfirmUserFreezeErr := tx.ConfirmUserFreeze(ctx, models.ConfirmUserFreezeParams{
		UserID:          freeze.UserID,
		Currency:        freeze.Currency,
		Amount:          freeze.Amount,
		ExpectedVersion: req.ExpectedUserVersion,
	})
	if ConfirmUserFreezeErr != nil {
		return nil, kerrors.New(500, "CONFIRM_FREEZE_FAILED", fmt.Sprintf("confirm freeze failed: %v", ConfirmUserFreezeErr.Error()))
	}
	// if rows == 0 {
	// 	return nil, kerrors.New(409, "USER_OPTIMISTIC_LOCK_FAILED", "user balance version mismatch")
	// }

	// 6. 获取商家ID（从订单ID）
	merchantId, err := b.getMerchantIDFromOrder(ctx, freeze.OrderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "MERCHANT_NOT_FOUND", "merchant not found for order")
		}
		return nil, kerrors.New(500, "GET_MERCHANT_ID_FAILED", fmt.Sprintf("get merchant id from order failed:%v", err))
	}

	// 7. 增加商家可用余额
	params := models.UpdateMerchantAvailableBalanceParams{
		MerchantID:      merchantId,
		Currency:        freeze.Currency,
		Amount:          freeze.Amount,
		ExpectedVersion: req.ExpectedMerchantVersion,
	}
	log.Debugf("UpdateMerchantAvailableBalanceParams:%+v", params)
	rows, err := tx.UpdateMerchantAvailableBalance(ctx, params)
	if err != nil {
		return nil, kerrors.New(500, "INCREASE_MERCHANT_BALANCE_FAILED", "increase merchant balance failed")
	}
	if rows == 0 {
		return nil, kerrors.New(409, "MERCHANT_OPTIMISTIC_LOCK_FAILED", "merchant balance version mismatch")
	}

	// 8. 创建交易记录
	// 构建额外信息JSON
	paymentExtra := map[string]interface{}{
		"freeze_id":       req.FreezeId,
		"idempotency_key": req.IdempotencyKey,
		"order_id":        freeze.OrderID,
	}
	paymentExtraJson, err := json.Marshal(paymentExtra)
	if err != nil {
		return nil, kerrors.New(500, "JSON_MARSHAL_FAILED", "marshal payment extra failed")
	}

	// 创建交易记录
	transactionId, err := tx.CreateTransaction(ctx, models.CreateTransactionParams{
		ID:                id.SnowflakeID(),
		Type:              string(constants.TransactionPayment),
		Amount:            freeze.Amount,
		Currency:          freeze.Currency,
		FromUserID:        freeze.UserID,
		ToMerchantID:      merchantId,
		PaymentMethodType: string(constants.PaymentMethodBalancer), // 使用余额支付
		PaymentAccount:    req.PaymentAccount,
		PaymentExtra:      paymentExtraJson,
		Status:            string(constants.PaymentPaid), // 已支付状态
	})
	if err != nil {
		return nil, kerrors.New(500, "CREATE_TRANSACTION_FAILED", fmt.Sprintf("create transaction failed: %v", err))
	}

	return &biz.ConfirmTransferReply{
		Success:            true,
		TransactionId:      transactionId,
		NewUserVersion:     req.ExpectedUserVersion + 1,     // 用户余额版本号+1
		NewMerchantVersion: req.ExpectedMerchantVersion + 1, // 商家余额版本号+1
	}, nil
}

// getMerchantIDFromOrder 从订单ID获取商家ID
func (b balanceRepo) getMerchantIDFromOrder(ctx context.Context, subOrderId int64) (uuid.UUID, error) {
	reply, err := b.data.merchantOrderv1.GetMerchantByOrderId(ctx, &merchantorderv1.GetMerchantByOrderIdReq{
		OrderId: subOrderId,
	})
	if err != nil {
		return uuid.Nil, err
	}
	merchantId, err := uuid.Parse(reply.MerchantId)
	if err != nil {
		return uuid.Nil, err
	}
	return merchantId, nil
}

func (b balanceRepo) getConsumerTransactions(ctx context.Context, req *biz.GetTransactionsRequest, page int64, size int64) (*biz.GetTransactionsReply, error) {
	transactions, err := b.data.db.GetConsumerTransactions(ctx, models.GetConsumerTransactionsParams{
		UserID:   req.UserId,
		Currency: req.Currency,
		Status:   string(req.PaymentStatus),
		Page:     page,
		PageSize: size,
	})
	if err != nil {
		return nil, err
	}
	transactionsReply := make([]*biz.Transaction, 0, len(transactions))
	// maps.Copy(transactionsReply, transactions)
	for _, t := range transactions {
		amount, err := types.NumericToFloat(t.Amount)
		if err != nil {
			return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to float64 failed")
		}
		transactionsReply = append(transactionsReply, &biz.Transaction{
			Id:                t.ID,
			Type:              constants.TransactionType(t.Type),
			Amount:            amount,
			Currency:          t.Currency,
			FromUserId:        t.FromUserID,
			ToMerchantId:      t.ToMerchantID,
			PaymentMethodType: constants.PaymentMethod(t.PaymentMethodType),
			PaymentAccount:    t.PaymentAccount,
			Status:            constants.PaymentStatus(t.Status),
			PaymentExtra:      t.PaymentExtra,
			CreatedAt:         t.CreatedAt,
			UpdatedAt:         t.UpdatedAt,
		})
	}

	return &biz.GetTransactionsReply{
		Transactions: transactionsReply,
	}, nil
}

func (b balanceRepo) getMerchantTransactions(ctx context.Context, req *biz.GetTransactionsRequest, page int64, size int64) (*biz.GetTransactionsReply, error) {
	transactions, err := b.data.db.GetMerchantTransactions(ctx, models.GetMerchantTransactionsParams{
		UserID:   req.UserId,
		Currency: req.Currency,
		Status:   string(req.PaymentStatus),
		Page:     page,
		PageSize: size,
	})
	if err != nil {
		return nil, err
	}

	transactionsReply := make([]*biz.Transaction, 0, len(transactions))
	// maps.Copy(transactionsReply, transactions)
	for _, t := range transactions {
		amount, err := types.NumericToFloat(t.Amount)
		if err != nil {
			return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to float64 failed")
		}
		transactionsReply = append(transactionsReply, &biz.Transaction{
			Id:                t.ID,
			Type:              constants.TransactionType(t.Type),
			Amount:            amount,
			Currency:          t.Currency,
			FromUserId:        t.FromUserID,
			ToMerchantId:      t.ToMerchantID,
			PaymentMethodType: constants.PaymentMethod(t.PaymentMethodType),
			PaymentAccount:    t.PaymentAccount,
			Status:            constants.PaymentStatus(t.Status),
			PaymentExtra:      t.PaymentExtra,
			CreatedAt:         t.CreatedAt,
			UpdatedAt:         t.UpdatedAt,
		})
	}

	return &biz.GetTransactionsReply{
		Transactions: transactionsReply,
	}, nil
}
