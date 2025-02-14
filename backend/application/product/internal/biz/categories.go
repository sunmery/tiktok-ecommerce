package biz

import (
	"context"
	"errors"
	"time"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Category struct {
	ID       uint32   `json:"id"`
	Name     string   `json:"name"`
	ParentID *uint32  `json:"parent_id"` // nil 表示顶级分类
	Children []Category `json:"children,omitempty"` // 非数据库字段，用于返回树形结构
	CreateAt time.Time `json:"create_at"`
	IsActive bool     `json:"is_active"`
}

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	ParentID *int32 `json:"parentID"`
	Owner   string `json:"owner"`
	Username string `json:"username"`
}

type CategoryReply struct {
	Category Category `json:"category"`
}

type ListCategoriesResp struct {
	Categories []Category `json:"categories"`
}

type GetCategoryChildrenResp struct {
	Categories []Category `json:"categories"`
}

func (P *ProductUsecase) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryReply, error) {
	P.log.WithContext(ctx).Infof("CreateCategory: %+v", req)
	if req.Name == "" {
		return nil, errors.New("category name can not be empty")
	}
	return P.repo.CreateCategory(ctx, req)
}

func (P *ProductUsecase) ListCategories(ctx context.Context,_ *emptypb.Empty) (*ListCategoriesResp, error) {
	P.log.WithContext(ctx).Infof("ListCategories")
	return P.repo.ListCategories(ctx, &emptypb.Empty{})
}

func (P *ProductUsecase) GetCategoryChildren(ctx context.Context, categoryID uint32) (*GetCategoryChildrenResp, error) {
	P.log.WithContext(ctx).Infof("GetCategoryChildren: %v", categoryID)
	categoryTree, _ := P.repo.ListCategories(ctx, &emptypb.Empty{})

	// **递归查找子分类**
	children := FindChildren(categoryTree.Categories, categoryID) // Pass categoryTree.Categories instead of categoryTree
	if len(children) == 0 {
		return nil, errors.New("no subcategories found")
	}

	return &GetCategoryChildrenResp{Categories: children}, nil
}

// 递归查找子分类
func FindChildren(categories []Category, parentID uint32) []Category {
	var result []Category
	for _, cat := range categories {
		if cat.ParentID != nil && *cat.ParentID == parentID {
			// **递归查找子分类**
			children := FindChildren(categories, cat.ID)
			if len(children) > 0 {
				cat.Children = children
			}
			result = append(result, cat)
		}
	}
	return result
}