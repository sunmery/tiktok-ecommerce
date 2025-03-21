package data

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"backend/pkg/types"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func (o *orderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error) {
	order, err := o.data.DB(ctx).CreateOrder(ctx, models.CreateOrderParams{
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
		type SubOrderItem struct {
			Item *biz.CartItem
			Cost float64
		}
		items := []SubOrderItem{
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

		subOrder, subOrderErr := o.data.DB(ctx).CreateSubOrder(ctx, models.CreateSubOrderParams{
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
			OrderId: order.ID.String(),
		},
	}, nil
}

func (o *orderRepo) ListOrder(ctx context.Context, req *biz.ListOrderReq) (*biz.ListOrderResp, error) {
	orders, err := o.data.DB(ctx).ListOrdersByUser(ctx, models.ListOrdersByUserParams{
		UserID: req.UserID,
		Limit:  int64(req.PageSize),
		Offset: int64((req.Page - 1) * req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	var respOrders []*biz.Order
	for _, order := range orders {
		respOrders = append(respOrders, &biz.Order{
			OrderID:       order.ID,
			UserID:        order.UserID,
			Currency:      order.Currency,
			Address:       nil,
			Email:         order.Email,
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
			SubOrders:     nil,
			PaymentStatus: "",
		})
	}
	return &biz.ListOrderResp{Orders: respOrders}, nil
}

// 自定义JSON解析
func parseSubOrders(data []byte) ([]*biz.SubOrder, error) {
	type dbSubOrder struct {
		ID          string
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

func (o *orderRepo) MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error) {
	// TODO implement me
	panic("implement me")
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
