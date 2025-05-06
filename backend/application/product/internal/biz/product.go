package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ProductStatusDraft    ProductStatus = iota // 商品草稿
	ProductStatusPending                       // 商品待审核。
	ProductStatusApproved                      // 商品审核通过。
	ProductStatusRejected                      // 商品审核未通过。
	ProductStatusSoldOut                       // 商品因某种原因不可购买。
)

const (
	Approved AuditAction = 1
	Rejected AuditAction = 2
)

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
	ProductStatusSoldOut: {
		ProductStatusSoldOut: true,
	},
}

type (
	AuditAction int
)

// AuditRecord 完善AuditRecord定义
type AuditRecord struct {
	ID         uuid.UUID
	ProductID  uuid.UUID
	OldStatus  ProductStatus
	NewStatus  ProductStatus
	Reason     string
	OperatorID uuid.UUID
	OperatedAt time.Time
}
type AuditInfo struct {
	AuditId    uuid.UUID // 审核记录ID
	Reason     string    // 审核意见/驳回原因
	OperatorId uuid.UUID // 操作人ID
	OperatedAt time.Time // 操作时间
}

// Product 商品领域模型
type (
	// ProductImage 商品图片信息
	ProductImage struct {
		URL       string
		IsPrimary bool
		SortOrder *int
	}
	// ProductStatus 商品状态
	ProductStatus int32
	// CategoryInfo 分类信息
	CategoryInfo struct {
		CategoryId   uint64
		CategoryName string
		SortOrder    int32
	}

	Product struct {
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
		Attributes  map[string]any
		Inventory   Inventory // 库存
	}
)

type GetProductsBatchRequest struct {
	ProductIds  []uuid.UUID
	MerchantIds []uuid.UUID
}

type SubmitAuditRequest struct {
	ProductID  uuid.UUID
	MerchantID uuid.UUID
	ID         uuid.UUID
	Reason     string
	OperatorID uuid.UUID
	OperatedAt time.Time
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
	Status     ProductStatus // 删除商品的状态, 默认为4, 但未来可能会根据需求进行修改
}

// GetProductRequest 完善GetProductRequest
type GetProductRequest struct {
	ID         uuid.UUID
	MerchantID uuid.UUID
}

