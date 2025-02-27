package service

import (
	"context"
	"errors"
	"fmt"

	"backend/application/category/internal/biz"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "backend/api/category/v1"
)

type CategoryServiceService struct {
	pb.UnimplementedCategoryServiceServer
	uc *biz.CategoryUsecase
}

func NewCategoryServiceService(uc *biz.CategoryUsecase) *CategoryServiceService {
	return &CategoryServiceService{uc: uc}
}

// CreateCategory 创建分类
// 接口文档：POST /v1/category
func (s *CategoryServiceService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	// 参数校验
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "分类名称不能为空")
	}
	if req.SortOrder < 0 {
		return nil, status.Error(codes.InvalidArgument, "排序序号不能为负数")
	}

	// 调用业务逻辑层
	category, err := s.uc.CreateCategory(ctx, &biz.CreateCategoryReq{
		ParentID:  req.ParentId,
		Name:      req.Name,
		SortOrder: int(req.SortOrder),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("create category failed: %v", err))
	}

	// 转换响应格式
	return &pb.Category{
		Id:        int64(category.ID),
		CreatedAt: timestamppb.New(category.CreatedAt),
		UpdatedAt: timestamppb.New(category.UpdatedAt),
	}, nil
}

// GetCategory 获取分类详情
// 接口文档：GET /v1/category/{id}
func (s *CategoryServiceService) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	// 参数校验
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}

	// 调用业务逻辑层
	category, err := s.uc.GetCategory(ctx, req.Id)
	if err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "分类不存在")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 转换响应格式
	return &pb.Category{
		Id:        int64(category.ID),
		ParentId:  int64(category.ParentID),
		Level:     int32(category.Level),
		Path:      category.Path,
		Name:      category.Name,
		SortOrder: int32(category.SortOrder),
		IsLeaf:    category.IsLeaf,
		CreatedAt: timestamppb.New(category.CreatedAt),
		UpdatedAt: timestamppb.New(category.UpdatedAt),
	}, nil
}

// UpdateCategory 更新分类名称
func (s *CategoryServiceService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*emptypb.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "分类名称不能为空")
	}

	if err := s.uc.UpdateCategoryName(ctx, &biz.Category{
		ID:   req.Id,
		Name: req.Name,
	}); err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "分类不存在")
		}
		if errors.Is(err, biz.ErrCategoryNameConflict) {
			return nil, status.Error(codes.AlreadyExists, "分类名称已存在")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// DeleteCategory 删除分类
func (s *CategoryServiceService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}

	if err := s.uc.DeleteCategory(ctx, req.Id); err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "分类不存在")
		}
		if errors.Is(err, biz.ErrCategoryHasChildren) {
			return nil, status.Error(codes.FailedPrecondition, "存在子分类不可删除")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// GetSubTree 获取子树
func (s *CategoryServiceService) GetSubTree(ctx context.Context, req *pb.GetSubTreeRequest) (*pb.Categories, error) {
	if req.RootId == 0 {
		return nil, status.Error(codes.InvalidArgument, "根分类ID不能为空")
	}

	categories, err := s.uc.GetSubTree(ctx, req.RootId)
	if err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "根分类不存在")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbCategories := make([]*pb.Category, 0, len(categories))
	for _, c := range categories {
		pbCategories = append(pbCategories, &pb.Category{
			Id:        int64(c.ID),
			ParentId:  int64(c.ParentID),
			Level:     int32(c.Level),
			Path:      c.Path,
			Name:      c.Name,
			SortOrder: int32(c.SortOrder),
			IsLeaf:    c.IsLeaf,
			CreatedAt: timestamppb.New(c.CreatedAt),
			UpdatedAt: timestamppb.New(c.UpdatedAt),
		})
	}

	return &pb.Categories{Categories: pbCategories}, nil
}

// GetCategoryPath 获取分类路径
func (s *CategoryServiceService) GetCategoryPath(ctx context.Context, req *pb.GetCategoryPathRequest) (*pb.Categories, error) {
	if req.CategoryId == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}

	pathCategories, err := s.uc.GetCategoryPath(ctx, req.CategoryId)
	if err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "分类不存在")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbCategories := make([]*pb.Category, 0, len(pathCategories))
	for _, c := range pathCategories {
		pbCategories = append(pbCategories, &pb.Category{
			Id:        int64(c.ID),
			ParentId:  int64(c.ParentID),
			Level:     int32(c.Level),
			Path:      c.Path,
			Name:      c.Name,
			SortOrder: int32(c.SortOrder),
			IsLeaf:    c.IsLeaf,
			CreatedAt: timestamppb.New(c.CreatedAt),
			UpdatedAt: timestamppb.New(c.UpdatedAt),
		})
	}

	return &pb.Categories{Categories: pbCategories}, nil
}

// GetLeafCategories 获取所有叶子分类
func (s *CategoryServiceService) GetLeafCategories(ctx context.Context, _ *emptypb.Empty) (*pb.Categories, error) {
	leafCategories, err := s.uc.GetLeafCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbCategories := make([]*pb.Category, len(leafCategories))
	for _, c := range leafCategories {
		pbCategories = append(pbCategories, &pb.Category{
			Id:        int64(c.ID),
			ParentId:  int64(c.ParentID),
			Level:     int32(c.Level),
			Path:      c.Path,
			Name:      c.Name,
			SortOrder: int32(c.SortOrder),
			IsLeaf:    c.IsLeaf,
			CreatedAt: timestamppb.New(c.CreatedAt),
			UpdatedAt: timestamppb.New(c.UpdatedAt),
		})
	}

	return &pb.Categories{Categories: pbCategories}, nil
}

// GetClosureRelations 获取闭包关系
func (s *CategoryServiceService) GetClosureRelations(ctx context.Context, req *pb.GetClosureRequest) (*pb.ClosureRelations, error) {
	if req.CategoryId == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}

	relations, err := s.uc.GetClosureRelations(ctx, req.CategoryId)
	if err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "分类不存在")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbRelations := make([]*pb.ClosureRelation, 0, len(relations))
	for _, r := range relations {
		pbRelations = append(pbRelations, &pb.ClosureRelation{
			Ancestor:   int64(r.Ancestor),
			Descendant: int64(r.Descendant),
			Depth:      int32(r.Depth),
		})
	}

	return &pb.ClosureRelations{Relations: pbRelations}, nil
}

// UpdateClosureDepth 更新闭包关系深度
func (s *CategoryServiceService) UpdateClosureDepth(ctx context.Context, req *pb.UpdateClosureDepthRequest) (*emptypb.Empty, error) {
	if req.CategoryId == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "深度变化值不能为0")
	}

	if err := s.uc.UpdateClosureDepth(ctx, &biz.UpdateClosureDepth{
		ID:   req.CategoryId,
		Name: req.Name,
	}); err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "分类不存在")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
