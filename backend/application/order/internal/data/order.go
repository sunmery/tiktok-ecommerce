package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"backend/application/payment/pkg"
	"backend/pkg/types"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func (o *orderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error) {
	// 生成雪花ID
	orderID := pkg.SnowflakeID()

	order, err := o.data.db.CreateOrder(ctx, models.CreateOrderParams{
		ID:            orderID,
		UserID:        req.UserId,
		Currency:      req.Currency,
		StreetAddress: req.Address.StreetAddress,
		City:          req.Address.City,
		State:         req.Address.State,
		Country:       req.Address.Country,
		ZipCode:       req.Address.ZipCode,
		Email:         req.Email,
	})
	fmt.Printf("order: %v", order)
	if err != nil {
		return nil, fmt.Errorf("创建主订单失败: %w", err)
	}

	// 分单
	for _, item := range req.OrderItems {
		// 序列化订单项
		items := []biz.OrderItem{
			{
				Item: item.Item,
				Cost: item.Cost,
			},
		}
		itemsJSON, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("序列化订单项失败: %w", err)
		}

		// 转换价格到pgtype.Numeric
		totalAmount, totalAmountErr := types.Float64ToNumeric(item.Cost)
		if totalAmountErr != nil {
			return nil, fmt.Errorf("invalid price format: %w", totalAmountErr)
		}
		fmt.Printf("totalAmount: %v", totalAmount)

		// 创建子订单ID
		subOrderID := pkg.SnowflakeID()

		subOrder, subOrderErr := o.data.db.CreateSubOrder(ctx, models.CreateSubOrderParams{
			ID:          subOrderID,
			OrderID:     order.ID,
			MerchantID:  item.Item.MerchantId,
			TotalAmount: totalAmount,
			Currency:    req.Currency,
			Status:      "created",
			Items:       itemsJSON,
		})
		if subOrderErr != nil {
			return nil, fmt.Errorf("创建子订单失败: %w", subOrderErr)
		}
		fmt.Printf("subOrder: %v", subOrder)

	}

	return &biz.PlaceOrderResp{
		Order: &biz.OrderResult{
			OrderId: order.ID,
		},
	}, nil
}

func (o *orderRepo) ListOrder(ctx context.Context, req *biz.ListOrderReq) (*biz.ListOrderResp, error) {
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
	orders, err := o.data.db.ListOrdersByUser(ctx, models.ListOrdersByUserParams{
		UserID: req.UserID,
		Limit:  int64(req.PageSize),
		Offset: int64((req.Page - 1) * req.PageSize),
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

	o.log.WithContext(ctx).Infof("Listed %d orders for user %s", len(respOrders), req.UserID)
	return &biz.ListOrderResp{Orders: respOrders}, nil
}

// 自定义JSON解析
func parseSubOrders(data []byte) ([]*biz.SubOrder, error) {
	type dbSubOrder struct {
		ID          int64
		MerchantID  string
		TotalAmount float64
		Currency    string
		Status      string
		Items       json.RawMessage
	}

	var dbSubs []dbSubOrder
	if err := json.Unmarshal(data, &dbSubs); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	subOrders := make([]*biz.SubOrder, 0, len(dbSubs))
	for _, d := range dbSubs {
		merchantID, err := uuid.Parse(d.MerchantID)
		if err != nil {
			return nil, fmt.Errorf("invalid merchant id: %w", err)
		}

		var items []*biz.OrderItem
		if err := json.Unmarshal(d.Items, &items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		subOrders = append(subOrders, &biz.SubOrder{
			ID:          d.ID,
			MerchantID:  merchantID,
			TotalAmount: d.TotalAmount,
			Currency:    d.Currency,
			Status:      d.Status,
			Items:       items,
		})
	}

	return subOrders, nil
}

// 获取订单的子订单信息
func (o *orderRepo) getSubOrders(ctx context.Context, orderID int64) ([]*biz.SubOrder, error) {
	// 创建独立上下文，设置合理超时（如5秒）
	subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查询子订单
	rows, err := o.data.db.QuerySubOrders(subCtx, orderID)
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
		amount, err := types.NumericToFloat(order.TotalAmount)
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

func (o *orderRepo) MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error) {
	o.log.WithContext(ctx).Infof("Marking order %d as paid for user %s", req.OrderId, req.UserId)
	tx := o.data.DB(ctx)

	// 获取订单信息，确认订单存在且属于该用户，使用FOR UPDATE锁定行
	var order struct {
		ID            int64
		UserID        uuid.UUID
		PaymentStatus string
	}
	updatePaymentStatusResult, err := tx.UpdatePaymentStatus(ctx, req.OrderId)
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to get order %d: %v", req.OrderId, err)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	fmt.Printf("updatePaymentStatusResult: %#v", updatePaymentStatusResult)

	// 验证订单所有者
	if order.UserID != req.UserId {
		o.log.WithContext(ctx).Warnf("Order %d does not belong to user %s", req.OrderId, req.UserId)
		return nil, fmt.Errorf("order does not belong to user")
	}

	// 检查订单当前支付状态
	if order.PaymentStatus == string(biz.PaymentPaid) {
		// 订单已经是已支付状态，直接返回成功
		o.log.WithContext(ctx).Infof("Order %d is already marked as paid", req.OrderId)
		return &biz.MarkOrderPaidResp{}, nil
	}

	o.log.WithContext(ctx).Infof("Updating order %d payment status from %s to %s",
		req.OrderId, order.PaymentStatus, string(biz.PaymentPaid))

	// 更新订单支付状态为已支付
	markOrderAsPaidResult, err := tx.MarkOrderAsPaid(ctx, models.MarkOrderAsPaidParams{
		PaymentStatus: string(biz.PaymentPaid),
		ID:            req.OrderId,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to update order payment status: %v", err)
		return nil, fmt.Errorf("failed to update order payment status: %w", err)
	}

	log.Debugf("markOrderAsPaidResult: %#v", markOrderAsPaidResult)

	// 更新子订单状态
	markSubOrderAsPaidResult, err := tx.MarkSubOrderAsPaid(ctx, models.MarkSubOrderAsPaidParams{
		PaymentStatus: string(biz.PaymentPaid),
		OrderID:       req.OrderId,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to update sub orders payment status: %v", err)
		return nil, fmt.Errorf("failed to update sub orders payment status: %w", err)
	}
	log.Debugf("markSubOrderAsPaidResult: %#v", markSubOrderAsPaidResult)

	// 获取更新的子订单数量
	// rowsAffected := markSubOrderAsPaidResult.
	// 	o.log.WithContext(ctx).Infof("Updated %d sub orders for order %s", rowsAffected, req.OrderId)

	o.log.WithContext(ctx).Infof("Successfully marked order %d as paid for user %s", req.OrderId, req.UserId)
	return &biz.MarkOrderPaidResp{}, nil
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
