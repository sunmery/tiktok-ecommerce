package data

import (
	"context"
	"encoding/json"
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

	o.log.WithContext(ctx).Infof("获取商户 %s 的订单列表, 页码: %d, 每页数量: %d", req.UserID, req.Page, req.PageSize)

	// 查询订单列表
	userId := types.ToPgUUID(req.UserID)
	pageSize := int64(req.PageSize)
	page := int64((req.Page - 1) * req.PageSize)
	orders, err := o.data.db.ListOrdersByUser(ctx, models.ListOrdersByUserParams{
		MerchantID: userId,
		PageSize:   &pageSize,
		Page:       &page,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("获取订单列表失败: %v", err)
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	if len(orders) == 0 {
		o.log.WithContext(ctx).Infof("商户 %s 没有订单记录", req.UserID)
		return &biz.GetMerchantOrdersReply{Orders: []*biz.SubOrder{}}, nil
	}

	// 创建一个映射以存储订单ID与子订单的关系
	var respOrders []*biz.SubOrder
	var orderIDs []int64

	// 收集所有主订单ID
	for _, order := range orders {
		orderIDs = append(orderIDs, order.OrderID)
	}

	// 去重订单ID
	uniqueOrderIDs := make(map[int64]bool)
	var uniqueIDs []int64
	for _, id := range orderIDs {
		if !uniqueOrderIDs[id] {
			uniqueOrderIDs[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	// 对每个唯一的订单ID获取子订单
	for _, orderID := range uniqueIDs {
		// 使用getSubOrders函数获取子订单
		subOrders, err := o.getSubOrders(ctx, orderID)
		if err != nil {
			o.log.WithContext(ctx).Errorf("获取订单 %d 的子订单失败: %v", orderID, err)
			// 继续处理其他订单，不因为一个订单失败而中断整个流程
			continue
		}

		// 添加子订单到结果
		for _, subOrder := range subOrders {
			// 查找对应的原始订单以获取支付状态
			paymentStatus := biz.PaymentPending
			for _, order := range orders {
				if order.ID == subOrder.ID {
					if order.Status != "" {
						paymentStatus = biz.PaymentStatus(order.Status)
					}
					break
				}
			}

			// 更新子订单的支付状态
			subOrder.Status = string(paymentStatus)
			respOrders = append(respOrders, subOrder)
		}
	}

	o.log.WithContext(ctx).Debugf("获取到 %d 个商户订单", len(respOrders))
	return &biz.GetMerchantOrdersReply{Orders: respOrders}, nil
}

// 获取订单的子订单信息
func (o *orderRepo) getSubOrders(ctx context.Context, orderID int64) ([]*biz.SubOrder, error) {
	// 使用父上下文创建子上下文，保留链路追踪信息
	subCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	o.log.WithContext(ctx).Debugf("查询订单 %d 的子订单信息", orderID)

	// 查询子订单
	rows, err := o.data.db.QuerySubOrders(subCtx, &orderID)
	if err != nil {
		// 检查是否是上下文取消或超时导致的错误
		if subCtx.Err() != nil {
			o.log.WithContext(ctx).Warnf("查询子订单时上下文取消: %v, 订单ID: %d", subCtx.Err(), orderID)
			return nil, fmt.Errorf("查询子订单超时: %w", subCtx.Err())
		}
		o.log.WithContext(ctx).Errorf("查询子订单失败: %v, 订单ID: %d", err, orderID)
		return nil, fmt.Errorf("查询子订单失败: %w", err)
	}

	if len(rows) == 0 {
		o.log.WithContext(ctx).Infof("订单 %d 没有子订单", orderID)
		return []*biz.SubOrder{}, nil
	}

	var subOrders []*biz.SubOrder
	for _, order := range rows {
		// 解析订单项 - 先解析为SubOrderItem结构
		type SubOrderItem struct {
			Item *biz.CartItem `json:"item"`
			Cost float64       `json:"cost"`
		}

		var subOrderItems []SubOrderItem
		if err := json.Unmarshal(order.Items, &subOrderItems); err != nil {
			o.log.WithContext(ctx).Errorf("解析子订单项失败: %v, 订单ID: %d, 子订单ID: %d", err, orderID, order.ID)

			// 尝试记录原始数据用于调试
			o.log.WithContext(ctx).Debugf("原始子订单数据: %s", string(order.Items))

			// 返回解析错误
			return nil, fmt.Errorf("解析子订单项失败: %w", err)
		}

		// 转换为biz.OrderItem
		var orderItems []*biz.OrderItem
		for _, item := range subOrderItems {
			// 验证商品ID和商家ID的有效性
			if item.Item == nil {
				o.log.WithContext(ctx).Warnf("子订单项缺少商品信息, 跳过此项, 订单ID: %d, 子订单ID: %d", orderID, order.ID)
				continue
			}

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
		var amount float64
		switch v := order.TotalAmount.(type) {
		case pgtype.Numeric:
			convertedAmount, err := types.NumericToFloat(v)
			if err != nil {
				o.log.WithContext(ctx).Errorf("转换金额失败: %v, 订单ID: %d, 子订单ID: %d", err, orderID, order.ID)
				return nil, fmt.Errorf("转换金额失败: %w", err)
			}
			amount = convertedAmount
		case float64:
			amount = v
		default:
			o.log.WithContext(ctx).Errorf("未知的金额类型: %T, 订单ID: %d, 子订单ID: %d", order.TotalAmount, orderID, order.ID)
			return nil, fmt.Errorf("未知的金额类型: %T", order.TotalAmount)
		}

		// 添加子订单到结果集
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

	o.log.WithContext(ctx).Debugf("获取到 %d 个子订单, 订单ID: %d", len(subOrders), orderID)
	return subOrders, nil
}
