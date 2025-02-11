package biz

import (
	"context"
	"errors"

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
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Picture     string  `json:"picture,omitempty"`
	Price       float32 `json:"price,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Owner       string  `json:"owner,omitempty"`
	Username	string  `json:"username,omitempty"`
}
type ListProductsResp struct {
	Products []Product `json:"products,omitempty"`
}

type ListProductsReq struct {
	Page         int32  `json:"page,omitempty"`
	PageSize     int64  `json:"page_size,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
}

type GetProductReq struct {
	Id uint32 `json:"id,omitempty"`
}

type GetProductResp struct {
	Product Product `json:"product,omitempty"`
}

type SearchProductsReq struct {
	Query string `json:"query,omitempty"`
}

type SearchProductsResp struct {
	Products []Product `json:"products,omitempty"`
}

// type CreateProductReq struct {
// 	Name        string  `json:"name,omitempty"`
// 	Description string  `json:"description,omitempty"`
// 	Picture     string  `json:"picture,omitempty"`
// 	Price       float32 `json:"price,omitempty"`
// 	Categories  []string `json:"categories,omitempty"`
// 	Owner       string  `json:"owner,omitempty"`
// 	Username	string  `json:"username,omitempty"`
// }

type ProductReply struct {
	Message string `json:"message,omitempty"`
	Code    int32  `json:"code,omitempty"`
}

type DeleteProductReq struct {
	Id uint32
}

type ProductRepo interface {
	ListProducts(ctx context.Context, req ListProductsReq) (*ListProductsResp, error)
	GetProduct(ctx context.Context, req GetProductReq) (*GetProductResp, error)
	SearchProducts(ctx context.Context, req SearchProductsReq) (*SearchProductsResp, error)
	CreateProduct(ctx context.Context, req Product) (*ProductReply, error)
	UpdateProduct(ctx context.Context, req Product) (*ProductReply, error)
	DeleteProduct(ctx context.Context, req DeleteProductReq) (*ProductReply, error)
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

func (s *ProductUsecase) ListProducts(ctx context.Context, req ListProductsReq) (*ListProductsResp, error) {
	s.log.WithContext(ctx).Debugf("ListProducts %v", req)

	if req.Page < 1 || req.PageSize < 1 {
		return nil, errors.New("invalid pagination parameters")
	}

	resp, err := s.repo.ListProducts(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to list products: %v", err)
		return nil, err
	}
	return resp, nil
}
func (s *ProductUsecase) GetProduct(ctx context.Context, req GetProductReq) (*GetProductResp, error) {
	s.log.WithContext(ctx).Debugf("GetProduct %v", req)
	
	resp, err := s.repo.GetProduct(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get product: %v", err)
		return nil, err
	}
	return resp, nil
}
func (s *ProductUsecase) SearchProducts(ctx context.Context, req SearchProductsReq) (*SearchProductsResp, error) {
	s.log.WithContext(ctx).Debugf("SearchProducts %v", req)
	
	if req.Query == "" {
		return nil, errors.New("search query cannot be empty")
	}

	resp, err := s.repo.SearchProducts(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to search products: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *ProductUsecase) CreateProduct(ctx context.Context, req Product) (*ProductReply, error) {
	s.log.WithContext(ctx).Debugf("CreateProduct %v", req)
	
	if req.Name == "" || req.Price <= 0 || len(req.Categories) == 0 || req.Categories[0] == "" || req.Picture == "" || req.Description == "" {
		return &ProductReply{}, errors.New("invalid product data")
	}
	resp, err := s.repo.CreateProduct(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to create product: %v", err)
		return &ProductReply{}, err
	}
	return resp, nil
}

func (s *ProductUsecase) UpdateProduct(ctx context.Context, req Product) (*ProductReply, error) {
	s.log.WithContext(ctx).Debugf("UpdateProduct %v", req)
	
	if req.Name == "" || req.Price <= 0 || len(req.Categories) == 0 || req.Categories[0] == "" || req.Picture == "" || req.Description == "" {
		return nil, errors.New("invalid product data")
	}

	resp, err := s.repo.UpdateProduct(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to update product: %v", err)
		return resp, err
	}
	return resp, nil
}

func (s *ProductUsecase) DeleteProduct(ctx context.Context, req DeleteProductReq) (*ProductReply, error) {
	s.log.WithContext(ctx).Debugf("DeleteProduct %v", req)
	
	if req.Id == 0 {
		return nil, errors.New("product ID cannot be empty")
	}

	resp, err := s.repo.DeleteProduct(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to delete product: %v", err)
		return resp, err
	}
	return resp, nil
}

