package service

import (
	pb "backend/api/product/v1"
	"backend/application/product/internal/biz"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductCatalogServiceService struct {
	pb.UnimplementedProductCatalogServiceServer
	pu *biz.ProductUsecase
}

func NewProductService(pu *biz.ProductUsecase) *ProductCatalogServiceService {
	return &ProductCatalogServiceService{pu: pu}
}

func (s *ProductCatalogServiceService) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsResp, error) {
	list, err := s.pu.ListProducts(ctx, &biz.ListProductsReq{
		Page:         uint(req.Page),
		PageSize:     uint(req.PageSize),
		CategoryId: req.CategoryId,
	})
	if err != nil {
		return nil, err
	}
	pbProduct := make([]*pb.Product, len(list.Product))
	for i, product := range list.Product {
		pbProduct[i] = &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
		}
	}
	return &pb.ListProductsResp{
		Products: pbProduct,
	}, nil
}

func (s *ProductCatalogServiceService) GetProduct(ctx context.Context, req *pb.GetProductReq) (*pb.ProductReply, error) {
	product, err := s.pu.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ProductReply{Product: &pb.Product{
		Id:          product.Product.Id,
		Name:        product.Product.Name,
		Description: product.Product.Description,
		Picture:     product.Product.Picture,
		Price:       product.Product.Price,
		CategoryId:  product.Product.CategoryId,
		TotalStock:  product.Product.TotalStock,
		AvailableStock: *product.Product.AvailableStock,
		ReservedStock: product.Product.ReservedStock,
		LowStockThreshold: product.Product.LowStockThreshold,
		AllowNegative: product.Product.AllowNegative,
		CreatedAt: &timestamppb.Timestamp{Seconds: product.Product.CreatedAt.Unix()},
		UpdatedAt: timestamppb.New(product.Product.UpdatedAt),
		Version: product.Product.Version,
	}}, nil
}

func (s *ProductCatalogServiceService) SearchProducts(ctx context.Context, req *pb.SearchProductsReq) (*pb.SearchProductsResp, error) {
	products, err := s.pu.SearchProducts(ctx, &biz.SearchProductsReq{Query: req.GetQuery()})
	if err != nil {
		return nil, err
	}
	pbProduct := make([]*pb.Product, len(products.Result))
	for i, product := range products.Result {
		pbProduct[i] = &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			CategoryId:  product.CategoryId,
			TotalStock:  product.TotalStock,
			AvailableStock: *product.AvailableStock,
			ReservedStock: product.ReservedStock,
			LowStockThreshold: product.LowStockThreshold,
			AllowNegative: product.AllowNegative,
			CreatedAt: &timestamppb.Timestamp{Seconds: product.CreatedAt.Unix()},
			UpdatedAt: timestamppb.New(product.UpdatedAt),
			Version: product.Version,
		}
	}
	return &pb.SearchProductsResp{
		Results: pbProduct,
	}, nil
}

func (s *ProductCatalogServiceService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductReply, error) {

	p, cErr := s.pu.UpdateProduct(ctx, &biz.UpdateProductRequest{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
		TotalStock:  req.TotalStock,
	})
	if cErr != nil {
		return nil, cErr
	}
	return &pb.ProductReply{
		Product: &pb.Product{
			Id:          p.Product.Id,
			Name:        p.Product.Name,
			Description: p.Product.Description,
			Picture:     p.Product.Picture,
			Price:       p.Product.Price,
			CategoryId:  p.Product.CategoryId,
			TotalStock:  p.Product.TotalStock,
		},
	}, nil	
}

func (s *ProductCatalogServiceService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductReply, error) {

	p, cErr := s.pu.CreateProduct(ctx, &biz.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
		TotalStock:  req.TotalStock,
	})
	if cErr != nil {
		return nil, cErr
	}
	return &pb.ProductReply{
		Product: &pb.Product{
			Id:          p.Product.Id,
			Name:        p.Product.Name,
			Description: p.Product.Description,
			Picture:     p.Product.Picture,
			Price:       p.Product.Price,
			CategoryId:  p.Product.CategoryId,
			TotalStock:  p.Product.TotalStock,
		},
	}, nil
}



func (s *ProductCatalogServiceService) DeleteProduct(ctx context.Context, req *pb.DeleteProductReq) (*pb.ProductReply, error) {
	p, cErr := s.pu.DeleteProduct(ctx, &biz.DeleteProductReq{
		Id:          req.Id,
	})
	if cErr != nil {
		return nil, cErr
	}
	return &pb.ProductReply{
		Product: &pb.Product{
			Id:          p.Product.Id,
			Name:        p.Product.Name,
			Description: p.Product.Description,
			Picture:     p.Product.Picture,
			Price:       p.Product.Price,
			CategoryId:  p.Product.CategoryId,
			TotalStock:  p.Product.TotalStock,
		},
	}, nil
}

