package service

import (
	"backend/application/product/internal/biz"
	"context"

	pb "backend/api/product/v1"
)

type ProductCatalogServiceService struct {
	pb.UnimplementedProductCatalogServiceServer
	pc *biz.ProductUsecase
}

func NewProductCatalogServiceService(pc *biz.ProductUsecase) *ProductCatalogServiceService {
	return &ProductCatalogServiceService{pc: pc}
}

func (s *ProductCatalogServiceService) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsResp, error) {

	return &pb.ListProductsResp{}, nil
}
func (s *ProductCatalogServiceService) GetProduct(ctx context.Context, req *pb.GetProductReq) (*pb.GetProductResp, error) {
	return &pb.GetProductResp{}, nil
}
func (s *ProductCatalogServiceService) SearchProducts(ctx context.Context, req *pb.SearchProductsReq) (*pb.SearchProductsResp, error) {
	return &pb.SearchProductsResp{}, nil
}
