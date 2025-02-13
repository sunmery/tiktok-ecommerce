package biz

import "context"

type Product struct {
	Id          uint32
	Name        string
	Description string
	Picture     string
	Price       float32
	Categories  []string
}

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	Price       float32  `json:"price"`
	Categories  []string `json:"categories"`

	Owner    string `json:"owner"`
	Username string `json:"username"`
}

type CreateProductReply struct {
	Product Product
}

type ListProductsReq struct {
	Page         uint   `json:"page"`
	PageSize     uint   `json:"pageSize"`
	CategoryName string `json:"categoryName"`
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

func (uc *ProductUsecase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductReply, error) {
	uc.log.WithContext(ctx).Infof("CreateProduct: %v", req)
	return uc.repo.CreateProduct(ctx, req)
}

func (uc *ProductUsecase) ListProducts(ctx context.Context, req *ListProductsReq) (*ListProductsResp, error) {
	uc.log.WithContext(ctx).Infof("ListProducts: %v", req)
	return uc.repo.ListProducts(ctx, req)
}

func (uc *ProductUsecase) GetProduct(ctx context.Context, id uint32) (*GetProductResp, error) {
	uc.log.WithContext(ctx).Infof("GetProductReq: %v", id)
	return uc.repo.GetProduct(ctx, id)
}

func (uc *ProductUsecase) SearchProducts(ctx context.Context, req *SearchProductsReq) (*SearchProductsResp, error) {
	uc.log.WithContext(ctx).Infof("SearchProducts: %v", req)
	return uc.repo.SearchProducts(ctx, req)
}
