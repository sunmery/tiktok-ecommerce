package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"google.golang.org/protobuf/types/known/structpb"

	"backend/pkg"

	"github.com/google/uuid"

	pb "backend/api/product/v1"
	"backend/application/product/internal/biz"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	uc *biz.ProductUsecase
}

func NewProductService(uc *biz.ProductUsecase) *ProductService {
	return &ProductService{uc: uc}
}

func (s *ProductService) UploadProductFile(ctx context.Context, req *pb.UploadProductFileRequest) (*pb.UploadProductFileReply, error) {
	result, err := s.uc.UploadProductFile(ctx, &biz.UploadProductFileRequest{
		Method:      biz.UploadMethod(req.Method),
		ContentType: req.ContentType,
		BucketName:  req.BucketName,
		FilePath:    req.FilePath,
		FileName:    req.FileName,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "生成预签名URL失败")
	}

	return &pb.UploadProductFileReply{
		UploadUrl:   result.UploadUrl,
		DownloadUrl: result.DownloadUrl,
		BucketName:  result.BucketName,
		ObjectName:  result.ObjectName,
		FormData:    result.FormData,
	}, nil
}

func (s *ProductService) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.UpdateInventoryReply, error) {
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}

	result, err := s.uc.UpdateInventory(ctx, &biz.UpdateInventoryRequest{
		ProductId:  productId,
		MerchantId: merchantId,
		Stock:      req.Stock,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateInventoryReply{
		ProductId:  result.ProductId.String(),
		MerchantId: result.MerchantId.String(),
		Stock:      result.Stock,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}

	created, createdErr := s.uc.CreateProduct(ctx, &biz.CreateProductRequest{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		MerchantId:  merchantId,
		Images:      convertPBImagesToBiz(req.Images),
		Status:      biz.ProductStatusApproved,
		Category: biz.CategoryInfo{
			CategoryId:   uint64(req.Category.CategoryId),
			CategoryName: req.Category.CategoryName,
		},
		Attributes: parseProtoValue(req.Attributes),
		Stock:      req.Stock,
	})
	if createdErr != nil {
		return nil, createdErr
	}
	return &pb.CreateProductReply{
		Id:        created.ID.String(),
		CreatedAt: timestamppb.New(created.CreatedAt),
		UpdatedAt: timestamppb.New(created.UpdatedAt),
	}, nil
}

// CreateProductBatch 批量创建商品（草稿状态）
func (s *ProductService) CreateProductBatch(ctx context.Context, req *pb.CreateProductBatchRequest) (*pb.CreateProductBatchReply, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}

	logger := log.NewHelper(log.With(log.GetLogger(), "module", "service/product"))
	logger.Infof("开始批量创建商品，商家ID: %v, 商品数量: %d", merchantId, len(req.Products))

	// 参数验证
	if len(req.Products) == 0 {
		return nil, status.Error(codes.InvalidArgument, "商品列表不能为空")
	}

	// 转换请求格式从proto到biz层
	bizProducts := make([]*biz.ProductDraft, 0, len(req.Products))
	for _, product := range req.Products {
		// 基本参数验证
		if product.Name == "" {
			return nil, status.Error(codes.InvalidArgument, "商品名称不能为空")
		}
		if product.Price <= 0 {
			return nil, status.Error(codes.InvalidArgument, "商品价格必须大于0")
		}

		bizProducts = append(bizProducts, &biz.ProductDraft{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			MerchantId:  merchantId,
			Status:      biz.ProductStatusPending,
			Category: biz.CategoryInfo{
				CategoryId:   uint64(product.Category.CategoryId),
				CategoryName: product.Category.CategoryName,
				// SortOrder:    product.Category.SortOrder,
			},
			Attributes: parseProtoValue(product.Attributes),
			Images:     convertPBImagesToBiz(product.Images),
		})
	}

	logger.Infof("调用业务层批量创建商品，商品数量: %d", len(bizProducts))
	created, createdErr := s.uc.CreateProductBatch(ctx, &biz.CreateProductBatchRequest{
		Products: bizProducts,
	})
	if createdErr != nil {
		logger.Errorf("批量创建商品失败: %v", createdErr)
		return nil, status.Errorf(codes.Internal, "批量创建商品失败: %v", createdErr)
	}

	// 处理错误结果
	errResults := make([]*pb.CreateProductBatchReply_BatchProductError, 0, len(created.Errors))
	for _, result := range created.Errors {
		errResults = append(errResults, &pb.CreateProductBatchReply_BatchProductError{
			Index:   uint32(result.Index),
			Message: result.Message,
		})
	}

	logger.Infof("批量创建商品完成，成功: %d, 失败: %d", created.SuccessCount, created.FailedCount)
	return &pb.CreateProductBatchReply{
		SuccessCount: created.SuccessCount,
		FailedCount:  created.FailedCount,
		Errors:       errResults,
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	productId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchant id")
	}

	product, err := s.uc.GetProduct(ctx, &biz.GetProductRequest{
		ID:         productId,
		MerchantID: merchantId,
	})
	if err != nil {
		return nil, err
	}

	return convertBizProductToPB(product), nil
}

