package service

import (
	"context"

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
		Owner:       profile.Owner,
		Name:        profile.Name,
		Avatar:      profile.Avatar,
		Email:       profile.Email,
		Id:          profile.Id.String(),
		Role:        profile.Role,
		DisplayName: profile.DisplayName,
		IsDeleted:   profile.IsDeleted,
		CreatedTime: profile.CreatedTime,
		UpdatedTime: profile.UpdatedTime,
	}, nil
}
