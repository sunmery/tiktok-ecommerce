package service

import (
	pb "backend/api/product/v1"
	"backend/application/product/internal/biz"
	"context"
	"errors"
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

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	bizProduct := convertPBToBizProduct(req.GetProduct())
	created, err := s.uc.CreateProduct(ctx, biz.CreateProductRequest{Product: *bizProduct})
	if err != nil {
		return nil, convertError(err)
	}
	return convertBizProductToPB(created), nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	updateReq := biz.UpdateProductRequest{
		ID:         req.Id,
		MerchantID: req.Product.MerchantId,
	}

	// 使用指针实现字段掩码
	pbProduct := req.GetProduct()
	if pbProduct.Name != "" {
		updateReq.Name = &pbProduct.Name
	}
	if pbProduct.Price > 0 {
		updateReq.Price = &pbProduct.Price
	}
	// if pbProduct.Stock >= 0 {
	// 	stock := int(pbProduct.Stock)
	// 	updateReq.Stock = &stock
	// }
	if pbProduct.Description != "" {
		updateReq.Description = pbProduct.Description
	}
	if pbProduct.Category != nil {
		updateReq.Category = biz.CategoryInfo{
			CategoryId:   pbProduct.Category.CategoryId,
			CategoryName: pbProduct.Category.CategoryName,
		}
	}

	updatedProduct, err := s.uc.UpdateProduct(ctx, updateReq)
	if err != nil {
		return nil, convertError(err)
	}

	return convertBizProductToPB(updatedProduct), nil
}

func (s *ProductService) SubmitForAudit(ctx context.Context, req *pb.SubmitAuditRequest) (*pb.AuditRecord, error) {
	bizReq := biz.SubmitAuditRequest{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
	}

	record, err := s.uc.SubmitForAudit(ctx, bizReq)
	if err != nil {
		return nil, convertError(err)
	}

	return &pb.AuditRecord{
		Id:         record.ID,
		ProductId:  record.ProductID,
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

	bizReq := biz.AuditProductRequest{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
		Action:     uint64(biz.AuditAction(req.Action)),
		Reason:     req.Reason,
		OperatorID: req.OperatorId,
	}

	record, err := s.uc.AuditProduct(ctx, bizReq)
	if err != nil {
		return nil, convertError(err)
	}

	return &pb.AuditRecord{
		Id:         record.ID,
		ProductId:  record.ProductID,
		OldStatus:  convertBizStatusToPB(record.OldStatus),
		NewStatus:  convertBizStatusToPB(record.NewStatus),
		Reason:     record.Reason,
		OperatorId: record.OperatorID,
		OperatedAt: timestamppb.New(record.OperatedAt),
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	bizReq := biz.GetProductRequest{
		ID:         req.Id,
		MerchantID: req.MerchantId,
	}

	product, err := s.uc.GetProduct(ctx, bizReq)
	if err != nil {
		return nil, convertError(err)
	}

	return convertBizProductToPB(product), nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	bizReq := biz.DeleteProductRequest{
		ID:         req.Id,
		MerchantID: req.MerchantId,
	}

	_, err := s.uc.DeleteProduct(ctx, bizReq)
	if err != nil {
		return nil, convertError(err)
	}

	return &emptypb.Empty{}, nil
}

// 辅助转换方法
func convertBizProductToPB(p *biz.Product) *pb.Product {
	pbProduct := &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		// Stock:       p.Stock,
		Status:      convertBizStatusToPB(p.Status),
		MerchantId:  p.MerchantId,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
		Category: &pb.CategoryInfo{
			CategoryId:   p.Category.CategoryId,
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

	if p.AuditInfo.AuditId > 0 {
		pbProduct.AuditInfo = &pb.AuditInfo{
			AuditId:    p.AuditInfo.AuditId,
			Reason:     p.AuditInfo.Reason,
			OperatorId: p.AuditInfo.OperatorId,
			OperatedAt: timestamppb.New(p.AuditInfo.OperatedAt),
		}
	}

	return pbProduct
}

func convertPBToBizProduct(p *pb.Product) *biz.Product {
	return &biz.Product{
		ID:          p.GetId(),
		Name:        p.GetName(),
		Description: p.GetDescription(),
		Price:       p.GetPrice(),
		// Stock:       p.GetStock(),
		Status:      convertPBStatusToBiz(p.GetStatus()),
		MerchantId:  p.GetMerchantId(),
		Category: biz.CategoryInfo{
			CategoryId:   p.GetCategory().GetCategoryId(),
			CategoryName: p.GetCategory().GetCategoryName(),
		},
		Images:    convertPBImagesToBiz(p.GetImages()),
		AuditInfo: convertPBAuditInfoToBiz(p.GetAuditInfo()),
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
	return biz.AuditInfo{
		AuditId:    pbInfo.AuditId,
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

// 错误转换
func convertError(err error) error {
	switch {
	case errors.Is(err, biz.ErrProductNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, biz.ErrInvalidStatus):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, biz.ErrAuditReasonMissing):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
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
