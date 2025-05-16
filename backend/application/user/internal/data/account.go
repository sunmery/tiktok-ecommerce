package data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	categoryv1 "backend/api/category/v1"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-kratos/kratos/v2/log"

	"backend/pkg/types"

	"backend/application/user/internal/data/models"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"

	"golang.org/x/sync/errgroup"

	"github.com/google/uuid"

	"backend/application/user/internal/biz"
)

// GetProfile 获取用户档案
func (u *userRepo) GetProfile(ctx context.Context, req *biz.GetProfileRequest) (*biz.GetProfileReply, error) {
	user, err := u.data.cs.GetUser(req.UserId.String())
	if err != nil {
		return nil, err
	}

	// 用户是否被注销
	if user.IsDeleted {
		return nil, fmt.Errorf(fmt.Sprintf("user %s is deleted", user.Name))
	}

	// 组装数据
	log.Debugf("Phone%+v", user.Phone)
	return &biz.GetProfileReply{
		Id:                req.UserId,
		Role:              user.Roles[0].Name,
		IsDeleted:         false,
		CreatedTime:       user.CreatedTime,
		UpdatedTime:       user.UpdatedTime,
		Owner:             user.Owner,
		SignupApplication: user.SignupApplication,
		Name:              user.Name,
		Email:             user.Email,
		Avatar:            user.Avatar,
		// DeletedTime:        user.DeletedTime,
		DisplayName: user.DisplayName,
		Phone:       user.Phone,
	}, nil
}