// CreateProductRequest 完善CreateProductRequest
type (
	CreateProductRequest struct {
		Name        string
		Price       float64
		Description string
		MerchantId  uuid.UUID
		Images      []*ProductImage
		Status      ProductStatus
		Category    CategoryInfo
		Attributes  map[string]any
		Stock       uint32
	}
	CreateProductReply struct {
		ID        uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

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

type ListRandomProductsRequest struct {
	Page     uint32
	PageSize uint32
	Status   uint32
}

// Products 批量商品
type Products struct {
	Items []*Product
}

type SearchProductsByNameRequest struct {
	Name     string
	Page     uint32
	PageSize uint32
	// 将自然语言文本转换为全文搜索查询条件（tsquery 类型），主要功能包括：
	// 词素标准化（Normalization）
	// 移除停用词（的、是、the、a 等）
	// 执行词干提取（running → run，dogs → dog）
	// 逻辑运算符转换
	// 自动用 &（AND）连接词汇
	// 示例："red apple" → 'red' & 'apple'
	// 安全过滤
	// 自动转义特殊字符（! : & 等）
	Query string
}
type (
	UploadMethod             int32
	UploadProductFileRequest struct {
		Method      UploadMethod
		ContentType *string
		BucketName  *string
		FilePath    *string
		FileName    *string
	}
	UploadProductFileReply struct {
		UploadUrl   string
		DownloadUrl string
		BucketName  *string
		ObjectName  string
		FormData    map[string]string
	}
)

// GetCategoryProducts 根据分类获取商品
type GetCategoryProducts struct {
	CategoryID uint32
	Status     uint32 // 商品状态机
	Page       int64
	PageSize   int64
}

// GetCategoryWithChildrenProducts 根据分类及其所有子分类获取商品
type GetCategoryWithChildrenProducts struct {
	CategoryID uint32
	Status     uint32 // 商品状态机
	Page       int64
	PageSize   int64
}

type (
	ProductDraft struct {
		Name           string
		Description    string
		Price          float64
		CurrentAuditID *string // 使用指针处理空值
		Stock          uint32
		MerchantId     uuid.UUID
		Status         ProductStatus
		Attributes     map[string]interface{}
		Category       CategoryInfo
		Images         []*ProductImage
	}

	BatchProductError struct {
		Index           int
		Message         string
		OriginalProduct *ProductDraft
	}

	CreateProductBatchRequest struct {
		Products []*ProductDraft
	}

	CreateProductBatchReply struct {
		SuccessCount uint32
		FailedCount  uint32
		Errors       []*BatchProductError
		ProductIds   []uuid.UUID
	}
)

// ProductRepo is a Greater repo.
type ProductRepo interface {
	UploadProductFile(ctx context.Context, req *UploadProductFileRequest) (*UploadProductFileReply, error)
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductReply, error)
	CreateProductBatch(ctx context.Context, req *CreateProductBatchRequest) (*CreateProductBatchReply, error)
	SubmitForAudit(ctx context.Context, req *SubmitAuditRequest) (*AuditRecord, error)
	AuditProduct(ctx context.Context, req *AuditProductRequest) (*AuditRecord, error)
	GetProduct(ctx context.Context, req *GetProductRequest) (*Product, error)
	GetProductBatch(ctx context.Context, req *GetProductsBatchRequest) (*Products, error)
	GetCategoryProducts(ctx context.Context, req *GetCategoryProducts) (*Products, error)
	GetCategoryWithChildrenProducts(ctx context.Context, req *GetCategoryWithChildrenProducts) (*Products, error)
	DeleteProduct(ctx context.Context, req *DeleteProductRequest) error
	ListRandomProducts(ctx context.Context, req *ListRandomProductsRequest) (*Products, error)
	SearchProductsByName(ctx context.Context, req *SearchProductsByNameRequest) (*Products, error)

	// UpdateInventory 更新库存
	UpdateInventory(ctx context.Context, req *UpdateInventoryRequest) (*UpdateInventoryReply, error)
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

func (p *ProductUsecase) UploadProductFile(ctx context.Context, req *UploadProductFileRequest) (*UploadProductFileReply, error) {
	return p.repo.UploadProductFile(ctx, req)
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductReply, error) {
	p.log.WithContext(ctx).Debugf("CreateProduct: %v", req)
	return p.repo.CreateProduct(ctx, req)
}

func (p *ProductUsecase) CreateProductBatch(ctx context.Context, req *CreateProductBatchRequest) (*CreateProductBatchReply, error) {
	p.log.WithContext(ctx).Debugf("CreateProductBatch: %v", req)
	return p.repo.CreateProductBatch(ctx, req)
}

func (p *ProductUsecase) SubmitForAudit(ctx context.Context, req *SubmitAuditRequest) (*AuditRecord, error) {
	return p.repo.SubmitForAudit(ctx, req)
}

func (p *ProductUsecase) AuditProduct(ctx context.Context, req *AuditProductRequest) (*AuditRecord, error) {
	return p.repo.AuditProduct(ctx, req)
}

func (p *ProductUsecase) GetProduct(ctx context.Context, req *GetProductRequest) (*Product, error) {
	return p.repo.GetProduct(ctx, req)
}

func (p *ProductUsecase) GetProductBatch(ctx context.Context, req *GetProductsBatchRequest) (*Products, error) {
	return p.repo.GetProductBatch(ctx, req)
}

func (p *ProductUsecase) DeleteProduct(ctx context.Context, req DeleteProductRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *ProductUsecase) ListRandomProducts(ctx context.Context, req *ListRandomProductsRequest) (*Products, error) {
	return p.repo.ListRandomProducts(ctx, req)
}

func (p *ProductUsecase) SearchProductsByName(ctx context.Context, req *SearchProductsByNameRequest) (*Products, error) {
	return p.repo.SearchProductsByName(ctx, req)
}

func (p *ProductUsecase) GetCategoryProducts(ctx context.Context, req *GetCategoryProducts) (*Products, error) {
	return p.repo.GetCategoryProducts(ctx, req)
}

func (p *ProductUsecase) GetCategoryWithChildrenProducts(ctx context.Context, req *GetCategoryWithChildrenProducts) (*Products, error) {
	return p.repo.GetCategoryWithChildrenProducts(ctx, req)
}

// UpdateInventory 更新库存
func (p *ProductUsecase) UpdateInventory(ctx context.Context, req *UpdateInventoryRequest) (*UpdateInventoryReply, error) {
	p.log.WithContext(ctx).Debugf("UpdateInventory: req:%+v", req)
	return p.repo.UpdateInventory(ctx, req)
}
