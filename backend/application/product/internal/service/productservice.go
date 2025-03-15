package service

import (
	"context"
	"errors"
	"fmt"

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

	created, err := s.uc.CreateProduct(ctx, &biz.CreateProductRequest{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		MerchantId:  merchantId,
		Images:      convertPBImagesToBiz(req.Images),
		Status:      biz.ProductStatusPending,
		Category: biz.CategoryInfo{
			CategoryId:   uint64(req.Category.CategoryId),
			CategoryName: req.Category.CategoryName,
		},
		Attributes: convertPBAttributes(req.Attributes),
		Stock:      req.Stock,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductReply{
		Id:        created.ID.String(),
		CreatedAt: timestamppb.New(created.CreatedAt),
		UpdatedAt: timestamppb.New(created.UpdatedAt),
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	merchantId, err := uuid.Parse(req.Product.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	updateReq := biz.UpdateProductRequest{
		ID:         id,
		MerchantID: merchantId,
	}

	// 使用指针实现字段掩码
	pbProduct := req.GetProduct()
	if pbProduct.Name != "" {
		updateReq.Name = &pbProduct.Name
	}
	if pbProduct.Price > 0 {
		updateReq.Price = &pbProduct.Price
	}

	if pbProduct.Description != "" {
		updateReq.Description = pbProduct.Description
	}
	if pbProduct.Category != nil {
		updateReq.Category = biz.CategoryInfo{
			CategoryId:   uint64(pbProduct.Category.CategoryId),
			CategoryName: pbProduct.Category.CategoryName,
		}
	}

	updatedProduct, err := s.uc.UpdateProduct(ctx, &updateReq)
	if err != nil {
		return nil, err
	}

	return convertBizProductToPB(updatedProduct), nil
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
		OldStatus:  convertBizStatusToPB(record.OldStatus),
		NewStatus:  convertBizStatusToPB(record.NewStatus),
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

	var userId string
	if md, ok := metadata.FromServerContext(ctx); ok {
		userId = md.Get("x-md-global-user-id")
	}
	operatorId, err := uuid.Parse(userId)
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
		OldStatus:  convertBizStatusToPB(record.OldStatus),
		NewStatus:  convertBizStatusToPB(record.NewStatus),
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

func (s *ProductService) GetMerchantProducts(ctx context.Context, _ *pb.GetMerchantProductRequest) (*pb.Products, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}
	products, err := s.uc.GetMerchantProducts(ctx, &biz.GetMerchantProducts{
		MerchantID: merchantId,
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
	var userId string
	if md, ok := metadata.FromServerContext(ctx); ok {
		userId = md.Get("x-md-global-user-id")
	}
	merchantId := uuid.MustParse(userId)
	bizReq := biz.DeleteProductRequest{
		ID:         id,
		MerchantID: merchantId,
		Status:     4,
	}

	_, err = s.uc.DeleteProduct(ctx, bizReq)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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
	var pbProducts []*pb.Product
	for _, product := range listRandomProducts.Items {
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
	pbProduct := &pb.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Status:      convertBizStatusToPB(p.Status),
		MerchantId:  p.MerchantId.String(),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
		Category: &pb.CategoryInfo{
			CategoryId:   uint32(p.Category.CategoryId),
			CategoryName: p.Category.CategoryName,
		},
	}

	for _, img := range p.Images {
		// 安全转换 SortOrder
		var sortOrder int32
		if img.SortOrder != nil {
			sortOrder = int32(*img.SortOrder) // 解引用并转换类型
		} else {
			sortOrder = 0 // 默认值
		}
		pbProduct.Images = append(pbProduct.Images, &pb.Image{
			Url:       img.URL,
			IsPrimary: img.IsPrimary,
			SortOrder: sortOrder,
		})
	}

	return pbProduct
}

// func convertPBToBizProduct(p *pb.Product) (*biz.Product, error) {
// 	merchantId, err := uuid.Parse(p.MerchantId)
// 	if err != nil {
// 		return nil, errors.New("invalid merchant ID")
// 	}
//
// 	return &biz.Product{
// 		MerchantId:  merchantId,
// 		Name:        p.GetName(),
// 		Price:       p.GetPrice(),
// 		Description: p.GetDescription(),
// 		Images:      convertPBImagesToBiz(p.GetImages()),
// 		Status:      convertPBStatusToBiz(p.GetStatus()),
// 		Category: biz.CategoryInfo{
// 			CategoryId:   uint64(p.GetCategory().GetCategoryId()),
// 			CategoryName: p.GetCategory().GetCategoryName(),
// 		},
// 		Attributes: convertPBObject(p.Attributes)
// 	},nil
// }

// 顶层转换函数
func convertPBToBizProduct(ctx context.Context, p *pb.Product) (*biz.Product, error) {
	var userId string
	if md, ok := metadata.FromServerContext(ctx); ok {
		userId = md.Get("x-md-global-user-id")
	}
	merchantID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid merchant ID: %w", err)
	}

	return &biz.Product{
		MerchantId:  merchantID,
		Name:        p.GetName(),
		Price:       p.GetPrice(),
		Description: p.GetDescription(),
		Images:      convertPBImagesToBiz(p.GetImages()),
		Status:      convertPBStatusToBiz(p.GetStatus()),
		Category:    convertPBCategoryToBiz(p.GetCategory()),
		Attributes:  convertPBAttributes(p.GetAttributes()),
	}, nil
}

// 递归处理属性转换
func convertPBAttributes(pbAttrs map[string]*pb.AttributeValue) map[string]*biz.AttributeValue {
	if pbAttrs == nil {
		return nil
	}

	bizAttrs := make(map[string]*biz.AttributeValue)
	for k, v := range pbAttrs {
		bizAttrs[k] = convertPBAttributeValue(v)
	}
	return bizAttrs
}

// 转换入口函数
func convertPBAttributeValue(pbVal *pb.AttributeValue) *biz.AttributeValue {
	if pbVal == nil {
		return nil
	}

	return &biz.AttributeValue{
		StringValue: pbVal.GetStringValue(),
		ArrayValue:  convertPBStringArray(pbVal.GetArrayValue()),
		ObjectValue: convertPBNestedObject(pbVal.GetObjectValue()),
	}
}

// 转换字符串数组
func convertPBStringArray(pbArr *pb.StringArray) *biz.ArrayValue {
	if pbArr == nil {
		return nil
	}

	// Protobuf 的 StringArray 应该包含 Items 字段
	return &biz.ArrayValue{
		Items: pbArr.GetItems(), // 直接映射到业务层的 ArrayValue.Items
	}
}

// 递归处理嵌套对象
func convertPBNestedObject(pbObj *pb.NestedObject) *biz.NestedObject {
	if pbObj == nil || pbObj.GetFields() == nil {
		return nil
	}

	fields := make(map[string]*biz.AttributeValue)
	for k, v := range pbObj.GetFields() {
		fields[k] = convertPBAttributeValue(v)
	}

	return &biz.NestedObject{
		Fields: fields,
	}
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
