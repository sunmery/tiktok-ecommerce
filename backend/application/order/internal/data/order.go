package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"backend/application/order/internal/pkg"

	"backend/constants"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	"github.com/jackc/pgx/v5"

	"backend/application/order/internal/pkg/id"

	productv1 "backend/api/product/v1"

	userv1 "backend/api/user/v1"

	cartv1 "backend/api/cart/v1"

	v1 "backend/api/order/v1"

	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"backend/pkg/types"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func (o *orderRepo) UpdateOrderStatus(ctx context.Context, req *biz.UpdateOrderStatusReq) (*biz.UpdateOrderStatusResp, error) {
	// err := o.data.db.UpdateOrderPaymentStatus(ctx, models.UpdateOrderPaymentStatusParams{
	// 	ID:            req.OrderId,
	// 	PaymentStatus: string(req.Status),
	// })
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return nil, kerrors.New(404, "ORDER_ID_NOT_FOUND", fmt.Sprintf("未找到该订单'%d'的商品", req.OrderId))
	// 	}
	// 	return nil, err
	// }
	return &biz.UpdateOrderStatusResp{}, nil
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
	if updatePaymentStatusResult.PaymentStatus == string(constants.PaymentPaid) {
		// 订单已经是已支付状态，直接返回成功
		o.log.WithContext(ctx).Infof("Order %d is already marked as paid", req.OrderId)
		return &biz.MarkOrderPaidResp{}, nil
	}

	o.log.WithContext(ctx).Infof("Updating order %d payment status from %s to %s",
		req.OrderId, updatePaymentStatusResult.PaymentStatus, string(constants.PaymentPaid))

	// 更新订单支付状态为已支付
	_, err = tx.MarkOrderAsPaid(ctx, models.MarkOrderAsPaidParams{
		PaymentStatus: string(constants.PaymentPaid),
		ID:            req.OrderId,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to update order payment status: %v", err)
		return nil, fmt.Errorf("failed to update order payment status: %w", err)
	}

	// 更新子订单状态
	_, err = tx.MarkSubOrderAsPaid(ctx, models.MarkSubOrderAsPaidParams{
		Status:  string(constants.PaymentPaid),
		OrderID: req.OrderId,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to update sub orders payment status: %v", err)
		return nil, fmt.Errorf("failed to update sub orders payment status: %w", err)
	}

	// 更新货运状态为等待操作
	err = tx.UpdateOrderShippingStatus(ctx, models.UpdateOrderShippingStatusParams{
		ShippingStatus: string(constants.ShippingWaitCommand),
		SubOrderID:     &req.OrderId,
	})

	o.log.WithContext(ctx).Infof("Successfully marked order %d as paid for user %s", req.OrderId, req.UserId)
	return &biz.MarkOrderPaidResp{}, nil
}

func (o *orderRepo) ConfirmReceived(ctx context.Context, req *biz.ConfirmReceivedReq) (*biz.ConfirmReceivedResp, error) {
	tx := o.data.DB(ctx)

	// 获取订单信息，确认订单存在且属于该用户
	order, err := tx.GetOrderByID(ctx, models.GetOrderByIDParams{
		UserID:  req.UserId,
		OrderID: req.OrderId,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to get order %d: %v", req.OrderId, err)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// 验证订单所有者
	if order.UserID.String() != req.UserId.String() {
		o.log.WithContext(ctx).Warnf("Order %d does not belong to user %s", req.OrderId, req.UserId)
		return nil, fmt.Errorf("order does not belong to user")
	}

	// 更新订单物流状态为已确认收货
	err = tx.UpdateOrderShippingStatus(ctx, models.UpdateOrderShippingStatusParams{
		ShippingStatus: string(constants.ShippingConfirmed),
		SubOrderID:     &req.OrderId,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to update order shipping status: %v", err)
		return nil, fmt.Errorf("failed to update order shipping status: %w", err)
	}

	return &biz.ConfirmReceivedResp{}, nil
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

	paymentStatus := constants.PaymentPending
	shippingStatus := constants.ShippingWaitCommand
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
		subOrder.PaymentStatus = paymentStatus
		subOrder.ShippingStatus = shippingStatus
	}

	// 构建主订单
	orderProto := &v1.Order{
		Items:      allItems,
		OrderId:    order.ID,
		SubOrderId: nil,
		UserId:     order.UserID.String(),
		Currency:   order.Currency,
		Address: &userv1.ConsumerAddress{
			StreetAddress: order.StreetAddress,
			City:          order.City,
			State:         order.State,
			Country:       order.Country,
			ZipCode:       order.ZipCode,
		},
		Email:          order.Email,
		CreatedAt:      timestamppb.New(order.CreatedAt),
		PaymentStatus:  pkg.MapPaymentStatusToProto(paymentStatus),
		ShippingStatus: pkg.MapShippingStatusToProto(shippingStatus),
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
	paymentStatus := constants.PaymentPending
	shippingStatus := constants.ShippingWaitCommand
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
		var subOrderId int64
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

			subOrderId = subOrder.ID // 设置子订单ID
			subOrder.PaymentStatus = paymentStatus
			subOrder.ShippingStatus = shippingStatus
		}

		// 构建主订单
		orderProto := &v1.Order{
			Items:    allItems,
			OrderId:  order.ID,
			UserId:   order.UserID.String(),
			Currency: order.Currency,
			Address: &userv1.ConsumerAddress{
				StreetAddress: order.StreetAddress,
				City:          order.City,
				State:         order.State,
				Country:       order.Country,
				ZipCode:       order.ZipCode,
			},
			Email:          order.Email,
			CreatedAt:      timestamppb.New(order.CreatedAt),
			PaymentStatus:  pkg.MapPaymentStatusToProto(paymentStatus),
			ShippingStatus: pkg.MapShippingStatusToProto(shippingStatus),
			SubOrderId:     &subOrderId,
		}

		orders = append(orders, orderProto)
	}

	return &biz.Orders{
		Orders: orders,
	}, nil
}

func (o *orderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error) {
	// 生成雪花ID
	params := models.CreateOrderParams{
		ID:            id.SnowflakeID(),
		UserID:        req.UserId,
		Currency:      req.Currency,
		StreetAddress: req.Address.StreetAddress,
		City:          req.Address.City,
		State:         req.Address.State,
		Country:       req.Address.Country,
		ZipCode:       req.Address.ZipCode,
		Email:         req.Email,
	}
	log.Debugf("params: %+v", params)
	order, _ := o.data.db.CreateOrder(ctx, params)

	// 分单
	for _, item := range req.OrderItems {
		// 序列化订单项
		items := []biz.OrderItem{
			{
				Item: item.Item,
				Cost: item.Cost,
			},
		}
		itemsJSON, marshalErr := json.Marshal(items)
		if marshalErr != nil {
			return nil, kerrors.New(400, "INVALID_ORDER_ITEM", fmt.Sprintf("序列化订单项失败: %v", marshalErr))
		}

		// 转换价格到pgtype.Numeric
		totalAmount, totalAmountErr := types.Float64ToNumeric(item.Cost)
		if totalAmountErr != nil {
			return nil, fmt.Errorf("invalid price format: %w", totalAmountErr)
		}

		// 创建子订单ID
		params := models.CreateSubOrderParams{
			ID:          id.SnowflakeID(),
			OrderID:     order.ID,
			MerchantID:  item.Item.MerchantId,
			TotalAmount: totalAmount,
			Currency:    req.Currency,
			Status:      string(constants.PaymentPending),
			Items:       itemsJSON,
		}
		log.Debugf("params: %+v", params)
		_, subOrderErr := o.data.db.CreateSubOrder(ctx, params)
		if subOrderErr != nil {
			return nil, fmt.Errorf("创建子订单失败: %w", subOrderErr)
		}

		// 预扣库存
		_, updateInventoryErr := o.data.productv1.UpdateInventory(ctx, &productv1.UpdateInventoryRequest{
			ProductId:  item.Item.ProductId.String(),
			MerchantId: item.Item.MerchantId.String(),
			Stock:      -int32(item.Item.Quantity),
		})
		if updateInventoryErr != nil {
			return nil, kerrors.New(500, "INVENTORY_UPDATE_FAILED", fmt.Sprintf("更新库存失败: %v", updateInventoryErr))
		}
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
		return &biz.GetAllOrdersReply{Orders: nil}, nil
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
	for _, orderId := range orderIDs {
		if !uniqueOrderIDs[orderId] {
			uniqueOrderIDs[orderId] = true
			uniqueIDs = append(uniqueIDs, orderId)
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
			for _, order := range orders {
				if order.ID == subOrder.ID {
					subOrder.PaymentStatus = constants.PaymentStatus(order.PaymentStatus.(string))
					subOrder.ShippingStatus = constants.ShippingStatus(order.ShippingStatus)
					break
				}
			}
			respOrders = append(respOrders, subOrder)
		}
	}

	return &biz.GetAllOrdersReply{Orders: respOrders}, nil
}

func (o *orderRepo) GetShipOrderStatus(ctx context.Context, req *biz.GetShipOrderStatusReq) (*biz.GetShipOrderStatusReply, error) {
	orderStatus, err := o.data.db.GetShipOrderStatus(ctx, req.SubOrderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "ORDER_ID_NOT_FOUND", fmt.Sprintf("未找到该订单'%d'的商品", req.SubOrderId))
		}
		return nil, err
	}

	shippingFee, err := types.NumericToFloat(orderStatus.ShippingFee)
	if err != nil {
		return nil, kerrors.New(400, "shipping_fee", "invalid shipping fee")
	}
	receiverAddress := make(map[string]any)
	err = json.Unmarshal(orderStatus.ReceiverAddress, &receiverAddress)

	shippingAddress := make(map[string]any)
	err = json.Unmarshal(orderStatus.ShippingAddress, &shippingAddress)
	if err != nil {
		return nil, kerrors.New(400, "shipping_fee", "invalid shipping fee")
	}

	return &biz.GetShipOrderStatusReply{
		Id:              orderStatus.ID,
		SubOrderId:      orderStatus.SubOrderID,
		TrackingNumber:  orderStatus.TrackingNumber,
		Carrier:         orderStatus.Carrier,
		ShippingStatus:  constants.ShippingStatus(orderStatus.ShippingStatus),
		Delivery:        orderStatus.Delivery,
		ShippingFee:     shippingFee,
		ReceiverAddress: receiverAddress,
		ShippingAddress: shippingAddress,
		CreatedAt:       orderStatus.CreatedAt,
		UpdatedAt:       orderStatus.UpdatedAt,
	}, nil
}

// 解析GetConsumerOrders返回的子订单JSON
func parseSubOrders(data []byte) ([]*biz.SubOrder, error) {
	if data == nil || len(data) == 0 || string(data) == "null" || string(data) == "[null]" {
		return []*biz.SubOrder{}, nil
	}

	type dbSubOrder struct {
		ID             int64           `json:"sub_order_id,omitempty"`
		MerchantID     string          `json:"merchant_id,omitempty"`
		TotalAmount    json.Number     `json:"total_amount,omitempty"`
		Currency       string          `json:"currency,omitempty"`
		PaymentStatus  string          `json:"status,omitempty"`
		ShippingStatus string          `json:"shipping_status,omitempty"`
		Items          json.RawMessage `json:"items,omitempty"`
		CreatedAt      time.Time       `json:"created_at"`
		UpdatedAt      time.Time       `json:"updated_at"`
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

		// 解析订单项
		type OrderItemData struct {
			Cost float64 `json:"cost"`
			Item struct {
				Name       string `json:"name"`
				Picture    string `json:"picture"`
				Quantity   int32  `json:"quantity"`
				ProductId  string `json:"productId"`
				MerchantId string `json:"merchantId"`
			} `json:"item"`
		}

		var itemsData []OrderItemData
		if err := json.Unmarshal(d.Items, &itemsData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		var orderItems []*biz.OrderItem
		for _, itemData := range itemsData {
			productId, err := uuid.Parse(itemData.Item.ProductId)
			if err != nil {
				return nil, fmt.Errorf("invalid product id: %w", err)
			}

			merchantId, err := uuid.Parse(itemData.Item.MerchantId)
			if err != nil {
				return nil, fmt.Errorf("invalid merchant id: %w", err)
			}

			orderItems = append(orderItems, &biz.OrderItem{
				Item: &biz.CartItem{
					MerchantId: merchantId,
					ProductId:  productId,
					Quantity:   uint32(itemData.Item.Quantity),
					Name:       itemData.Item.Name,
					Picture:    itemData.Item.Picture,
				},
				Cost: itemData.Cost,
			})
		}

		subOrders = append(subOrders, &biz.SubOrder{
			ID:             d.ID,
			MerchantID:     merchantID,
			TotalAmount:    totalAmount,
			Currency:       d.Currency,
			PaymentStatus:  constants.PaymentStatus(d.PaymentStatus),
			ShippingStatus: constants.ShippingStatus(d.ShippingStatus),
			Items:          orderItems,
			CreatedAt:      d.CreatedAt,
			UpdatedAt:      d.UpdatedAt,
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
		amount, err := types.NumericToFloat(order.TotalAmount.(pgtype.Numeric))
		if err != nil {
			return nil, fmt.Errorf("failed to convert amount: %w", err)
		}

		subOrders = append(subOrders, &biz.SubOrder{
			ID:             order.ID,
			MerchantID:     order.MerchantID,
			TotalAmount:    amount,
			Currency:       order.Currency,
			PaymentStatus:  constants.PaymentStatus(order.PaymentStatus.(string)),
			ShippingStatus: constants.ShippingStatus(order.ShippingStatus.(string)),
			Items:          orderItems,
			CreatedAt:      order.CreatedAt.Time,
			UpdatedAt:      order.UpdatedAt.Time,
		})
	}

	return subOrders, nil
}
