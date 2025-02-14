package service

import (
	pb "backend/api/product/v1"
	"backend/application/product/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
	"backend/pkg/token"
	"context"
)


func (s *ProductCatalogServiceService) CreateCategory(ctx context.Context, req *pb.CreateCategoryReq) (*pb.CategoryReply, error){
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}
	p, cErr := s.pu.CreateCategory(ctx, &biz.CreateCategoryRequest{
		Owner:       payload.Owner,
		Username:    payload.Name,
		Name:        req.Name,
		ParentID:    &req.ParentId,
	})
	if cErr != nil {
		return nil, cErr
	}
	return &pb.CategoryReply{
		Category: &pb.Category{
			Name:        p.Category.Name,
			ParentId:   int32(*p.Category.ParentID),
		},
	}, nil
}

func (s *ProductCatalogServiceService) ListCategories(ctx context.Context, _ *emptypb.Empty) (*pb.ListCategoriesResp, error) {
	p, cErr := s.pu.ListCategories(ctx, &emptypb.Empty{})
	if cErr != nil {
		return nil, cErr
	}

	var categories []*pb.Category
	for _, c := range p.Categories {
		category := &pb.Category{
			Name:     c.Name,
			ParentId: int32(*c.ParentID),
			Children: []*pb.Category{},
			Id:	   int32(c.ID),
		}
		categories = append(categories, category)
	}

	return &pb.ListCategoriesResp{
		Categories: categories,
	}, nil
}

func (s *ProductCatalogServiceService) GetCategoryChildren(ctx context.Context, req *pb.GetCategoryChildrenReq) (*pb.GetCategoryChildrenResp, error){
	p, cErr := s.pu.GetCategoryChildren(ctx, uint32(req.Id))
	if cErr != nil {
		return nil, cErr
	}

	var children []*pb.Category
	for _, c := range p.Categories {
		category := &pb.Category{
			Name:     c.Name,
			ParentId: int32(*c.ParentID),
			Children: []*pb.Category{},
			Id:       int32(c.ID),
		}
		children = append(children, category)
	}

	return &pb.GetCategoryChildrenResp{
		Categories: children,
	}, nil
}

