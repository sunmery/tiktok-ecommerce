package data

import (
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

type productRepo struct {
	data *Data
	log  *log.Helper
}

func (p productRepo) ListProducts(ctx context.Context, req biz.ListProductsReq) (*biz.ListProductsResp, error) {
	products, err := p.data.db.ListProducts(ctx, models.ListProductsParams{
		Page:     int64(req.Page),
		PageSize: req.PageSize,
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
			Categories:  product.Categories,
		}
	}
	return &biz.ListProductsResp{
		Products: productsResp,
	}, nil

}

func (p productRepo) GetProduct(ctx context.Context, req biz.GetProductReq) (*biz.GetProductResp, error) {
	// TODO implement me
	panic("implement me")
}

func (p productRepo) SearchProducts(ctx context.Context, req biz.SearchProductsReq) (*biz.SearchProductsResp, error) {
	// TODO implement me
	panic("implement me")
}
