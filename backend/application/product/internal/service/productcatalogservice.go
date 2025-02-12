package service

import (
	pb "backend/api/product/v1"
	"backend/application/product/internal/biz"
	"backend/pkg/token"
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
		Product: &pb.Product{
			Id:          p.Product.Id,
			Name:        p.Product.Name,
			Description: p.Product.Description,
			Picture:     p.Product.Picture,
			Price:       p.Product.Price,
			Categories:  p.Product.Categories,
		},
	}, nil
}

// func (s *ProductCatalogServiceService) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.ProductReply, error) {
// 	payload, err := token.ExtractPayload(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Printf("✅ req: %+v\n", req)
//     fmt.Printf("✅ owner: %+v\n username: %+v\n", req.Owner, req.Username)

// 	if req.Owner != payload.Owner || req.Username != payload.Name {
// 		return nil, errors.New("invalid token")
// 	}
	
// 	result, err := s.pc.UpdateProduct(ctx, biz.Product{ // Pass the address of the struct
// 		Name:        req.Name,
// 		Description: req.Description,
// 		Picture:     req.Picture,
// 		Price:       req.Price,
// 		Categories:  req.Categories,
// 		Id:          req.Id,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.ProductReply{
// 		Message: result.Message,
// 		Code:   result.Code,
// 	}, nil
// }

// func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductReq) (*pb.ProductReply, error) {
// 	payload, err := token.ExtractPayload(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Printf("✅ req: %+v\n", req)
//     fmt.Printf("✅ owner: %+v\n username: %+v\n", req.Owner, req.Username)

// 	if req.Owner != payload.Owner || req.Username != payload.Name {
// 		return nil, errors.New("invalid token")
// 	}

// 	result, err := s.pc.DeleteProduct(ctx, biz.DeleteProductReq{
// 		Id: req.Id,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.ProductReply{
// 		Message: result.Message,
// 		Code:   result.Code,
// 	}, nil
// }

// func (s *ProductCatalogServiceService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
// 	return &pb.DeleteProductReply{}, nil
// }
