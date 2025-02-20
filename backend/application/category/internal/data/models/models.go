// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"time"
)

// 商品分类主表（ltree路径+闭包表双重优化）
type CategoriesCategories struct {
	ID        int64     `json:"id"`
	ParentID  int64     `json:"parent_id"`
	Level     int16     `json:"level"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	SortOrder int16     `json:"sort_order"`
	IsLeaf    bool      `json:"is_leaf"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 分类闭包关系表（存储所有层级关系）
type CategoriesCategoryClosure struct {
	Ancestor   int64 `json:"ancestor"`
	Descendant int64 `json:"descendant"`
	Depth      int16 `json:"depth"`
}
