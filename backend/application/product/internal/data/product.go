package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"

	category "backend/api/category/v1"
	v1 "backend/api/product/v1"
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"backend/pkg/types"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
)

const (
	defaultExpiryTime = time.Second * 24 * 60 * 60 // 1 day
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

func (p *productRepo) GetCategoryProducts(ctx context.Context, req *biz.GetCategoryProducts) (*biz.Products, error) {
	products, err := p.data.DB(ctx).GetCategoryProducts(ctx, models.GetCategoryProductsParams{
		CategoryID: int64(req.CategoryID),
		Status:     int16(req.Status),
		Limit:      req.PageSize,
		Offset:     (req.Page - 1) * req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// 收集所有不同的分类ID
	categoryIDs := make([]int64, 0)
	categoryIDMap := make(map[int64]bool)
	for _, product := range products {
		if !categoryIDMap[product.CategoryID] {
			categoryIDMap[product.CategoryID] = true
			categoryIDs = append(categoryIDs, product.CategoryID)
		}
	}

	// 从分类服务获取分类信息
	var categoryMap map[int64]*category.Category
	if len(categoryIDs) > 0 {
		categoriesResp, err := p.data.categoryClient.BatchGetCategories(ctx, &category.BatchGetCategoriesRequest{
			Ids: categoryIDs,
		})
		if err != nil {
			p.log.WithContext(ctx).Warnf("failed to get categories: %v", err)
		} else {
			categoryMap = make(map[int64]*category.Category)
			for _, cat := range categoriesResp.Categories {
				categoryMap[cat.Id] = cat
			}
		}
	}

	items := make([]*biz.Product, 0)
	for _, product := range products {
		var images []*biz.ProductImage
		if len(product.Images) > 0 {
			if err := json.Unmarshal(product.Images, &images); err != nil {
				// 处理错误或记录日志
				p.log.WithContext(ctx).Warnf("unmarshal images error: %v", err)
			}
		}

		price, err := types.NumericToFloat(product.Price)
		if err != nil {
			p.log.WithContext(ctx).Warnf("unmarshal price error: %v", err)
		}

		// 构建分类信息
		categoryInfo := biz.CategoryInfo{
			CategoryId: uint64(product.CategoryID),
		}

		// 如果找到了分类信息，则设置分类名称
		if c, ok := categoryMap[product.CategoryID]; ok {
			categoryInfo.CategoryName = c.Name
			categoryInfo.SortOrder = c.SortOrder
		}

		// 处理商品属性
		var attributes map[string]any
		if len(product.Attributes) > 0 {
			if err := json.Unmarshal(product.Attributes, &attributes); err != nil {
				p.log.WithContext(ctx).Warnf("unmarshal attributes error: %v", err)
			}
		}

		items = append(items, &biz.Product{
			ID:          product.ID,
			MerchantId:  product.MerchantID,
			Name:        product.Name,
			Price:       price,
			Description: *product.Description,
			Images:      images,
			Status:      biz.ProductStatus(product.Status),
			Category:    categoryInfo,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			Attributes:  attributes,
		})
	}

	return &biz.Products{Items: items}, err
}

func (p *productRepo) UploadProductFile(ctx context.Context, req *biz.UploadProductFileRequest) (*biz.UploadProductFileReply, error) {
	expiry := defaultExpiryTime

	policy := minio.NewPostPolicy()
	_ = policy.SetBucket(*req.BucketName)
	_ = policy.SetKey(*req.FileName)
	_ = policy.SetExpires(time.Now().UTC().Add(expiry))
	presignedURL, formData, err := p.data.minio.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}

	url, err := p.data.minio.PresignedPutObject(ctx, *req.BucketName, *req.FileName, expiry)
	if err != nil {
		return nil, err
	}

	return &biz.UploadProductFileReply{
		UploadUrl:   presignedURL.String(),
		DownloadUrl: url.String(),
		BucketName:  req.BucketName,
		ObjectName:  *req.FileName,
		FormData:    formData,
	}, nil
}

func (p *productRepo) CreateProduct(ctx context.Context, req *biz.CreateProductRequest) (_ *biz.CreateProductReply, err error) {
	var (
		result            models.CreateProductRow
		categoryID        uint64
		createdCategoryID uint64 // 记录新创建的分类 ID（用于补偿）
	)
	// 获取事务版 DB 操作
	db := p.data.DB(ctx)

	// Step 1: 获取或创建分类（跨服务操作）
	getCategory, err := p.data.categoryClient.GetCategory(ctx, &category.GetCategoryRequest{
		Id: req.Category.CategoryId,
	})
	if status.Code(err) == codes.NotFound {
		// 创建新分类（跨服务操作）
		newCategory, createErr := p.data.categoryClient.CreateCategory(ctx, &category.CreateCategoryRequest{
			Name:      req.Category.CategoryName,
			SortOrder: req.Category.SortOrder,
		})
		if createErr != nil {
			return nil, fmt.Errorf("create category failed: %w", createErr)
		}
		categoryID = uint64(newCategory.Id)
		createdCategoryID = categoryID // 记录新创建的分类 ID
	} else if err != nil {
		return nil, fmt.Errorf("get category failed: %w", err)
	}
	fmt.Println("getCategory", getCategory)
	fmt.Println("createdCategoryID", createdCategoryID)

	// else {
	// 	categoryID = uint64(getCategory.Id)
	// }

	// Step 2: 执行本地事务（商品相关操作）

	// // 执行本地数据库事务（包裹商品、图片、属性、库存）
	// txErr := p.data.ExecTx(ctx, func(ctx context.Context) error {
	// 注意：这里使用 ctx 作为上下文

	// 1. 创建商品
	price, err := types.Float64ToNumeric(req.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	result, err = db.CreateProduct(ctx, models.CreateProductParams{
		Name:        req.Name,
		Description: &req.Description,
		Price:       price,
		Status:      int16(req.Status),
		MerchantID:  req.MerchantId,
		CategoryID:  int64(req.Category.CategoryId), // 假设分类 ID 已存在

	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// 2. 并行创建图片、属性、库存
	var eg errgroup.Group

	// 插入图片
	eg.Go(func() error {
		if len(req.Images) > 0 {
			return p.createProductImages(ctx, result.ID, req.MerchantId, req.Images)
		}
		return nil
	})

	// 插入属性
	eg.Go(func() error {
		attributes, err := json.Marshal(req.Attributes)
		if err != nil {
			return fmt.Errorf("marshal attributes failed: %w", err)
		}
		return db.CreateProductAttribute(ctx, models.CreateProductAttributeParams{
			MerchantID: req.MerchantId,
			ProductID:  result.ID,
			Attributes: attributes,
		})
	})

	// 插入库存
	eg.Go(func() error {
		_, err = db.CreateInventory(ctx, models.CreateInventoryParams{
			ProductID:  result.ID,
			MerchantID: req.MerchantId,
			Stock:      int32(req.Stock),
		})
		return err
	})

	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	return &biz.CreateProductReply{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (p *productRepo) SubmitForAudit(ctx context.Context, req *biz.SubmitAuditRequest) (*biz.AuditRecord, error) {
	db := p.data.DB(ctx)

	// 获取当前产品状态
	current, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         req.ProductID,
		MerchantID: req.MerchantID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, v1.ErrorProductNotFound("查询不到该商品")
		}
	}

	// 创建审核记录
	auditRecord, err := db.GetLatestAudit(ctx, models.GetLatestAuditParams{
		MerchantID: req.MerchantID,
		ProductID:  req.ProductID,
		OldStatus:  current.Status,
		NewStatus:  int16(biz.ProductStatusPending),
		Reason:     nil,
		OperatorID: req.OperatorID, // 从上下文中获取实际操作人
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create audit record: %w", err)
	}

	// 更新产品状态
	if err := db.UpdateProductStatus(ctx, models.UpdateProductStatusParams{
		ID:             req.ProductID,
		Status:         int16(biz.ProductStatusPending),
		CurrentAuditID: types.ToPgUUID(auditRecord.ID),
		MerchantID:     req.MerchantID,
	}); err != nil {
		return nil, fmt.Errorf("failed to update product status: %w", err)
	}

	return &biz.AuditRecord{
		ID:         auditRecord.ID,
		ProductID:  req.ProductID,
		OldStatus:  biz.ProductStatus(current.Status),
		NewStatus:  biz.ProductStatusPending,
		OperatedAt: auditRecord.CreatedAt,
	}, nil
}

// AuditProduct 审核商品
func (p *productRepo) AuditProduct(ctx context.Context, req *biz.AuditProductRequest) (*biz.AuditRecord, error) {
	db := p.data.DB(ctx)

	// 获取当前产品状态
	fmt.Printf("req%+v", req)
	current, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         req.ProductID,
		MerchantID: req.MerchantID,
	})
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, v1.ErrorProductNotFound("查询不到该商品")
		}
		return nil, v1.ErrorInvalidStatus("data/product: AuditProduct: %+v", err)
	}

	// 确定新状态
	var newStatus biz.ProductStatus
	switch biz.AuditAction(req.Action) { // 添加类型转换
	case biz.Approved:
		newStatus = biz.ProductStatusApproved
	case biz.Rejected:
		newStatus = biz.ProductStatusRejected
	default:
		p.log.Warnf("非法的Action行为: %v", req.Action)
		return nil, v1.ErrorInvalidAuditAction("AuditProduct: 非法的req.Action参数, 1为通过, 2为驳回")
	}

	// 创建审核记录
	fmt.Printf("current%+v", current)
	auditRecord, err := db.CreateAuditRecord(ctx, models.CreateAuditRecordParams{
		MerchantID: req.MerchantID,
		ProductID:  req.ProductID,
		OldStatus:  current.Status,
		NewStatus:  int16(newStatus),
		Reason:     &req.Reason,
		OperatorID: req.OperatorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create audit record: %w", err)
	}

	// 更新产品状态
	currentAuditID := types.ToPgUUID(auditRecord.ID)
	if err := db.UpdateProductStatus(ctx, models.UpdateProductStatusParams{
		ID:             req.ProductID,
		Status:         int16(newStatus),
		CurrentAuditID: currentAuditID,
		MerchantID:     req.MerchantID,
	}); err != nil {
		return nil, fmt.Errorf("failed to update product status: %w", err)
	}

	return &biz.AuditRecord{
		ID:         auditRecord.ID,
		ProductID:  req.ProductID,
		OldStatus:  biz.ProductStatus(current.Status),
		NewStatus:  newStatus,
		Reason:     req.Reason,
		OperatorID: req.OperatorID,
		OperatedAt: auditRecord.CreatedAt,
	}, nil
}

func (p *productRepo) ListRandomProducts(ctx context.Context, req *biz.ListRandomProductsRequest) (*biz.Products, error) {
	offset := (req.Page - 1) * req.PageSize
	listRandomProducts, err := p.data.DB(ctx).ListRandomProducts(ctx, models.ListRandomProductsParams{
		Status: int16(req.Status),
		Limit:  int64(req.PageSize),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, err
	}

	// 收集所有不同的分类ID
	categoryIDs := make([]int64, 0)
	categoryIDMap := make(map[int64]bool)
	for _, product := range listRandomProducts {
		if !categoryIDMap[product.CategoryID] {
			categoryIDMap[product.CategoryID] = true
			categoryIDs = append(categoryIDs, product.CategoryID)
		}
	}

	// 从分类服务获取分类信息
	var categoryMap map[int64]*category.Category
	if len(categoryIDs) > 0 {
		categoriesResp, err := p.data.categoryClient.BatchGetCategories(ctx, &category.BatchGetCategoriesRequest{
			Ids: categoryIDs,
		})
		if err != nil {
			p.log.WithContext(ctx).Warnf("failed to get categories: %v", err)
		} else {
			categoryMap = make(map[int64]*category.Category)
			for _, cat := range categoriesResp.Categories {
				categoryMap[cat.Id] = cat
			}
		}
	}

	items := make([]*biz.Product, 0)
	for _, product := range listRandomProducts {
		var images []*biz.ProductImage
		if len(product.Images) > 0 {
			if err := json.Unmarshal(product.Images, &images); err != nil {
				// 处理错误或记录日志
				p.log.WithContext(ctx).Warnf("unmarshal images error: %v", err)
			}
		}

		price, err := types.NumericToFloat(product.Price)
		if err != nil {
			p.log.WithContext(ctx).Warnf("unmarshal price error: %v", err)
		}

		// 构建分类信息
		categoryInfo := biz.CategoryInfo{
			CategoryId: uint64(product.CategoryID),
		}

		// 如果找到了分类信息，则设置分类名称
		if c, ok := categoryMap[product.CategoryID]; ok {
			categoryInfo.CategoryName = c.Name
			categoryInfo.SortOrder = c.SortOrder
		}

		// 处理商品属性
		var attributes map[string]any
		if len(product.Attributes) > 0 {
			if err := json.Unmarshal(product.Attributes, &attributes); err != nil {
				p.log.WithContext(ctx).Warnf("unmarshal attributes error: %v", err)
				attributes = nil
			}
		} else {
			attributes = nil
		}

		items = append(items, &biz.Product{
			ID:          product.ID,
			MerchantId:  product.MerchantID,
			Name:        product.Name,
			Price:       price,
			Description: *product.Description,
			Images:      images,
			Status:      biz.ProductStatus(product.Status),
			Category:    categoryInfo,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			Attributes:  attributes,
			Inventory: biz.Inventory{
				ProductId:  product.ID,
				MerchantId: product.MerchantID,
				Stock:      uint32(product.Stock),
			},
		})
	}

	return &biz.Products{Items: items}, err
}

func (p *productRepo) SearchProductsByName(ctx context.Context, req *biz.SearchProductsByNameRequest) (*biz.Products, error) {
	// 参数转换
	page := (req.Page - 1) * req.PageSize
	params := models.SearchFullProductsByNameParams{
		Name:     &req.Name,
		Query:    req.Query,
		PageSize: int64(req.PageSize),
		Page:     int64(page),
	}

	// 执行数据库查询
	productsByNameRows, err := p.data.DB(ctx).SearchFullProductsByName(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	// 注意：SearchFullProductsByNameRow结构体不包含CategoryID字段，所以无法获取分类信息

	// 处理商品数据
	g, ctx := errgroup.WithContext(ctx)
	var (
		products = make([]*biz.Product, 0, len(productsByNameRows))
		pMu      sync.Mutex
	)

	for _, product := range productsByNameRows {
		product := product // 创建局部变量
		g.Go(func() error {
			// 处理商品图片
			var images []*biz.ProductImage
			if len(product.Images) > 0 {
				if err := json.Unmarshal(product.Images, &images); err != nil {
					p.log.WithContext(ctx).Warnf("unmarshal images error: %v", err)
				}
			}

			// 处理商品属性
			var attributes map[string]any
			if len(product.Attributes) > 0 {
				if err := json.Unmarshal(product.Attributes, &attributes); err != nil {
					p.log.WithContext(ctx).Warnf("unmarshal attributes error: %v", err)
				}
			}

			// 处理价格
			price, err := types.NumericToFloat(product.Price)
			if err != nil {
				p.log.WithContext(ctx).Warnf("price conversion error: %v", err)
				price = 0 // 设置默认值
			}

			// 构建商品对象
			productData := &biz.Product{
				ID:          product.ID,
				MerchantId:  product.MerchantID,
				Name:        product.Name,
				Price:       price,
				Description: *product.Description,
				Images:      images,
				Status:      biz.ProductStatus(product.Status),
				CreatedAt:   product.CreatedAt,
				UpdatedAt:   product.UpdatedAt,
				Attributes:  attributes,
				Category:    biz.CategoryInfo{},
				Inventory: biz.Inventory{
					ProductId:  product.ID,
					MerchantId: product.MerchantID,
					Stock:      uint32(product.Stock),
				},
			}

			// 安全添加至结果集
			pMu.Lock()
			products = append(products, productData)
			pMu.Unlock()
			return nil
		})
	}

	// 等待所有goroutine完成
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("parallel processing failed: %w", err)
	}

	return &biz.Products{
		Items: products,
	}, nil
}

func (p *productRepo) GetProduct(ctx context.Context, req *biz.GetProductRequest) (*biz.Product, error) {
	db := p.data.DB(ctx)

	// 获取基础信息
	product, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         req.ID,
		MerchantID: req.MerchantID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 从分类服务获取分类信息
	var categoryInfo *category.Category
	if product.CategoryID > 0 {
		categoryResp, err := p.data.categoryClient.GetCategory(ctx, &category.GetCategoryRequest{
			Id: uint64(product.CategoryID),
		})
		if err != nil {
			p.log.WithContext(ctx).Warnf("failed to get category: %v", err)
		} else {
			categoryInfo = categoryResp
		}
	}

	result, err := p.fullProductData(ctx, product)
	if err != nil {
		return nil, err
	}

	// 补充分类信息
	if categoryInfo != nil {
		result.Category.CategoryName = categoryInfo.Name
		result.Category.SortOrder = categoryInfo.SortOrder
	}

	return result, nil
}

func (p *productRepo) GetProductBatch(ctx context.Context, req *biz.GetProductsBatchRequest) (*biz.Products, error) {
	db := p.data.DB(ctx)

	// 获取基础信息
	products, err := db.GetProductsBatch(ctx, models.GetProductsBatchParams{
		ProductIds:  req.ProductIds,
		MerchantIds: req.MerchantIds,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &biz.Products{Items: nil}, nil
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 收集所有不同的分类ID
	categoryIDs := make([]int64, 0)
	categoryIDMap := make(map[int64]bool)
	for _, product := range products {
		if !categoryIDMap[product.CategoryID] {
			categoryIDMap[product.CategoryID] = true
			categoryIDs = append(categoryIDs, product.CategoryID)
		}
	}

	// 从分类服务获取分类信息
	var categoryMap map[int64]*category.Category
	if len(categoryIDs) > 0 {
		categoriesResp, err := p.data.categoryClient.BatchGetCategories(ctx, &category.BatchGetCategoriesRequest{
			Ids: categoryIDs,
		})
		if err != nil {
			p.log.WithContext(ctx).Warnf("failed to get categories: %v", err)
		} else {
			categoryMap = make(map[int64]*category.Category)
			for _, cat := range categoriesResp.Categories {
				categoryMap[cat.Id] = cat
			}
		}
	}

	items := make([]*biz.Product, 0)
	for _, product := range products {
		var images []*biz.ProductImage
		if len(product.Images) > 0 {
			if err := json.Unmarshal(product.Images, &images); err != nil {
				// 处理错误或记录日志
				p.log.WithContext(ctx).Warnf("unmarshal images error: %v", err)
			}
		}

		price, err := types.NumericToFloat(product.Price)
		if err != nil {
			p.log.WithContext(ctx).Warnf("unmarshal price error: %v", err)
		}

		// 构建分类信息
		categoryInfo := biz.CategoryInfo{
			CategoryId: uint64(product.CategoryID),
		}

		// 如果找到了分类信息，则设置分类名称
		if c, ok := categoryMap[product.CategoryID]; ok {
			categoryInfo.CategoryName = c.Name
			categoryInfo.SortOrder = c.SortOrder
		}
		// 处理商品属性
		var attributes map[string]any
		if len(product.Attributes) > 0 {
			if err := json.Unmarshal(product.Attributes, &attributes); err != nil {
				p.log.WithContext(ctx).Warnf("unmarshal attributes error: %v", err)
			}
		}
		items = append(items, &biz.Product{
			ID:          product.ID,
			MerchantId:  product.MerchantID,
			Name:        product.Name,
			Price:       price,
			Description: *product.Description,
			Images:      images,
			Status:      biz.ProductStatus(product.Status),
			Category:    categoryInfo,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			Attributes:  attributes,
			Inventory: biz.Inventory{
				ProductId:  product.ID,
				MerchantId: product.MerchantID,
				Stock:      uint32(product.Stock),
			},
		})
	}

	return &biz.Products{Items: items}, err
}

func (p *productRepo) DeleteProduct(ctx context.Context, req *biz.DeleteProductRequest) error {
	db := p.data.DB(ctx)
	result, err := db.SoftDeleteProduct(ctx, models.SoftDeleteProductParams{
		ID:         req.ID,
		MerchantID: req.MerchantID,
		Status:     int16(req.Status), // 下架状态
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		return err
	}
	log.Debugf("result: %+v", result)
	return nil
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
	price, err := types.NumericToFloat(product.Price)
	if err != nil {
		p.log.WithContext(ctx).Warnf("unmarshal price error: %v", err)
	}

	// 处理属性
	var attributes map[string]interface{}
	if len(product.Attributes) > 0 {
		if err := json.Unmarshal(product.Attributes, &attributes); err != nil {
			p.log.WithContext(ctx).Warnf("unmarshal attributes error: %v", err)
			attributes = nil
		}
	} else {
		attributes = nil
	}

	// 组装返回结果
	return &biz.Product{
		ID:          product.ID,
		MerchantId:  product.MerchantID,
		Name:        product.Name,
		Price:       price,
		Description: *product.Description,
		Images:      convertImages(images),
		Status:      biz.ProductStatus(product.Status),
		Category: biz.CategoryInfo{
			CategoryId: uint64(product.CategoryID),
		},
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
		Attributes: attributes,
		Inventory: biz.Inventory{
			ProductId:  product.ID,
			MerchantId: product.MerchantID,
			Stock:      uint32(product.Stock),
		},
	}, nil
}

// 转换图片数据
func convertImages(images []models.ProductsProductImages) []*biz.ProductImage {
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
func (p *productRepo) createProductImages(ctx context.Context, productID uuid.UUID, merchantID uuid.UUID, images []*biz.ProductImage) error {
	if len(images) == 0 {
		return nil
	}

	bulkParams := models.BulkCreateProductImagesParams{
		MerchantIds: make([]uuid.UUID, len(images)),
		ProductIds:  make([]uuid.UUID, len(images)),
		Urls:        make([]string, len(images)),
		IsPrimary:   make([]bool, len(images)),
		SortOrders:  make([]int16, len(images)), // 数据库字段类型为 SMALLINT
	}

	for i, img := range images {
		bulkParams.MerchantIds[i] = merchantID
		bulkParams.ProductIds[i] = productID
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
