package data

import (
	"backend/application/category/internal/biz"
	"backend/application/category/internal/data/models"
	"context"
	"fmt"
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

func (c *categoryRepo) CreateCategory(ctx context.Context, req *biz.CreateCategoryRequest) (*biz.Category, error) {
	// TODO 可以使用不存在即创建
	// var category models.CreateCategoryRow
	// categoryName, err := c.data.db.GetCategoryByName(ctx, &req.Name)
	// if err != nil {
	// 	if !errors.Is(err, sql.ErrNoRows) {
	//
	// 	}
	// 	return nil, err
	// }
	// fmt.Println("categoryName:", categoryName)
	fmt.Printf("CreateCategory: %+v\n", req)
	category, err := c.data.DB(ctx).CreateCategory(ctx, models.CreateCategoryParams{
		Name:  req.Name,
		Level: int32(req.Level),
	})
	if err != nil {
		return nil, fmt.Errorf("category creation failed: %w", err)
	}
	return &biz.Category{
		ID:   category.ID,
		Name: category.Name,
		// CreatedAt: category.CreatedAt,
		// UpdatedAt: category.UpdatedAt,
	}, err
}

func (c *categoryRepo) DeleteCategory(context.Context, *biz.DeleteCategoryRequest) (*biz.DeleteCategoryReply, error) {
	panic("todo")
}
func (c *categoryRepo) GetCategoryTree(context.Context, *biz.GetCategoryTreeRequest) (*biz.CategoryTree, error) {
	panic("todo")
}
func (c *categoryRepo) UpdateCategory(context.Context, *biz.UpdateCategoryRequest) (*biz.Category, error) {
	panic("todo")
}
