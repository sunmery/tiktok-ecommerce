package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	balancerv1 "backend/api/balancer/v1"

	"github.com/jackc/pgx/v5/pgtype"

	"backend/constants"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	"github.com/jackc/pgx/v5"

	"backend/application/consumer/internal/pkg/id"

	productv1 "backend/api/product/v1"

	"backend/application/consumer/internal/biz"
	"backend/application/consumer/internal/data/models"
	"backend/pkg/types"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type consumerOrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewConsumerOrderRepo(data *Data, logger log.Logger) biz.ConsumerOrderRepo {
	return &consumerOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (o *consumerOrderRepo) GetConsumerOrders(ctx context.Context, req *biz.GetConsumerOrdersRequest) (*biz.GetConsumerOrdersReply, error) {
	rows, err := o.data.db.GetConsumerOrders(ctx, models.GetConsumerOrdersParams{
		UserID:   req.UserId,
		PageSize: int64(req.PageSize),
		Page:     int64((req.Page - 1) * req.PageSize),
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("获取用户订单列表失败: %v", err)
		return nil, fmt.Errorf("获取用户订单列表失败: %w", err)
	}

	orders := make([]*biz.ConsumerOrder, 0, len(rows))
	for _, row := range rows {
		var items []*biz.OrderItem
		err := json.Unmarshal(row.Items, &items)
		if err != nil {
			log.Warnf("failed to unmarshal items: %v, item order_id: '%+v' sub_order_id:%+v", err, row.OrderID, row.SubOrderID)
			continue
		}

		orderItems := make([]*biz.ConsumerOrderItem, 0, len(items))
		for _, i := range items {
			orderItems = append(orderItems, &biz.ConsumerOrderItem{
				Cost: i.Cost,
				Item: &biz.CartItem{
					MerchantId: i.Item.MerchantId,
					ProductId:  i.Item.ProductId,
					Quantity:   i.Item.Quantity,
					Name:       i.Item.Name,
					Picture:    i.Item.Picture,
				},
			})
		}

		orders = append(orders, &biz.ConsumerOrder{
			OrderId:    row.OrderID,
			SubOrderID: *row.SubOrderID,
			Items:      orderItems,
			Address: biz.Address{
				StreetAddress: row.StreetAddress,
				City:          row.City,
				State:         row.State,
				Country:       row.Country,
				ZipCode:       row.ZipCode,
			},
			Currency:       *row.Currency,
			PaymentStatus:  constants.PaymentStatus(*row.PaymentStatus),
			ShippingStatus: constants.ShippingStatus(*row.ShippingStatus),
			Email:          row.Email,
			CreatedAt:      row.CreatedAt.Time,
			UpdatedAt:      row.UpdatedAt.Time,
		})
	}

	return &biz.GetConsumerOrdersReply{
		SubOrders: orders,
	}, nil
}

func (o *consumerOrderRepo) GetConsumerOrdersWithSuborders(ctx context.Context, req *biz.GetConsumerOrdersWithSubordersRequest) (*biz.GetConsumerOrdersWithSubordersReply, error) {
	suborders, err := o.data.db.GetConsumerOrdersWithSuborders(ctx, models.GetConsumerOrdersWithSubordersParams{
		UserID:  req.UserId,
		OrderID: req.OrderId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(suborders) == 0 {
		return nil, nil
	}

	orders := make([]*biz.SubOrder, 0, len(suborders))
	for _, s := range suborders {
		totalAmount, err := types.NumericToFloat(s.TotalAmount)
		if err != nil {
			return nil, err
		}
		var allItems []*biz.OrderItem
		if s.Items != nil {
			if err := json.Unmarshal(s.Items, &allItems); err != nil {
				return nil, fmt.Errorf("failed to unmarshal items: %w", err)
			}
			for _, item := range allItems {
				allItems = append(allItems, &biz.OrderItem{
					Item: &biz.CartItem{
						MerchantId: item.Item.MerchantId,
						ProductId:  item.Item.ProductId,
						Quantity:   item.Item.Quantity,
						Name:       item.Item.Name,
						Picture:    item.Item.Picture,
					},
					Cost: item.Cost,
				})
			}
		}

		orders = append(orders, &biz.SubOrder{
			OrderId:        s.ID,
			SubOrderId:     *s.SubOrderID,
			StreetAddress:  s.StreetAddress,
			City:           s.City,
			State:          s.State,
			Country:        s.Country,
			ZipCode:        s.ZipCode,
			Email:          s.Email,
			MerchantId:     s.MerchantID.String(),
			TotalAmount:    totalAmount,
			PaymentStatus:  constants.PaymentStatus(s.PaymentStatus),
			ShippingStatus: constants.ShippingStatus(*(s.ShippingStatus)),
			Currency:       *s.Currency,
			Items:          allItems,
			CreatedAt:      s.CreatedAt,
			UpdatedAt:      s.UpdatedAt,
		})
	}
	return &biz.GetConsumerOrdersWithSubordersReply{
		Orders: orders,
	}, nil
}

func (o *consumerOrderRepo) ConfirmReceived(ctx context.Context, req *biz.ConfirmReceivedRequest) (*biz.ConfirmReceivedReply, error) {
	tx := o.data.DB(ctx)

	// 获取订单信息，确认订单存在且属于该用户
	order, err := tx.GetSubOrderByID(ctx, models.GetSubOrderByIDParams{
		UserID:  req.UserId,
		OrderID: req.OrderId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.NotFound("order", fmt.Sprintf("order '%d' not found", req.OrderId))
		}
		o.log.WithContext(ctx).Errorf("Failed to get order %d: %v", req.OrderId, err)
		return nil, fmt.Errorf("getOrderById failed to get order: %w", err)
	}

	// 验证订单所有者
	if order.UserID.String() != req.UserId.String() {
		o.log.WithContext(ctx).Warnf("Order %d does not belong to user %s", req.OrderId, req.UserId)
		return nil, fmt.Errorf("order does not belong to user")
	}

	// 更新订单物流状态为已确认收货
	shippingConfirmed := string(constants.ShippingConfirmed)
	params := models.UpdateOrderShippingStatusParams{
		ShippingStatus: &shippingConfirmed,
		SubOrderID:     &req.OrderId,
	}
	log.Debugf("params: %+v", params)
	err = tx.UpdateOrderShippingStatus(ctx, params)
	if err != nil {
		o.log.WithContext(ctx).Errorf("Failed to update order shipping status: %v", err)
		return nil, fmt.Errorf("failed to update order shipping status: %w", err)
	}

	return &biz.ConfirmReceivedReply{}, nil
}

func (o *consumerOrderRepo) GetConsumerSubOrderDetail(ctx context.Context, req *biz.GetConsumerSubOrderDetailRequest) (*biz.ConsumerOrder, error) {
	order, err := o.data.db.GetOrderByID(ctx, models.GetOrderByIDParams{
		UserID:     req.UserId,
		SubOrderID: req.SubOrderId,
	})
	if err != nil {
		return nil, err
	}

	var items []*biz.ConsumerOrderItem
	err = json.Unmarshal(order.Items, &items)
	if err != nil {
		return nil, err
	}

	return &biz.ConsumerOrder{
		OrderId:    order.OrderID,
		SubOrderID: order.SubOrderID,
		Items:      items,
		Address: biz.Address{
			StreetAddress: order.StreetAddress,
			City:          order.City,
			State:         order.State,
			Country:       order.Country,
			ZipCode:       order.ZipCode,
		},
		Currency:       order.Currency,
		PaymentStatus:  constants.PaymentStatus(order.PaymentStatus),
		ShippingStatus: constants.ShippingStatus(order.ShippingStatus),
		Email:          order.Email,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}, nil
}

func (o *consumerOrderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderRequest) (*biz.PlaceOrderReply, error) {
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

	// 计算订单总金额
	var totalOrderAmount float64
	var subOrders []struct {
		merchantId uuid.UUID
		amount     float64
		item       biz.OrderItem
	}

	// 第一步：计算总金额并收集子订单信息
	for _, item := range req.OrderItems {
		totalOrderAmount += item.Cost

		subOrders = append(subOrders, struct {
			merchantId uuid.UUID
			amount     float64
			item       biz.OrderItem
		}{
			merchantId: item.Item.MerchantId,
			amount:     item.Cost,
			item: biz.OrderItem{
				Item: &biz.CartItem{
					MerchantId: item.Item.MerchantId,
					ProductId:  item.Item.ProductId,
					Quantity:   item.Item.Quantity,
					Name:       item.Item.Name,
					Picture:    item.Item.Picture,
				},
				Cost: item.Cost,
			},
		})
	}

	// 第二步：获取用户余额并一次性冻结总金额
	userBalance, err := o.data.balancerv1.GetUserBalance(ctx, &balancerv1.GetUserBalanceRequest{
		UserId:   req.UserId.String(),
		Currency: req.Currency,
	})
	if err != nil {
		return nil, err
	}

	if userBalance.Available < totalOrderAmount {
		return nil, kerrors.New(400, "INSUFFICIENT_BALANCE", "可用余额不足")
	}

	// 一次性冻结总金额
	freezeBalance, freezeBalanceErr := o.data.balancerv1.FreezeBalance(ctx, &balancerv1.FreezeBalanceRequest{
		UserId:          req.UserId.String(),
		OrderId:         order.ID,
		Amount:          totalOrderAmount,
		Currency:        req.Currency,
		IdempotencyKey:  strconv.FormatInt(order.ID, 10),
		ExpectedVersion: userBalance.Version,
	})
	if freezeBalanceErr != nil {
		return nil, kerrors.New(500, "FREEZE_BALANCE_FAILED", fmt.Sprintf("冻结余额失败: %v", freezeBalanceErr))
	}

	// 第三步：创建子订单并预扣库存
	merchantVersions := make(map[string]int64)
	for _, subOrder := range subOrders {
		totalAmount, err := types.Float64ToNumeric(subOrder.amount)
		log.Debugf("405 totalAmount: %v", totalAmount)
		// 序列化订单项
		items := []biz.OrderItem{subOrder.item}
		itemsJSON, marshalErr := json.Marshal(items)
		if marshalErr != nil {
			return nil, kerrors.New(400, "INVALID_ORDER_ITEM", fmt.Sprintf("序列化订单项失败: %v", marshalErr))
		}

		// 创建子订单ID
		params := models.CreateSubOrderParams{
			ID:          id.SnowflakeID(),
			OrderID:     order.ID,
			MerchantID:  subOrder.merchantId,
			TotalAmount: totalAmount,
			Status:      string(constants.PaymentPending),
			Items:       itemsJSON,
		}
		log.Debugf("params: %+v", params)

		createSubOrderResult, subOrderErr := o.data.db.CreateSubOrder(ctx, params)
		if subOrderErr != nil {
			return nil, fmt.Errorf("创建子订单失败: %w", subOrderErr)
		}

		// 预扣库存
		_, updateInventoryErr := o.data.productv1.UpdateInventory(ctx, &productv1.UpdateInventoryRequest{
			ProductId:  subOrder.item.Item.ProductId.String(),
			MerchantId: subOrder.item.Item.MerchantId.String(),
			Stock:      -int32(subOrder.item.Item.Quantity),
		})
		if updateInventoryErr != nil {
			if errors.Is(updateInventoryErr, pgx.ErrNoRows) {
				return nil, kerrors.New(500, "INVENTORY_UPDATE_FAILED", fmt.Sprintf("查找该商家 id:%v的商品 id%v, 失败%v", subOrder.item.Item.MerchantId.String(), subOrder.item.Item.ProductId.String(), updateInventoryErr))
			}
			log.Errorf("商家%v 商品id%v 数量%v", subOrder.item.Item.ProductId.String(), subOrder.item.Item.MerchantId.String(), subOrder.item.Item.Quantity)
			return nil, kerrors.New(500, "INVENTORY_UPDATE_FAILED", fmt.Sprintf("更新库存失败: %v", updateInventoryErr))
		}

		// 创建物流信息为等待操作
		shippingAddress := map[string]any{}
		ShippingAddressJSON, err := json.Marshal(shippingAddress)
		if err != nil {
			return nil, kerrors.New(400, "INVALID_SHIPPING_ADDRESS", fmt.Sprintf("序列化物流地址失败: %v", err))
		}

		receiverAddress := map[string]any{}
		receiverAddressJSON, err := json.Marshal(receiverAddress)
		if err != nil {
			return nil, kerrors.New(400, "INVALID_SHIPPING_ADDRESS", fmt.Sprintf("序列化物流地址失败: %v", err))
		}
		shippingFee, err := types.Float64ToNumeric(0)
		if err != nil {
			return nil, fmt.Errorf("invalid price format: %w", err)
		}
		createOrderShippingParams := models.CreateOrderShippingParams{
			ID:              id.SnowflakeID(),
			MerchantID:      createSubOrderResult.MerchantID,
			SubOrderID:      createSubOrderResult.ID,
			ShippingStatus:  string(constants.ShippingWaitCommand),
			TrackingNumber:  "0",                  // 默认值, 等待商家发货时更新
			Carrier:         "unknown",            // 默认值, 等待商家发货时更新
			Delivery:        pgtype.Timestamptz{}, // 默认值, 等待用户确认收货触发送达时间
			ShippingAddress: ShippingAddressJSON,  // 默认值, 等待商家发货时更新
			ReceiverAddress: receiverAddressJSON,  // 默认值, 等待商家发货时更新
			ShippingFee:     shippingFee,          // 默认值, 等待商家发货时更新
		}
		createOrderShippingErr := o.data.db.CreateOrderShipping(ctx, createOrderShippingParams)
		if createOrderShippingErr != nil {
			return nil, kerrors.New(500, "CREATE_ORDER_SHIPPING_FAILED", fmt.Sprintf("创建物流信息失败: %v", createOrderShippingErr))
		}

		// 获取商家余额信息（仅用于返回版本号）
		merchantId := subOrder.merchantId.String()
		if _, exists := merchantVersions[merchantId]; !exists {
			merchantBalance, merchantBalanceErr := o.data.balancerv1.GetMerchantBalance(ctx, &balancerv1.GetMerchantBalanceRequest{
				MerchantId: merchantId,
				Currency:   req.Currency,
			})
			if merchantBalanceErr != nil {
				return nil, kerrors.New(500, "GET_MERCHANT_BALANCE_NOT_FOUND", fmt.Sprintf("查询商家余额失败: %v", merchantBalanceErr))
			}
			merchantVersions[merchantId] = int64(merchantBalance.Version)
		}
	}

	// 返回订单结果
	var merchantVersion []int64
	for _, version := range merchantVersions {
		merchantVersion = append(merchantVersion, version)
	}

	return &biz.PlaceOrderReply{
		Order: &biz.OrderResult{
			OrderId:         order.ID,
			FreezeId:        freezeBalance.FreezeId,
			ConsumerVersion: int64(freezeBalance.NewVersion),
			MerchantVersion: merchantVersion,
		},
	}, nil
}

func (o *consumerOrderRepo) GetShipOrderStatus(ctx context.Context, req *biz.GetShipOrderStatusRequest) (*biz.GetShipOrderStatusReply, error) {
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
	receiverAddress := biz.ReceiverAddress{}
	err = json.Unmarshal(orderStatus.ReceiverAddress, &receiverAddress)
	if err != nil {
		return nil, kerrors.New(400, "shipping_fee", "invalid shipping fee")
	}

	shippingAddress := biz.ShippingAddress{}
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
		Delivery:        orderStatus.Delivery.Time,
		ShippingFee:     shippingFee,
		ReceiverAddress: receiverAddress,
		ShippingAddress: shippingAddress,
		CreatedAt:       orderStatus.CreatedAt,
		UpdatedAt:       orderStatus.UpdatedAt,
	}, nil
}

// 解析GetOrders返回的子订单JSON
func parseSubOrders(data []byte) ([]*biz.SubOrder, error) {
	if data == nil || len(data) == 0 || string(data) == "null" || string(data) == "[null]" {
		return []*biz.SubOrder{}, nil
	}

	type dbSubOrder struct {
		ID             int64           `json:"sub_order_id,omitempty"`
		MerchantID     string          `json:"merchant_id,omitempty"`
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
			OrderId: d.ID,
			// MerchantID:     merchantID,
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
