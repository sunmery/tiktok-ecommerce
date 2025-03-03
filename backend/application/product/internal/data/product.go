package data

import (
	category "backend/api/category/v1"
	v1 "backend/api/product/v1"
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"backend/pkg/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

func (p *productRepo) CreateProduct(ctx context.Context, req *biz.CreateProductRequest) (*biz.CreateProductReply, error) {
	db := p.data.DB(ctx)

	var (
		eg         errgroup.Group
		categoryId uint64
		result     models.CreateProductRow
		attributes []byte
	)

	// 获取分类ID
	eg.Go(func() (err error) {
		getCategory, getCategoryErr := p.data.categoryClient.GetCategory(ctx, &category.GetCategoryRequest{
			// _, err = p.data.categoryClient.GetCategory(ctx, &category.GetCategoryRequest{
			Id: req.Category.CategoryId,
		})
		newCategory := &category.Category{}
		if getCategoryErr != nil {
			// 明确处理"未找到分类"的情况
			if status.Code(getCategoryErr) == codes.NotFound {
				// 创建分类时需要指定父分类（示例使用根分类）
				newCategory, err = p.data.categoryClient.CreateCategory(ctx, &category.CreateCategoryRequest{
					ParentId:  1, // 默认挂载到根分类
					Name:      req.Category.CategoryName,
					SortOrder: req.Category.SortOrder,
				})
				if err != nil {
					return fmt.Errorf("create category failed: %w", err)
				}
				categoryId = uint64(newCategory.Id)
			} else {
				return fmt.Errorf("get category failed: %w", getCategoryErr)
			}
		} else {
			categoryId = uint64(getCategory.Id)
		}

		if getCategory != nil {
			fmt.Printf("getCategory%+v", categoryId)
			categoryId = uint64(getCategory.Id)
		}
		if newCategory != nil {
			fmt.Printf("newCategory%+v", categoryId)
			categoryId = uint64(newCategory.Id)
		}
		return nil
	})

	// 创建商品
	eg.Go(func() (createErr error) {
		// 转换价格到pgtype.Numeric
		price, err := types.Float64ToNumeric(req.Price)
		if err != nil {
			return fmt.Errorf("invalid price format: %w", err)
		}

		result, createErr = db.CreateProduct(ctx, models.CreateProductParams{
			Name:        req.Name,
			Description: &req.Description,
			Price:       price,
			CategoryID:  int64(categoryId),
			Status:      int16(req.Status),
			MerchantID:  req.MerchantId,
		})
		if createErr != nil {
			return fmt.Errorf("failed to create product: %w", createErr)
		}
		return
	})

	// 创建商品图片
	eg.Go(func() error {
		// 创建图片记录
		if len(req.Images) > 0 {
			if err := p.createProductImages(ctx, result.ID, req.MerchantId, req.Images); err != nil {
				p.log.Warnf("created product but failed to create images: %v", err)
			}
		}
		return nil
	})

	// 创建属性记录
	eg.Go(func() (err error) {
		// 转成JSON
		attributes, err = json.Marshal(req.Attributes)
		if err != nil {
			return err
		}
		createProductAttributeErr := db.CreateProductAttribute(ctx, models.CreateProductAttributeParams{
			MerchantID: req.MerchantId,
			ProductID:  result.ID,
			Attributes: attributes,
		})
		if createProductAttributeErr != nil {
			return fmt.Errorf("failed to create product attribute: %w", createProductAttributeErr)
		}
		return nil
	})

	// 创建库存记录
	eg.Go(func() error {
		inventory, err := p.data.DB(ctx).CreateInventory(ctx, models.CreateInventoryParams{
			ProductID:  result.ID,
			MerchantID: req.MerchantId,
			Stock:      int32(req.Stock),
		})
		fmt.Printf("inventory%+v", inventory)
		return err
	})

	return &biz.CreateProductReply{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (p *productRepo) UpdateProduct(ctx context.Context, req *biz.UpdateProductRequest) (*biz.Product, error) {
	db := p.data.DB(ctx)

	// 获取当前版本
	current, err := db.GetProduct(ctx, models.GetProductParams{
		ID:         req.ID,
		MerchantID: req.MerchantID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, v1.ErrorProductNotFound("查询不到该商品")
		}
		return nil, v1.ErrorInvalidStatus("failed to get product: %w", err)
	}

	// 准备更新参数
	params := models.UpdateProductParams{
		ID:         req.ID,
		MerchantID: req.MerchantID,
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
		ID:         req.ID,
		MerchantID: req.MerchantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated product: %w", err)
	}

	return p.fullProductData(ctx, updated)
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

	// TODO 从分类服务获取分类信息

	var items = make([]*biz.Product, 0)
	for _, product := range listRandomProducts {
		var images []*biz.ProductImage
		if len(product.Images) > 0 {
			if err := json.Unmarshal(product.Images, &images); err != nil {
				// 处理错误或记录日志
				p.log.WithContext(ctx).Warnf("unmarshal images error: %v", err)
			}
		}

		var attributes map[string]*biz.AttributeValue
		if len(product.Attributes) > 0 {
			if err := json.Unmarshal(product.Attributes, &attributes); err != nil {
				// 处理错误或记录日志
				p.log.WithContext(ctx).Warnf("unmarshal attributes error: %v", err)
			}
		}

		price, err := types.NumericToFloat(product.Price)
		if err != nil {
			p.log.WithContext(ctx).Warnf("unmarshal price error: %v", err)
		}

		items = append(items, &biz.Product{
			ID:          product.ID,
			MerchantId:  product.MerchantID,
			Name:        product.Name,
			Price:       price,
			Description: *product.Description,
			Images:      images,
			Status:      biz.ProductStatus(product.Status),
			Category: biz.CategoryInfo{
				CategoryId: uint64(product.CategoryID),
				// CategoryName: product.,
				// SortOrder:    0,
			},
			CreatedAt:  product.CreatedAt,
			UpdatedAt:  product.UpdatedAt,
			Attributes: attributes,
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

	// 第一阶段：收集所有分类ID
	var (
		categoryIDs = make([]int64, 0, len(productsByNameRows))
		mu          sync.Mutex
	)


	// 第二阶段：批量获取分类信息
	g, ctx := errgroup.WithContext(ctx)
	var categoryMap map[int]*category.Category // 使用int类型作为key

	g.Go(func() error {
		// 调用分类微服务批量接口
		resp, err := p.data.categoryClient.BatchGetCategories(ctx, &category.BatchGetCategoriesRequest{
			Ids: categoryIDs,
		})
		if err != nil {
			return fmt.Errorf("category service failed: %w", err)
		}

		// 构建分类映射表
		categoryMap = make(map[int]*category.Category, len(resp.Categories))
		for _, c := range resp.Categories {
			categoryMap[int(c.Id)] = c // 确保类型转换正确
		}
		return nil
	})

	// 第三阶段：并行处理商品数据
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
			var attributes map[string]*biz.AttributeValue
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

			// 获取分类信息
			var cg *category.Category
			if categoryMap != nil {
				cg = categoryMap[int(product.CategoryID)] // 类型转换确保匹配
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
				Attributes: attributes,
				Inventory: biz.Inventory{
					ProductId:  product.ID,
					MerchantId: product.MerchantID,
					Stock:      product.Stock,
				},
			}

			// 添加分类信息
			if cg != nil {
				productData.Category = biz.CategoryInfo{
					CategoryId:   uint64(cg.Id),    // int -> uint64
					CategoryName: cg.Name,
					SortOrder:    cg.SortOrder,
				}
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
		Items:      products,
	}, nil
}

func (p *productRepo) GetProduct(ctx context.Context, req *biz.GetProductRequest) (*biz.Product, error) {
	db := p.data.DB(ctx)

	// 获取基础信息
	product, err := db.GetProduct(ctx, models.GetProductParams{
		ID: req.ID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, v1.ErrorProductNotFound("查询不到该商品")
		}
		return nil, v1.ErrorInvalidStatus("GetProduct 内部错误")
	}

	return p.fullProductData(ctx, product)
}

func (p *productRepo) DeleteProduct(ctx context.Context, req *biz.DeleteProductRequest) error {
	db := p.data.DB(ctx)
	err := db.SoftDeleteProduct(ctx, models.SoftDeleteProductParams{
		ID:         req.ID,
		MerchantID: req.MerchantID,
		Status:     int16(req.Status), // 下架状态
	})
	if err != nil {
		return err
	}
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
	price, _ := decimal.NewFromString(fmt.Sprintf("%s%d", product.Price.Int.String(), product.Price.Exp))

	// 组装返回结果
	return &biz.Product{
		ID:          product.ID,
		MerchantId:  product.MerchantID,
		Name:        product.Name,
		Description: *product.Description,
		Price:       float64(price.IntPart()),
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
