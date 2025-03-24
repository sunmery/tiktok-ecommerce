package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewInventoryUsecase, NewProductUsecase)

const (
	ProductStatusDraft    ProductStatus = iota // 商品草稿
	ProductStatusPending                       // 商品待审核。
	ProductStatusApproved                      // 商品审核通过。
	ProductStatusRejected                      // 商品审核未通过。
	ProductStatusSoldOut                       // 商品因某种原因不可购买。
)

// InventoryRepo 库存域方法
type InventoryRepo interface {
	// GetProductStock 获取产品库存
	GetProductStock(ctx context.Context, req *GetProductStockRequest) (*GetProductStockResponse, error)
	// UpdateProductStock 更新产品库存
	UpdateProductStock(ctx context.Context, req *UpdateProductStockRequest) (*UpdateProductStockResponse, error)
	// SetStockAlert 设置库存警报阈值
	SetStockAlert(ctx context.Context, req *SetStockAlertRequest) (*SetStockAlertResponse, error)
	// GetStockAlerts 获取库存警报配置
	GetStockAlerts(ctx context.Context, req *GetStockAlertsRequest) (*GetStockAlertsResponse, error)
	// GetLowStockProducts 获取低库存产品列表
	GetLowStockProducts(ctx context.Context, req *GetLowStockProductsRequest) (*GetLowStockProductsResponse, error)
	// RecordStockAdjustment 记录库存调整
	RecordStockAdjustment(ctx context.Context, req *RecordStockAdjustmentRequest) (*RecordStockAdjustmentResponse, error)
	// GetStockAdjustmentHistory 获取库存调整历史
	GetStockAdjustmentHistory(ctx context.Context, req *GetStockAdjustmentHistoryRequest) (*GetStockAdjustmentHistoryResponse, error)
}

// ProductRepo 商品域方法
type ProductRepo interface {
	// GetMerchantProducts 获取商家自身商品列表
	GetMerchantProducts(ctx context.Context, req *GetMerchantProducts) (*Products, error)
}

type InventoryUsecase struct {
	repo InventoryRepo
	log  *log.Helper
}

type ProductUsecase struct {
	repo ProductRepo
	log  *log.Helper
}

func NewInventoryUsecase(repo InventoryRepo, logger log.Logger) *InventoryUsecase {
	return &InventoryUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func NewProductUsecase(repo ProductRepo, logger log.Logger) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
