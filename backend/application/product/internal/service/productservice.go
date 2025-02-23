package service

import (
	"context"
	"fmt"

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

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
	bizProduct := convertPBToBizProduct(req.GetProduct())
	created, err := s.uc.CreateProduct(ctx, &biz.CreateProductRequest{Product: *bizProduct})
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
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
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
		OperatorId: record.OperatorID,
		OperatedAt: timestamppb.New(record.OperatedAt),
	}, nil
}

func (s *ProductService) AuditProduct(ctx context.Context, req *pb.AuditProductRequest) (*pb.AuditRecord, error) {
	if req.Action == pb.AuditAction_AUDIT_ACTION_REJECT && req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "reject reason required")
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	operatorId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
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
		OperatorId: record.OperatorID,
		OperatedAt: timestamppb.New(record.OperatedAt),
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	// TODO service 层获取用户ID
	// 层层传递下去
	var userId string
	if md, ok := metadata.FromServerContext(ctx); ok {
		userId = md.Get("x-md-global-user-id")
	}
	fmt.Printf("x-md-global-user-id %s\n", userId)

	id, err := uuid.Parse(req.Id)
	product, err := s.uc.GetProduct(ctx, &biz.GetProductRequest{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	return convertBizProductToPB(product), nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	bizReq := biz.DeleteProductRequest{
		ID:         id,
		MerchantID: merchantId,
	}

	_, err = s.uc.DeleteProduct(ctx, bizReq)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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
		pbProduct.Images = append(pbProduct.Images, &pb.Product_Image{
			Url:       img.URL,
			IsPrimary: img.IsPrimary,
			SortOrder: sortOrder,
		})
	}

	if p.AuditInfo.AuditId != uuid.Nil {
		pbProduct.AuditInfo = &pb.AuditInfo{
			AuditId:    p.AuditInfo.AuditId.String(),
			Reason:     p.AuditInfo.Reason,
			OperatorId: p.AuditInfo.OperatorId,
			OperatedAt: timestamppb.New(p.AuditInfo.OperatedAt),
		}
	}

	return pbProduct
}

func convertPBToBizProduct(p *pb.Product) *biz.Product {
	id, err := uuid.Parse(p.Id)
	if err != nil {
		return nil
	}
	merchantId, err := uuid.Parse(p.MerchantId)
	if err != nil {
		return nil
	}

	return &biz.Product{
		ID:          id,
		MerchantId:  merchantId,
		Name:        p.GetName(),
		Price:       p.GetPrice(),
		Description: p.GetDescription(),
		Images:      convertPBImagesToBiz(p.GetImages()),
		Status:      convertPBStatusToBiz(p.GetStatus()),
		Category: biz.CategoryInfo{
			CategoryId:   uint64(p.GetCategory().GetCategoryId()),
			CategoryName: p.GetCategory().GetCategoryName(),
		},
		Attributes: nil,
		AuditInfo:  convertPBAuditInfoToBiz(p.GetAuditInfo()),
	}
}

func convertPBImagesToBiz(pbImages []*pb.Product_Image) []*biz.ProductImage {
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

func convertPBAuditInfoToBiz(pbInfo *pb.AuditInfo) biz.AuditInfo {
	if pbInfo == nil {
		return biz.AuditInfo{}
	}
	auditId, err := uuid.Parse(pbInfo.AuditId)
	if err != nil {
		return biz.AuditInfo{}
	}
	return biz.AuditInfo{
		AuditId:    auditId,
		Reason:     pbInfo.Reason,
		OperatorId: pbInfo.OperatorId,
		OperatedAt: pbInfo.OperatedAt.AsTime(),
	}
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
