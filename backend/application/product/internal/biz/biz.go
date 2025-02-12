package biz

import (
	"context"
	"github.com/google/wire"

	"github.com/go-kratos/kratos/v2/log"
)

var ProviderSet = wire.NewSet(NewProductUsecase)

type ProductRepo interface {
	ListProducts(ctx context.Context, req *ListProductsReq) (*ListProductsResp, error)
	GetProduct(ctx context.Context, id uint32) (*GetProductResp, error)
	SearchProducts(ctx context.Context, req *SearchProductsReq) (*SearchProductsResp, error)
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductReply, error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*ProductReply, error)
	DeleteProduct(ctx context.Context, req *DeleteProductReq) (*ProductReply, error)
}

// ProductUsecase is a Product usecase.
type ProductUsecase struct {
	repo ProductRepo
	log  *log.Helper
}

// NewProductUsecase new a Product usecase.
func NewProductUsecase(repo ProductRepo, logger log.Logger) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}