func (s *ProductService) GetProductsBatch(ctx context.Context, req *pb.GetProductsBatchRequest) (*pb.Products, error) {
	var productIds []uuid.UUID
	for _, id := range req.ProductIds {
		productId, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid product id")
		}
		productIds = append(productIds, productId)
	}

	var merchantIds []uuid.UUID
	for _, id := range req.MerchantIds {
		merchantId, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid merchant id")
		}
		merchantIds = append(merchantIds, merchantId)
	}

	products, err := s.uc.GetProductBatch(ctx, &biz.GetProductsBatchRequest{
		ProductIds:  productIds,
		MerchantIds: merchantIds,
	})
	if err != nil {
		return nil, err
	}
	var pbProducts []*pb.Product
	for _, product := range products.Items {
		pbProducts = append(pbProducts, convertBizProductToPB(product))
	}
	return &pb.Products{
		Items: pbProducts,
	}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	// 从网关获取用户ID, 这里的用户是商户, 只有商户角色才能删除商品
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}

	log.Debugf("DeleteProduct userid:%+v", userId)
	bizReq := biz.DeleteProductRequest{
		ID:         id,
		MerchantID: userId,
		Status:     4,
	}

	_, err = s.uc.DeleteProduct(ctx, bizReq)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, err
}

// ListRandomProducts 随机返回商品数据
func (s *ProductService) ListRandomProducts(ctx context.Context, req *pb.ListRandomProductsRequest) (*pb.Products, error) {
	listRandomProducts, err := s.uc.ListRandomProducts(ctx, &biz.ListRandomProductsRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	if listRandomProducts == nil {
		return &pb.Products{
			Items: []*pb.Product{},
		}, nil
	}
	if listRandomProducts.Items == nil {
		return &pb.Products{
			Items: []*pb.Product{},
		}, nil
	}
	var pbProducts []*pb.Product
	for _, product := range listRandomProducts.Items {
		if product != nil {
			pbProducts = append(pbProducts, convertBizProductToPB(product))
		}
	}
	return &pb.Products{
		Items: pbProducts,
	}, nil
}

func (s *ProductService) GetCategoryProducts(ctx context.Context, req *pb.GetCategoryProductsRequest) (*pb.Products, error) {
	// 设置默认分页参数
	page := uint32(1)
	pageSize := uint32(10)

	// 使用请求中的参数，如果有提供的话
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	listRandomProducts, err := s.uc.GetCategoryProducts(ctx, &biz.GetCategoryProducts{
		CategoryID: req.CategoryId,
		Status:     req.Status,
		Page:       int64(page),
		PageSize:   int64(pageSize),
	})
	if err != nil {
		return nil, err
	}
	var pbProducts []*pb.Product
	for _, product := range listRandomProducts.Items {
		pbProducts = append(pbProducts, convertBizProductToPB(product))
	}
	return &pb.Products{
		Items: pbProducts,
	}, nil
}

func (s *ProductService) GetCategoryWithChildrenProducts(ctx context.Context, req *pb.GetCategoryProductsRequest) (*pb.Products, error) {
	// 设置默认分页参数
	page := uint32(1)
	pageSize := uint32(10)

	// 使用请求中的参数，如果有提供的话
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	products, err := s.uc.GetCategoryWithChildrenProducts(ctx, &biz.GetCategoryWithChildrenProducts{
		CategoryID: req.CategoryId,
		Status:     req.Status,
		Page:       int64(page),
		PageSize:   int64(pageSize),
	})
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, product := range products.Items {
		pbProducts = append(pbProducts, convertBizProductToPB(product))
	}

	return &pb.Products{
		Items: pbProducts,
	}, nil
}

// SearchProductsByName 根据商品名称模糊查询
func (s *ProductService) SearchProductsByName(ctx context.Context, req *pb.SearchProductRequest) (*pb.Products, error) {
	products, err := s.uc.SearchProductsByName(context.Background(), &biz.SearchProductsByNameRequest{
		Name:     req.Name,
		Page:     req.Page,
		PageSize: req.PageSize,
		Query:    req.Query,
	})
	if err != nil {
		return nil, err
	}
	var pbProducts []*pb.Product
	for _, product := range products.Items {
		pbProducts = append(pbProducts, convertBizProductToPB(product))
	}
	return &pb.Products{
		Items: pbProducts,
	}, nil
}

// 辅助转换方法
func convertBizProductToPB(p *biz.Product) *pb.Product {
	if p == nil {
		return nil
	}

	// 转换图片
	images := make([]*pb.Image, 0)
	if p.Images != nil {
		for _, img := range p.Images {
			if img != nil {
				sortOrder := int32(0)
				if img.SortOrder != nil {
					sortOrder = int32(*img.SortOrder)
				}
				images = append(images, &pb.Image{
					Url:       img.URL,
					IsPrimary: img.IsPrimary,
					SortOrder: sortOrder,
				})
			}
		}
	}

	// 转换商品属性
	var attributes *structpb.Value
	if p.Attributes != nil && len(p.Attributes) > 0 {
		protoStruct, err := structpb.NewStruct(p.Attributes)
		if err != nil {
			log.Warn("Error creating struct: %w", err)
			attributes = nil
		} else {
			attributes = structpb.NewStructValue(protoStruct)
		}
	}

	// 构建返回结果
	result := &pb.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Status:      uint32(p.Status),
		MerchantId:  p.MerchantId.String(),
		Images:      images,
		Attributes:  attributes,
		Category: &pb.CategoryInfo{
			CategoryId:   uint32(p.Category.CategoryId),
			CategoryName: p.Category.CategoryName,
		},
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}

	// 添加库存信息
	if p.Inventory.ProductId != uuid.Nil {
		result.Inventory = &pb.Inventory{
			ProductId:  p.Inventory.ProductId.String(),
			MerchantId: p.Inventory.MerchantId.String(),
			Stock:      p.Inventory.Stock,
		}
	}

	return result
}

