package service

import (
	pb "backend/api/category/v1"
	"backend/application/category/internal/biz"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryServiceService struct {
	pb.UnimplementedCategoryServiceServer

	uc *biz.CategoryUsecase
}

func NewCategoryServiceService(uc *biz.CategoryUsecase) *CategoryServiceService {
	return &CategoryServiceService{uc: uc}
}

func (s *CategoryServiceService) GetCategoryTree(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	return &pb.Category{}, nil
}
func (s *CategoryServiceService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := s.uc.CreateCategory(ctx, &biz.Category{
		ParentID:  req.ParentId,
		Name:      req.Name,
		SortOrder: int(req.SortOrder),
	})
	if err != nil {
		return nil, err
	}

	return &pb.Category{
		Id:        category.ID,
		Name:      category.Name,
		SortOrder: int32(category.SortOrder),
	}, nil
}
func (s *CategoryServiceService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
	return &pb.Category{}, nil
}
func (s *CategoryServiceService) DeleteCategory(context.Context, *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	return nil, nil
}
