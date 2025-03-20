package biz

import "github.com/google/uuid"

type (
	UpsertItemReq struct {
		UserId     uuid.UUID
		MerchantId uuid.UUID
		ProductId  uuid.UUID
		Quantity   uint32
	}
	UpsertItemResp struct {
		Success bool
	}
)

type (
	EmptyCartReq struct {
		UserId uuid.UUID
	}
	EmptyCartResp struct {
		Success bool
	}
)

type (
	GetCartReq struct {
		UserId uuid.UUID
	}
	CartInfo struct {
		MerchantId uuid.UUID
		ProductId  uuid.UUID
		Quantity   uint32
		Price      float64
		Name       string
		Picture    string
	}
	GetCartRelpy struct {
		Items []*CartInfo
	}
)

type (
	RemoveCartItemReq struct {
		UserId     uuid.UUID
		MerchantId uuid.UUID
		ProductId  uuid.UUID
	}
	RemoveCartItemResp struct {
		Success bool
	}
)
