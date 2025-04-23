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
