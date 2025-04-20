package service

import (
	"context"
	"errors"
	"fmt"

	"backend/pkg"

	"github.com/google/uuid"

	"backend/application/merchant/internal/biz"

	v1 "backend/api/merchant/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InventoryService struct {
	v1.UnimplementedInventoryServer
	ic *biz.InventoryUsecase
}

func NewInventoryService(ic *biz.InventoryUsecase) *InventoryService {
	return &InventoryService{ic: ic}
}

func (uc *InventoryService) GetProductStock(ctx context.Context, req *v1.GetProductStockRequest) (*v1.GetProductStockResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("invalid product id: '%s' error: %v", req.ProductId, err)
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
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid product id: '%s' error: %v", req.ProductId, err))
	}

	result, err := uc.ic.UpdateProductStock(ctx, &biz.UpdateProductStockRequest{
		ProductId:  productId,
		MerchantId: userId,
		Quantity:   req.Stock,
		Reason:     req.Reason,
		OperatorId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateProductStockResponse{
		Success: result.Success,
		Message: result.Message,
	}, nil
}

func (uc *InventoryService) SetStockAlert(ctx context.Context, req *v1.SetStockAlertRequest) (*v1.SetStockAlertResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid product id: '%s' error: %v", req.ProductId, err))
	}

	result, err := uc.ic.SetStockAlert(ctx, &biz.SetStockAlertRequest{
		ProductId:  productId,
		MerchantId: userId,
		Threshold:  req.Threshold,
	})
	if err != nil {
		return nil, err
	}

	return &v1.SetStockAlertResponse{
		Success: result.Success,
		Message: result.Message,
	}, nil
}

func (uc *InventoryService) GetStockAlerts(ctx context.Context, req *v1.GetStockAlertsRequest) (*v1.GetStockAlertsResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := uc.ic.GetStockAlerts(ctx, &biz.GetStockAlertsRequest{
		MerchantId: userId,
		Page:       int64(req.Page),
		PageSize:   int64(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	var alerts []*v1.StockAlert
	for _, alert := range result.Alerts {
		alerts = append(alerts, &v1.StockAlert{
			ProductId:    alert.ProductId.String(),
			MerchantId:   alert.MerchantId.String(),
			ProductName:  alert.ProductName,
			CurrentStock: alert.CurrentStock,
			Threshold:    alert.Threshold,
			CreatedAt:    timestamppb.New(alert.CreatedAt),
			UpdatedAt:    timestamppb.New(alert.UpdatedAt),
		})
	}

	return &v1.GetStockAlertsResponse{
		Alerts: alerts,
		Total:  result.Total,
	}, nil
}

func (uc *InventoryService) GetLowStockProducts(ctx context.Context, req *v1.GetLowStockProductsRequest) (*v1.GetLowStockProductsResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	if req.Threshold <= 0 {
		req.Threshold = 10 // 使用系统默认阈值
	}

	result, err := uc.ic.GetLowStockProducts(ctx, &biz.GetLowStockProductsRequest{
		MerchantId: userId,
		Page:       int64(req.Page),
		PageSize:   int64(req.PageSize),
		Threshold:  int32(req.Threshold),
	})
	if err != nil {
		return nil, err
	}

	var products []*v1.LowStockProduct
	for _, product := range result.Products {
		products = append(products, &v1.LowStockProduct{
			ProductId:    product.ProductId.String(),
			MerchantId:   product.MerchantId.String(),
			ProductName:  product.ProductName,
			CurrentStock: product.CurrentStock,
			Threshold:    product.Threshold,
			ImageUrl:     product.ImageUrl,
		})
	}

	return &v1.GetLowStockProductsResponse{
		Products: products,
		Total:    uint32(result.Total),
	}, nil
}

func (uc *InventoryService) RecordStockAdjustment(ctx context.Context, req *v1.RecordStockAdjustmentRequest) (*v1.RecordStockAdjustmentResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid product id: '%s' error: %v", req.ProductId, err))
	}
	operatorId, err := uuid.Parse(req.OperatorId)
	if err != nil {
		operatorId = userId // 如果操作者ID无效，使用当前用户ID
	}

	result, err := uc.ic.RecordStockAdjustment(ctx, &biz.RecordStockAdjustmentRequest{
		ProductId:  productId,
		MerchantId: userId,
		Quantity:   req.Quantity,
		Reason:     req.Reason,
		OperatorId: operatorId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RecordStockAdjustmentResponse{
		Success:      result.Success,
		Message:      result.Message,
		AdjustmentId: result.AdjustmentId.String(),
	}, nil
}

func (uc *InventoryService) GetStockAdjustmentHistory(ctx context.Context, req *v1.GetStockAdjustmentHistoryRequest) (*v1.GetStockAdjustmentHistoryResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	// productId, err := uuid.Parse(req.ProductId)
	// if err != nil {
	// 	return nil, errors.New(fmt.Sprintf("invalid product id: '%s' error: %v", req.ProductId, err))
	// }

	result, err := uc.ic.GetStockAdjustmentHistory(ctx, &biz.GetStockAdjustmentHistoryRequest{
		// ProductId:  productId,
		MerchantId: userId,
		Page:       int64(req.Page),
		PageSize:   int64(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	var adjustments []*v1.StockAdjustment
	for _, adjustment := range result.Adjustments {
		adjustments = append(adjustments, &v1.StockAdjustment{
			Id:          adjustment.Id.String(),
			ProductId:   adjustment.ProductId.String(),
			MerchantId:  adjustment.MerchantId.String(),
			ProductName: adjustment.ProductName,
			Quantity:    adjustment.Quantity,
			Reason:      adjustment.Reason,
			OperatorId:  adjustment.OperatorId.String(),
			CreatedAt:   timestamppb.New(adjustment.CreatedAt),
		})
	}

	return &v1.GetStockAdjustmentHistoryResponse{
		Adjustments: adjustments,
		Total:       uint32(result.Total),
	}, nil
}
