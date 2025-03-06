package pkg

import (
	"context"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/google/uuid"
)

const UserId = "x-md-global-user-id"

// GetMetadataUesrID 从网关获取用户ID
func GetMetadataUesrID(ctx context.Context) (uuid.UUID, error) {
	var userIdStr string
	if md, ok := metadata.FromServerContext(ctx); ok {
		userIdStr = md.Get(UserId)
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}
