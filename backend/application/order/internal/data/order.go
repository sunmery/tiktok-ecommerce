package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"backend/pkg/types"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// 数据访问层实现
const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
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
	for _, items := range req.OrderItems {

		// 序列化订单项
		itemsJSON, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("序列化订单项失败: %w", err)
		}

		// 转换价格到pgtype.Numeric
		totalAmount, err := types.Float64ToNumeric(items.Cost)
		if err != nil {
			return nil, fmt.Errorf("invalid price format: %w", err)
		}
		fmt.Printf("totalAmount: %v", totalAmount)

		subOrder, subOrderErr := o.data.DB(ctx).CreateSubOrder(ctx, models.CreateSubOrderParams{
			OrderID:     order.ID,
			MerchantID:  items.Item.MerchantId,
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
	// 参数校验
	if req.Page < 1 {
		req.Page = defaultPage
	}
	if req.PageSize < 1 || req.PageSize > maxPageSize {
		req.PageSize = defaultPageSize
	}

	// 计算时间范围
	start, end, err := o.calculateTimeRange(req)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	// 执行数据库查询
	orders, err := o.queryOrders(ctx, req.UserID, start, end, req.Page, req.PageSize)
	// orders, err := o.data.DB(ctx).ListOrdersByUserWithDate(ctx, models.ListOrdersByUserWithDateParams{
	// 	UserID: req.UserID,
	// 	Limit:  int64(req.PageSize),
	// 	Offset: int64((req.Page - 1) * req.PageSize),
	// })
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	return &biz.ListOrderResp{Orders: orders}, nil
}

func (o *orderRepo) MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error) {
	// TODO implement me
	panic("implement me")
}

func (o *orderRepo) calculateTimeRange(req *biz.ListOrderReq) (start, end time.Time, err error) {
	switch req.DateRangeType {
	case "today":
		return time.Now().Truncate(24 * time.Hour), time.Now(), nil
	case "week":
		return time.Now().AddDate(0, 0, -7), time.Now(), nil
	case "custom":
		return req.StartTime, req.EndTime, nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date range type")
	}
}

func (o *orderRepo) queryOrders(ctx context.Context, userID uuid.UUID, start, end time.Time, page, pageSize int) ([]*biz.Order, error) {
	offset := (page - 1) * pageSize

	rows, err := o.data.DB(ctx).ListOrdersByUserWithDate(ctx, models.ListOrdersByUserWithDateParams{
		UserID:    userID,
		StartTime: start,
		EndTime:   end,
		Offsets:   int64(pageSize),
		Limits:    int64(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var orders []*biz.Order
	for _, dbOrder := range rows {
		// 解析子订单
		subOrders, err := parseSubOrders(dbOrder.SubOrders)
		if err != nil {
			return nil, fmt.Errorf("parse sub_orders failed: %w", err)
		}

		orders = append(orders, &biz.Order{
			OrderID:  dbOrder.ID.String(),
			UserID:   dbOrder.UserID,
			Currency: dbOrder.Currency,
			Address: &biz.Address{
				StreetAddress: dbOrder.StreetAddress,
				City:          dbOrder.City,
				State:         dbOrder.State,
				Country:       dbOrder.Country,
				ZipCode:       dbOrder.ZipCode,
			},
			Email:         dbOrder.Email,
			CreatedAt:     dbOrder.CreatedAt,
			SubOrders:     subOrders,
			PaymentStatus: dbOrder.PaymentStatus,
		})
	}

	return orders, nil
}

// 自定义JSON解析（使用标准库）
func parseSubOrders(data []byte) ([]*biz.SubOrder, error) {
	var dbSubOrders []struct {
		ID          string    `json:"id"`
		MerchantID  uuid.UUID `json:"merchant_id"`
		TotalAmount string    `json:"total_amount"`
		Currency    string    `json:"currency"`
		Status      string    `json:"status"`
		Items       []struct {
			MerchantId uuid.UUID
			// 商品ID
			ProductId uuid.UUID
			// 商品数量
			Quantity uint32
		} `json:"items"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if err := json.Unmarshal(data, &dbSubOrders); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	subOrders := make([]*biz.SubOrder, 0, len(dbSubOrders))
	for _, so := range dbSubOrders {
		items := make([]biz.OrderItem, 0, len(so.Items))
		for _, item := range so.Items {
			items = append(items, biz.OrderItem{
				Item: &biz.CartItem{
					MerchantId: item.MerchantId,
					ProductId:  item.ProductId,
					Quantity:   item.Quantity,
				},
				Cost: 0,
			})
		}

		subOrders = append(subOrders, &biz.SubOrder{
			ID:          so.ID,
			MerchantID:  so.MerchantID,
			TotalAmount: so.TotalAmount,
			Currency:    so.Currency,
			Status:      so.Status,
			Items:       items,
			CreatedAt:   so.CreatedAt,
			UpdatedAt:   so.UpdatedAt,
		})
	}

	return subOrders, nil
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
