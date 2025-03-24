package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GetProductStockRequest 获取产品库存请求
type GetProductStockRequest struct {
	ProductId  uuid.UUID
	MerchantId uuid.UUID
}

// 获取产品库存响应
type GetProductStockResponse struct {
	ProductId      uuid.UUID
	MerchantId     uuid.UUID
	Stock          int32
	AlertThreshold int32
	IsLowStock     bool
}

// 更新产品库存请求
type UpdateProductStockRequest struct {
	ProductId  uuid.UUID
	MerchantId uuid.UUID
	Quantity   int32 // 调整数量（正数增加，负数减少）
	Reason     string
	OperatorId uuid.UUID
}

// 更新产品库存响应
type UpdateProductStockResponse struct {
	Success bool
	Message string
}

// 设置库存警报阈值请求
type SetStockAlertRequest struct {
	ProductId  uuid.UUID
	MerchantId uuid.UUID
	Threshold  int32
}

// 设置库存警报阈值响应
type SetStockAlertResponse struct {
	Success bool
	Message string
}

// 获取库存警报配置请求
type GetStockAlertsRequest struct {
	MerchantId uuid.UUID
	Page       int32
	PageSize   int32
}

// 库存警报配置
type StockAlert struct {
	Id           uuid.UUID
	ProductId    uuid.UUID
	MerchantId   uuid.UUID
	ProductName  string
	CurrentStock int32
	Threshold    int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// 获取库存警报配置响应
type GetStockAlertsResponse struct {
	Alerts []StockAlert
	Total  uint32
}

// 获取低库存产品请求
type GetLowStockProductsRequest struct {
	MerchantId uuid.UUID
	Threshold  *int32
	Page       int64
	PageSize   int64
}

// 低库存产品信息
type LowStockProduct struct {
	ProductId    uuid.UUID
	MerchantId   uuid.UUID
	ProductName  string
	CurrentStock int32
	Threshold    int32
	ImageUrl     string
}

// 获取低库存产品响应
type GetLowStockProductsResponse struct {
	Products []LowStockProduct
	Total    int64
}

// 记录库存调整请求
type RecordStockAdjustmentRequest struct {
	ProductId  uuid.UUID
	MerchantId uuid.UUID
	Quantity   int32
	Reason     string
	OperatorId uuid.UUID
}

// 记录库存调整响应
type RecordStockAdjustmentResponse struct {
	Success      bool
	Message      string
	AdjustmentId uuid.UUID
}

// 获取库存调整历史请求
type GetStockAdjustmentHistoryRequest struct {
	ProductId  uuid.UUID
	MerchantId uuid.UUID
	Page       int64
	PageSize   int64
}

// 库存调整记录
type StockAdjustment struct {
	Id          uuid.UUID
	ProductId   uuid.UUID
	MerchantId  uuid.UUID
	ProductName string
	Quantity    int32
	Reason      string
	OperatorId  uuid.UUID
	CreatedAt   time.Time
}

// 获取库存调整历史响应
type GetStockAdjustmentHistoryResponse struct {
	Adjustments []StockAdjustment
	Total       uint32
}

// 库存
type Inventory struct {
	ProductId  uuid.UUID
	MerchantId uuid.UUID
	Stock      int32
}

// GetProductStock 获取产品库存
func (uc *InventoryUsecase) GetProductStock(ctx context.Context, req *GetProductStockRequest) (*GetProductStockResponse, error) {
	return uc.repo.GetProductStock(ctx, req)
}

// UpdateProductStock 更新产品库存
func (uc *InventoryUsecase) UpdateProductStock(ctx context.Context, req *UpdateProductStockRequest) (*UpdateProductStockResponse, error) {
	return uc.repo.UpdateProductStock(ctx, req)
}

// SetStockAlert 设置库存警报阈值
func (uc *InventoryUsecase) SetStockAlert(ctx context.Context, req *SetStockAlertRequest) (*SetStockAlertResponse, error) {
	uc.log.WithContext(ctx).Debugf("设置库存警报阈值: %v", req)
	return uc.repo.SetStockAlert(ctx, req)
}

// GetStockAlerts 获取库存警报配置
func (uc *InventoryUsecase) GetStockAlerts(ctx context.Context, req *GetStockAlertsRequest) (*GetStockAlertsResponse, error) {
	return uc.repo.GetStockAlerts(ctx, req)
}

// GetLowStockProducts 获取低库存产品列表
func (uc *InventoryUsecase) GetLowStockProducts(ctx context.Context, req *GetLowStockProductsRequest) (*GetLowStockProductsResponse, error) {
	return uc.repo.GetLowStockProducts(ctx, req)
}

// RecordStockAdjustment 记录库存调整
func (uc *InventoryUsecase) RecordStockAdjustment(ctx context.Context, req *RecordStockAdjustmentRequest) (*RecordStockAdjustmentResponse, error) {
	return uc.repo.RecordStockAdjustment(ctx, req)
}

// GetStockAdjustmentHistory 获取库存调整历史
func (uc *InventoryUsecase) GetStockAdjustmentHistory(ctx context.Context, req *GetStockAdjustmentHistoryRequest) (*GetStockAdjustmentHistoryResponse, error) {
	return uc.repo.GetStockAdjustmentHistory(ctx, req)
}
