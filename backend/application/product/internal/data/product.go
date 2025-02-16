package data

import (
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

func (p *productRepo) CreateProduct(ctx context.Context, req *biz.Product) (*biz.Product, error) {
	product, err := p.data.DB(ctx).CreateProduct(ctx, models.CreateProductParams{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       float32(req.Price),
		Stock:       req.Stock,
		// CategoryID:  req.CategoryID,
	})
	if err != nil {
		return nil, err
	}
	return &biz.Product{
		ID:          product.ID,
		Name:        product.Name,
		Price:       float64(product.Price),
		Picture:     product.Picture,
		Description: product.Description,
		Stock:       product.Stock,
		// CategoryID:  product.CategoryID,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, err
}

func (p *productRepo) AddProduct(ctx context.Context, req *biz.AddProductRequest) (*biz.Product, error) {
	// TODO implement me
	panic("implement me")
}

func (p *productRepo) ListProducts(ctx context.Context, req *biz.ListProductsReq) (*biz.ListProductsResp, error) {
	// products, err := p.data.db.ListProducts(ctx, models.ListProductsParams{
	// 	Page:     int64((req.Page - 1) * req.PageSize),
	// 	PageSize: int64(req.PageSize),
	// })
	// if err != nil {
	// 	return nil, err
	// }
	//
	// productsResp := make([]*biz.Product, len(products))
	// for i, product := range products {
	// 	productsResp[i] = &biz.Product{
	// 		ID:         uuid.UUID{},
	// 		Name:       "",
	// 		Price:      decimal.Decimal{},
	// 		MainImage:  "",
	// 		Images:     nil,
	// 		Stock:      0,
	// 		CategoryID: uuid.UUID{},
	// 		CreatedAt:  time.Time{},
	// 		UpdatedAt:  time.Time{},
	// 	}
	// }
	// return &biz.ListProductsResp{
	// 	Product: productsResp,
	// }, nil
	// TODO implement me
	panic("implement me")
}

func (p *productRepo) GetProduct(ctx context.Context, id uint32) (*biz.GetProductResp, error) {
	// product, err := p.data.db.GetProduct(ctx, int32(id))
	// if err != nil {
	// 	return nil, err
	// }
	// return &biz.GetProductResp{Product: &biz.Product{
	// 	Id:          uint32(product.ID),
	// 	Name:        product.Name,
	// 	Description: product.Description,
	// 	Picture:     product.Picture,
	// 	Price:       product.Price,
	// 	Categories:  product.Categories,
	// }}, nil
	// TODO implement me
	panic("implement me")
}

func (p *productRepo) SearchProducts(ctx context.Context, req *biz.SearchProductsReq) (*biz.SearchProductsResp, error) {

	// products, err := p.data.db.SearchProducts(ctx, &req.Query)
	// if err != nil {
	// 	return nil, err
	// }
	// productsResp := make([]*biz.Product, len(products))
	// for i, product := range products {
	// 	productsResp[i] = &biz.Product{
	// 		Id:          uint32(product.ID),
	// 		Name:        product.Name,
	// 		Description: product.Description,
	// 		Picture:     product.Picture,
	// 		Price:       product.Price,
	// 		Categories:  product.Categories,
	// 	}
	// }
	// return &biz.SearchProductsResp{
	// 	Result: productsResp,
	// }, nil
	// TODO implement me
	panic("implement me")
}

// NewProductRepo .
func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
