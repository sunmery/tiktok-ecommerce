package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgtype"

	"backend/constants"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	"backend/application/merchant/internal/pkg/id"

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

func (o *orderRepo) ShipOrder(ctx context.Context, req *biz.ShipOrderReq) (*biz.ShipOrderResp, error) {
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
		return nil, kerrors.New(400, "receiver_address", "invalid receiver address")
	}

	ship, err := tx.CreateShip(ctx, models.CreateShipParams{
		ID:             &snowflakeID,
		MerchantID:     merchantID,
		SubOrderID:     &req.SubOrderId,
		TrackingNumber: &req.TrackingNumber,
		Carrier:        &req.Carrier,
		// Delivery:        req.Delivery,
		ShippingAddress: req.ShippingAddress, // 使用JSON格式的地址
		ReceiverAddress: receiverAddressJSON, // 使用JSON格式的地址
		ShippingFee:     shippingFee,
	})
	if err != nil {
		return nil, kerrors.New(500, "shipId", "create shipping error")
	}

	return &biz.ShipOrderResp{
		Id:        ship.ID,
		CreatedAt: ship.CreatedAt.Local(),
	}, nil
}

// GetMerchantByOrderId 根据订单 id 查询商家
func (o *orderRepo) GetMerchantByOrderId(ctx context.Context, req *biz.GetMerchantByOrderIdReq) (*biz.GetMerchantByOrderIdReply, error) {
	merchantId, err := o.data.db.GetMerchantByOrderId(ctx, &req.OrderId)
	if err != nil {
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
			o.log.WithContext(ctx).Errorf("解析子订单项失败: %v, 订单ID: %d", err, order.ID)
			continue
		}

		// 转换为biz.OrderItem
		var orderItems []*biz.OrderItem
		for _, item := range subOrderItems {
			if item.Item == nil {
				o.log.WithContext(ctx).Warnf("子订单项缺少商品信息, 跳过此项, 订单ID: %d", order.ID)
				continue
			}

			orderItems = append(orderItems, &biz.OrderItem{
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

		// 转换金额
		var amount float64
		switch v := order.TotalAmount.(type) {
		case pgtype.Numeric:
			convertedAmount, err := types.NumericToFloat(v)
			if err != nil {
				o.log.WithContext(ctx).Errorf("转换金额失败: %v, 订单ID: %d", err, order.ID)
				continue
			}
			amount = convertedAmount
		case float64:
			amount = v
		default:
			o.log.WithContext(ctx).Errorf("未知的金额类型: %T, 订单ID: %d", order.TotalAmount, order.ID)
			continue
		}

		// 添加子订单到结果
		subOrder := &biz.SubOrder{
			OrderID:        order.OrderID,
			SubOrderID:     order.ID,
			UserID:         order.UserID,
			MerchantID:     order.MerchantID,
			TotalAmount:    amount,
			Currency:       order.Currency,
			Status:         constants.PaymentStatus(order.PaymentStatus),
			ShippingStatus: constants.ShippingStatus(order.ShippingStatus),
			Items:          orderItems,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
			StreetAddress:  order.StreetAddress,
			City:           order.City,
			State:          order.State,
			Country:        order.Country,
			ZipCode:        order.ZipCode,
			Email:          order.Email,
		}
		respOrders = append(respOrders, subOrder)
	}

	o.log.WithContext(ctx).Debugf("获取到 %d 个商户订单", len(respOrders))
	return &biz.GetMerchantOrdersReply{Orders: respOrders}, nil
}

func (o *orderRepo) UpdateOrderShippingStatus(ctx context.Context, req *biz.UpdateOrderShippingStatusReq) (*biz.UpdateOrderShippingStatusResply, error) {
	shippingStatus := string(req.ShippingStatus)
	err := o.data.db.UpdateOrderShippingStatus(ctx, models.UpdateOrderShippingStatusParams{
		ShippingStatus: &shippingStatus,
		SubOrderID:     &req.SubOrderId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "ORDER_ID_NOT_FOUND", fmt.Sprintf("未找到该子订单'%d'的商品", req.SubOrderId))
		}
		return nil, kerrors.New(500, "UPDATE_ORDER_SHIPPING_STATUS", fmt.Sprintf("更新子订单'%d'的物流状态失败: %v", req.SubOrderId, err))
	}
	return &biz.UpdateOrderShippingStatusResply{}, nil
}
