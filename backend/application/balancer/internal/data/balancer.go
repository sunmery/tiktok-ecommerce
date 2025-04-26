package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	merchantorderv1 "backend/api/merchant/order/v1"

	globalPkg "backend/pkg"

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

type balancerRepo struct {
	data *Data
	log  *log.Helper
}

func NewBalancerRepo(data *Data, logger log.Logger) biz.BalancerRepo {
	return &balancerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (b balancerRepo) GetTransactions(ctx context.Context, req *biz.GetTransactionsRequest) (*biz.GetTransactionsReply, error) {
	page := (req.Page - 1) * req.PageSize
	pageSize := req.PageSize

	role, err := globalPkg.GetMetadataRole(ctx)
	if err != nil {
		return nil, err
	}

	transactions := make([]models.BalancesTransactions, 0, pageSize)
	switch role {
	case constants.Consumer:
		transactions, err = b.data.db.GetConsumerTransactions(ctx, models.GetConsumerTransactionsParams{
			UserID:   req.UserId,
			Currency: req.Currency,
			Status:   string(req.PaymentStatus),
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			return nil, err
		}
	case constants.Merchant:
		transactions, err = b.data.db.GetMerchantTransactions(ctx, models.GetMerchantTransactionsParams{
			MerchantID: req.UserId,
			Currency:   req.Currency,
			Status:     string(req.PaymentStatus),
			Page:       page,
			PageSize:   pageSize,
		})
	}

	bizTransactions := make([]*biz.Transactions, 0, len(transactions))
	for _, t := range transactions {
		amount, err := types.NumericToFloat(t.Amount)
		if err != nil {
			return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to float64 failed")
		}

		fromUserID, err := uuid.Parse(t.FromUserID.String())
		if err != nil {
			return nil, kerrors.New(500, "PARSE_USER_ID_FAILED", "parse user id failed")
		}

		toMerchantID, err := uuid.Parse(t.ToMerchantID.String())
		if err != nil {
			return nil, kerrors.New(500, "PARSE_MERCHANT_ID_FAILED", "parse merchant id failed")
		}

		bizTransaction := &biz.Transactions{
			Id:                t.ID,
			Type:              constants.TransactionType(t.Type),
			Amount:            amount,
			Currency:          t.Currency,
			FromUserId:        fromUserID,
			ToMerchantId:      toMerchantID,
			PaymentMethodType: constants.PaymentMethod(t.PaymentMethodType),
			PaymentAccount:    t.PaymentAccount,
			PaymentExtra:      t.PaymentExtra,
			Status:            constants.PaymentStatus(t.Status),
			CreatedAt:         t.CreatedAt,
			UpdatedAt:         t.UpdatedAt,
		}

		bizTransactions = append(bizTransactions, bizTransaction)
	}

	return &biz.GetTransactionsReply{
		Transactions: bizTransactions,
	}, nil
}

// FreezeBalance 冻结余额
func (b balancerRepo) FreezeBalance(ctx context.Context, req *biz.FreezeBalanceRequest) (*biz.FreezeBalanceReply, error) {
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
func (b balancerRepo) CancelFreeze(ctx context.Context, req *biz.CancelFreezeRequest) (*biz.CancelFreezeReply, error) {
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

func (b balancerRepo) GetMerchantBalance(ctx context.Context, req *biz.GetMerchantBalanceRequest) (*biz.BalanceReply, error) {
	balance, err := b.data.db.GetMerchantBalance(ctx, models.GetMerchantBalanceParams{
		MerchantID: req.MerchantId,
		Currency:   string(req.Currency),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "BALANCE_NOT_FOUND", "balance record not found")
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

func (b balancerRepo) CreateConsumerBalance(ctx context.Context, req *biz.CreateConsumerBalanceRequest) (*biz.CreateConsumerBalanceReply, error) {
	tx := b.data.DB(ctx)

	initialBalance, err := types.Float64ToNumeric(req.InitialBalance)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	CreateUserPaymentMethodsErr := tx.CreateConsumerPaymentMethods(ctx, models.CreateConsumerPaymentMethodsParams{
		ID:             id.SnowflakeID(),
		UserID:         req.UserId,
		Type:           req.BalancerType,
		IsDefault:      req.IsDefault,
		AccountDetails: req.AccountDetails,
	})
	if CreateUserPaymentMethodsErr != nil {
		return nil, err
	}

	reply, err := tx.CreateConsumerBalance(ctx, models.CreateConsumerBalanceParams{
		UserID:    req.UserId,
		Currency:  string(req.Currency),
		Available: initialBalance,
	})
	if err != nil {
		return nil, err
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

func (b balancerRepo) CreateMerchantBalance(ctx context.Context, req *biz.CreateMerchantBalanceRequest) (*biz.CreateMerchantBalanceReply, error) {
	tx := b.data.DB(ctx)

	initialBalance, err := types.Float64ToNumeric(req.InitialBalance)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	CreateUserPaymentMethodsErr := tx.CreateMerchantPaymentMethods(ctx, models.CreateMerchantPaymentMethodsParams{
		ID:             id.SnowflakeID(),
		MerchantID:     req.MerchantId,
		Type:           req.BalancerType,
		IsDefault:      req.IsDefault,
		AccountDetails: req.AccountDetails,
	})
	if CreateUserPaymentMethodsErr != nil {
		return nil, err
	}

	reply, err := tx.CreateMerchantBalance(ctx, models.CreateMerchantBalanceParams{
		MerchantID: req.MerchantId,
		Currency:   string(req.Currency),
		Available:  initialBalance,
	})
	if err != nil {
		return nil, err
	}

	available, err := types.NumericToFloat(reply.Available)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AVAILABLE_FAILED", "convert available to float64 failed")
	}

	return &biz.CreateMerchantBalanceReply{
		UserId:    reply.MerchantID,
		Currency:  constants.Currency(reply.Currency),
		Available: available,
	}, nil
}

func (b balancerRepo) GetUserBalance(ctx context.Context, req *biz.GetUserBalanceRequest) (*biz.BalanceReply, error) {
	balance, err := b.data.db.GetUserBalance(ctx, models.GetUserBalanceParams{
		UserID:   req.UserId,
		Currency: string(req.Currency),
	})
	if err != nil {
		return nil, err
	}
	available, err := types.NumericToFloat(balance.Available)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AVAILABLE_FAILED", "convert available to float64 failed")
	}
	frozen, err := types.NumericToFloat(balance.Frozen)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "BALANCE_NOT_FOUND", "balance record not found")
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

func (b balancerRepo) RechargeBalance(ctx context.Context, req *biz.RechargeBalanceRequest) (*biz.RechargeBalanceReply, error) {
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

	// 3. 创建交易记录
	merchantId := req.MerchantId

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
		ID:                id.SnowflakeID(),
		Type:              string(constants.TransactionRecharge),
		Amount:            amount,
		Currency:          string(req.Currency),
		FromUserID:        req.UserId,
		ToMerchantID:      merchantId,
		PaymentMethodType: req.PaymentMethodType,
		PaymentAccount:    req.PaymentAccount,
		PaymentExtra:      paymentExtraJson,
		Status:            "PAID", // 充值通常是已支付状态
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

// WithdrawBalance 用户提现
func (b balancerRepo) WithdrawBalance(ctx context.Context, req *biz.WithdrawBalanceRequest) (*biz.WithdrawBalanceReply, error) {
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
func (b balancerRepo) ConfirmTransfer(ctx context.Context, req *biz.ConfirmTransferRequest) (*biz.ConfirmTransferReply, error) {
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
	rows, ConfirmUserFreezeErr := tx.ConfirmUserFreeze(ctx, models.ConfirmUserFreezeParams{
		UserID:          freeze.UserID,
		Currency:        freeze.Currency,
		Amount:          freeze.Amount,
		ExpectedVersion: req.ExpectedUserVersion,
	})
	if ConfirmUserFreezeErr != nil {
		return nil, kerrors.New(500, "CONFIRM_FREEZE_FAILED", fmt.Sprintf("confirm freeze failed: %v", ConfirmUserFreezeErr.Error()))
	}
	if rows == 0 {
		return nil, kerrors.New(409, "USER_OPTIMISTIC_LOCK_FAILED", "user balance version mismatch")
	}

	// 6. 获取商家ID（从订单ID）
	merchantId, err := b.getMerchantIDFromOrder(ctx, freeze.OrderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "MERCHANT_NOT_FOUND", "merchant not found for order")
		}
		return nil, kerrors.New(500, "GET_MERCHANT_ID_FAILED", fmt.Sprintf("get merchant id from order failed:%v", err))
	}

	// 7. 增加商家可用余额
	rows, err = tx.IncreaseMerchantAvailableBalance(ctx, models.IncreaseMerchantAvailableBalanceParams{
		MerchantID:      merchantId,
		Currency:        freeze.Currency,
		Amount:          freeze.Amount,
		ExpectedVersion: req.ExpectedMerchantVersion,
	})
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
		ToMerchantID:      req.MerchantId,
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
func (b balancerRepo) getMerchantIDFromOrder(ctx context.Context, orderId int64) (uuid.UUID, error) {
	// 这里可以根据订单ID查询数据库获取商家ID
	reply, err := b.data.merchantOrderv1.GetMerchantByOrderId(ctx, &merchantorderv1.GetMerchantByOrderIdReq{
		OrderId: orderId,
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
