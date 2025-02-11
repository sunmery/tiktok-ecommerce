package data

import (
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"context"
	"net/http"
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

func (p *productRepo) CreateProduct(ctx context.Context, req biz.Product) (*biz.ProductReply, error) {
	db := p.data.DB(ctx)

	_, err := db.CreateProduct(ctx, models.CreateProductParams{
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

func (p productRepo) ListProducts(ctx context.Context, req biz.ListProductsReq) (*biz.ListProductsResp, error) {
	products, err := p.data.db.ListProducts(ctx, models.ListProductsParams{
		Page:     int64(req.Page),
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	productsResp := make([]biz.Product, len(products))
	for i, product := range products {
		productsResp[i] = biz.Product{
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
	productInfo, err := p.data.db.GetProduct(ctx, int32(req.Id))

	if err != nil {
		return nil, err
	}

	return &biz.GetProductResp{
		Product: biz.Product{
			Id:          uint32(productInfo.ID),
			Name:        productInfo.Name,
			Description: productInfo.Description,
			Picture:     productInfo.Picture,
			Price:       productInfo.Price,
			Categories:  productInfo.Categories,
		},
	}, nil
}

func (p productRepo) SearchProducts(ctx context.Context, req biz.SearchProductsReq) (*biz.SearchProductsResp, error) {
	products, err := p.data.db.SearchProducts(ctx, &req.Query)
	if err != nil {
		return nil, err
	}
	productsResp := make([]biz.Product, len(products))
	for i, product := range products {
		productsResp[i] = biz.Product{
			Id:          uint32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			Categories:  product.Categories,
		}
	}
	return &biz.SearchProductsResp{
		Products: productsResp,
	}, nil
}
