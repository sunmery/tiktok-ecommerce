package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	userv1 "backend/api/user/v1"

	cartv1 "backend/api/cart/v1"

	v1 "backend/api/order/v1"

	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"backend/application/order/internal/pkg"
	"backend/pkg/types"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	tx := o.data.DB(ctx)

	// 获取订单信息，确认订单存在且属于该用户，使用FOR UPDATE锁定行
	updatePaymentStatusResult, err := tx.UpdatePaymentStatus(ctx, req.OrderId)
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to get order %d: %v", req.OrderId, err)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	fmt.Printf("updatePaymentStatusResult: %#v", updatePaymentStatusResult)

	// 验证订单所有者
	if updatePaymentStatusResult.UserID.String() != req.UserId.String() {
		o.log.WithContext(ctx).Warnf("Order %d does not belong to user %s, order.UserID:%s", req.OrderId, req.UserId.String(), updatePaymentStatusResult.UserID.String())
		return nil, fmt.Errorf("order does not belong to user")
	}

	// 检查订单当前支付状态
	if updatePaymentStatusResult.PaymentStatus == string(biz.PaymentPaid) {
		// 订单已经是已支付状态，直接返回成功
		o.log.WithContext(ctx).Infof("Order %d is already marked as paid", req.OrderId)
		return &biz.MarkOrderPaidResp{}, nil
	}

	o.log.WithContext(ctx).Infof("Updating order %d payment status from %s to %s",
		req.OrderId, updatePaymentStatusResult.PaymentStatus, string(biz.PaymentPaid))

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

