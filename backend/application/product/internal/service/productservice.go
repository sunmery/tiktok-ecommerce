package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	"google.golang.org/protobuf/types/known/structpb"

	"backend/pkg"

	"github.com/google/uuid"

	pb "backend/api/product/v1"
	"backend/application/product/internal/biz"

	"github.com/go-kratos/kratos/v2/metadata"
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
		// Images:      convertPBImagesToBiz(req.Images),
		Status: biz.ProductStatusPending,
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

func (s *ProductService) SubmitForAudit(ctx context.Context, req *pb.SubmitAuditRequest) (*pb.AuditRecord, error) {
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	var userId string
	if md, ok := metadata.FromServerContext(ctx); ok {
		userId = md.Get("x-md-global-user-id")
	}
	merchantId, err := uuid.Parse(userId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}

	bizReq := biz.SubmitAuditRequest{
		ProductID:  productId,
		MerchantID: merchantId,
	}

	record, err := s.uc.SubmitForAudit(ctx, &bizReq)
	if err != nil {
		return nil, err
	}
	return &pb.AuditRecord{
		Id:         record.ID.String(),
		ProductId:  productId.String(),
		OldStatus:  uint32(record.OldStatus),
		NewStatus:  uint32(record.NewStatus),
		Reason:     record.Reason,
		OperatorId: record.OperatorID.String(),
		OperatedAt: timestamppb.New(record.OperatedAt),
	}, nil
}

func (s *ProductService) AuditProduct(ctx context.Context, req *pb.AuditProductRequest) (*pb.AuditRecord, error) {
	if req.Action == pb.AuditAction_AUDIT_ACTION_REJECT && req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "reject reason required")
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid productId ID")
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}

	operatorId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid operatorId ID")
	}

	bizReq := biz.AuditProductRequest{
		ProductID:  productId,
		MerchantID: merchantId,
		Action:     uint64(biz.AuditAction(req.Action)),
		Reason:     req.Reason,
		OperatorID: operatorId,
	}

	record, err := s.uc.AuditProduct(ctx, &bizReq)
	if err != nil {
		return nil, err
	}

	return &pb.AuditRecord{
		Id:         record.ID.String(),
		ProductId:  record.ProductID.String(),
		OldStatus:  uint32(record.OldStatus),
		NewStatus:  uint32(record.NewStatus),
		Reason:     record.Reason,
		OperatorId: record.OperatorID.String(),
		OperatedAt: timestamppb.New(record.OperatedAt),
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
		fmt.Println(id)
		productId, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid product id")
		}
		productIds = append(productIds, productId)
	}

	var merchantIds []uuid.UUID
	for _, id := range req.MerchantIds {
		fmt.Println(id)
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

// 其他辅助转换函数
// func convertPBImagesToBiz(pbImages []*pb.Product_Image) []*biz.ProductImage {
// 	var images []*biz.ProductImage
// 	for _, img := range pbImages {
// 		images = append(images, &biz.ProductImage{
// 			URL:       img.GetUrl(),
// 			IsPrimary: img.GetIsPrimary(),
// 			SortOrder: int(img.GetSortOrder()),
// 		})
// 	}
// 	return images
// }

// func convertPBStatusToBiz(pbStatus pb.ProductStatus) biz.ProductStatus {
// 	return biz.ProductStatus(pbStatus)
// }

func convertPBCategoryToBiz(pbCategory *pb.CategoryInfo) biz.CategoryInfo {
	if pbCategory == nil {
		return biz.CategoryInfo{}
	}
	return biz.CategoryInfo{
		CategoryId:   uint64(pbCategory.GetCategoryId()),
		CategoryName: pbCategory.GetCategoryName(),
	}
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

func convertPBAuditInfoToBiz(pbInfo *pb.AuditInfo) (*biz.AuditInfo, error) {
	if pbInfo == nil {
		return nil, errors.New("audit info is nil")
	}
	auditId, err := uuid.Parse(pbInfo.AuditId)
	if err != nil {
		return nil, errors.New("invalid audit ID")
	}

	operatorId, err := uuid.Parse(pbInfo.OperatorId)
	if err != nil {
		return nil, errors.New("invalid operator ID")
	}
	return &biz.AuditInfo{
		AuditId:    auditId,
		Reason:     pbInfo.Reason,
		OperatorId: operatorId,
		OperatedAt: pbInfo.OperatedAt.AsTime(),
	}, nil
}

func convertBizStatusToPB(s biz.ProductStatus) pb.ProductStatus {
	switch s {
	case biz.ProductStatusDraft:
		return pb.ProductStatus_PRODUCT_STATUS_DRAFT
	case biz.ProductStatusPending:
		return pb.ProductStatus_PRODUCT_STATUS_PENDING
	case biz.ProductStatusApproved:
		return pb.ProductStatus_PRODUCT_STATUS_APPROVED
	case biz.ProductStatusRejected:
		return pb.ProductStatus_PRODUCT_STATUS_REJECTED
	default:
		return pb.ProductStatus_PRODUCT_STATUS_DRAFT
	}
}

func convertPBStatusToBiz(s pb.ProductStatus) biz.ProductStatus {
	switch s {
	case pb.ProductStatus_PRODUCT_STATUS_DRAFT:
		return biz.ProductStatusDraft
	case pb.ProductStatus_PRODUCT_STATUS_PENDING:
		return biz.ProductStatusPending
	case pb.ProductStatus_PRODUCT_STATUS_APPROVED:
		return biz.ProductStatusApproved
	case pb.ProductStatus_PRODUCT_STATUS_REJECTED:
		return biz.ProductStatusRejected
	default:
		return biz.ProductStatusDraft
	}
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
