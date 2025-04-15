package service

import (
	"context"

	pb "backend/api/product/v1"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"

	v1 "backend/api/user/v1"
	"backend/application/user/internal/biz"
	"backend/pkg"
)

func (s *UserService) GetUserProfile(ctx context.Context, _ *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	owner, err := pkg.GetMetadataOwner(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := s.uc.GetProfile(ctx, &biz.GetProfileRequest{
		Owner:  owner,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponse{
		Owner:             profile.Owner,
		Name:              profile.Name,
		Avatar:            profile.Avatar,
		Email:             profile.Email,
		Id:                profile.Id.String(),
		Role:              profile.Role,
		DisplayName:       profile.DisplayName,
		IsDeleted:         profile.IsDeleted,
		CreatedTime:       profile.CreatedTime,
		UpdatedTime:       profile.UpdatedTime,
		SignupApplication: profile.SignupApplication,
	}, nil
}

func (s *UserService) GetUsers(ctx context.Context, _ *v1.GetUsersRequest) (*v1.GetUsersResponse, error) {
	adminId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}

	log.Debugf("get adminId: %v", adminId)
	users, getUsersErr := s.uc.GetUsers(ctx, &biz.GetUsersRequest{
		AdminId: adminId,
	})
	if getUsersErr != nil {
		return nil, getUsersErr
	}

	resp := make([]*v1.GetProfileResponse, 0, len(users.Users))
	if resp == nil {
		return &v1.GetUsersResponse{}, nil
	}

	for _, user := range users.Users {
		resp = append(resp, &v1.GetProfileResponse{
			Owner:             user.Owner,
			Name:              user.Name,
			Avatar:            user.Avatar,
			Email:             user.Email,
			Id:                user.Id.String(),
			Role:              user.Role,
			DisplayName:       user.DisplayName,
			IsDeleted:         user.IsDeleted,
			CreatedTime:       user.CreatedTime,
			UpdatedTime:       user.UpdatedTime,
			SignupApplication: user.SignupApplication,
		})
	}

	return &v1.GetUsersResponse{
		Users: resp,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	result, deleteUsersErr := s.uc.DeleteUser(ctx, &biz.DeleteUserRequest{
		Owner:  req.Owner,
		UserId: userId,
		Name:   req.Name,
	})
	if deleteUsersErr != nil {
		return nil, deleteUsersErr
	}

	return &v1.DeleteUserResponse{
		Status: result.Status,
		Code:   result.Code,
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	result, err := s.uc.UpdateUser(ctx, &biz.UpdateUserRequest{
		Owner:             req.Owner,
		UserId:            userId,
		Name:              req.Name,
		Email:             req.Email,
		Avatar:            req.Avatar,
		DisplayName:       req.DisplayName,
		SignupApplication: req.SignupApplication,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateUserResponse{
		Status: result.Status,
		Code:   result.Code,
	}, nil
}

// GetFavorites 获取用户的全部收藏
func (s *UserService) GetFavorites(ctx context.Context, req *v1.GetFavoritesRequest) (*v1.Favorites, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	favorites, err := s.uc.GetFavorites(ctx, &biz.GetFavoritesRequest{
		UserId:   userId,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	if favorites == nil {
		return &v1.Favorites{
			Items: []*pb.Product{},
		}, nil
	}
	if favorites.Items == nil {
		return &v1.Favorites{
			Items: []*pb.Product{},
		}, nil
	}
	var pbProducts []*pb.Product
	for _, product := range favorites.Items {
		if product != nil {
			pbProducts = append(pbProducts, convertBizProductToPB(product))
		}
	}

	return &v1.Favorites{
		Items: pbProducts,
	}, nil
}

// DeleteFavorites 删除收藏
func (s *UserService) DeleteFavorites(ctx context.Context, req *v1.UpdateFavoritesRequest) (*v1.UpdateFavoritesResply, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, err
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, err
	}
	resply, err := s.uc.DeleteFavorites(ctx, &biz.UpdateFavoritesRequest{
		UserId:     userId,
		ProductId:  productId,
		MerchantId: merchantId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateFavoritesResply{
		Message: resply.Message,
		Code:    resply.Code,
	}, nil
}

// SetFavorites 添加收藏
func (s *UserService) SetFavorites(ctx context.Context, req *v1.UpdateFavoritesRequest) (*v1.UpdateFavoritesResply, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, err
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, err
	}
	resply, err := s.uc.SetFavorites(ctx, &biz.UpdateFavoritesRequest{
		UserId:     userId,
		ProductId:  productId,
		MerchantId: merchantId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateFavoritesResply{
		Message: resply.Message,
		Code:    resply.Code,
	}, nil
}

func convertBizProductToPB(p *biz.Product) *pb.Product {
	if p == nil {
		return nil
	}

	// 转换图片
	images := make([]*pb.Image, 0)
	if p.Images != nil {
		for _, img := range p.Images {
			if img != nil {
				sortOrder := int32(0)
				if img.SortOrder != nil {
					sortOrder = int32(*img.SortOrder)
				}
				images = append(images, &pb.Image{
					Url:       img.URL,
					IsPrimary: img.IsPrimary,
					SortOrder: sortOrder,
				})
			}
		}
	}

	// 转换商品属性
	var attributes *structpb.Value
	if p.Attributes != nil && len(p.Attributes) > 0 {
		protoStruct, err := structpb.NewStruct(p.Attributes)
		if err != nil {
			log.Warn("Error creating struct: %w", err)
			attributes = nil
		} else {
			attributes = structpb.NewStructValue(protoStruct)
		}
	}

	// 构建返回结果
	result := &pb.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Status:      uint32(p.Status),
		MerchantId:  p.MerchantId.String(),
		Images:      images,
		Attributes:  attributes,
		Category: &pb.CategoryInfo{
			CategoryId:   uint32(p.Category.CategoryId),
			CategoryName: p.Category.CategoryName,
		},
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}

	// 添加库存信息
	if p.Inventory.ProductId != uuid.Nil {
		result.Inventory = &pb.Inventory{
			ProductId:  p.Inventory.ProductId.String(),
			MerchantId: p.Inventory.MerchantId.String(),
			Stock:      p.Inventory.Stock,
		}
	}

	return result
}
