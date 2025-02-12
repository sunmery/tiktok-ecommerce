package data

import (
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"context"
	"net/http"
	"github.com/go-kratos/kratos/v2/log"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

func (p *productRepo) DeleteProduct(ctx context.Context, req biz.DeleteProductReq) (*biz.ProductReply, error) {
	_, err := p.data.db.DeleteProduct(ctx, models.DeleteProductParams{
		ID: int32(req.Id),
	})
	if err != nil {
		return nil, err
	}
	return &biz.ProductReply{
		Message: "OK",
		Code:    http.StatusOK,
	}, nil
}

func (p *productRepo) UpdateProduct(ctx context.Context, req biz.Product) (*biz.ProductReply, error) {
	_, err := p.data.db.UpdateProduct(ctx, models.UpdateProductParams{
		ID:          int32(req.Id),
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  req.Categories,
	})
	if err != nil {
		return nil, err
	}
	return &biz.ProductReply{
		Message: "OK",
		Code:    http.StatusOK,
	}, nil
}

func (p *productRepo) CreateProduct(ctx context.Context, req *biz.CreateProductRequest) (*biz.CreateProductReply, error) {
	// 通过 data.DB(ctx) 自动获取事务或普通连接
	db := p.data.DB(ctx)

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

	// 审计日志自动使用同一事务(相同的 ctx)
	_, err = db.CreateAuditLog(ctx, models.CreateAuditLogParams{
		Action:    "CREATE",
		ProductID: product.ID,
		Owner:     req.Owner,
		Name:      req.Username,
	})
	if err != nil {
		return nil, err
	}

	return &biz.CreateProductReply{
		Product: biz.Product{
			Id:          uint32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			CategoryId:  product.CategoryId,
		},
	}, nil
}

func (p *productRepo) ListProducts(ctx context.Context, req *biz.ListProductsReq) (*biz.ListProductsResp, error) {
	products, err := p.data.db.ListProducts(ctx, models.ListProductsParams{
		Offset: int64((req.Page - 1) * req.PageSize),
		Limit:  int64(req.PageSize),
		CategoryID: req.CategoryId,
	})
	if err != nil {
		return nil, err
	}

	productsResp := make([]*biz.Product, len(products))
	for i, product := range products {
		productsResp[i] = &biz.Product{
			Id:          uint32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			TotalStock: product.TotalStock,
			AvailableStock: product.AvailableStock,
			ReservedStock: product.ReservedStock,
			LowStockThreshold: product.LowStockThreshold,
			AllowNegative: product.AllowNegative,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			Version: product.Version,
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
		Id:          uint32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Picture:     product.Picture,
		Price:       product.Price,
		CategoryId:  product.CategoryID,
		TotalStock: product.TotalStock,
		AvailableStock: product.AvailableStock,
		ReservedStock: product.ReservedStock,
		LowStockThreshold: product.LowStockThreshold,
		AllowNegative: product.AllowNegative,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		Version: product.Version,
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
			Id:          uint32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			CategoryId:  product.CategoryID,
			TotalStock: product.TotalStock,
			AvailableStock: product.AvailableStock,
			ReservedStock: product.ReservedStock,
			LowStockThreshold: product.LowStockThreshold,
			AllowNegative: product.AllowNegative,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			Version: product.Version,
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