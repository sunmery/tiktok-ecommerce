package service

import (
	cartv1 "backend/api/cart/v1"
	orderv1 "backend/api/order/v1"
	"backend/application/order/internal/biz"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewOrderService)

type OrderService struct {
	orderv1.UnimplementedOrderServiceServer
	cartClient cartv1.CartServiceClient
	oc         *biz.OrderUsecase
}

func NewOrderService(oc *biz.OrderUsecase) *OrderService {
	return &OrderService{oc: oc}
}
