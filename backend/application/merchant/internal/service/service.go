package service

import (
	inventoryv1 "backend/api/merchant/inventory/v1"
	orderv1 "backend/api/merchant/order/v1"
	productv1 "backend/api/merchant/product/v1"
	"backend/application/merchant/internal/biz"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewInventoryService, NewProductService, NewOrderService)

type InventoryService struct {
	inventoryv1.UnimplementedInventoryServer
	ic *biz.InventoryUsecase
}
type ProductService struct {
	productv1.UnimplementedProductServer
	pc *biz.ProductUsecase
}

type OrderServiceService struct {
	orderv1.UnimplementedOrderServer
	oc *biz.OrderUsecase
}

func NewInventoryService(ic *biz.InventoryUsecase) *InventoryService {
	return &InventoryService{ic: ic}
}

func NewProductService(pc *biz.ProductUsecase) *ProductService {
	return &ProductService{pc: pc}
}

func NewOrderService(oc *biz.OrderUsecase) *OrderServiceService {
	return &OrderServiceService{oc: oc}
}
