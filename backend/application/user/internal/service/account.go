package service

import (
	"backend/application/user/internal/biz"
	"backend/pkg"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "backend/api/user/v1"
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
		CreatedTime: timestamppb.New(profile.CreatedTime),
		UpdatedTime: timestamppb.New(profile.UpdatedTime),
	}, nil
}
