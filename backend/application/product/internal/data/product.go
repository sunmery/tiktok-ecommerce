package data

import (
	v1 "backend/api/category/v1"
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"context"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

func (p *productRepo) CreateProduct(ctx context.Context, req biz.CreateProductRequest) (*biz.Product, error) {
	db := p.data.DB(ctx)

	// 转换价格到pgtype.Numeric
	price, err := decimal.NewFromString(fmt.Sprintf("%.2f", req.Product.Price))
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	// TODO 创建分类
	categoryID, err := strconv.ParseInt(req.Product.Category.CategoryId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID format: %w", err)
	}

	category, err := p.data.categoryClient.GetCategory(ctx, &v1.GetCategoryRequest{
		Id: categoryID,
	})
	if err != nil{
		return nil,err
	}

	if category == nil {
		category, err = p.data.categoryClient.CreateCategory(ctx, &v1.CreateCategoryRequest{
			Name:      req.Product.Category.CategoryName,
			ParentId:  0,
			SortOrder: 0,
		})
	}

	// 执行创建
	result, err := db.CreateProduct(ctx, models.CreateProductParams{
		Name:        req.Product.Name,
		Description: &req.Product.Description,
		Price:       pgtype.Numeric{Int: price.Coefficient(), Exp: price.Exponent()},
		Status:      int16(req.Product.Status),
		MerchantID:  int64(req.Product.MerchantId),
		CategoryID:  category.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// 转换基础信息
	product := biz.Product{
		ID: uint64(result.ID),
		// MerchantId: uint64(result.MerchantID),
		Name:      req.Product.Name,
		Price:     req.Product.Price,
		Stock:     req.Product.Stock,
		Status:    req.Product.Status,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	// 创建图片记录
	if len(req.Product.Images) > 0 {
		if err := p.createProductImages(ctx, uint64(result.ID), req.Product.MerchantId, req.Product.Images); err != nil {
			p.log.Warnf("created product but failed to create images: %v", err)
		}
	}

	return &product, nil
}

func (p *productRepo) UpdateProduct(ctx context.Context, req biz.UpdateProductRequest) (*biz.Product, error) {
	db := p.data.DB(ctx)

	// 获取当前版本
	current, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         int64(req.ID),
		MerchantID: int64(req.MerchantID),
	})
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// 准备更新参数
	params := models.UpdateProductParams{
		ID:         int64(req.ID),
		MerchantID: int64(req.MerchantID),
		UpdatedAt:  pgtype.Timestamptz{Time: current.UpdatedAt, Valid: true},
	}

	// 字段掩码处理
	if req.Name != nil {
		params.Name = *req.Name
	} else {
		params.Name = current.Name
	}

	if req.Price != nil {
		price, err := decimal.NewFromString(fmt.Sprintf("%.2f", *req.Price))
		if err != nil {
			return nil, fmt.Errorf("invalid price format: %w", err)
		}
		params.Price = pgtype.Numeric{Int: price.Coefficient(), Exp: price.Exponent()}
	} else {
		params.Price = current.Price
	}

	// params.Stock = current.Stock
	// stock := int32(req.Stock)
	// if req.Stock != nil {
	// 	params.Stock = &stock
	// } else {
	// 	params.Stock = current.Stock
	// }

	if req.Description != "" {
		params.Description = &req.Description
	} else {
		params.Description = current.Description
	}

	// 执行更新
	if err := db.UpdateProduct(ctx, params); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// 获取更新后的数据
	updated, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         int64(req.ID),
		MerchantID: int64(req.MerchantID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated product: %w", err)
	}

	return p.fullProductData(ctx, updated)
}

func (p *productRepo) SubmitForAudit(ctx context.Context, req biz.SubmitAuditRequest) (*biz.AuditRecord, error) {
	db := p.data.DB(ctx)

	// 获取当前产品状态
	current, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         int64(req.ProductID),
		MerchantID: int64(req.MerchantID),
	})
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// 创建审核记录
	auditRecord, err := db.GetLatestAudit(ctx, models.GetLatestAuditParams{
		MerchantID: int64(req.MerchantID),
		ProductID:  int64(req.ProductID),
		OldStatus:  current.Status,
		NewStatus:  int16(biz.ProductStatusPending),
		Reason:     nil,
		OperatorID: 0, // 从上下文中获取实际操作人
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create audit record: %w", err)
	}

	// 更新产品状态
	if err := db.UpdateProductStatus(ctx, models.UpdateProductStatusParams{
		ID:             int64(req.ProductID),
		Status:         int16(biz.ProductStatusPending),
		CurrentAuditID: &auditRecord.ID,
		MerchantID:     int64(req.MerchantID),
	}); err != nil {
		return nil, fmt.Errorf("failed to update product status: %w", err)
	}

	return &biz.AuditRecord{
		ID:         uint64(auditRecord.ID),
		ProductID:  req.ProductID,
		OldStatus:  biz.ProductStatus(current.Status),
		NewStatus:  biz.ProductStatusPending,
		OperatedAt: auditRecord.CreatedAt,
	}, nil
}

func (p *productRepo) AuditProduct(ctx context.Context, req biz.AuditProductRequest) (*biz.AuditRecord, error) {
	db := p.data.DB(ctx)

	// 获取当前产品状态
	current, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         int64(req.ProductID),
		MerchantID: int64(req.MerchantID),
	})
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// 确定新状态
	var newStatus biz.ProductStatus
	switch biz.AuditAction(req.Action) { // 添加类型转换
	case biz.AuditActionApprove:
		newStatus = biz.ProductStatusApproved
	case biz.AuditActionReject:
		newStatus = biz.ProductStatusRejected
	default:
		return nil, biz.ErrInvalidAuditAction
	}

	// 创建审核记录
	auditRecord, err := db.GetLatestAudit(ctx, models.GetLatestAuditParams{
		MerchantID: int64(req.MerchantID),
		ProductID:  int64(req.ProductID),
		OldStatus:  current.Status,
		NewStatus:  int16(newStatus),
		Reason:     &req.Reason,
		OperatorID: int64(req.OperatorID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create audit record: %w", err)
	}

	// 更新产品状态
	if err := db.UpdateProductStatus(ctx, models.UpdateProductStatusParams{
		ID:             int64(req.ProductID),
		Status:         int16(newStatus),
		CurrentAuditID: &auditRecord.ID,
		MerchantID:     int64(req.MerchantID),
	}); err != nil {
		return nil, fmt.Errorf("failed to update product status: %w", err)
	}

	return &biz.AuditRecord{
		ID:         uint64(auditRecord.ID),
		ProductID:  req.ProductID,
		OldStatus:  biz.ProductStatus(current.Status),
		NewStatus:  newStatus,
		Reason:     req.Reason,
		OperatorID: req.OperatorID,
		OperatedAt: auditRecord.CreatedAt,
	}, nil
}

func (p *productRepo) GetProduct(ctx context.Context, req biz.GetProductRequest) (*biz.Product, error) {
	db := p.data.DB(ctx)

	// 获取基础信息
	product, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         int64(req.ID),
		MerchantID: int64(req.MerchantID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return p.fullProductData(ctx, product)
}

func (p *productRepo) DeleteProduct(ctx context.Context, req biz.DeleteProductRequest) error {
	db := p.data.DB(ctx)

	return db.SoftDeleteProduct(ctx, models.SoftDeleteProductParams{
		ID:         int64(req.ID),
		MerchantID: int64(req.MerchantID),
	})
}

// 辅助方法：完整产品数据组装
func (p *productRepo) fullProductData(ctx context.Context, product models.GetProductRow) (*biz.Product, error) {
	// 获取图片
	images, err := p.data.DB(ctx).GetProductImages(ctx, models.GetProductImagesParams{
		MerchantID: product.MerchantID,
		ProductID:  product.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}

	// 转换价格
	price, _ := decimal.NewFromString(fmt.Sprintf("%s%d", product.Price.Int.String(), product.Price.Exp))

	// 组装返回结果
	return &biz.Product{
		ID:          uint64(product.ID),
		MerchantId:  uint64(product.MerchantID),
		Name:        product.Name,
		Description: *product.Description,
		Price:       float64(price.IntPart()),
		// Stock:       *product.Stock,
		Status:      biz.ProductStatus(product.Status),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Images:      p.convertImages(images),
		// 其他字段根据实际需求补充
	}, nil
}

// 转换图片数据
func (p *productRepo) convertImages(images []models.ProductsProductImages) []*biz.ProductImage {
	result := make([]*biz.ProductImage, 0, len(images))
	for _, img := range images {
		sortOrder := 0
		result = append(result, &biz.ProductImage{
			URL:       img.Url,
			IsPrimary: img.IsPrimary,
			SortOrder: &sortOrder,
		})
	}
	return result
}

// 创建产品图片
func (p *productRepo) createProductImages(ctx context.Context, productID uint64, merchantID uint64, images []*biz.ProductImage) error {
	if len(images) == 0 {
		return nil
	}

	bulkParams := models.BulkCreateProductImagesParams{
		MerchantIds: make([]int64, len(images)),
		ProductIds:  make([]int64, len(images)),
		Urls:        make([]string, len(images)),
		IsPrimary:   make([]bool, len(images)),
		SortOrders:  make([]int16, len(images)), // 数据库字段类型为 SMALLINT
	}

	for i, img := range images {
		bulkParams.MerchantIds[i] = int64(merchantID)
		bulkParams.ProductIds[i] = int64(productID)
		bulkParams.Urls[i] = img.URL
		bulkParams.IsPrimary[i] = img.IsPrimary

		// 处理 SortOrder
		if img.SortOrder != nil {
			bulkParams.SortOrders[i] = int16(*img.SortOrder) // 解引用指针并转换类型
		} else {
			bulkParams.SortOrders[i] = 0 // 默认值
		}
	}

	return p.data.DB(ctx).BulkCreateProductImages(ctx, bulkParams)
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
