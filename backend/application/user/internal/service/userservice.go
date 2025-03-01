package service

import (
	"backend/application/user/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "backend/api/user/v1"
)

func (s *UserService) GetUserProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var (
		owner string
		id    string
	)
	if md, ok := metadata.FromServerContext(ctx); ok {
		owner = md.Get("x-md-global-owner")
		id = md.Get("x-md-global-user-id")
	}
	userId, err := uuid.Parse(id)
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

	return &pb.GetProfileResponse{
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
