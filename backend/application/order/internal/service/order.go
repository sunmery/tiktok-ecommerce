package service

import (
	v1 "backend/api/cart/v1"
	pb "backend/api/order/v1"
	"backend/application/order/internal/biz"
	"backend/application/order/pkg/convert"
	"backend/application/order/pkg/token"
	"context"
	"log"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *OrderService) PlaceOrder(ctx context.Context, req *pb.PlaceOrderReq) (*pb.PlaceOrderResp, error) {

	// 从上下文获取荷载
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	UserId, err := strconv.ParseUint(payload.ID, 10, 32)
	if err != nil {
		// 处理转换错误，例如返回错误信息给调用者
		return nil, status.Error(codes.Internal, "用户ID转换失败")
	}

	// 调用购物车服务获取购物车商品
	cartResp, err := s.cartClient.GetCart(ctx, &v1.GetCartReq{
		UserId: payload.ID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "获取购物车失败: "+err.Error())
	}

	// 将购物车商品转换为 biz.OrderItem
	var items []biz.OrderItem
	for _, cartItem := range cartResp.Cart.Items {
		items = append(items, biz.OrderItem{
			Id:        int32(cartItem.ProductId),
			Name:      "dorr",
			Price:     133,
			Quantity:  cartItem.Quantity,
			OrderId:   0,
			ProductId: int32(cartItem.ProductId),
		})
	}

	result, err := s.oc.PlaceOrder(ctx, &biz.PlaceOrderReq{
		UserId:       uint32(UserId),
		UserCurrency: req.UserCurrency,
		Address: biz.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			ZipCode:       req.Address.ZipCode,
		},
		Items: items,
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &pb.PlaceOrderResp{
		Order: &pb.OrderResult{
			OrderId: string(result.Order.OrderId),
		},
	}, nil

}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrderReq) (*pb.ListOrderResp, error) {

	bizResp, err := s.oc.ListOrders(ctx, &biz.ListOrderReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, order := range bizResp.Orders {
		pbOrders = append(pbOrders, &pb.Order{
			OrderId:      order.OrderId,
			UserId:       order.UserId,
			UserCurrency: order.UserCurrency,
			Email:        order.Email,
			CreatedAt:    uint32(order.CreatedAt),
			OrderItems:   convert.ToPbOrderItems(order.OrderItems),
			Address: &pb.Address{
				StreetAddress: order.Address.StreetAddress,
				City:          order.Address.City,
				State:         order.Address.State,
				ZipCode:       order.Address.ZipCode,
			},
		})
	}

	// 返回响应
	return &pb.ListOrderResp{
		Orders: pbOrders,
	}, nil
}

func (s *OrderService) MarkOrderPaid(ctx context.Context, req *pb.MarkOrderPaidReq) (*pb.MarkOrderPaidResp, error) {

	log.Printf("MarkOrderPaid called with OrderId: %s, UserId: %d", req.OrderId, req.UserId)

	if req.OrderId == "" {
		log.Println("OrderId is empty")
		return nil, status.Error(codes.InvalidArgument, "订单ID不能为空")
	}

	_, err := s.oc.MarkOrderPaid(ctx, &biz.MarkOrderPaidReq{
		UserId:  uint32(req.UserId),
		OrderId: req.OrderId,
	})

	if err != nil {
		log.Printf("MarkOrderPaid failed: %v", err)
		return nil, status.Error(codes.Internal, "订单支付失败")
	}

	log.Println("MarkOrderPaid succeeded")

	return &pb.MarkOrderPaidResp{}, nil
}
