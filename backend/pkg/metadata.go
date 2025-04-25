package pkg

import (
	"context"
	"errors"
	"fmt"

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
		return uuid.Nil, errors.New(fmt.Sprintf("invalid user id: '%s' error: %v", userIdStr, err))
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

// GetMetadataRole 从网关获取用户角色
func GetMetadataRole(ctx context.Context) (constants.RoleType, error) {
	var role string
	if md, ok := metadata.FromServerContext(ctx); ok {
		role = md.Get(constants.Role)
	}
	return getRole(role), nil
}

func getRole(role string) constants.RoleType {
	switch role {
	case "consumer":
		return constants.Consumer
	case "merchant":
		return constants.Merchant
	case "admin":
		return constants.Admin
	case "guest":
		return constants.Guest
	default:
		return constants.Guest
	}
}
