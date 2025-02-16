package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

var (
	ErrInvalidparentIdArgument = errors.New("category: invalid parent_id argument")
)

// Category 分类
type Category struct {
	ID        int64
	ParentID  int64
	Level     int
	Path      string
	Name      string
	SortOrder int
	IsLeaf    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeleteCategoryReply struct {
	Success bool
}

type CategoryRepo interface {
	CreateCategory(ctx context.Context, req *Category) (*Category, error)
	UpdateCategory(ctx context.Context, req *Category) (*Category, error)
	DeleteCategory(ctx context.Context, id string) (*DeleteCategoryReply, error)
	GetCategoryTree(ctx context.Context, id uint32) ([]*Category, error)
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

func (uc *CategoryUsecase) CreateCategory(ctx context.Context, req *Category) (*Category, error) {
	uc.log.WithContext(ctx).Infof("CreateCategory request: %+v", req)
	return uc.repo.CreateCategory(ctx, req)
}

func (uc *CategoryUsecase) DeleteCategory(ctx context.Context, id string) (*DeleteCategoryReply, error) {
	uc.log.WithContext(ctx).Infof("DeleteCategory request: %+v", id)
	return uc.repo.DeleteCategory(ctx, id)
}
func (uc *CategoryUsecase) GetCategoryTree(ctx context.Context, id uint32) ([]*Category, error) {
	uc.log.WithContext(ctx).Info("GetCategoryTree", id)
	return uc.repo.GetCategoryTree(ctx, id)

}
func (uc *CategoryUsecase) UpdateCategory(ctx context.Context, req *Category) (*Category, error) {
	uc.log.WithContext(ctx).Infof("UpdateCategory request: %+v", req)
	return uc.repo.UpdateCategory(ctx, req)
}
