package service

import (
	"backend/application/product/internal/biz"

	"context"

	pb "backend/api/product/v1"
)

type ProductCatalogServiceService struct {
	pb.UnimplementedProductCatalogServiceServer

	uc *biz.ProductUsecase
}

func NewProductCatalogServiceService(uc *biz.ProductUsecase) *ProductCatalogServiceService {
	return &ProductCatalogServiceService{
		uc: uc,
	}
}

func (s *ProductCatalogServiceService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductReply, error) {

	p, cErr := s.uc.CreateProduct(ctx, &biz.CreateProductRequest{
		Name:      req.Name,
		Price:     float64(req.Price),
		MainImage: req.Picture,
		Images:    nil,
		Stock:     req.Stock,
		// Category:  req.Categories,
	})
	if cErr != nil {
		return nil, cErr
	}
	return &pb.ProductReply{
		Product: &pb.Product{

			// Id:   p.ID,
			Name: p.Name,
			// Description: p,
			Picture: p.MainImage,
			// Price:       p.Price,
			Categories: nil,
		},
	}, nil
}

func (s *ProductCatalogServiceService) GetProduct(context.Context, *pb.GetProductReq) (*pb.ProductReply, error) {
	panic("todo")
}
func (s *ProductCatalogServiceService) ListProducts(context.Context, *pb.ListProductsReq) (*pb.ListProductsResp, error) {
	panic("todo")
}
func (s *ProductCatalogServiceService) SearchProducts(context.Context, *pb.SearchProductsReq) (*pb.SearchProductsResp, error) {
	panic("todo")
}
func (s *ProductCatalogServiceService) UpdateProduct(context.Context, *pb.Product) (*pb.ProductReply, error) {
	panic("todo")
}
