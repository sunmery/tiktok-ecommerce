package service

import (
	inventoryv1 "backend/api/merchant/inventory/v1"
	productv1 "backend/api/merchant/product/v1"
	"backend/application/merchant/internal/biz"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewInventoryService, NewProductService)

type InventoryService struct {
	inventoryv1.UnimplementedInventoryServer
	ic *biz.InventoryUsecase
}
type ProductService struct {
	productv1.UnimplementedProductServer
	pc *biz.ProductUsecase
}

func NewInventoryService(ic *biz.InventoryUsecase, pc *biz.ProductUsecase) *InventoryService {
	return &InventoryService{ic: ic}
}

func NewProductService(pc *biz.ProductUsecase) *ProductService {
	return &ProductService{pc: pc}
}
