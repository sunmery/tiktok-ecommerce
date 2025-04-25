package data

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"

	"backend/constants"

	"github.com/google/uuid"

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

func (b balancerRepo) FreezeBalance(ctx context.Context, req *biz.FreezeBalanceRequest) (*biz.FreezeBalanceResponse, error) {
	// 1. 开始事务
	tx := b.data.DB(ctx)

	// 2. 冻结用户余额
	amount, err := types.Float64ToNumeric(req.Amount)
	if err != nil {
		return nil, kerrors.New(500, "CONVERT_AMOUNT_FAILED", "convert amount to numeric failed")
	}

	// 检查是否已存在相同订单的冻结记录（幂等性检查）
	orderUUID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, kerrors.New(400, "INVALID_ORDER_ID", "invalid order id format")
	}

	// 尝试获取已存在的冻结记录
	existingFreeze, err := tx.GetFreezeByOrderForUser(ctx, models.GetFreezeByOrderForUserParams{
		UserID:  req.UserId,
		OrderID: orderUUID,
	})

	// 如果已存在冻结记录，直接返回
	if err == nil {
		// 已存在冻结记录，返回现有记录ID
		return &biz.FreezeBalanceResponse{
			FreezeId:   strconv.FormatInt(existingFreeze.ID, 10),
			NewVersion: req.ExpectedVersion, // 保持版本不变，因为没有执行冻结操作
		}, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		// 如果是其他错误，则返回错误
		return nil, kerrors.New(500, "GET_FREEZE_FAILED", "failed to check existing freeze")
	}

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
		return nil, kerrors.New(409, "OPTIMISTIC_LOCK_FAILED", "balance version mismatch or insufficient funds")
	}

	// 3. 创建冻结记录
	freezeId, err := tx.CreateFreeze(ctx, models.CreateFreezeParams{
		UserID:   req.UserId,
		OrderID:  orderUUID,
		Currency: string(req.Currency),
		Amount:   amount,
		// ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		return nil, kerrors.New(500, "CREATE_FREEZE_FAILED", "create freeze record failed")
	}

	return &biz.FreezeBalanceResponse{
		FreezeId:   strconv.FormatInt(freezeId, 10),
		NewVersion: req.ExpectedVersion + 1, // 版本号+1
	}, nil
}

