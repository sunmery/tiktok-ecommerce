package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewProductUsecase)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type Product struct {
	Id          uint32
	Name        string
	Description string
	Picture     string
	Price       float32
	Categories  []string
}
type ListProductsResp struct {
	Products []Product
}

type GetProductReq struct {
}

type GetProductResp struct {
}

type SearchProductsReq struct {
}

type SearchProductsResp struct {
}

type ProductRepo interface {
	ListProducts(ctx context.Context, req ListProductsReq) (*ListProductsResp, error)
	GetProduct(ctx context.Context, req GetProductReq) (*GetProductResp, error)
	SearchProducts(ctx context.Context, req SearchProductsReq) (*SearchProductsResp, error)
}

type ProductUsecase struct {
	repo ProductRepo
	log  *log.Helper
}

func NewProductUsecase(repo ProductRepo, logger log.Logger) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

type ListProductsReq struct {
	Page         int32
	PageSize     int64
	CategoryName string
}

func (s *ProductUsecase) ListProducts(ctx context.Context, req ListProductsReq) (*ListProductsResp, error) {
	s.log.WithContext(ctx).Debugf("ListProducts %v", req)
	return s.repo.ListProducts(ctx, req)
}
func (s *ProductUsecase) GetProduct(ctx context.Context, req GetProductReq) (*GetProductResp, error) {
	s.log.WithContext(ctx).Debugf("GetProduct %v", req)
	return s.repo.GetProduct(ctx, req)
}
func (s *ProductUsecase) SearchProducts(ctx context.Context, req SearchProductsReq) (*SearchProductsResp, error) {
	s.log.WithContext(ctx).Debugf("SearchProducts %v", req)
	return s.repo.SearchProducts(ctx, req)
}
