package biz

import (
	"github.com/google/uuid"
	"time"
)

// Category 分类
type Category struct {
	ID        uuid.UUID
	Name      string
	Level     int
	ParentID  string
	Children  []*Category
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetCategoryTreeRequest struct {
	Query string
}

type CategoryTree struct {
	Category []Category
}

type UpdateCategoryRequest struct {
	ID   string
	Name string
}

type DeleteCategoryRequest struct {
	ID string
}

type DeleteCategoryReply struct {
	Message string
}

type CreateCategoryRequest struct {
	Name     string
	Level    uint32
	ParentID *string
	// Children  []*Category
}