func (b balancerRepo) CancelFreeze(ctx context.Context, req *biz.CancelFreezeRequest) (*biz.CancelFreezeResponse, error) {
	// 1. 开始事务
	tx := b.data.DB(ctx)

	// 2. 获取冻结记录
	freezeId, err := strconv.ParseInt(req.FreezeId, 10, 64)
	if err != nil {
		return nil, kerrors.New(400, "INVALID_FREEZE_ID", "invalid freeze id format")
	}

	freeze, err := tx.GetFreeze(ctx, freezeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "FREEZE_NOT_FOUND", "freeze record not found")
		}
		return nil, kerrors.New(500, "GET_FREEZE_FAILED", "get freeze record failed")
	}

	// 3. 检查冻结状态
	if freeze.Status != "FROZEN" {
		return nil, kerrors.New(400, "INVALID_FREEZE_STATUS", "freeze is not in FROZEN status")
	}

	// 4. 更新冻结记录状态为取消
	rows, err := tx.UpdateFreezeStatus(ctx, models.UpdateFreezeStatusParams{
		Status:        "CANCELED",
		ID:            freezeId,
		CurrentStatus: "FROZEN",
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

	return &biz.CancelFreezeResponse{
		Success:    true,
		NewVersion: req.ExpectedVersion + 1, // 版本号+1
	}, nil
}

func (b balancerRepo) GetMerchantBalance(ctx context.Context, req *biz.GetMerchantBalanceRequest) (*biz.BalanceResponse, error) {
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

	return &biz.BalanceResponse{
		Available: available,
		Frozen:    0, // 商家账户没有冻结余额
		Currency:  constants.Currency(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (b balancerRepo) GetUserBalance(ctx context.Context, req *biz.GetUserBalanceRequest) (*biz.BalanceResponse, error) {
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

	return &biz.BalanceResponse{
		Available: available,
		Frozen:    frozen,
		Currency:  constants.Currency(balance.Currency),
		Version:   balance.Version,
	}, nil
}

func (b balancerRepo) RechargeBalance(ctx context.Context, req *biz.RechargeBalanceRequest) (*biz.RechargeBalanceResponse, error) {
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
	// 为了简化，我们使用一个假的商家ID（实际应用中可能是平台账户）
	merchantId := uuid.New()

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
		Type:              "RECHARGE",
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
		return nil, kerrors.New(500, "CREATE_TRANSACTION_FAILED", "create transaction record failed")
	}

	return &biz.RechargeBalanceResponse{
		Success:       true,
		TransactionId: strconv.FormatInt(transactionId, 10),
		NewVersion:    req.ExpectedVersion + 1, // 版本号+1
	}, nil
}

func (b balancerRepo) WithdrawBalance(ctx context.Context, req *biz.WithdrawBalanceRequest) (*biz.WithdrawBalanceResponse, error) {
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
		return nil, kerrors.New(409, "OPTIMISTIC_LOCK_FAILED", "balance version mismatch or insufficient funds")
	}

	// 4. 创建交易记录
	// 为了简化，我们使用一个假的商家ID（实际应用中可能是平台账户）
	merchantId := uuid.New()

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
		Type:              "WITHDRAW",
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
		return nil, kerrors.New(500, "CREATE_TRANSACTION_FAILED", "create transaction record failed")
	}

	return &biz.WithdrawBalanceResponse{
		Success:       true,
		TransactionId: strconv.FormatInt(transactionId, 10),
		NewVersion:    req.ExpectedVersion + 1, // 版本号+1
	}, nil
}

func (b balancerRepo) ConfirmTransfer(ctx context.Context, req *biz.ConfirmTransferRequest) (*biz.ConfirmTransferResponse, error) {
	// 1. 开始事务
	tx := b.data.DB(ctx)

	// 2. 获取冻结记录
	freezeId, err := strconv.ParseInt(req.FreezeId, 10, 64)
	if err != nil {
		return nil, kerrors.New(400, "INVALID_FREEZE_ID", "invalid freeze id format")
	}

	freeze, err := tx.GetFreeze(ctx, freezeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "FREEZE_NOT_FOUND", "freeze record not found")
		}
		return nil, kerrors.New(500, "GET_FREEZE_FAILED", "get freeze record failed")
	}

	// 3. 检查冻结状态
	if freeze.Status != "FROZEN" {
		return nil, kerrors.New(400, "INVALID_FREEZE_STATUS", "freeze is not in FROZEN status")
	}

	// 4. 更新冻结记录状态为确认
	rows, err := tx.UpdateFreezeStatus(ctx, models.UpdateFreezeStatusParams{
		Status:        "CONFIRMED",
		ID:            freezeId,
		CurrentStatus: "FROZEN",
	})
	if err != nil {
		return nil, kerrors.New(500, "UPDATE_FREEZE_STATUS_FAILED", "update freeze status failed")
	}
	if rows == 0 {
		return nil, kerrors.New(409, "FREEZE_STATUS_CHANGED", "freeze status has been changed")
	}

	// 5. 确认用户冻结（减少冻结余额）
	rows, err = tx.ConfirmUserFreeze(ctx, models.ConfirmUserFreezeParams{
		UserID:          freeze.UserID,
		Currency:        freeze.Currency,
		Amount:          freeze.Amount,
		ExpectedVersion: req.ExpectedUserVersion,
	})
	if err != nil {
		return nil, kerrors.New(500, "CONFIRM_FREEZE_FAILED", "confirm freeze failed")
	}
	if rows == 0 {
		return nil, kerrors.New(409, "USER_OPTIMISTIC_LOCK_FAILED", "user balance version mismatch")
	}

	// 6. 获取商家ID（从订单ID）
	orderUUID, err := uuid.Parse(freeze.OrderID.String())
	if err != nil {
		return nil, kerrors.New(500, "INVALID_ORDER_ID", "invalid order id format")
	}

	merchantId, err := b.getMerchantIDFromOrder(ctx, orderUUID)
	if err != nil {
		return nil, kerrors.New(500, "GET_MERCHANT_ID_FAILED", "get merchant id from order failed")
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
		"order_id":        freeze.OrderID.String(),
	}
	paymentExtraJson, err := json.Marshal(paymentExtra)
	if err != nil {
		return nil, kerrors.New(500, "JSON_MARSHAL_FAILED", "marshal payment extra failed")
	}

	// 创建交易记录
	transactionId, err := tx.CreateTransaction(ctx, models.CreateTransactionParams{
		Type:              "PAYMENT",
		Amount:            freeze.Amount,
		Currency:          freeze.Currency,
		FromUserID:        freeze.UserID,
		ToMerchantID:      merchantId,
		PaymentMethodType: "BALANCE", // 使用余额支付
		PaymentAccount:    "internal_balance",
		PaymentExtra:      paymentExtraJson,
		Status:            "PAID", // 已支付状态
	})
	if err != nil {
		return nil, kerrors.New(500, "CREATE_TRANSACTION_FAILED", "create transaction record failed")
	}

	return &biz.ConfirmTransferResponse{
		Success:            true,
		TransactionId:      strconv.FormatInt(transactionId, 10),
		NewUserVersion:     req.ExpectedUserVersion + 1,     // 用户余额版本号+1
		NewMerchantVersion: req.ExpectedMerchantVersion + 1, // 商家余额版本号+1
	}, nil
}

// 辅助方法：从订单ID获取商家ID
func (b balancerRepo) getMerchantIDFromOrder(ctx context.Context, orderID uuid.UUID) (uuid.UUID, error) {
	// 实际实现中，可能需要调用订单服务获取商家ID
	// 这里简化处理，返回一个假的商家ID
	// 在实际项目中，应该替换为真实的实现
	return uuid.New(), nil
}
