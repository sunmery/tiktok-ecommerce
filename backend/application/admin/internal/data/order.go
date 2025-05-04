package data

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/pkg/types"

	"github.com/jackc/pgx/v5/pgtype"

	"backend/constants"

	"backend/application/admin/internal/data/models"

	"backend/application/admin/internal/biz"
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
		return nil, nil
	}

	var respOrders []*biz.SubOrder
	for _, o := range orders {
		var subOrderItem []*biz.SubOrderItem
		unmarshalSubOrderItemsErr := json.Unmarshal(o.Items, &subOrderItem)
		if unmarshalSubOrderItemsErr != nil {
			return nil, fmt.Errorf("解析订单商品列表失败: %w", unmarshalSubOrderItemsErr)
		}

		for _, item := range subOrderItem {
			subOrderItem = append(subOrderItem, &biz.SubOrderItem{
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

		var consumerAddress biz.ConsumerAddress
		err := json.Unmarshal(o.ConsumerAddress, &consumerAddress)
		if err != nil {
			return nil, fmt.Errorf("解析地址信息失败: %w", err)
		}

		var totalAmount float64
		_, ok := o.TotalAmount.(pgtype.Numeric)
		if ok {
			totalAmount, err = types.NumericToFloat(o.TotalAmount.(pgtype.Numeric))
			if err != nil {
				return nil, fmt.Errorf("解析订单金额失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("解析订单金额失败: %w", err)
		}

		respOrders = append(respOrders, &biz.SubOrder{
			OrderID:     o.OrderID,
			SubOrderID:  o.SubOrderID,
			TotalAmount: totalAmount,
			ConsumerId:  o.UserID,
			ConsumerAddress: &biz.ConsumerAddress{
				City:          consumerAddress.City,
				State:         consumerAddress.State,
				Country:       consumerAddress.Country,
				ZipCode:       consumerAddress.ZipCode,
				StreetAddress: consumerAddress.StreetAddress,
			},
			ConsumerEmail:  o.Email,
			Currency:       constants.Currency(o.Currency),
			PaymentStatus:  constants.PaymentStatus(o.PaymentStatus),
			ShippingStatus: constants.ShippingStatus(o.ShippingStatus),
			SubOrderItems:  subOrderItem,
			CreatedAt:      o.CreatedAt.Time,
			UpdatedAt:      o.UpdatedAt.Time,
		})
	}

	o.log.WithContext(ctx).Debugf("获取到 %d 个订单", len(respOrders))
	return &biz.GetAllOrdersReply{Orders: respOrders}, nil
}
