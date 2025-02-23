package biz

import "github.com/google/uuid"

type CartItem struct {
	ProductId int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type UpsertItemReq struct {
	UserId uuid.UUID
	Item   CartItem `json:"item"`
}

type UpsertItemResp struct {
	Success bool `json:"success"`
}

type EmptyCartReq struct {
	UserId uuid.UUID
}

type EmptyCartResp struct {
	Success bool `json:"success"`
}

type GetCartReq struct {
	UserId uuid.UUID
}

type GetCartResp struct {
	Cart Cart `json:"cart"`
}

type Cart struct {
	UserId uuid.UUID

	Items []CartItem `json:"items"`
}

type RemoveCartItemReq struct {
	UserId    uuid.UUID
	ProductId uint32 `json:"product_id"`
}

type RemoveCartItemResp struct {
	Success bool `json:"success"`
}
