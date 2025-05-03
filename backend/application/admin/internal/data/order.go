package data

import (
	"context"
	"encoding/json"
	"fmt"

	cartv1 "backend/api/cart/v1"
	userv1 "backend/api/user/v1"
	"backend/application/admin/internal/data/models"

	"backend/application/admin/internal/biz"
	"backend/constants"

	"github.com/go-kratos/kratos/v2/log"
)

type adminOrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewAdminOrderRepo(data *Data, logger log.Logger) biz.AdminOrderRepo {
	return &adminOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (o *adminOrderRepo) GetAllOrders(ctx context.Context, req *biz.GetAllOrdersReq) (*biz.GetAllOrdersReply, error) {
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

	page := int64((req.Page - 1) * req.PageSize)
	pageSize := int64(req.PageSize)
	orders, err := o.data.db.GetAllOrders(ctx, models.GetAllOrdersParams{
		PageSize: &pageSize,
		Page:     &page,
	})
	if err != nil {
		o.log.WithContext(ctx).Errorf("获取订单列表失败: %v", err)
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	if len(orders) == 0 {
		o.log.WithContext(ctx).Infof("没有订单记录")
		return &biz.GetAllOrdersReply{Orders: nil}, nil
	}

	var respOrders []*biz.SubOrder
	for _, order := range orders {
		// 解析订单项
		var subOrderItems []biz.OrderItem
		if err := json.Unmarshal(order.Items, &subOrderItems); err != nil {
			o.log.WithContext(ctx).Errorf("解析子订单项失败: %v, 订单ID: %d", err, order.OrderID)
			continue
		}

		var orderItems []*biz.OrderItem
		for _, item := range subOrderItems {
			if item.Item == nil {
				o.log.WithContext(ctx).Warnf("子订单项缺少商品信息, 跳过此项, 订单ID: %d", order.OrderID)
				continue
			}

			log.Debugf("item: %+v", item)
			orderItems = append(orderItems, &biz.OrderItem{
				Cost: item.Cost,
				Item: &cartv1.CartItem{
					MerchantId: item.Item.MerchantId,
					ProductId:  item.Item.ProductId,
					Quantity:   item.Item.Quantity,
					Name:       item.Item.Name,
					Picture:    item.Item.Picture,
				},
				Email: item.Email,
				ConsumerAddress: userv1.ConsumerAddress{
					StreetAddress: item.ConsumerAddress.StreetAddress,
					City:          item.ConsumerAddress.City,
					State:         item.ConsumerAddress.State,
					Country:       item.ConsumerAddress.Country,
					ZipCode:       item.ConsumerAddress.ZipCode,
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
			// ID:             order.ID,
			// MerchantID:     order.MerchantID,
			TotalAmount:    order.TotalAmount.(float64),
			Currency:       order.Currency,
			PaymentStatus:  constants.PaymentStatus(order.PaymentStatus),
			ShippingStatus: constants.ShippingStatus(order.ShippingStatus),
			Items:          orderItems,
			CreatedAt:      order.CreatedAt.Time,
			UpdatedAt:      order.UpdatedAt.Time,
		}
		respOrders = append(respOrders, subOrder)
	}

	o.log.WithContext(ctx).Debugf("获取到 %d 个商户订单", len(respOrders))
	return &biz.GetAllOrdersReply{Orders: respOrders}, nil
}
