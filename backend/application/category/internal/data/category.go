package data

import (
	"backend/application/category/internal/biz"
	"backend/application/category/internal/data/models"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type categoryRepo struct {
	data *Data
	log  *log.Helper
}

func (r *categoryRepo) GetCategories(ctx context.Context, ids []int64) ([]*biz.Category, error) {
	// 执行SQL查询
	rows, err := r.data.DB(ctx).BatchGetCategories(ctx, ids)
	if err != nil {
		return nil, status.Error(codes.Internal, "数据库查询失败")
	}
	var categories []*biz.Category
	for _, row := range rows {
		categories = append(categories, &biz.Category{
			ID:        uint64(row.ID),
			ParentID:  uint64(row.ParentID),
			Level:     int(row.Level),
			Path:      row.Path,
			Name:      row.Name,
			SortOrder: int(row.SortOrder),
			IsLeaf:    row.IsLeaf,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}
	return categories, nil
}

func NewCategoryRepo(data *Data, logger log.Logger) biz.CategoryRepo {
	return &categoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// GetCategory 获取单个分类详情
func (r *categoryRepo) GetCategory(ctx context.Context, id int64) (*biz.Category, error) {
	dbCategory, err := r.data.DB(ctx).GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
		}
		return nil, fmt.Errorf("get category failed: %w", err)
	}
	return convertDBGetCategoryRow(dbCategory), nil
}

// DeleteCategory 删除分类及关联关系
// 流程：
// 1. 开启事务
// 2. 删除所有关联闭包记录
// 3. 删除分类记录
// 4. 更新父节点的is_leaf状态（如果没有其他子节点则设为true）
// 数据访问层实现
func (r *categoryRepo) DeleteCategory(ctx context.Context, id uint64) error {

	qtx := r.data.DB(ctx)

	// 1. 获取被删分类的path
	category, err := qtx.GetCategory(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return biz.ErrCategoryNotFound
		}
		return fmt.Errorf("get category failed: %w", err)
	}

	// 2. 执行级联删除
	err = qtx.DeleteCategory(ctx, models.DeleteCategoryParams{
		ID:   &category.ID,
		Path: &category.Path,
	})
	if err != nil {
		return fmt.Errorf("delete operation failed: %w", err)
	}

	return nil
}

// GetSubTree 获取指定分类的子树
// GetSubTree 获取子树（使用WITH RECURSIVE查询优化）
func (r *categoryRepo) GetSubTree(ctx context.Context, rootId uint64) ([]*biz.Category, error) {
	dbCategories, err := r.data.DB(ctx).GetSubTree(ctx, int64(rootId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, biz.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get subtree failed: %w", err)
	}

	result := make([]*biz.Category, 0, len(dbCategories))
	for _, c := range dbCategories {
		result = append(result, convertDBCategory(models.CategoriesCategories{
			ID:        c.ID,
			ParentID:  &c.ParentID,
			Level:     c.Level,
			Path:      c.CPath,
			Name:      c.Name,
			SortOrder: c.SortOrder,
			IsLeaf:    c.IsLeaf,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}))
	}
	return result, nil
}

// GetCategoryPath 获取分类路径（包含层级排序）
func (r *categoryRepo) GetCategoryPath(ctx context.Context, categoryID uint64) ([]*biz.Category, error) {
	dbCategories, err := r.data.DB(ctx).GetCategoryPath(ctx, int64(categoryID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, biz.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get category path failed: %w", err)
	}

	result := make([]*biz.Category, 0, len(dbCategories))
	for _, c := range dbCategories {
		result = append(result, convertDBCategory(models.CategoriesCategories{
			ID:        c.ID,
			ParentID:  &c.ParentID,
			Level:     c.Level,
			Path:      c.CPath,
			Name:      c.Name,
			SortOrder: c.SortOrder,
			IsLeaf:    c.IsLeaf,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}))
	}
	return result, nil
}

// GetLeafCategories 获取所有叶子分类
func (r *categoryRepo) GetLeafCategories(ctx context.Context) ([]*biz.Category, error) {
	dbCategories, err := r.data.DB(ctx).GetLeafCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("get leaf categories failed: %w", err)
	}

	var result []*biz.Category
	for _, c := range dbCategories {
		result = append(result, convertDBGetLeafCategoriesRow(c))
	}
	return result, nil
}

// GetClosureRelations 获取分类闭包关系
func (r *categoryRepo) GetClosureRelations(ctx context.Context, categoryId uint64) ([]*biz.ClosureRelation, error) {
	dbRelations, err := r.data.DB(ctx).GetClosureRelations(ctx, int64(categoryId))
	if err != nil {
		return nil, fmt.Errorf("get closure relations failed: %w", err)
	}

	result := make([]*biz.ClosureRelation, 0, len(dbRelations))
	for _, r := range dbRelations {
		result = append(result, &biz.ClosureRelation{
			Ancestor:   uint64(r.Ancestor),
			Descendant: uint64(r.Descendant),
			Depth:      uint32(r.Depth),
		})
	}
	return result, nil
}

// UpdateClosureDepth 更新闭包关系深度
func (r *categoryRepo) UpdateClosureDepth(ctx context.Context, req *biz.UpdateClosureDepth) error {
	tx := r.data.DB(ctx)

	// 验证分类存在
	if _, err := tx.GetCategory(ctx, req.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return biz.ErrCategoryNotFound
		}
		return fmt.Errorf("get category failed: %w", err)
	}

	// 执行深度更新（示例实现，实际可能需要更复杂的逻辑）
	err := tx.UpdateCategory(ctx, models.UpdateCategoryParams{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return err
	}
	return nil
}

// CreateCategory 创建分类
// 自动处理根分类：
// 通过 root_check 确保根分类（id=0）存在，无需手动初始化。
// 层级验证：
// 父分类层级 ≥n 时禁止插入（new_level 会设为 NULL，跳过插入）。
// 路径生成：
// 根据父分类的 path 生成唯一路径（使用 gen_random_uuid() 避免冲突）。
// 更新父分类状态：
// 若父分类之前是叶子节点（is_leaf=TRUE），插入子节点后将其设为非叶子节点。
// 闭包表维护：
// 自动维护 category_closure 表，记录分类的祖先 - 后代关系
func (r *categoryRepo) CreateCategory(ctx context.Context, req *biz.CreateCategoryReq) (*biz.Category, error) {
	qtx := r.data.DB(ctx)
	// 验证根分类存在
	// ParentID为 0 即视为创建一个根分类
	// 不为 0 就去查询该根类
	if req.ParentID != 0 {
		parentID := int32(req.ParentID)
		_, err := r.data.db.GetCategory(ctx, int64(parentID))
		if err != nil {
			return nil, biz.ErrParentIdUnprocessableEntiy
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

	category, err := qtx.CreateCategory(ctx, params)
	if err != nil {
		return nil, errors.New(500, "Failed", "failed to create category")
	}

	// fmt.Printf("id:%+v", category)
	// dbCategory, err := qtx.GetCategoryByID(ctx, category.ID)
	// if err != nil {
	// 	return nil, fmt.Errorf("get created category failed: %w", err)
	// }

	return convertDBCreateCategory(category), nil
}

func (r *categoryRepo) UpdateCategoryName(ctx context.Context, req *biz.Category) error {
	// params := models.UpdateCategoryNameParams{
	// 	ID:   req.ID,
	// 	Name: req.Name,
	// }
	//
	// err := r.data.DB(ctx).UpdateCategoryName(ctx, params)
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return biz.ErrCategoryNameNotFound
	// 	}
	// 	if strings.Contains(err.Error(), "unique constraint") {
	// 		return biz.ErrCategoryNameConflict
	// 	}
	// 	return biz.ErrCategoryFailed
	// }
	//
	// return err
	panic("TODO")
}

// convertDBCategory 转换数据库模型到业务模型
func convertDBCreateCategory(dbCategory models.CreateCategoryRow) *biz.Category {
	return &biz.Category{
		ID:        uint64(*dbCategory.ParentID),
		ParentID:  uint64(*dbCategory.ParentID),
		Level:     int(dbCategory.Level),
		Path:      dbCategory.Path,
		Name:      dbCategory.Name,
		SortOrder: int(dbCategory.SortOrder),
		IsLeaf:    dbCategory.IsLeaf,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}
}

func convertDBCategory(dbCategory models.CategoriesCategories) *biz.Category {
	return &biz.Category{
		ID:        uint64(dbCategory.ID),
		ParentID:  uint64(*dbCategory.ParentID),
		Level:     int(dbCategory.Level),
		Path:      dbCategory.Path,
		Name:      dbCategory.Name,
		SortOrder: int(dbCategory.SortOrder),
		IsLeaf:    dbCategory.IsLeaf,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}
}
func convertDBGetCategoryRow(dbCategory models.GetCategoryRow) *biz.Category {
	return &biz.Category{
		ID:        uint64(dbCategory.ID),
		ParentID:  uint64(dbCategory.ParentID),
		Level:     int(dbCategory.Level),
		Path:      dbCategory.Path,
		Name:      dbCategory.Name,
		SortOrder: int(dbCategory.SortOrder),
		IsLeaf:    dbCategory.IsLeaf,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}
}
func convertDBGetLeafCategoriesRow(dbCategory models.GetLeafCategoriesRow) *biz.Category {
	return &biz.Category{
		ID:        uint64(dbCategory.ID),
		ParentID:  uint64(dbCategory.ParentID),
		Level:     int(dbCategory.Level),
		Path:      dbCategory.Path,
		Name:      dbCategory.Name,
		SortOrder: int(dbCategory.SortOrder),
		IsLeaf:    dbCategory.IsLeaf,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}
}
