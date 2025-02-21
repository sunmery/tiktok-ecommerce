package service

import (
	"backend/application/category/internal/biz"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"

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
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 转换响应格式
	return &pb.Category{
		Id:        category.ID,
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
		Id:        category.ID,
		ParentId:  category.ParentID,
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
// 接口文档：PUT /v1/category/{id}
func (s *CategoryServiceService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*emptypb.Empty, error) {
	// 参数校验
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "分类名称不能为空")
	}

	// 构造业务对象
	category := &biz.Category{
		ID:   req.Id,
		Name: req.Name,
	}

	// 调用业务逻辑层
	if err := s.uc.UpdateCategoryName(ctx, category); err != nil {
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
// 接口文档：DELETE /v1/category/{id}
func (s *CategoryServiceService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	// 参数校验
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}

	// 调用业务逻辑层
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
// 接口文档：GET /v1/category/{root_id}/subtree
func (s *CategoryServiceService) GetSubTree(req *pb.GetSubTreeRequest, stream pb.CategoryService_GetSubTreeServer) error {
	// 参数校验
	if req.RootId == 0 {
		return status.Error(codes.InvalidArgument, "根分类ID不能为空")
	}

	// 获取子树数据
	categories, err := s.uc.GetSubTree(stream.Context(), req.RootId)
	if err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return status.Error(codes.NotFound, "根分类不存在")
		}
		return status.Error(codes.Internal, err.Error())
	}

	// 流式返回结果
	for _, category := range categories {
		if err := stream.Send(&pb.Category{
			Id:        category.ID,
			ParentId:  category.ParentID,
			Level:     int32(category.Level),
			Path:      category.Path,
			Name:      category.Name,
			SortOrder: int32(category.SortOrder),
			IsLeaf:    category.IsLeaf,
			CreatedAt: timestamppb.New(category.CreatedAt),
			UpdatedAt: timestamppb.New(category.UpdatedAt),
		}); err != nil {
			// 处理客户端断开连接的情况
			if errors.Is(err, io.EOF) {
				return nil
			}
			return status.Error(codes.Internal, "流式传输中断")
		}
	}
	return nil
}

// GetCategoryPath 获取分类路径
// 接口文档：GET /v1/category/{category_id}/path
func (s *CategoryServiceService) GetCategoryPath(req *pb.GetCategoryPathRequest, stream grpc.ServerStreamingServer[pb.Category]) error {
	// 参数校验
	if req.CategoryId == 0 {
		return status.Error(codes.InvalidArgument, "分类ID不能为空")
	}

	// 获取路径数据
	pathCategories, err := s.uc.GetCategoryPath(stream.Context(), req.CategoryId)
	if err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return status.Error(codes.NotFound, "分类不存在")
		}
		return status.Error(codes.Internal, err.Error())
	}

	// 流式返回路径节点（从根到当前）
	for _, category := range pathCategories {
		if err := stream.Send(&pb.Category{
			Id:        category.ID,
			ParentId:  category.ParentID,
			Level:     int32(category.Level),
			Path:      category.Path,
			Name:      category.Name,
			SortOrder: int32(category.SortOrder),
			IsLeaf:    category.IsLeaf,
			CreatedAt: timestamppb.New(category.CreatedAt),
			UpdatedAt: timestamppb.New(category.UpdatedAt),
		}); err != nil {
			return status.Error(codes.Internal, "流式传输中断")
		}
	}
	return nil
}

// GetLeafCategories 获取所有叶子分类
// 接口文档：GET /v1/category/leaves
// func (s *CategoryServiceService) GetLeafCategories(req *pb.GetLeafCategoriesRequest, stream grpc.ServerStreamingServer[pb.Category]) error {
func (s *CategoryServiceService) GetLeafCategories(ctx context.Context, _ *emptypb.Empty) (*pb.Categorys, error) {

	// 获取叶子分类
	leafCategories, err := s.uc.GetLeafCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fmt.Printf("GetLeafCategories res: %+v", leafCategories)
	var categories = make([]*pb.Category, 0, len(leafCategories))
	// 流式返回结果
	for _, category := range leafCategories {
		categories = append(categories, &pb.Category{
			Id:        category.ID,
			ParentId:  category.ParentID,
			Level:     int32(category.Level),
			Path:      category.Path,
			Name:      category.Name,
			SortOrder: int32(category.SortOrder),
			IsLeaf:    category.IsLeaf,
			CreatedAt: timestamppb.New(category.CreatedAt),
			UpdatedAt: timestamppb.New(category.UpdatedAt),
		})
	}

	return &pb.Categorys{
		Categorys: categories,
	}, nil
}

// GetClosureRelations 获取闭包关系
// 接口文档：GET /v1/category/{category_id}/closure
func (s *CategoryServiceService) GetClosureRelations(req *pb.GetClosureRequest, stream pb.CategoryService_GetClosureRelationsServer) error {
	// 参数校验
	if req.CategoryId == 0 {
		return status.Error(codes.InvalidArgument, "分类ID不能为空")
	}

	// 获取闭包关系
	relations, err := s.uc.GetClosureRelations(stream.Context(), req.CategoryId)
	if err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return status.Error(codes.NotFound, "分类不存在")
		}
		return status.Error(codes.Internal, err.Error())
	}

	// 流式返回闭包关系
	for _, rel := range relations {
		if err := stream.Send(&pb.ClosureRelation{
			Ancestor:   rel.Ancestor,
			Descendant: rel.Descendant,
			Depth:      rel.Depth,
		}); err != nil {
			return status.Error(codes.Internal, "流式传输中断")
		}
	}
	return nil
}

// UpdateClosureDepth 更新闭包关系深度
// 接口文档：PATCH /v1/category/{category_id}/closure
func (s *CategoryServiceService) UpdateClosureDepth(ctx context.Context, req *pb.UpdateClosureDepthRequest) (*emptypb.Empty, error) {
	// 参数校验
	if req.CategoryId == 0 {
		return nil, status.Error(codes.InvalidArgument, "分类ID不能为空")
	}
	if req.Delta == 0 {
		return nil, status.Error(codes.InvalidArgument, "深度变化值不能为0")
	}

	// 调用业务逻辑层
	if err := s.uc.UpdateClosureDepth(ctx, req.CategoryId, int32(int(req.Delta))); err != nil {
		if errors.Is(err, biz.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "分类不存在")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
