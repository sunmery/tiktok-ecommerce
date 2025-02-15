package data

import (
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}



func (p *productRepo) DeleteProduct(ctx context.Context, req *biz.DeleteProductReq) (*biz.ProductReply, error) {
	db := p.data.DB(ctx)

	product, err := db.DeleteProduct(ctx, int32(req.Id))


	if err != nil {
		return nil, err
	}

	var userID uuid.UUID
	if md, ok := metadata.FromServerContext(ctx); ok {
		userIDStr := md.Get("x-md-global-userid")
		parsedUserID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		userID = parsedUserID
	}

	_, err = db.UpdateAuditLog(ctx, models.UpdateAuditLogParams{
		ChangeReason: "DELETE",
		NewStock:     0,
		UserID:       userID,
		ProductID:    int32(req.Id),
	})

	if err != nil {
		return nil, err
	}

	return &biz.ProductReply{
		Product: biz.Product{
			Id:   uint32(product.ID),
			Name: product.Name,
		},
	}, nil
}

func (p *productRepo) UpdateProduct(ctx context.Context, req *biz.UpdateProductRequest) (*biz.ProductReply, error) {
	db := p.data.DB(ctx)

	product, err := db.UpdateProduct(ctx, models.UpdateProductParams{
		ID:          int32(req.Id),
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		CategoryID:  req.CategoryId,
		TotalStock:  req.TotalStock,
	})
	if err != nil {
		return nil, err
	}

	
	var userID uuid.UUID
	if md, ok := metadata.FromServerContext(ctx); ok {
		userIDStr := md.Get("x-md-global-userid")
		parsedUserID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		userID = parsedUserID
	}

	_, err = db.UpdateAuditLog(ctx, models.UpdateAuditLogParams{
		ChangeReason: "UPDATE",
		ProductID:    int32(req.Id),
		NewStock:     req.TotalStock,
		UserID:       userID,
	})
	if err != nil {
		return nil, err
	}

	return &biz.ProductReply{
		Product: biz.Product{
			Id:          uint32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			CategoryId:  product.CategoryID,
			TotalStock:  product.TotalStock,
		},
	}, nil
}

func (p *productRepo) CreateProduct(ctx context.Context, req *biz.CreateProductRequest) (*biz.ProductReply, error) {
	// 通过 data.DB(ctx) 自动获取事务或普通连接
	db := p.data.DB(ctx)

	// 生成分类树
	resp, err := p.ListCategories(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	categoryTree := resp.Categories

	// 将分类转换为 map，方便查找
	categoryMap := make(map[uint32]*biz.Category)
	var flatten func([]biz.Category)
	flatten = func(categories []biz.Category) {
		for _, cat := range categories {
			categoryMap[cat.ID] = &cat
			flatten(cat.Children) // 递归存储所有分类
		}
	}
	flatten(categoryTree)

	// 遍历商品的所有分类，检查是否都是叶子分类
	for _, categoryID := range req.CategoryId {
		category, exists := categoryMap[uint32(categoryID)]
		if !exists {
			return nil, fmt.Errorf("分类 ID %d 不存在", categoryID)
		}

		// 不是叶子分类则报错
		if len(category.Children) > 0 {
			return nil, fmt.Errorf("分类 ID %d 不是叶子分类，无法添加商品", categoryID)
		}
	}

	product, err := db.CreateProduct(ctx, models.CreateProductParams{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		CategoryID:  req.CategoryId,
		TotalStock:  req.TotalStock,
	})
	if err != nil {
		return nil, err
	}

	
	var userID uuid.UUID
	if md, ok := metadata.FromServerContext(ctx); ok {
		userIDStr := md.Get("x-md-global-userid")
		parsedUserID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		userID = parsedUserID
	}

	// 审计日志自动使用同一事务(相同的 ctx)
	_, err = db.CreateAuditLog(ctx, models.CreateAuditLogParams{
		ChangeReason: "CREATE",
		ProductID:    product.ID,
		NewStock:     product.TotalStock,
		OldStock:     product.TotalStock,
		UserID:       userID,
	})
	if err != nil {
		return nil, err
	}

	return &biz.ProductReply{
		Product: biz.Product{
			Id:          uint32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			CategoryId:  product.CategoryID,
			TotalStock:  product.TotalStock,
		},
	}, nil
}

func (p *productRepo) ListProducts(ctx context.Context, req *biz.ListProductsReq) (*biz.ListProductsResp, error) {
	products, err := p.data.db.ListProducts(ctx, models.ListProductsParams{
		Offset:     int64((req.Page - 1) * req.PageSize),
		Limit:      int64(req.PageSize),
		CategoryID: req.CategoryId,
	})
	if err != nil {
		return nil, err
	}

	productsResp := make([]*biz.Product, len(products))
	for i, product := range products {
		productsResp[i] = &biz.Product{
			Id:                uint32(product.ID),
			Name:              product.Name,
			Description:       product.Description,
			Picture:           product.Picture,
			Price:             product.Price,
			TotalStock:        product.TotalStock,
			AvailableStock:    product.AvailableStock,
			ReservedStock:     product.ReservedStock,
			LowStockThreshold: product.LowStockThreshold,
			AllowNegative:     product.AllowNegative,
			CreatedAt:         product.CreatedAt,
			UpdatedAt:         product.UpdatedAt,
			Version:           product.Version,
		}
	}
	return &biz.ListProductsResp{
		Product: productsResp,
	}, nil
}

func (p *productRepo) GetProduct(ctx context.Context, id uint32) (*biz.GetProductResp, error) {
	product, err := p.data.db.GetProduct(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return &biz.GetProductResp{Product: &biz.Product{
		Id:                uint32(product.ID),
		Name:              product.Name,
		Description:       product.Description,
		Picture:           product.Picture,
		Price:             product.Price,
		CategoryId:        product.CategoryID,
		TotalStock:        product.TotalStock,
		AvailableStock:    product.AvailableStock,
		ReservedStock:     product.ReservedStock,
		LowStockThreshold: product.LowStockThreshold,
		AllowNegative:     product.AllowNegative,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
		Version:           product.Version,
	}}, nil
}

func (p *productRepo) SearchProducts(ctx context.Context, req *biz.SearchProductsReq) (*biz.SearchProductsResp, error) {

	products, err := p.data.db.SearchProducts(ctx, &req.Query)
	if err != nil {
		return nil, err
	}
	productsResp := make([]*biz.Product, len(products))
	for i, product := range products {
		productsResp[i] = &biz.Product{
			Id:                uint32(product.ID),
			Name:              product.Name,
			Description:       product.Description,
			Picture:           product.Picture,
			Price:             product.Price,
			CategoryId:        product.CategoryID,
			TotalStock:        product.TotalStock,
			AvailableStock:    product.AvailableStock,
			ReservedStock:     product.ReservedStock,
			LowStockThreshold: product.LowStockThreshold,
			AllowNegative:     product.AllowNegative,
			CreatedAt:         product.CreatedAt,
			UpdatedAt:         product.UpdatedAt,
			Version:           product.Version,
		}
	}
	return &biz.SearchProductsResp{
		Result: productsResp,
	}, nil
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
