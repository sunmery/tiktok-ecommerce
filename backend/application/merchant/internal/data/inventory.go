package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"backend/pkg/types"

	"backend/application/merchant/internal/data/models"

	"backend/application/merchant/internal/biz"
)

// GetProductStock 获取产品库存
func (i *inventoryRepo) GetProductStock(ctx context.Context, req *biz.GetProductStockRequest) (*biz.GetProductStockResponse, error) {
	productId := types.ToPgUUID(req.ProductId)
	merchantId := types.ToPgUUID(req.MerchantId)
	stock, err := i.data.DB(ctx).GetProductStock(ctx, models.GetProductStockParams{
		ProductID:  productId,
		MerchantID: merchantId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			i.log.WithContext(ctx).Infof(" 查询不到该商家'%s'的ID为'%s'的商品", req.MerchantId, req.ProductId)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product stock: %w", err)
	}

	return &biz.GetProductStockResponse{
		ProductId:      stock.ProductID,
		MerchantId:     stock.MerchantID,
		Stock:          stock.Stock,
		AlertThreshold: *stock.AlertThreshold,
		IsLowStock:     *stock.IsLowStock,
	}, nil
}

// UpdateProductStock 更新产品库存
func (i *inventoryRepo) UpdateProductStock(ctx context.Context, req *biz.UpdateProductStockRequest) (*biz.UpdateProductStockResponse, error) {
	// 更新库存
	productId := types.ToPgUUID(req.ProductId)
	merchantId := types.ToPgUUID(req.MerchantId)
	_, err := i.data.DB(ctx).UpdateProductStock(ctx, models.UpdateProductStockParams{
		Stock:      &req.Quantity,
		ProductID:  productId,
		MerchantID: merchantId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update product stock: %w", err)
	}

	// 记录库存调整
	_, err = i.data.DB(ctx).RecordStockAdjustment(ctx, models.RecordStockAdjustmentParams{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
		Quantity:   req.Quantity,
		Reason:     &req.Reason,
		OperatorID: req.OperatorId,
	})
	if err != nil {
		i.log.Warnf("failed to record stock adjustment: %v", err)
		// 不影响主流程，继续执行
	}

	return &biz.UpdateProductStockResponse{
		Success: true,
		Message: "Stock updated successfully",
	}, nil
}

// SetStockAlert 设置库存警报阈值
func (i *inventoryRepo) SetStockAlert(ctx context.Context, req *biz.SetStockAlertRequest) (*biz.SetStockAlertResponse, error) {
	_, err := i.data.DB(ctx).SetStockAlert(ctx, models.SetStockAlertParams{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
		Threshold:  req.Threshold,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &biz.SetStockAlertResponse{
				Success: true,
				Message: "未找到该商品，不执行任何操作",
			}, nil
		}
		return nil, fmt.Errorf("failed to set stock alert: %w", err)
	}

	return &biz.SetStockAlertResponse{
		Success: true,
		Message: "Stock alert threshold set successfully",
	}, nil
}

// GetStockAlerts 获取库存警报配置
func (i *inventoryRepo) GetStockAlerts(ctx context.Context, req *biz.GetStockAlertsRequest) (*biz.GetStockAlertsResponse, error) {
	// 获取警报配置列表
	merchantId := types.ToPgUUID(req.MerchantId)
	pageSize := (req.Page - 1) * req.PageSize
	alerts, err := i.data.DB(ctx).GetStockAlerts(ctx, models.GetStockAlertsParams{
		MerchantID: merchantId,
		Page:       &req.Page,
		PageSize:   &pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get stock alerts: %w", err)
	}

	// 获取总数
	count, err := i.data.DB(ctx).CountStockAlerts(ctx, req.MerchantId)
	if err != nil {
		return nil, fmt.Errorf("failed to count stock alerts: %w", err)
	}

	// 转换为业务模型
	result := make([]biz.StockAlert, 0, len(alerts))
	for _, alert := range alerts {
		result = append(result, biz.StockAlert{
			Id:           alert.ID,
			ProductId:    alert.ProductID,
			MerchantId:   alert.MerchantID,
			ProductName:  alert.ProductName,
			CurrentStock: alert.CurrentStock,
			Threshold:    alert.Threshold,
			CreatedAt:    alert.CreatedAt,
			UpdatedAt:    alert.UpdatedAt,
		})
	}

	return &biz.GetStockAlertsResponse{
		Alerts: result,
		Total:  uint32(count),
	}, nil
}

// GetLowStockProducts 获取低库存产品列表
func (i *inventoryRepo) GetLowStockProducts(ctx context.Context, req *biz.GetLowStockProductsRequest) (*biz.GetLowStockProductsResponse, error) {
	// 获取低库存产品列表
	merchantId := types.ToPgUUID(req.MerchantId)
	pageSize := (req.Page - 1) * req.PageSize
	products, err := i.data.DB(ctx).GetLowStockProducts(ctx, models.GetLowStockProductsParams{
		MerchantID: merchantId,
		Page:       &req.PageSize,
		PageSize:   &pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}

	// 获取总数
	count, err := i.data.DB(ctx).CountLowStockProducts(ctx, models.CountLowStockProductsParams{
		Threshold:  req.Threshold,
		MerchantID: merchantId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to count low stock products: %w", err)
	}

	// 转换为业务模型
	result := make([]biz.LowStockProduct, 0, len(products))
	for _, product := range products {
		result = append(result, biz.LowStockProduct{
			ProductId:    product.ProductID,
			MerchantId:   product.MerchantID,
			ProductName:  product.ProductName,
			CurrentStock: product.CurrentStock,
			Threshold:    *product.Threshold,
			ImageUrl:     *product.ImageUrl,
		})
	}

	return &biz.GetLowStockProductsResponse{
		Products: result,
		Total:    *count,
	}, nil
}

// RecordStockAdjustment 记录库存调整
func (i *inventoryRepo) RecordStockAdjustment(ctx context.Context, req *biz.RecordStockAdjustmentRequest) (*biz.RecordStockAdjustmentResponse, error) {
	// 记录库存调整
	adjustment, err := i.data.DB(ctx).RecordStockAdjustment(ctx, models.RecordStockAdjustmentParams{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
		Quantity:   req.Quantity,
		Reason:     &req.Reason,
		OperatorID: req.OperatorId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to record stock adjustment: %w", err)
	}

	productId := types.ToPgUUID(req.ProductId)
	merchantId := types.ToPgUUID(req.MerchantId)
	// 更新产品库存
	_, err = i.data.DB(ctx).UpdateProductStock(ctx, models.UpdateProductStockParams{
		Stock:      &req.Quantity,
		ProductID:  productId,
		MerchantID: merchantId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update product stock: %w", err)
	}

	return &biz.RecordStockAdjustmentResponse{
		Success:      true,
		Message:      "Stock adjustment recorded successfully",
		AdjustmentId: adjustment.ID,
	}, nil
}

// GetStockAdjustmentHistory 获取库存调整历史
func (i *inventoryRepo) GetStockAdjustmentHistory(ctx context.Context, req *biz.GetStockAdjustmentHistoryRequest) (*biz.GetStockAdjustmentHistoryResponse, error) {
	// 获取库存调整历史
	productId := types.ToPgUUID(req.ProductId)
	merchantId := types.ToPgUUID(req.MerchantId)
	pageSize := (req.Page - 1) * req.PageSize
	adjustments, err := i.data.DB(ctx).GetStockAdjustmentHistory(ctx, models.GetStockAdjustmentHistoryParams{
		ProductID:  productId,
		MerchantID: merchantId,
		Page:       &req.Page,
		PageSize:   &pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get stock adjustment history: %w", err)
	}

	// 获取总数
	count, err := i.data.DB(ctx).CountStockAdjustmentHistory(ctx, models.CountStockAdjustmentHistoryParams{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to count stock adjustment history: %w", err)
	}

	// 转换为业务模型
	result := make([]biz.StockAdjustment, 0, len(adjustments))
	for _, adjustment := range adjustments {
		result = append(result, biz.StockAdjustment{
			Id:          adjustment.ID,
			ProductId:   adjustment.ProductID,
			MerchantId:  adjustment.MerchantID,
			ProductName: adjustment.ProductName,
			Quantity:    adjustment.Quantity,
			Reason:      *adjustment.Reason,
			OperatorId:  adjustment.OperatorID,
			CreatedAt:   adjustment.CreatedAt,
		})
	}

	return &biz.GetStockAdjustmentHistoryResponse{
		Adjustments: result,
		Total:       uint32(count),
	}, nil
}
