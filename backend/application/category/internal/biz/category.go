package biz

import (
	"context"
	"time"

	v1 "backend/api/category/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrParentIdUnprocessableEntiy PARENT_ID不符合业务规则
	ErrParentIdUnprocessableEntiy = errors.BadRequest(v1.ErrorReason_PARENT_ID_UNPROCESSABLE_ENTITY.String(), "category: invalid parent_id argument")
	// ErrCategoryNameNotFound 找不到该分类名称
	ErrCategoryNameNotFound = errors.NotFound(v1.ErrorReason_CATEGORY_NAME_NOT_FOUND.String(), "category: category name not found")
	// ErrCategoryNotFound 找不到该分类
	ErrCategoryNotFound = errors.NotFound(v1.ErrorReason_CATEGORY_NOT_FOUND.String(), "category: category not found")
	// ErrCategoryNameConflict 分类已存在
	ErrCategoryNameConflict = errors.New(409, "Already Exists ", "category name exists")
	// ErrCategoryFailed 内部错误
	ErrCategoryFailed = errors.New(500, "Failed", "failed category")
	// ErrCategoryHasChildren 存在子分类不可删除
	ErrCategoryHasChildren = errors.New(403, "Forbidden", "存在子分类不可删除")
)

// Category 分类
type Category struct {
	ID        uint64
	ParentID  uint64
	Level     int
	Path      string
	Name      string
	SortOrder int
	IsLeaf    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCategoryReq struct {
	ParentID  int64
	Name      string
	SortOrder int
}

type DeleteCategoryReply struct {
	Success bool
}

type CategoryRepo interface {
	// CreateCategory 基础操作
	CreateCategory(ctx context.Context, req *CreateCategoryReq) (*Category, error)
	GetCategory(ctx context.Context, id int64) (*Category, error)
	GetCategories(ctx context.Context, ids []int64) ([]*Category, error)
	UpdateCategoryName(ctx context.Context, req *Category) error
	DeleteCategory(ctx context.Context, id uint64) error

	// GetSubTree 树形操作
	GetSubTree(ctx context.Context, rootID uint64) ([]*Category, error)
	GetCategoryPath(ctx context.Context, categoryId uint64) ([]*Category, error)
	GetLeafCategories(ctx context.Context) ([]*Category, error)
	// GetDirectSubCategories 获取指定分类的直接子分类（只返回下一级）
	GetDirectSubCategories(ctx context.Context, parentId uint64) ([]*Category, error)

	// GetClosureRelations 闭包关系
	GetClosureRelations(ctx context.Context, categoryId uint64) ([]*ClosureRelation, error)
	UpdateClosureDepth(ctx context.Context, req *UpdateClosureDepth) error
}

type CategoryUsecase struct {
	repo CategoryRepo
	log  *log.Helper
}

func NewCategoryUsecase(repo CategoryRepo, logger log.Logger) *CategoryUsecase {
	return &CategoryUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// ClosureRelation 闭包关系业务对象
type ClosureRelation struct {
	Ancestor   uint64
	Descendant uint64
	Depth      uint32
}
type UpdateClosureDepth struct {
	ID      int64
	Name string
}

// CreateCategory 创建分类
func (uc *CategoryUsecase) CreateCategory(ctx context.Context, req *CreateCategoryReq) (*Category, error) {
	uc.log.WithContext(ctx).Debugf("CreateCategory request: %+v", req)
	return uc.repo.CreateCategory(ctx, req)
}

// UpdateCategoryName 更新分类名称
func (uc *CategoryUsecase) UpdateCategoryName(ctx context.Context, req *Category) error {
	uc.log.WithContext(ctx).Debugf("UpdateCategory request: %+v", req)
	return uc.repo.UpdateCategoryName(ctx, req)
}

// GetCategory 获取单个分类详情
func (uc *CategoryUsecase) GetCategory(ctx context.Context, id int64) (*Category, error) {
	uc.log.WithContext(ctx).Debugf("GetCategory request: %d", id)
	return uc.repo.GetCategory(ctx, id)
}

// UpdateCategory 更新分类信息
func (uc *CategoryUsecase) UpdateCategory(ctx context.Context, req *Category) error {
	uc.log.WithContext(ctx).Debugf("UpdateCategory request: %+v", req)
	return uc.repo.UpdateCategoryName(ctx, req)
}

// DeleteCategory 删除分类（包含子树）
func (uc *CategoryUsecase) DeleteCategory(ctx context.Context, id uint64) error {
	uc.log.WithContext(ctx).Debugf("DeleteCategory request: %d", id)
	return uc.repo.DeleteCategory(ctx, id)
}

// GetSubTree 获取子树
func (uc *CategoryUsecase) GetSubTree(ctx context.Context, rootID uint64) ([]*Category, error) {
	uc.log.WithContext(ctx).Debugf("GetSubTree request: %d", rootID)
	return uc.repo.GetSubTree(ctx, rootID)
}

// GetCategoryPath 获取分类路径
func (uc *CategoryUsecase) GetCategoryPath(ctx context.Context, categoryID uint64) ([]*Category, error) {
	uc.log.WithContext(ctx).Debugf("GetCategoryPath request: %d", categoryID)
	return uc.repo.GetCategoryPath(ctx, categoryID)
}

// GetLeafCategories 获取所有叶子分类
func (uc *CategoryUsecase) GetLeafCategories(ctx context.Context) ([]*Category, error) {
	// uc.log.WithContext(ctx).Debugf("GetLeafCategories request: %d", level)
	return uc.repo.GetLeafCategories(ctx)
}

// GetDirectSubCategories 获取指定分类的直接子分类（只返回下一级）
func (uc *CategoryUsecase) GetDirectSubCategories(ctx context.Context, parentId uint64) ([]*Category, error) {
	uc.log.WithContext(ctx).Debugf("GetDirectSubCategories request: %d", parentId)
	return uc.repo.GetDirectSubCategories(ctx, parentId)
}

// GetClosureRelations 获取闭包关系
func (uc *CategoryUsecase) GetClosureRelations(ctx context.Context, categoryID uint64) ([]*ClosureRelation, error) {
	uc.log.WithContext(ctx).Debugf("GetClosureRelations request: %d", categoryID)
	return uc.repo.GetClosureRelations(ctx, categoryID)
}

// UpdateClosureDepth 更新闭包深度
func (uc *CategoryUsecase) UpdateClosureDepth(ctx context.Context, req *UpdateClosureDepth) error {
	uc.log.WithContext(ctx).Debugf("UpdateClosureDepth request: %d delta:%d", req)
	return uc.repo.UpdateClosureDepth(ctx, req)
}
func (uc *CategoryUsecase) GetCategories(ctx context.Context, ids []int64) ([]*Category, error){
	uc.log.WithContext(ctx).Debugf("GetCategories request: %d", ids)
	return uc.repo.GetCategories(ctx, ids)
}
