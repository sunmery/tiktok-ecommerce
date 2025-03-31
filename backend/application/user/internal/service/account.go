package service

import (
	"context"

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
