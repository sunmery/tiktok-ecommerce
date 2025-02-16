package data

import (
	"context"
	"fmt"
	"strings"

	"backend/application/category/internal/biz"
	"backend/application/category/internal/data/models"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type categoryRepo struct {
	data *Data
	log  *log.Helper
}

func NewCategoryRepo(data *Data, logger log.Logger) biz.CategoryRepo {
	return &categoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateCategory 创建分类
// 自动处理根分类：
// 通过 root_check 确保根分类（id=0）存在，无需手动初始化。
// 层级验证：
// 父分类层级 ≥3 时禁止插入（new_level 会设为 NULL，跳过插入）。
// 路径生成：
// 根据父分类的 path 生成唯一路径（使用 gen_random_uuid() 避免冲突）。
// 更新父分类状态：
// 若父分类之前是叶子节点（is_leaf=TRUE），插入子节点后将其设为非叶子节点。
// 闭包表维护：
// 自动维护 category_closure 表，记录分类的祖先 - 后代关系
func (r *categoryRepo) CreateCategory(ctx context.Context, req *biz.Category) (*biz.Category, error) {
	qtx := r.data.DB(ctx)
	// 验证根分类存在
	// ParentID为 0 即视为创建一个根分类
	// 不为 0 就去查询该根类
	if req.ParentID != 0 {
		_, err := r.data.db.GetCategoryByID(ctx, req.ParentID)
		if err != nil {
			return nil, biz.ErrInvalidparentIdArgument
		}
	}

	// 执行创建
	sortOrder := int16(req.SortOrder)
	// 参数
	// ParentID: 根分类 ID
	// Name: 分类名称
	// SortOrder: 分类的排序权重
	params := models.CreateCategoryParams{
		ParentID:  req.ParentID,
		Name:      req.Name,
		SortOrder: sortOrder,
	}

	id, err := qtx.CreateCategory(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, status.Errorf(codes.AlreadyExists, "category name exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
	}

	dbCategory, err := r.data.DB(ctx).GetCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get created category failed: %w", err)
	}

	return convertDBCategory(dbCategory), nil
}

func (r *categoryRepo) UpdateCategory(ctx context.Context, req *biz.Category) (*biz.Category, error) {
	// params := models.UpdateCategoryParams{
	// 	ID:   req.ID,
	// 	Name: req.Name,
	// }
	//
	// updated, err := r.q.UpdateCategoryName(ctx, params)
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return nil, status.Errorf(codes.NotFound, "category not found")
	// 	}
	// 	if strings.Contains(err.Error(), "unique constraint") {
	// 		return nil, status.Errorf(codes.AlreadyExists, "category name exists")
	// 	}
	// 	return nil, status.Errorf(codes.Internal, "failed to update category: %v", err)
	// }
	//
	// return &biz.Category{
	// 	ID:        updated.ID,
	// 	ParentID:  updated.ParentID,
	// 	Level:     int(updated.Level),
	// 	Path:      updated.Path,
	// 	Name:      updated.Name,
	// 	SortOrder: int(updated.SortOrder),
	// 	IsLeaf:    updated.IsLeaf,
	// 	CreatedAt: updated.CreatedAt.Time,
	// 	UpdatedAt: updated.UpdatedAt.Time,
	// }, nil
	panic("implement me")
}

func (r *categoryRepo) DeleteCategory(ctx context.Context, id string) (*biz.DeleteCategoryReply, error) {
	// categoryID, err := dbtypes.ParseInt64(id)
	// if err != nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, "invalid category id")
	// }
	//
	// if err := r.q.DeleteCategory(ctx, categoryID); err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return nil, status.Errorf(codes.NotFound, "category not found")
	// 	}
	// 	return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
	// }
	//
	// return &biz.DeleteCategoryReply{Success: true}, nil
	panic("implement me")
}

func (r *categoryRepo) GetCategoryTree(ctx context.Context, id uint32) ([]*biz.Category, error) {
	// dbCategory, err := r.data.DB(ctx).GetCategoryByID(ctx, int64(id))
	// if err != nil {
	// 	return nil, fmt.Errorf("get category failed: %w", err)
	// }
	// var tree []*biz.Category
	// return convertDBCategory(dbCategory), nil
	panic("implement me")
}

// convertDBCategory 转换数据库模型到业务模型
func convertDBCategory(dbCategory models.CategoriesCategories) *biz.Category {
	return &biz.Category{
		ID:        dbCategory.ID,
		ParentID:  dbCategory.ParentID,
		Level:     int(dbCategory.Level),
		Path:      dbCategory.Path,
		Name:      dbCategory.Name,
		SortOrder: int(dbCategory.SortOrder),
		IsLeaf:    dbCategory.IsLeaf,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}
}
