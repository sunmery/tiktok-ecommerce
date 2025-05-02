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

func (o *orderRepo) GetConsumerOrders(ctx context.Context, req *biz.GetConsumerOrdersReq) (*biz.GetConsumerOrdersReply, error) {
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
	var OrderID int64
	for _, row := range rows {
		OrderID = row.OrderID
		subOrders := make([]*biz.ConsumerOrder, 0, len(row.SubOrders))
		err := json.Unmarshal(row.SubOrders, &subOrders)
		if err != nil {
			return nil, err
		}

		for _, order := range subOrders {
			items := make([]*biz.ConsumerOrderItem, 0, len(order.Items))
			for _, item := range order.Items {
				items = append(items, &biz.ConsumerOrderItem{
					Cost: item.Cost,
					Item: &biz.CartItem{
						MerchantId: item.Item.MerchantId,
						ProductId:  item.Item.ProductId,
						Quantity:   item.Item.Quantity,
						Name:       item.Item.Name,
						Picture:    item.Item.Picture,
					},
				})
			}
			orders = append(orders, &biz.ConsumerOrder{
				Items: items,
				Address: biz.Address{
					StreetAddress: order.Address.StreetAddress,
					City:          order.Address.City,
					State:         order.Address.State,
					Country:       order.Address.Country,
					ZipCode:       order.Address.ZipCode,
				},
				SubOrderID:     order.SubOrderID,
				Currency:       order.Currency,
				PaymentStatus:  order.PaymentStatus,
				ShippingStatus: order.ShippingStatus,
				Email:          order.Email,
				CreatedAt:      order.CreatedAt,
				UpdatedAt:      order.UpdatedAt,
			})
		}
	}

	return &biz.GetConsumerOrdersReply{
		SubOrders: orders,
		OrderId:   OrderID,
	}, nil
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (o *orderRepo) GetUserOrdersWithSuborders(ctx context.Context, req *biz.GetUserOrdersWithSubordersReq) (*biz.GetUserOrdersWithSubordersReply, error) {
	suborders, err := o.data.db.GetUserOrdersWithSuborders(ctx, models.GetUserOrdersWithSubordersParams{
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

	orders := make([]*biz.Suborder, 0, len(suborders))
	for _, s := range suborders {

		totalAmount, err := types.NumericToFloat(s.TotalAmount)
		if err != nil {
			return nil, err
		}
		var allItems []*biz.OrderItem
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

		orders = append(orders, &biz.Suborder{
			OrderId:        s.ID,
			SubOrderId:     *s.SubOrderID,
			StreetAddress:  s.StreetAddress,
			City:           s.City,
			State:          s.State,
			Country:        s.Country,
			ZipCode:        s.ZipCode,
			Email:          s.Email,
			MerchantId:     s.MerchantID.String(),
			PaymentStatus:  constants.PaymentStatus(s.PaymentStatus),
			ShippingStatus: constants.ShippingStatus(*(s.ShippingStatus)),
			TotalAmount:    totalAmount,
			Currency:       *s.Currency,
			Items:          allItems,
			CreatedAt:      s.CreatedAt,
			UpdatedAt:      s.UpdatedAt,
		})
	}
	return &biz.GetUserOrdersWithSubordersReply{
		Orders: orders,
	}, nil
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

func (o *orderRepo) ConfirmReceived(ctx context.Context, req *biz.ConfirmReceivedReq) (*biz.ConfirmReceivedResp, error) {
	tx := o.data.DB(ctx)

	// 获取订单信息，确认订单存在且属于该用户
	order, err := tx.GetOrderByID(ctx, models.GetOrderByIDParams{
		UserID:  req.UserId,
		OrderID: req.OrderId,
	})
	if err != nil {
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
	err = tx.UpdateOrderShippingStatus(ctx, models.UpdateOrderShippingStatusParams{
		ShippingStatus: &shippingConfirmed,
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
		// 序列化订单项
		items := []biz.OrderItem{subOrder.item}
		itemsJSON, marshalErr := json.Marshal(items)
		if marshalErr != nil {
			return nil, kerrors.New(400, "INVALID_ORDER_ITEM", fmt.Sprintf("序列化订单项失败: %v", marshalErr))
		}

		// 转换价格到pgtype.Numeric
		totalAmount, totalAmountErr := types.Float64ToNumeric(subOrder.amount)
		if totalAmountErr != nil {
			return nil, fmt.Errorf("invalid price format: %w", totalAmountErr)
		}

		// 创建子订单ID
		params := models.CreateSubOrderParams{
			ID:          id.SnowflakeID(),
			OrderID:     order.ID,
			MerchantID:  subOrder.merchantId,
			TotalAmount: totalAmount,
			Currency:    req.Currency,
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

	return &biz.PlaceOrderResp{
		Order: &biz.OrderResult{
			OrderId:         order.ID,
			FreezeId:        freezeBalance.FreezeId,
			ConsumerVersion: int64(freezeBalance.NewVersion),
			MerchantVersion: merchantVersion,
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
					subOrder.PaymentStatus = constants.PaymentStatus(order.PaymentStatus)
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
			PaymentStatus:  constants.PaymentStatus(order.PaymentStatus),
			ShippingStatus: constants.ShippingStatus(order.ShippingStatus),
			Items:          orderItems,
			CreatedAt:      order.CreatedAt.Time,
			UpdatedAt:      order.UpdatedAt.Time,
		})
	}

	return subOrders, nil
}
