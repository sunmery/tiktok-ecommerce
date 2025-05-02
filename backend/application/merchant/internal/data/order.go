package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"backend/application/merchant/internal/pkg/id"

	"backend/constants"

	"github.com/jackc/pgx/v5"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/log"

	"backend/pkg/types"

	"backend/application/merchant/internal/data/models"

	"backend/application/merchant/internal/biz"
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

func (o *orderRepo) CreateOrderShip(ctx context.Context, req *biz.CreateOrderShipReq) (*biz.CreateOrderShipResp, error) {
	tx := o.data.DB(ctx)

	merchantID := types.ToPgUUID(req.MerchantID)
	snowflakeID := id.SnowflakeID()
	shippingFee, err := types.Float64ToNumeric(req.ShippingFee)
	if err != nil {
		return nil, kerrors.New(400, "shipping_fee", "invalid shipping fee")
	}

	receiverAddress, err := tx.GetConsumerAddress(ctx, &req.SubOrderId)
	if err != nil {
		return nil, kerrors.New(400, "receiver_address", "invalid receiver address")
	}

	// TODO 这个receiverAddress存储了处理地址以外的字段, 可能在数据库去除或者 go 结构体提取
	receiverAddressJSON, err := json.Marshal(receiverAddress)
	if err != nil {
		return nil, kerrors.New(400, "receiver_address", fmt.Sprintf("invalid receiver address: %v", err))
	}
	shippingStatus := string(constants.ShippingShipped)
	ship, err := tx.CreateShip(ctx, models.CreateShipParams{
		ID:             &snowflakeID,
		MerchantID:     merchantID,
		SubOrderID:     &req.SubOrderId,
		ShippingStatus: &shippingStatus,
		TrackingNumber: &req.TrackingNumber,
		Carrier:        &req.Carrier,
		// Delivery:        req.Delivery,
		ShippingAddress: req.ShippingAddress, // 使用JSON格式的地址
		ReceiverAddress: receiverAddressJSON, // 使用JSON格式的地址
		ShippingFee:     shippingFee,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "CREATE_ORDER_SHIP_ORDER_NOT_FOUND", "order not found")
		}
		return nil, kerrors.New(500, "CREATE_ORDER_SHIP_ORDER_INTERNAL_ERROR", fmt.Sprintf("create ship failed: %v", err))
	}

	return &biz.CreateOrderShipResp{
		Id:        ship.ID,
		CreatedAt: ship.CreatedAt,
	}, nil
}

