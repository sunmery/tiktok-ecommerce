package data

import (
	cartv1 "backend/api/cart/v1"
	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"backend/application/order/pkg/convert"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
)

func (o *orderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error) {
	var extra string
	if md, ok := metadata.FromServerContext(ctx); ok {
		extra = md.Get("x-md-global-userid")
	}
	fmt.Println(extra)
	// 获取购物车商品
	cartItems, err := o.data.cartClient.GetCart(ctx, &cartv1.GetCartReq{
		Owner: "test",
		Name:  "test",
	})

	if err != nil {
		return nil, err
	}
	// 将购物车商品映射到订单商品项
	orderItems := make([]biz.OrderItem, len(cartItems.Cart.Items))
	for i, cartItem := range cartItems.Cart.Items {
		orderItems[i] = biz.OrderItem{
			Id:        int32(cartItem.ProductId),
			OrderId:   int32(req.OrderId),
			ProductId: int32(cartItem.ProductId),
			Name:      "test",
			Quantity:  cartItem.Quantity,
			Price:     6.32,
		}
	}

	// 创建订单
	orderID, err := o.data.db.CreateOrder(ctx, models.CreateOrderParams{
		Owner:         fmt.Sprintf("%d", req.UserId),
		Name:          req.Name,
		Email:         req.Email,
		StreetAddress: req.Address.StreetAddress,
		City:          req.Address.City,
		State:         req.Address.State,
		ZipCode:       string(req.Address.ZipCode),
		Currency:      req.UserCurrency,
	})
	if err != nil {
		return nil, err
	}

	// 创建订单商品
	for _, item := range req.Items {
		price, err := convert.Float32ToNumeric(item.Price)
		if err != nil {
			return nil, err
		}

		_, err = o.data.db.CreateOrderItems(ctx, models.CreateOrderItemsParams{
			OrderID:   orderID,
			ProductID: int32(item.Id),
			Name:      item.Name,
			Price:     price,
			Quantity:  item.Quantity,
		})
		if err != nil {
			return nil, err
		}
	}

	// 返回响应
	return &biz.PlaceOrderResp{
		Order: biz.OrderResult{
			OrderId: orderID,
		},
	}, nil
}

func (o *orderRepo) ListOrders(ctx context.Context, req *biz.ListOrderReq) (*biz.ListOrderResp, error) {

	// 从数据库获取订单数据
	dbOrders, err := o.data.db.ListOrders(ctx, models.ListOrdersParams{
		Owner: fmt.Sprintf("%d", req.UserId),
	})
	if err != nil {
		return nil, err
	}

	var orderSummaries []biz.OrderSummary
	for _, dbOrder := range dbOrders {
		items, err := o.data.db.ListOrderItems(ctx, dbOrder.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list order items for order %s: %w", dbOrder.ID, err)
		}

		// 将 items 转换为 biz.OrderItem 类型
		var orderItems []biz.OrderItem
		for _, item := range items {
			orderItems = append(orderItems, biz.OrderItem{
				Id:        item.ID,
				Name:      item.Name,
				OrderId:   item.OrderID,
				ProductId: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			})
		}

		orderSummaries = append(orderSummaries, biz.OrderSummary{
			OrderId:   string(dbOrder.ID),
			CreatedAt: int32(dbOrder.CreatedAt.Unix()),
			Address: biz.Address{
				StreetAddress: dbOrder.StreetAddress,
				City:          dbOrder.City,
				State:         dbOrder.State,
				Country:       dbOrder.Country,
				ZipCode:       int32(dbOrder.ZipCode),
			},
			Status:       dbOrder.Status,
			UserCurrency: dbOrder.Currency,
			Email:        dbOrder.Email,
			OrderItems:   orderItems, // 使用转换后的 orderItems
		})
	}

	// 返回包含订单列表的响应
	return &biz.ListOrderResp{
		Orders: orderSummaries,
	}, nil

}

func (o *orderRepo) MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error) {
	// 解析订单ID
	orderId, err := strconv.ParseInt(req.OrderId, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID格式: %w", err)
	}

	// 调用 casdoor 接口标记订单为已支付
	_, err = o.data.db.MarkOrderPaid(ctx, models.MarkOrderPaidParams{
		ID:    int32(orderId),
		Owner: fmt.Sprintf("%d", req.UserId),
		Name:  "test",
	})

	if err != nil {
		// 根据错误类型返回具体的错误信息
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("订单不存在")
		}
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	return &biz.MarkOrderPaidResp{
		Success: true,
	}, nil

}

func NewOrderrRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