func convertPBImagesToBiz(pbImages []*pb.Image) []*biz.ProductImage {
	var images []*biz.ProductImage
	for _, img := range pbImages {
		var sortOrderPtr *int
		if img.SortOrder != 0 { // 0 视为未设置
			sortOrderValue := int(img.SortOrder)
			sortOrderPtr = &sortOrderValue
		}
		images = append(images, &biz.ProductImage{
			URL:       img.Url,
			IsPrimary: img.IsPrimary,
			SortOrder: sortOrderPtr,
		})
	}
	return images
}

var validTransitions = map[biz.ProductStatus]map[biz.ProductStatus]bool{
	biz.ProductStatusDraft: {
		biz.ProductStatusPending:  true,
		biz.ProductStatusRejected: true,
	},
	biz.ProductStatusPending: {
		biz.ProductStatusApproved: true,
		biz.ProductStatusRejected: true,
	},
	biz.ProductStatusRejected: {
		biz.ProductStatusDraft: true,
	},
	biz.ProductStatusApproved: {
		// 已审核状态不允许修改
	},
}

func parseProtoValue(v *structpb.Value) map[string]any {
	if v == nil {
		return nil
	}
	if v.GetStructValue() == nil {
		return nil
	}
	return v.GetStructValue().AsMap()
}
