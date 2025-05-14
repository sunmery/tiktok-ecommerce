package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// 用户档案
type (
	GetProfileRequest struct {
		Owner  string
		UserId uuid.UUID
	}
	GetProfileReply struct {
		Id                uuid.UUID // 用户唯一标识
		Role              string
		IsDeleted         bool
		Owner             string // 用户所在的组织
		SignupApplication string // 用户所在的应用
		Name              string // 用户姓名
		Email             string // 用户邮箱
		Avatar            string // 用户头像 string
		CreatedTime       string
		UpdatedTime       string
		DeletedTime       time.Time
		DisplayName       string // 显示的用户名(可作为昵称)
		Phone             string // 手机号
	}
)

type (
	GetUsersRequest struct {
		AdminId uuid.UUID
	}

	GetUsersReply struct {
		Users []*GetProfileReply
	}
)

type (
	DeleteUserRequest struct {
		Owner  string
		UserId uuid.UUID
		Name   string
	}
	DeleteUserReply struct {
		Code   int32
		Status string
	}
)

type (
	UpdateUserRequest struct {
		Owner             string
		UserId            uuid.UUID
		Name              string
		Email             string
		Avatar            string
		DisplayName       string
		SignupApplication string
	}
	UpdateUserReply struct {
		Code   int32
		Status string
	}
)

// Product 商品领域模型
type (
	// ProductImage 商品图片信息
	ProductImage struct {
		URL       string
		IsPrimary bool
		SortOrder *int
	}
	// ProductStatus 商品状态
	ProductStatus int32
	// CategoryInfo 分类信息
	CategoryInfo struct {
		CategoryId   uint64
		CategoryName string
		SortOrder    int32
	}

	// Inventory 库存
	Inventory struct {
		ProductId  uuid.UUID
		MerchantId uuid.UUID
		Stock      uint32
	}

	Product struct {
		ID          uuid.UUID
		MerchantId  uuid.UUID
		Name        string
		Price       float64
		Description string
		Images      []*ProductImage
		Status      ProductStatus
		Category    CategoryInfo
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Attributes  map[string]any
		Inventory   Inventory // 库存
	}
)

func (cc *UserUsecase) GetProfile(ctx context.Context, req *GetProfileRequest) (*GetProfileReply, error) {
	cc.log.WithContext(ctx).Debugf("GetProfile: %+v", req)
	return cc.repo.GetProfile(ctx, req)
}

func (cc *UserUsecase) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersReply, error) {
	cc.log.WithContext(ctx).Debugf("GetUsers: %+v", req)
	return cc.repo.GetUsers(ctx, req)
}

func (cc *UserUsecase) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserReply, error) {
	cc.log.WithContext(ctx).Debugf("DeleteUser: %+v", req)
	return cc.repo.DeleteUser(ctx, req)
}

func (cc *UserUsecase) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserReply, error) {
	cc.log.WithContext(ctx).Debugf("UpdateUser: %+v", req)
	return cc.repo.UpdateUser(ctx, req)
}

func (cc *UserUsecase) GetFavorites(ctx context.Context, req *GetFavoritesRequest) (*Favorites, error) {
	cc.log.WithContext(ctx).Debugf("GetFavorites req:%+v", req)
	return cc.repo.GetFavorites(ctx, req)
}

func (cc *UserUsecase) DeleteFavorites(ctx context.Context, req *UpdateFavoritesRequest) (*UpdateFavoritesResply, error) {
	cc.log.WithContext(ctx).Debugf("DeleteFavorites req:%+v", req)
	return cc.repo.DeleteFavorites(ctx, req)
}

func (cc *UserUsecase) SetFavorites(ctx context.Context, req *UpdateFavoritesRequest) (*UpdateFavoritesResply, error) {
	cc.log.WithContext(ctx).Debugf("SetFavorites req:%+v", req)
	return cc.repo.SetFavorites(ctx, req)
}
