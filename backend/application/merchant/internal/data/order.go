package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"backend/pkg/types"

	"backend/application/merchant/internal/data/models"

	"backend/application/merchant/internal/biz"
)

func (o *orderRepo) GetMerchantOrders(ctx context.Context, req *biz.GetMerchantOrdersReq) (*biz.GetMerchantOrdersReply, error) {
	// 设置默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	// 限制最大页面大小
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	o.log.WithContext(ctx).Infof("Listing orders for user %s, page %d, page size %d", req.UserID, req.Page, req.PageSize)

	// 查询订单列表
	userId := types.ToPgUUID(req.UserID)
	pageSize := int64(req.PageSize)
	page := int64((req.Page - 1) * req.PageSize)
	orders, err := o.data.db.ListOrdersByUser(ctx, models.ListOrdersByUserParams{
		UserID:   userId,
		PageSize: &pageSize,
		Page:     &page,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to list orders: %v", err)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// 构建响应
	var respOrders []*biz.Order
	for _, order := range orders {
		// 构建地址信息
		address := &biz.Address{
			StreetAddress: order.StreetAddress,
			City:          order.City,
			State:         order.State,
			Country:       order.Country,
			ZipCode:       order.ZipCode,
		}

		// 获取子订单信息
		var subOrders []*biz.SubOrder
		var err error

		// 先检查上下文是否已取消
		select {
		case <-ctx.Done():
			o.log.WithContext(ctx).Warnf("Context canceled before fetching sub orders for order %d", order.ID)
			// 上下文已取消，不执行查询，继续处理其他字段
		default:
			// 上下文未取消，执行查询
			subOrders, err = o.getSubOrders(ctx, order.ID)
			if err != nil {
				if errors.Is(ctx.Err(), context.Canceled) {
					o.log.WithContext(ctx).Warnf("Context canceled during fetching sub orders for order %d", order.ID)
				} else {
					o.log.WithContext(ctx).Errorf("Failed to get sub orders for order %d: %v", order.ID, err)
				}
				// 错误发生时继续处理，不中断整个列表查询
			}
		}

		// 解析支付状态
		paymentStatus := biz.PaymentPending
		if order.PaymentStatus != "" {
			paymentStatus = biz.PaymentStatus(order.PaymentStatus)
		}

		respOrders = append(respOrders, &biz.Order{
			OrderID:       order.ID,
			UserID:        order.UserID,
			Currency:      order.Currency,
			Address:       address,
			Email:         order.Email,
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
			SubOrders:     subOrders,
			PaymentStatus: paymentStatus,
		})
	}

	return &biz.GetMerchantOrdersReply{Orders: respOrders}, nil
}

// 获取订单的子订单信息
func (o *orderRepo) getSubOrders(ctx context.Context, orderID int64) ([]*biz.SubOrder, error) {
	// 创建独立上下文，设置合理超时（如5秒）
	subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查询子订单
	rows, err := o.data.db.QuerySubOrders(subCtx, &orderID)
	if err != nil {
		// 检查是否是上下文取消导致的错误
		if ctx.Err() != nil {
			o.log.WithContext(ctx).Warnf("Context canceled during database query for order %d", orderID)
			return nil, fmt.Errorf("failed to query sub orders: %w", ctx.Err())
		}
		return nil, fmt.Errorf("failed to query sub orders: %w", err)
	}

	var subOrders []*biz.SubOrder
	for _, order := range rows {
		// 解析订单项 - 先解析为SubOrderItem结构
		type SubOrderItem struct {
			Item *biz.CartItem
			Cost float64
		}
		var subOrderItems []SubOrderItem
		if err := json.Unmarshal(order.Items, &subOrderItems); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sub order items: %w", err)
		}

		// 转换为biz.OrderItem
		var orderItems []*biz.OrderItem
		for _, item := range subOrderItems {
			// 确保CartItem中的MerchantId和ProductId正确映射
			cartItem := &biz.CartItem{
				MerchantId: item.Item.MerchantId,
				ProductId:  item.Item.ProductId,
				Quantity:   item.Item.Quantity,
			}

			orderItems = append(orderItems, &biz.OrderItem{
				Item: cartItem,
				Cost: item.Cost,
			})
		}

		// 转换金额
		amount, err := types.NumericToFloat(order.TotalAmount.(pgtype.Numeric))
		if err != nil {
			return nil, fmt.Errorf("failed to convert amount: %w", err)
		}

		subOrders = append(subOrders, &biz.SubOrder{
			ID:          order.ID,
			MerchantID:  order.MerchantID,
			TotalAmount: amount,
			Currency:    order.Currency,
			Status:      order.Status,
			Items:       orderItems,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		})
	}

	return subOrders, nil
}
