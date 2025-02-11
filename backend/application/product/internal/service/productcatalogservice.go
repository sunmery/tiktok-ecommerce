package service

import (
	pb "backend/api/product/v1"
	"backend/application/product/internal/biz"
	"backend/pkg/token"
	"context"
	"errors"
	"fmt"
)

type ProductService struct {
	pb.UnimplementedProductCatalogServiceServer
	pu *biz.ProductUsecase
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
	products, err := s.pu.SearchProducts(ctx, &biz.SearchProductsReq{Query: req.GetQuery()})
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

func (s *ProductCatalogServiceService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductReply, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	p, cErr := s.pu.CreateProduct(ctx, &biz.CreateProductRequest{
		Owner:       payload.Owner,
		Username:    payload.Name,
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  req.Categories,
	})
	if cErr != nil {
		return nil, cErr
	}
	return &pb.ProductReply{
		Message: result.Message,
		Code:   result.Code,
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.ProductReply, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("✅ req: %+v\n", req)
    fmt.Printf("✅ owner: %+v\n username: %+v\n", req.Owner, req.Username)

	if req.Owner != payload.Owner || req.Username != payload.Name {
		return nil, errors.New("invalid token")
	}
	
	result, err := s.pc.UpdateProduct(ctx, biz.Product{ // Pass the address of the struct
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  req.Categories,
		Id:          req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ProductReply{
		Message: result.Message,
		Code:   result.Code,
	}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductReq) (*pb.ProductReply, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("✅ req: %+v\n", req)
    fmt.Printf("✅ owner: %+v\n username: %+v\n", req.Owner, req.Username)

	if req.Owner != payload.Owner || req.Username != payload.Name {
		return nil, errors.New("invalid token")
	}

	result, err := s.pc.DeleteProduct(ctx, biz.DeleteProductReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ProductReply{
		Message: result.Message,
		Code:   result.Code,
	}, nil
}

// func (s *ProductCatalogServiceService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
// 	return &pb.DeleteProductReply{}, nil
// }
