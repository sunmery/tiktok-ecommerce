package service

import (
	pb "backend/api/product/v1"
	"backend/application/product/internal/biz"
	"context"
)

type ProductService struct {
	pb.UnimplementedProductCatalogServiceServer

	pc *biz.ProductUsecase
}

func NewProductService(pc *biz.ProductUsecase) *ProductService {
	return &ProductService{pc: pc}
}

func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsResp, error) {
	result, err := s.pc.ListProducts(ctx, biz.ListProductsReq{
		Page:         req.Page,
		PageSize:     req.PageSize,
		CategoryName: req.CategoryName,
	})
	if err != nil {
		return nil, err
	}

	productList := make([]*pb.Product, len(result.Products))
	for i, product := range result.Products {
		productList[i] = &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			Categories:  product.Categories,
		}
	}
	return &pb.ListProductsResp{
		Products: productList,
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductReq) (*pb.GetProductResp, error) {
	result, err := s.pc.GetProduct(ctx, biz.GetProductReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResp{
		Product: &pb.Product{
			Id:          result.Product.Id,
			Name:        result.Product.Name,
			Description: result.Product.Description,
			Picture:     result.Product.Picture,
			Price:       result.Product.Price,
			Categories:  result.Product.Categories,
		},
	}, nil
}

func (s *ProductService) SearchProducts(ctx context.Context, req *pb.SearchProductsReq) (*pb.SearchProductsResp, error) {
	result, err := s.pc.SearchProducts(ctx, biz.SearchProductsReq{
		Query: req.Query,
	})
	if err != nil {
		return nil, err
	}

	products := make([]*pb.Product, len(result.Products))
	for i, product := range result.Products {
		products[i] = &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			Categories:  product.Categories,
		}
	}

	return &pb.SearchProductsResp{
		Products: products,
	}, nil
}

// func (s *ProductCatalogServiceService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
// 	return &pb.CreateProductReply{}, nil
// }
// func (s *ProductCatalogServiceService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductReply, error) {
// 	return &pb.UpdateProductReply{}, nil
// }
// func (s *ProductCatalogServiceService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
// 	return &pb.DeleteProductReply{}, nil
// }
