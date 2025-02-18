package data

import (
	"context"

	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"

	"backend/application/category/internal/biz"
	"backend/application/category/internal/data/models"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
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

// GetCategory 获取单个分类详情
func (r *categoryRepo) GetCategory(ctx context.Context, id int64) (*biz.Category, error) {
	dbCategory, err := r.data.DB(ctx).GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, biz.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get category failed: %w", err)
	}
	return convertDBCategory(dbCategory), nil
}

// DeleteCategory 删除分类及关联关系
// 流程：
// 1. 开启事务
// 2. 删除所有关联闭包记录
// 3. 删除分类记录
// 4. 更新父分类叶子状态
func (r *categoryRepo) DeleteCategory(ctx context.Context, id int64) error {
	// tx, err := r.data.pool.BeginTx(ctx, pgx.TxOptions{
	// 	IsoLevel: pgx.Serializable, // 需要最高级别隔离
	// })
	qtx := r.data.DB(ctx)

	// 获取被删分类信息用于后续处理
	category, err := qtx.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return biz.ErrCategoryNotFound
		}
		return fmt.Errorf("get category failed: %w", err)
	}

	// 删除闭包关系
	if err := qtx.DeleteClosureRelations(ctx, &id); err != nil {
		return fmt.Errorf("delete closure relations failed: %w", err)
	}

	// 删除分类记录
	if err := qtx.DeleteCategory(ctx, &id); err != nil {
		return fmt.Errorf("delete category failed: %w", err)
	}

	// 更新父分类的叶子状态
	if category.ParentID != 0 { // 根分类无父分类
		if err := qtx.UpdateParentLeafStatus(ctx, &category.ParentID); err != nil {
			return fmt.Errorf("update parent leaf status failed: %w", err)
		}
	}

	return err
}

// GetSubTree 获取指定分类的子树
func (r *categoryRepo) GetSubTree(ctx context.Context, rootID int64) ([]*biz.Category, error) {
	dbCategories, err := r.data.DB(ctx).GetSubTree(ctx, &rootID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, biz.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get subtree failed: %w", err)
	}

	result := make([]*biz.Category, 0, len(dbCategories))
	for _, c := range dbCategories {
		result = append(result, convertDBCategory(c))
	}
	return result, nil
}

// GetCategoryPath 获取分类的完整路径
func (r *categoryRepo) GetCategoryPath(ctx context.Context, categoryID int64) ([]*biz.Category, error) {
	dbCategories, err := r.data.DB(ctx).GetCategoryPath(ctx, categoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, biz.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get category path failed: %w", err)
	}

	result := make([]*biz.Category, 0, len(dbCategories))
	for _, c := range dbCategories {
		result = append(result, convertDBCategory(c))
	}
	return result, nil
}

// GetLeafCategories 获取所有叶子分类（三级分类）
func (r *categoryRepo) GetLeafCategories(ctx context.Context) ([]*biz.Category, error) {
	dbCategories, err := r.data.DB(ctx).GetLeafCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("get leaf categories failed: %w", err)
	}

	result := make([]*biz.Category, 0, len(dbCategories))
	for _, c := range dbCategories {
		result = append(result, convertDBCategory(c))
	}
	return result, nil
}

// GetClosureRelations 获取分类闭包关系
func (r *categoryRepo) GetClosureRelations(ctx context.Context, categoryID int64) ([]*biz.ClosureRelation, error) {
	dbRelations, err := r.data.DB(ctx).GetClosureRelations(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("get closure relations failed: %w", err)
	}

	result := make([]*biz.ClosureRelation, 0, len(dbRelations))
	for _, r := range dbRelations {
		result = append(result, &biz.ClosureRelation{
			Ancestor:   r.Ancestor,
			Descendant: r.Descendant,
			Depth:      int32(int(r.Depth)),
		})
	}
	return result, nil
}

// UpdateClosureDepth 更新闭包关系深度
func (r *categoryRepo) UpdateClosureDepth(ctx context.Context, categoryID int64, delta int32) error {
	var deltaType = int16(delta)
	if err := r.data.DB(ctx).UpdateClosureDepth(ctx, models.UpdateClosureDepthParams{
		CategoryID: &categoryID,
		Delta:      &deltaType,
	}); err != nil {
		return fmt.Errorf("update closure depth failed: %w", err)
	}
	return nil
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
func (r *categoryRepo) CreateCategory(ctx context.Context, req *biz.CreateCategoryReq) (*biz.Category, error) {
	qtx := r.data.DB(ctx)
	// 验证根分类存在
	// ParentID为 0 即视为创建一个根分类
	// 不为 0 就去查询该根类
	if req.ParentID != 0 {
		_, err := r.data.db.GetCategoryByID(ctx, req.ParentID)
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

	id, err := qtx.CreateCategory(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, biz.ErrCategoryNameNotFound
		}
		return nil, errors.New(500, "Failed", "failed to create category")
	}

	dbCategory, err := qtx.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get created category failed: %w", err)
	}

	return convertDBCategory(dbCategory), nil
}

func (r *categoryRepo) UpdateCategoryName(ctx context.Context, req *biz.Category) error {
	params := models.UpdateCategoryNameParams{
		ID:   req.ID,
		Name: req.Name,
	}

	err := r.data.DB(ctx).UpdateCategoryName(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return biz.ErrCategoryNameNotFound
		}
		if strings.Contains(err.Error(), "unique constraint") {
			return biz.ErrCategoryNameConflict
		}
		return biz.ErrCategoryFailed
	}

	return err
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