// GetMerchantByOrderId 根据订单 id 查询商家
func (o *orderRepo) GetMerchantByOrderId(ctx context.Context, req *biz.GetMerchantByOrderIdReq) (*biz.GetMerchantByOrderIdReply, error) {
	merchantId, err := o.data.db.GetMerchantByOrderId(ctx, &req.OrderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "ORDER_NOT_FOUND", "order not found")
		}
		return nil, fmt.Errorf("获取商家失败: %w", err)
	}
	return &biz.GetMerchantByOrderIdReply{
		MerchantId: merchantId,
	}, nil
}

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
	// orders, err := o.data.orderv1.GetMerchantOrders(ctx, &v1.GetMerchantOrdersReq{
	// 	MerchantId: req.UserID.String(),
	// 	Page:       (req.Page - 1) * req.PageSize,
	// 	PageSize:   req.PageSize,
	// })
	userID := types.ToPgUUID(req.UserID)
	page := int64((req.Page - 1) * req.PageSize)
	pageSize := int64(req.PageSize)
	orders, err := o.data.db.GetMerchantOrders(ctx, models.GetMerchantOrdersParams{
		MerchantID: userID,
		PageSize:   &pageSize,
		Page:       &page,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("获取订单列表失败: %v", err)
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	if len(orders) == 0 {
		o.log.WithContext(ctx).Infof("商户 %s 没有订单记录", req.UserID)
		return &biz.GetMerchantOrdersReply{Orders: nil}, nil
	}

	var respOrders []*biz.SubOrder
	for _, order := range orders {
		// 解析订单项
		var subOrderItems []biz.OrderItem
		if err := json.Unmarshal(order.Items, &subOrderItems); err != nil {
			o.log.WithContext(ctx).Errorf("解析子订单项失败: %v, 订单ID: %d", err, order.OrderID)
			continue
		}

		// 转换为biz.OrderItem
		var orderItems []*biz.OrderItem
		for _, item := range subOrderItems {
			if item.Item == nil {
				o.log.WithContext(ctx).Warnf("子订单项缺少商品信息, 跳过此项, 订单ID: %d", order.OrderID)
				continue
			}

			log.Debugf("item: %+v", item)
			orderItems = append(orderItems, &biz.OrderItem{
				Cost: item.Cost,
				Item: &biz.CartItem{
					MerchantId: item.Item.MerchantId,
					ProductId:  item.Item.ProductId,
					Quantity:   item.Item.Quantity,
					Name:       item.Item.Name,
					Picture:    item.Item.Picture,
				},
				Email: item.Email,
				Address: biz.Address{
					StreetAddress: item.Address.StreetAddress,
					City:          item.Address.City,
					State:         item.Address.State,
					Country:       item.Address.Country,
					ZipCode:       item.Address.ZipCode,
				},
				UserID:         item.UserID,
				SubOrderID:     item.SubOrderID,
				TotalAmount:    item.TotalAmount,
				Currency:       item.Currency,
				PaymentStatus:  item.PaymentStatus,
				ShippingStatus: item.ShippingStatus,
				CreatedAt:      item.CreatedAt,
				UpdatedAt:      item.UpdatedAt,
			})
		}

		// 添加子订单到结果
		subOrder := &biz.SubOrder{
			OrderID:   order.OrderID,
			Items:     orderItems,
			CreatedAt: order.CreatedAt,
		}
		respOrders = append(respOrders, subOrder)
	}

	o.log.WithContext(ctx).Debugf("获取到 %d 个商户订单", len(respOrders))
	return &biz.GetMerchantOrdersReply{Orders: respOrders}, nil
}

func (o *orderRepo) UpdateOrderShippingStatus(ctx context.Context, req *biz.UpdateOrderShippingStatusReq) (*biz.UpdateOrderShippingStatusResply, error) {
	tx := o.data.DB(ctx)
	shippingStatus := string(req.ShippingStatus)
	merchantID := types.ToPgUUID(req.MerchantID)
	receiverAddress, err := tx.GetConsumerAddress(ctx, &req.SubOrderId)
	if err != nil {
		return nil, kerrors.New(400, "receiver_address", "invalid receiver address")
	}
	receiverAddressJSON, err := json.Marshal(receiverAddress)
	if err != nil {
		return nil, kerrors.New(400, "receiver_address", fmt.Sprintf("invalid receiver address: %v", err))
	}
	shippingFee, err := types.Float64ToNumeric(req.ShippingFee)
	if err != nil {
		return nil, kerrors.New(400, "shipping_fee", "invalid shipping fee")
	}
	params := models.UpdateOrderShippingStatusParams{
		MerchantID:     merchantID,
		SubOrderID:     &req.SubOrderId,
		ShippingStatus: &shippingStatus,
		TrackingNumber: &req.TrackingNumber,
		Carrier:        &req.Carrier,
		// Delivery:        req.Delivery,
		ShippingAddress: req.ShippingAddress, // 使用JSON格式的地址
		ReceiverAddress: receiverAddressJSON, // 使用JSON格式的地址
		ShippingFee:     shippingFee,
	}
	log.Debugf("params:%+v", params)
	result, err := tx.UpdateOrderShippingStatus(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "ORDER_ID_NOT_FOUND", fmt.Sprintf("未找到该子订单'%d'的商品", req.SubOrderId))
		}
		return nil, kerrors.New(500, "UPDATE_ORDER_SHIPPING_STATUS", fmt.Sprintf("更新子订单'%d'的物流状态失败: %v", req.SubOrderId, err))
	}
	return &biz.UpdateOrderShippingStatusResply{
		ID:        result.ID,
		UpdatedAt: result.UpdatedAt,
	}, nil
}
