package data

import (
	"backend/application/product/internal/biz"
	"backend/application/product/internal/data/models"
	"google.golang.org/protobuf/types/known/emptypb"
	"context"
)

func (c *productRepo) CreateCategory(ctx context.Context, req *biz.CreateCategoryRequest) (*biz.CategoryReply, error) {
	db := c.data.DB(ctx)
	category, err := db.CreateCategories(ctx, models.CreateCategoriesParams{
		Name: req.Name,
		ParentID: req.ParentID,
	})
	if err != nil {
		return nil, err
	}
	return &biz.CategoryReply{
		Category: biz.Category{
			ID:   uint32(category.ID),
			Name: category.Name,
		},
	}, nil
}

func (c *productRepo) ListCategories(ctx context.Context,_ *emptypb.Empty) (*biz.ListCategoriesResp, error) {
	db := c.data.DB(ctx)
	var categories []biz.Category
	productsCategories, err := db.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	for _, pc := range productsCategories {
		var parentID *uint32
		if pc.ParentID != nil {
			convertedParentID := uint32(*pc.ParentID)
			parentID = &convertedParentID
		}
		category := biz.Category{
			ID:       uint32(pc.ID),
			Name:     pc.Name,
			ParentID: parentID,
		}
		categories = append(categories, category)
	}

	return &biz.ListCategoriesResp{
		Categories: BuildCategoryTree(categories, nil),
	}, nil
}

// 构建树形结构
func BuildCategoryTree(categories []biz.Category, parentID *uint32) []biz.Category {
	var tree []biz.Category
	for _, cat := range categories {
		if (cat.ParentID == nil && parentID == nil) || (cat.ParentID != nil && parentID != nil && uint32(*cat.ParentID) == *parentID) {
			convertedID := uint32(cat.ID)
			children := BuildCategoryTree(categories, &convertedID)
			if len(children) > 0 {
				cat.Children = children
			}
			tree = append(tree, cat)
		}
	}
	return tree
}


