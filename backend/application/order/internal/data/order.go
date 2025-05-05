package data

import (
	"context"
	"fmt"

	"backend/constants"

	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"github.com/go-kratos/kratos/v2/log"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (o *orderRepo) MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error) {
	o.log.WithContext(ctx).Infof("Marking order %d as paid for user %s", req.OrderId, req.UserId)
	tx := o.data.db

	// 获取订单信息，确认订单存在且属于该用户，使用FOR UPDATE锁定行
	updatePaymentStatusResult, err := o.data.db.UpdatePaymentStatus(ctx, req.OrderId)
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to get order %d: %v", req.OrderId, err)
		return nil, fmt.Errorf("updatePaymentStatusResult failed to get order: %w", err)
	}
	fmt.Printf("updatePaymentStatusResult: %#v", updatePaymentStatusResult)

	// 验证订单所有者
	log.Debugf("r%+v, req:%+v", updatePaymentStatusResult.UserID.String(), req.UserId.String())
	if updatePaymentStatusResult.UserID.String() != req.UserId.String() {
		o.log.WithContext(ctx).Warnf("Order %d does not belong to user %s, order.UserID:%s", req.OrderId, req.UserId.String(), updatePaymentStatusResult.UserID.String())
		return nil, fmt.Errorf("order does not belong to user")
	}

	// 检查订单当前支付状态
	if updatePaymentStatusResult.PaymentStatus == string(constants.PaymentPaid) {
		// 订单已经是已支付状态，直接返回成功
		o.log.WithContext(ctx).Infof("Order %d is already marked as paid", req.OrderId)
		return &biz.MarkOrderPaidResp{}, nil
	}

	o.log.WithContext(ctx).Infof("Updating order %d payment status from %s to %s",
		req.OrderId, updatePaymentStatusResult.PaymentStatus, string(constants.PaymentPaid))

	// 更新订单支付状态为已支付
	_, markOrderAsPaidErr := tx.MarkOrderAsPaid(ctx, models.MarkOrderAsPaidParams{
		PaymentStatus: string(constants.PaymentPaid),
		ID:            req.OrderId,
	})
	if markOrderAsPaidErr != nil {
		o.log.WithContext(ctx).Errorf("Failed to update order payment status: %v", markOrderAsPaidErr)
		return nil, fmt.Errorf("failed to update order payment status: %w", markOrderAsPaidErr)
	}

	// 更新子订单状态
	_, markSubOrderAsPaidErr := tx.MarkSubOrderAsPaid(ctx, models.MarkSubOrderAsPaidParams{
		Status:  string(constants.PaymentPaid),
		OrderID: req.OrderId,
	})
	if markSubOrderAsPaidErr != nil {
		o.log.WithContext(ctx).Errorf("Failed to update sub orders payment status: %v", markSubOrderAsPaidErr)
		return nil, fmt.Errorf("failed to update sub orders payment status: %w", markSubOrderAsPaidErr)
	}

	o.log.WithContext(ctx).Infof("Successfully marked order %d as paid for user %s", req.OrderId, req.UserId)
	return &biz.MarkOrderPaidResp{}, nil
}
