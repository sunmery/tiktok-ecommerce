package biz

import "github.com/google/uuid"

type CartItem struct {
	MerchantId uuid.UUID
	ProductId  uuid.UUID // 商品ID
	Quantity   uint32
	Price      float64
	Name       string
	Picture    string
}

type UpsertItemReq struct {
	UserId     uuid.UUID
	MerchantId uuid.UUID
	ProductId  uuid.UUID
	Quantity   uint32
}

type UpsertItemResp struct {
	Success bool
}

type EmptyCartReq struct {
	UserId uuid.UUID
}

type EmptyCartResp struct {
	Success bool
}

type GetCartReq struct {
	UserId uuid.UUID
}

type GetCartResp struct {
	Cart Cart `json:"cart"`
}

type Cart struct {
	UserId string
	Items  []CartItem
}

type RemoveCartItemReq struct {
	UserId     uuid.UUID
	MerchantId uuid.UUID
	ProductId  uuid.UUID
}

type RemoveCartItemResp struct {
	Success bool
}

type CheckCartItemReq struct {
	UserId     uuid.UUID
	MerchantId uuid.UUID
	ProductId  uuid.UUID
}

type CheckCartItemResp struct {
	Success bool
}

type UncheckCartItemReq struct {
	UserId     uuid.UUID
	MerchantId uuid.UUID
	ProductId  uuid.UUID
}

type UncheckCartItemResp struct {
	Success bool
}

type CreateOrderReq struct {
	UserId uuid.UUID
}

type CreateOrderResp struct {
	Success bool
	Items   []CartItem
}

type ListCartsReq struct {
	UserId uuid.UUID
}
