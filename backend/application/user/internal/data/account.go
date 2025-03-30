package data

import (
	"context"
	"fmt"
	"net/http"
	"sync"

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
				return fmt.Errorf("failed to get user %s: %w", user.Id, err)
			}
			if fullUserProfile == nil {
				return fmt.Errorf("user %s not found", user.Id)
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
		CreatedTime:       user.CreatedTime,
		UpdatedTime:       user.UpdatedTime,
		Owner:             user.Owner,
		SignupApplication: user.SignupApplication,
		Name:              user.Name,
		Email:             user.Email,
		Avatar:            user.Avatar,
		DisplayName:       user.DisplayName,
	}, nil
}
