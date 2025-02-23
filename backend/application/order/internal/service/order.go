package service

import (
	pb "backend/api/order/v1"
	"backend/application/order/internal/biz"
	"backend/application/order/pkg/convert"
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *OrderService) PlaceOrder(ctx context.Context, req *pb.PlaceOrderReq) (*pb.PlaceOrderResp, error) {
	// 将请求中的OrderItems转换为业务层需要的结构
	var bizItems []biz.OrderItem
	for _, pbItem := range req.OrderItems {
		if pbItem.Item == nil {
			return nil, status.Error(codes.InvalidArgument, "购物车商品信息不完整")
		}

		bizItems = append(bizItems, biz.OrderItem{
			ProductId: int32(pbItem.Item.ProductId), // 假设CartItem包含ProductId字段
			Name:      pbItem.Item.Name,             // 假设CartItem包含Name字段
			Price:     pbItem.Item.Price,            // 假设CartItem包含Price字段
			Quantity:  pbItem.Item.Quantity,         // 假设CartItem包含Quantity字段
			// OrderId 会在业务层自动生成，无需赋值
		})
	}

	// 调用业务层创建订单
	result, err := s.oc.PlaceOrder(ctx, &biz.PlaceOrderReq{
		UserId:       "123456", // 使用请求中的用户ID
		UserCurrency: req.UserCurrency,
		Address: biz.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country, // 添加国家字段
			ZipCode:       req.Address.ZipCode,
		},
		Items: bizItems,
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
