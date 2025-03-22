package service

import (
	"context"
	"errors"
	"fmt"

	"backend/pkg"

	"github.com/google/uuid"

	"backend/application/merchant/internal/biz"

	v1 "backend/api/merchant/inventory/v1"
)

func (uc *InventoryService) GetProductStock(ctx context.Context, req *v1.GetProductStockRequest) (*v1.GetProductStockResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid product id: '%s' error: %v", req.ProductId, err))
	}

	stock, err := uc.ic.GetProductStock(ctx, &biz.GetProductStockRequest{
		ProductId:  productId,
		MerchantId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.GetProductStockResponse{
		ProductId:      stock.ProductId.String(),
		MerchantId:     stock.MerchantId.String(),
		Stock:          stock.Stock,
		AlertThreshold: stock.AlertThreshold,
		IsLowStock:     stock.IsLowStock,
	}, nil
}

func (uc *InventoryService) UpdateProductStock(ctx context.Context, req *v1.UpdateProductStockRequest) (*v1.UpdateProductStockResponse, error) {
	return &v1.UpdateProductStockResponse{}, nil
}

func (uc *InventoryService) SetStockAlert(ctx context.Context, req *v1.SetStockAlertRequest) (*v1.SetStockAlertResponse, error) {
	return &v1.SetStockAlertResponse{}, nil
}

func (uc *InventoryService) GetStockAlerts(ctx context.Context, req *v1.GetStockAlertsRequest) (*v1.GetStockAlertsResponse, error) {
	return &v1.GetStockAlertsResponse{}, nil
}

func (uc *InventoryService) GetLowStockProducts(ctx context.Context, req *v1.GetLowStockProductsRequest) (*v1.GetLowStockProductsResponse, error) {
	return &v1.GetLowStockProductsResponse{}, nil
}

func (uc *InventoryService) RecordStockAdjustment(ctx context.Context, req *v1.RecordStockAdjustmentRequest) (*v1.RecordStockAdjustmentResponse, error) {
	return &v1.RecordStockAdjustmentResponse{}, nil
}

func (uc *InventoryService) GetStockAdjustmentHistory(ctx context.Context, req *v1.GetStockAdjustmentHistoryRequest) (*v1.GetStockAdjustmentHistoryResponse, error) {
	return &v1.GetStockAdjustmentHistoryResponse{}, nil
}
