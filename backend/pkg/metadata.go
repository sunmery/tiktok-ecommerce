package pkg

import (
	"context"
	"errors"

	"backend/constants"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/google/uuid"
)

// GetMetadataUesrID 从网关获取用户ID
func GetMetadataUesrID(ctx context.Context) (uuid.UUID, error) {
	var userIdStr string
	if md, ok := metadata.FromServerContext(ctx); ok {
		userIdStr = md.Get(constants.UserId)
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil || userId == uuid.Nil {
		return uuid.Nil, errors.New("无效的用户ID")
	}

	return userId, nil
}

// GetMetadataOwner 从网关获取用户组织
func GetMetadataOwner(ctx context.Context) (string, error) {
	var owner string
	if md, ok := metadata.FromServerContext(ctx); ok {
		owner = md.Get(constants.Owner)
	}
	return owner, nil
}
