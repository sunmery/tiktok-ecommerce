package biz

import "github.com/google/uuid"

// 库存
type (
	Inventory struct {
		ProductId  uuid.UUID
		MerchantId uuid.UUID
		Stock      uint32
	}
	UpdateInventoryRequest struct {
		ProductId  uuid.UUID
		MerchantId uuid.UUID
		Stock      int32
	}
	UpdateInventoryReply struct {
		ProductId  uuid.UUID
		MerchantId uuid.UUID
		Stock      uint32
	}
)
