package biz

import (
	"context"
	"errors"
	"time"
)

type Product struct {
	Id          uint32	`json:"id"`
	Name        string	`json:"name"`
	Description string` json:"description"`
	Picture     string	`json:"picture"`
	Price       float32	`json:"price"`
	CategoryId  []int32		`json:"categoryId"`
	TotalStock        int32     `json:"totalStock"`
	AvailableStock    *int32    `json:"availableStock"`
	ReservedStock     int32     `json:"reservedStock"`
	LowStockThreshold int32     `json:"lowStockThreshold"`
	AllowNegative     bool      `json:"allowNegative"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Version           int32     `json:"version"`
}

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	Price       float32  `json:"price"`
	CategoryId []int32 `json:"categoryId"`
	TotalStock  int32   `json:"totalStock"`

	Owner    string `json:"owner"`
	Username string `json:"username"`
}

type CreateProductReply struct {
	Product Product
}

type ListProductsReq struct {
	Page         uint   `json:"page"`
	PageSize     uint   `json:"pageSize"`
	CategoryId int32 `json:"categoryId"`
}

type ListProductsResp struct {
	Product []*Product `json:"product"`
}

type GetProductResp struct {
	Product *Product `json:"product"`
}

type SearchProductsReq struct {
	Query string `json:"query"`
}
type SearchProductsResp struct {
	Result []*Product `json:"result"`
}

func (s *ProductUsecase) ListProducts(ctx context.Context, req *ListProductsReq) (*ListProductsResp, error) {
	s.log.WithContext(ctx).Infof("ListProducts %v", req)

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

func (s *ProductUsecase) GetProduct(ctx context.Context, id uint32) (*GetProductResp, error) {
	s.log.WithContext(ctx).Infof("GetProduct %v", id)
	
	resp, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get product: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *ProductUsecase) SearchProducts(ctx context.Context, req *SearchProductsReq) (*SearchProductsResp, error) {
	s.log.WithContext(ctx).Infof("SearchProducts %v", req)
	
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

func (s *ProductUsecase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductReply, error) {
	s.log.WithContext(ctx).Infof("CreateProduct %v", req)
	return s.repo.CreateProduct(ctx, req)
}

func (s *ProductUsecase) UpdateProduct(ctx context.Context, req Product) (*ProductReply, error) {
	s.log.WithContext(ctx).Infof("UpdateProduct %v", req)
	return s.repo.UpdateProduct(ctx, req)
}

func (s *ProductUsecase) DeleteProduct(ctx context.Context, req DeleteProductReq) (*ProductReply, error) {
	s.log.WithContext(ctx).Infof("DeleteProduct %v", req)
	
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