func (o *orderRepo) GetOrder(ctx context.Context, req *biz.GetOrderReq) (*v1.Order, error) {
	order, err := o.data.db.GetOrderByID(ctx, models.GetOrderByIDParams{
		UserID:  req.UserId,
		OrderID: req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	// 解析子订单JSON
	var subOrdersData []byte
	if order.SubOrders != nil {
		subOrdersData = order.SubOrders
	}

	// 解析子订单并转换为OrderItem
	subOrders, err := parseSubOrders(subOrdersData)
	if err != nil {
		o.log.WithContext(ctx).Warnf("解析订单 %d 的子订单失败: %v", order.ID, err)
	}

	// 将所有子订单的OrderItem合并到一个列表
	var allItems []*v1.OrderItem
	for _, subOrder := range subOrders {
		for _, item := range subOrder.Items {
			allItems = append(allItems, &v1.OrderItem{
				Item: &cartv1.CartItem{
					MerchantId: item.Item.MerchantId.String(),
					ProductId:  item.Item.ProductId.String(),
					Quantity:   item.Item.Quantity,
					Name:       item.Item.Name,
					Picture:    item.Item.Picture,
				},
				Cost: item.Cost,
			})
		}
	}

	// 构建主订单
	orderProto := &v1.Order{
		OrderId:  order.ID,
		UserId:   order.UserID.String(),
		Currency: order.Currency,
		Address: &userv1.Address{
			StreetAddress: order.StreetAddress,
			City:          order.City,
			State:         order.State,
			Country:       order.Country,
			ZipCode:       order.ZipCode,
		},
		Email:     order.Email,
		CreatedAt: timestamppb.New(order.CreatedAt),
		Items:     allItems,
	}

	// 设置支付状态（如果有）
	if order.PaymentStatus != "" {
		switch order.PaymentStatus {
		case string(biz.PaymentPending):
			orderProto.PaymentStatus = v1.PaymentStatus_NOT_PAID
		case string(biz.PaymentProcessing):
			orderProto.PaymentStatus = v1.PaymentStatus_PROCESSING
		case string(biz.PaymentPaid):
			orderProto.PaymentStatus = v1.PaymentStatus_PAID
		case string(biz.PaymentFailed):
			orderProto.PaymentStatus = v1.PaymentStatus_FAILED
		case string(biz.PaymentCancelled):
			orderProto.PaymentStatus = v1.PaymentStatus_CANCELLED
		default:
			orderProto.PaymentStatus = v1.PaymentStatus_NOT_PAID
		}
	} else {
		orderProto.PaymentStatus = v1.PaymentStatus_NOT_PAID
	}

	return orderProto, nil
}

func (o *orderRepo) GetConsumerOrders(ctx context.Context, req *biz.GetConsumerOrdersReq) (*biz.Orders, error) {
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

	consumerOrders, err := o.data.db.GetConsumerOrders(ctx, models.GetConsumerOrdersParams{
		UserID:   req.UserId,
		PageSize: int64(req.PageSize),
		Page:     int64((req.Page - 1) * req.PageSize),
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("获取用户订单列表失败: %v", err)
		return nil, fmt.Errorf("获取用户订单列表失败: %w", err)
	}

	orders := make([]*v1.Order, 0, len(consumerOrders))
	for _, order := range consumerOrders {
		// 解析子订单JSON
		var subOrdersData []byte
		if order.SubOrders != nil {
			subOrdersData = order.SubOrders
		}

		// 解析子订单并转换为OrderItem
		subOrders, err := parseSubOrders(subOrdersData)
		if err != nil {
			o.log.WithContext(ctx).Warnf("解析订单 %d 的子订单失败: %v", order.ID, err)
			// 继续处理其他订单，不中断
			continue
		}

		// 将所有子订单的OrderItem合并到一个列表
		var allItems []*v1.OrderItem
		for _, subOrder := range subOrders {
			for _, item := range subOrder.Items {
				allItems = append(allItems, &v1.OrderItem{
					Item: &cartv1.CartItem{
						MerchantId: item.Item.MerchantId.String(),
						ProductId:  item.Item.ProductId.String(),
						Quantity:   item.Item.Quantity,
						Name:       item.Item.Name,
						Picture:    item.Item.Picture,
					},
					Cost: item.Cost,
				})
			}
		}

		// 构建主订单
		orderProto := &v1.Order{
			OrderId:  order.ID,
			UserId:   order.UserID.String(),
			Currency: order.Currency,
			Address: &userv1.Address{
				StreetAddress: order.StreetAddress,
				City:          order.City,
				State:         order.State,
				Country:       order.Country,
				ZipCode:       order.ZipCode,
			},
			Email:     order.Email,
			CreatedAt: timestamppb.New(order.CreatedAt),
			Items:     allItems,
		}

		// 设置支付状态（如果有）
		if order.PaymentStatus != "" {
			switch order.PaymentStatus {
			case string(biz.PaymentPending):
				orderProto.PaymentStatus = v1.PaymentStatus_NOT_PAID
			case string(biz.PaymentProcessing):
				orderProto.PaymentStatus = v1.PaymentStatus_PROCESSING
			case string(biz.PaymentPaid):
				orderProto.PaymentStatus = v1.PaymentStatus_PAID
			case string(biz.PaymentFailed):
				orderProto.PaymentStatus = v1.PaymentStatus_FAILED
			case string(biz.PaymentCancelled):
				orderProto.PaymentStatus = v1.PaymentStatus_CANCELLED
			default:
				orderProto.PaymentStatus = v1.PaymentStatus_NOT_PAID
			}
		} else {
			orderProto.PaymentStatus = v1.PaymentStatus_NOT_PAID
		}

		orders = append(orders, orderProto)
	}

	return &biz.Orders{
		Orders: orders,
	}, nil
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

func (o *orderRepo) GetAllOrders(ctx context.Context, req *biz.GetAllOrdersReq) (*biz.GetAllOrdersReply, error) {
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

	// 查询订单列表
	pageSize := int64(req.PageSize)
	page := int64((req.Page - 1) * req.PageSize)
	orders, err := o.data.db.ListOrders(ctx, models.ListOrdersParams{
		PageSize: pageSize,
		Page:     page,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("获取订单列表失败: %v", err)
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	if len(orders) == 0 {
		return &biz.GetAllOrdersReply{Orders: []*biz.SubOrder{}}, nil
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
					if order.PaymentStatus != "" {
						paymentStatus = biz.PaymentStatus(order.PaymentStatus)
					}
					break
				}
			}

			// 更新子订单的支付状态
			subOrder.Status = string(paymentStatus)
			respOrders = append(respOrders, subOrder)
		}
	}

	return &biz.GetAllOrdersReply{Orders: respOrders}, nil
}

// 解析GetConsumerOrders返回的子订单JSON
func parseSubOrders(data []byte) ([]*biz.SubOrder, error) {
	if data == nil || len(data) == 0 || string(data) == "null" || string(data) == "[null]" {
		return []*biz.SubOrder{}, nil
	}

	type dbSubOrder struct {
		ID          int64           `json:"id"`
		MerchantID  string          `json:"merchant_id"`
		TotalAmount json.Number     `json:"total_amount"`
		Currency    string          `json:"currency"`
		Status      string          `json:"status"`
		Items       json.RawMessage `json:"items"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at"`
	}

	var dbSubs []dbSubOrder
	if err := json.Unmarshal(data, &dbSubs); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	subOrders := make([]*biz.SubOrder, 0, len(dbSubs))
	for _, d := range dbSubs {
		// 检查是否为空值
		if d.ID == 0 {
			continue
		}

		merchantID, err := uuid.Parse(d.MerchantID)
		if err != nil {
			return nil, fmt.Errorf("invalid merchant id: %w", err)
		}

		totalAmount, err := d.TotalAmount.Float64()
		if err != nil {
			return nil, fmt.Errorf("invalid total amount: %w", err)
		}

		var orderItems []*biz.OrderItem
		if err := json.Unmarshal(d.Items, &orderItems); err != nil {
			// 尝试另一种格式
			type OrderItemWrapper struct {
				Items []*biz.OrderItem `json:"items"`
			}
			var wrapper OrderItemWrapper
			if wrapErr := json.Unmarshal(d.Items, &wrapper); wrapErr != nil {
				return nil, fmt.Errorf("failed to unmarshal items: %w (original error: %v)", wrapErr, err)
			}
			orderItems = wrapper.Items
		}

		subOrders = append(subOrders, &biz.SubOrder{
			ID:          d.ID,
			MerchantID:  merchantID,
			TotalAmount: totalAmount,
			Currency:    d.Currency,
			Status:      d.Status,
			Items:       orderItems,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
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