func (u *userRepo) GetUsers(ctx context.Context, _ *biz.GetUsersRequest) (*biz.GetUsersReply, error) {
	users, err := u.data.cs.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var (
		g    errgroup.Group
		mu   sync.Mutex
		resp = make([]*biz.GetProfileReply, 0, len(users))
	)
	if resp == nil {
		return &biz.GetUsersReply{}, nil
	}
	// if resp == nil {
	// 	return &biz.GetUsersReply{}, nil
	// }

	// 并发控制
	g.SetLimit(10)

	for _, user := range users {
		g.Go(func() error {
			fullUserProfile, err := u.data.cs.GetUser(user.Id)
			if err != nil {
				return fmt.Errorf("failed to get user '%s': %w", user.Id, err)
			}
			if fullUserProfile == nil {
				return fmt.Errorf("user '%s' not found", user.Id)
			}

			users, err := convertUserToProfile(fullUserProfile)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()

			resp = append(resp, users)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &biz.GetUsersReply{
		Users: resp,
	}, nil
}

func (u *userRepo) DeleteUser(ctx context.Context, req *biz.DeleteUserRequest) (*biz.DeleteUserReply, error) {
	ok, err := u.data.cs.DeleteUser(&casdoorsdk.User{
		Id:    req.UserId.String(),
		Owner: req.Owner,
		Name:  req.Name,
	})
	if err != nil || !ok {
		return nil, errors.New(500, "InternalServerError", "delete user failed")
	}

	return &biz.DeleteUserReply{
		Status: "ok",
		Code:   http.StatusOK,
	}, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, req *biz.UpdateUserRequest) (*biz.UpdateUserReply, error) {
	ok, err := u.data.cs.UpdateUser(&casdoorsdk.User{
		Id:                req.UserId.String(),
		Owner:             req.Owner,
		Name:              req.Name,
		Email:             req.Email,
		Avatar:            req.Avatar,
		DisplayName:       req.DisplayName,
		SignupApplication: req.SignupApplication,
	})
	if err != nil || !ok {
		return nil, errors.New(500, "InternalServerError", "delete user failed")
	}

	return &biz.UpdateUserReply{
		Status: "ok",
		Code:   http.StatusOK,
	}, nil
}

func (u *userRepo) GetFavorites(ctx context.Context, req *biz.GetFavoritesRequest) (*biz.Favorites, error) {
	userID := types.ToPgUUID(req.UserId)
	page := (req.Page - 1) * req.PageSize
	favorites, err := u.data.db.GetFavorites(ctx, models.GetFavoritesParams{
		UserID:   userID,
		Page:     &page,
		PageSize: &req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	log.Debug("asda:", len(favorites))

	// products := make([]*productv1.Product, 0, len(favorites))
	// 收集所有不同的分类ID
	categoryIDs := make([]int64, 0)
	categoryIDMap := make(map[int64]bool)
	for _, product := range favorites {
		if !categoryIDMap[product.CategoryID] {
			categoryIDMap[product.CategoryID] = true
			categoryIDs = append(categoryIDs, product.CategoryID)
		}
	}

	// 从分类服务获取分类信息
	var categoryMap map[int64]*categoryv1.Category
	if len(categoryIDs) > 0 {
		categoriesResp, err := u.data.categoryClient.BatchGetCategories(ctx, &categoryv1.BatchGetCategoriesRequest{
			Ids: categoryIDs,
		})
		if err != nil {
			u.log.WithContext(ctx).Warnf("failed to get categories: %v", err)
		} else {
			categoryMap = make(map[int64]*categoryv1.Category)
			for _, cat := range categoriesResp.Categories {
				categoryMap[cat.Id] = cat
			}
		}
	}

	items := make([]*biz.Product, 0, len(favorites))
	for _, product := range favorites {
		var images []*biz.ProductImage
		if len(product.Images) > 0 {
			if err := json.Unmarshal(product.Images, &images); err != nil {
				// 处理错误或记录日志
				u.log.WithContext(ctx).Warnf("unmarshal images error: %v", err)
			}
		}

		price, err := types.NumericToFloat(product.Price.(pgtype.Numeric))
		if err != nil {
			u.log.WithContext(ctx).Warnf("unmarshal price error: %v", err)
		}

		// 构建分类信息
		categoryInfo := biz.CategoryInfo{
			CategoryId: uint64(product.CategoryID),
		}

		// 如果找到了分类信息，则设置分类名称
		if c, ok := categoryMap[product.CategoryID]; ok {
			categoryInfo.CategoryName = c.Name
			categoryInfo.SortOrder = c.SortOrder
		}

		// 处理商品属性
		var attributes map[string]any
		if len(product.Attributes) > 0 {
			if err := json.Unmarshal(product.Attributes, &attributes); err != nil {
				u.log.WithContext(ctx).Warnf("unmarshal attributes error: %v", err)
				attributes = nil
			}
		} else {
			attributes = nil
		}

		log.Debugf("Status: %+v", product.Status)
		items = append(items, &biz.Product{
			ID:          product.ID,
			MerchantId:  product.MerchantID,
			Name:        product.Name,
			Price:       price,
			Description: *product.Description,
			Images:      images,
			Status:      biz.ProductStatus(product.Status),
			Category:    categoryInfo,
			// CreatedAt:   product.CreatedAt,
			// UpdatedAt:   product.UpdatedAt,
			Attributes: attributes,
			Inventory: biz.Inventory{
				ProductId:  product.ID,
				MerchantId: product.MerchantID,
				Stock:      uint32(product.Stock),
			},
		})
	}

	return &biz.Favorites{
		Items: items,
	}, nil
}

func (u *userRepo) DeleteFavorites(ctx context.Context, req *biz.UpdateFavoritesRequest) (*biz.UpdateFavoritesResply, error) {
	// 调用数据库删除收藏
	err := u.data.db.DeleteFavorites(ctx, models.DeleteFavoritesParams{
		UserID:     req.UserId,
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
	})
	if err != nil {
		return nil, err
	}

	return &biz.UpdateFavoritesResply{
		Message: "收藏删除成功",
		Code:    http.StatusOK,
	}, nil
}

func (u *userRepo) SetFavorites(ctx context.Context, req *biz.UpdateFavoritesRequest) (*biz.UpdateFavoritesResply, error) {
	// 调用数据库添加收藏
	err := u.data.db.SetFavorites(ctx, models.SetFavoritesParams{
		UserID:     req.UserId,
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
	})
	if err != nil {
		return nil, err
	}

	return &biz.UpdateFavoritesResply{
		Message: "收藏添加成功",
		Code:    http.StatusOK,
	}, nil
}

func convertUserToProfile(user *casdoorsdk.User) (*biz.GetProfileReply, error) {
	if user == nil {
		return nil, fmt.Errorf("nil user provided to convertUserToProfile")
	}
	userId, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %s", user.Id)
	}

	// 安全获取角色
	var role string
	if len(user.Roles) > 0 {
		role = user.Roles[0].Name
	} else {
		role = "guest" // 访客角色
		// return nil, errors.New("user has no role assigned")
	}
	return &biz.GetProfileReply{
		Id:                userId,
		Role:              role,
		IsDeleted:         user.IsDeleted,
		Owner:             user.Owner,
		SignupApplication: user.SignupApplication,
		Name:              user.Name,
		Email:             user.Email,
		Avatar:            user.Avatar,
		CreatedTime:       user.CreatedTime,
		UpdatedTime:       user.UpdatedTime,
		DisplayName:       user.DisplayName,
		Phone:             user.Phone,
	}, nil
}
