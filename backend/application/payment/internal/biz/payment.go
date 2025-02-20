package biz

import "context"

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type CreateRequest struct {
}

type CreateReply struct {
}

type PaymentRepo interface {
	Create(context.Context, *CreateRequest) (*CreateReply, error)
}

func (cc *PaymentUsecase) Create(ctx context.Context, req *CreateRequest) (*CreateReply, error) {
	cc.log.WithContext(ctx).Debugf("Create request: %+v", req)
	return cc.repo.Create(ctx, req)
}
