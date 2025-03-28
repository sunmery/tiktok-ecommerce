// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"time"
)

type CategoriesCategories struct {
	ID        int64
	ParentID  *int64
	Level     int16
	Path      string
	Name      string
	SortOrder int16
	IsLeaf    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoriesCategoryClosure struct {
	Ancestor   int64
	Descendant int64
	Depth      int16
}
