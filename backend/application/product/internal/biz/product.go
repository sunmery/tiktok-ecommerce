package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	pb "backend/api/product/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductStatus uint

const (
	ProductStatusDraft ProductStatus = iota
	ProductStatusPending
	ProductStatusApproved
	ProductStatusRejected
)

// 补充状态映射
var pbStatusMapping = map[ProductStatus]pb.ProductStatus{
	ProductStatusDraft:    pb.ProductStatus_PRODUCT_STATUS_DRAFT,
	ProductStatusPending:  pb.ProductStatus_PRODUCT_STATUS_PENDING,
	ProductStatusApproved: pb.ProductStatus_PRODUCT_STATUS_APPROVED,
	ProductStatusRejected: pb.ProductStatus_PRODUCT_STATUS_REJECTED,
}

var validTransitions = map[ProductStatus]map[ProductStatus]bool{
	ProductStatusDraft: {
		ProductStatusPending: true,
	},
	ProductStatusPending: {
		ProductStatusApproved: true,
		ProductStatusRejected: true,
	},
	ProductStatusRejected: {
		ProductStatusDraft: true,
	},
}

// AuditAction 添加AuditAction类型
type (
	AuditAction    int
	AttributeValue struct{}
)

const (
	AuditActionApprove AuditAction = 0
	AuditActionReject  AuditAction = 1
)

// AuditRecord 完善AuditRecord定义
type AuditRecord struct {
	ID         uuid.UUID
	ProductID  uuid.UUID
	OldStatus  ProductStatus
	NewStatus  ProductStatus
	Reason     string
	OperatorID uint64
	OperatedAt time.Time
}
type AuditInfo struct {
	AuditId    uuid.UUID // 审核记录ID
	Reason     string    // 审核意见/驳回原因
	OperatorId uint64    // 操作人ID
	OperatedAt time.Time // 操作时间
}

type CategoryInfo struct {
	CategoryId   uint64
	CategoryName string
	SortOrder    int32
}
type ProductImage struct {
	URL       string
	IsPrimary bool
	SortOrder *int
}

// Product 商品领域模型
type Product struct {
	ID          uuid.UUID
	MerchantId  uuid.UUID
	Name        string
	Price       float64
	Description string
	Images      []*ProductImage
	Status      ProductStatus
	Category    CategoryInfo
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Attributes  map[string]*AttributeValue
	AuditInfo   AuditInfo
}

type CreateProductReply struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SubmitAuditRequest struct {
	ProductID  uuid.UUID
	MerchantID uuid.UUID
	ID         uuid.UUID
	Reason     string
	OperatorID uuid.UUID
	OperatedAt time.Time
}

// UpdateProductRequest 更新商品请求结构体
type UpdateProductRequest struct {
	ID          uuid.UUID
	MerchantID  uuid.UUID // 添加缺失字段
	Name        *string
	Price       *float64
	Description string
	Stock       *int
	Category    CategoryInfo
}

type ListProductsReq struct {
	Page         uint   `json:"page"`
	PageSize     uint   `json:"pageSize"`
	CategoryName string `json:"categoryName"`
}

type ListProductsResp struct {
	Product []*Product `json:"product"`
}

type GetProductResp struct {
	Product *Product `json:"product"`
}

type SearchProductsReq struct {
	Query string `json:"query"`
}
type SearchProductsResp struct {
	Result []*Product `json:"result"`
}

type AuditProductRequest struct {
	ProductID  uuid.UUID
	MerchantID uuid.UUID
	Action     uint64
	Reason     string
	OperatorID uuid.UUID
}

// DeleteProductRequest 完善DeleteProductRequest
type DeleteProductRequest struct {
	ID         uuid.UUID
	MerchantID uuid.UUID
}

// GetProductRequest 完善GetProductRequest
type GetProductRequest struct {
	ID uuid.UUID
}

// CreateProductRequest 完善CreateProductRequest
type CreateProductRequest struct {
	Product Product
}

type ImageModel struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint64
	URL       string
	IsPrimary bool
	SortOrder int32
}

type AttributeModel struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint64
	Key       string
	Type      string // "string", "array", "object"
	Value     string
}

type AuditInfoModel struct {
	AuditID    uint64
	Reason     string
	OperatorID uint64
	OperatedAt time.Time
}

// ProductRepo is a Greater repo.
type ProductRepo interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductReply, error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*Product, error)
	SubmitForAudit(ctx context.Context, req *SubmitAuditRequest) (*AuditRecord, error)
	AuditProduct(ctx context.Context, req *AuditProductRequest) (*AuditRecord, error)
	GetProduct(ctx context.Context, req *GetProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, req *DeleteProductRequest) error
}

// CanTransitionTo 添加状态转换方法
func (p *Product) CanTransitionTo(newStatus ProductStatus) bool {
	return validTransitions[p.Status][newStatus]
}

func (p *Product) ChangeStatus(newStatus ProductStatus) error {
	if !validTransitions[p.Status][newStatus] {
		return fmt.Errorf("invalid status transition from %d to %d", p.Status, newStatus)
	}
	p.Status = newStatus
	return nil
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductReply, error) {
	p.log.Debugf("CreateProduct: %+v", req)
	return p.repo.CreateProduct(ctx, req)
}

func (p *ProductUsecase) UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*Product, error) {
	return p.repo.UpdateProduct(ctx, req)
}

func (p *ProductUsecase) SubmitForAudit(ctx context.Context, req *SubmitAuditRequest) (*AuditRecord, error) {
	return p.repo.SubmitForAudit(ctx, req)
}

func (p *ProductUsecase) AuditProduct(ctx context.Context, req *AuditProductRequest) (*AuditRecord, error) {
	return p.repo.AuditProduct(ctx, req)
}

func (p *ProductUsecase) GetProduct(ctx context.Context, req *GetProductRequest) (*Product, error) {
	p.log.Debugf("GetProduct: %+v", req)
	return p.repo.GetProduct(ctx, req)
}

func (p *ProductUsecase) DeleteProduct(ctx context.Context, req DeleteProductRequest) (*emptypb.Empty, error) {
	p.log.Debugf("DeleteProduct: %+v", req)
	return &emptypb.Empty{}, nil
}

// 辅助验证函数
func validateProduct(p *Product) error {
	if p.Name == "" {
		return errors.New(403, "", "product name required")
	}
	if p.Price <= 0 {
		return errors.New(403, "", "invalid price")
	}
	return nil
}
