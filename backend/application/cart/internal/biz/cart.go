package biz

import "github.com/google/uuid"

type CartItem struct {
	MerchantId uuid.UUID `json:"merchant_id"` // 新增字段，表示商家ID
	ProductId  uuid.UUID `json:"product_id"`  // 商品ID
	Quantity   int32  `json:"quantity"`    // 商品数量
}

type UpsertItemReq struct {
	UserId uuid.UUID   `json:"user_id"` // 新增字段，表示用户ID
	Item   CartItem `json:"item"`
}

type UpsertItemResp struct {
	Success bool `json:"success"`
}

type EmptyCartReq struct {
	UserId uuid.UUID `json:"user_id"` // 新增字段，表示用户ID
}

type EmptyCartResp struct {
	Success bool `json:"success"`
}

type GetCartReq struct {
	UserId uuid.UUID `json:"user_id"` // 新增字段，表示用户ID
}

type GetCartResp struct {
	Cart Cart `json:"cart"`
}

type Cart struct {
	UserId string     `json:"user_id"` // 新增字段，表示用户ID
	Items  []CartItem `json:"items"`   // 购物车商品列表
}

type RemoveCartItemReq struct {
	UserId     uuid.UUID `json:"user_id"`     // 新增字段，表示用户ID
	MerchantId uuid.UUID `json:"merchant_id"` // 新增字段，表示商家ID
	ProductId  uuid.UUID `json:"product_id"`
}

type RemoveCartItemResp struct {
	Success bool `json:"success"`
}

type CheckCartItemReq struct {
	UserId     uuid.UUID `json:"user_id"`     // 新增字段，表示用户ID
	MerchantId uuid.UUID `json:"merchant_id"` // 新增字段，表示商家ID
	ProductId  uuid.UUID `json:"product_id"`
}

type CheckCartItemResp struct {
	Success bool `json:"success"`
}

type UncheckCartItemReq struct {
	UserId     uuid.UUID `json:"user_id"`     // 新增字段，表示用户ID
	MerchantId uuid.UUID `json:"merchant_id"` // 新增字段，表示商家ID
	ProductId  uuid.UUID `json:"product_id"`
}

type UncheckCartItemResp struct {
	Success bool `json:"success"`
}

type CreateOrderReq struct {
	UserId uuid.UUID `json:"user_id"` // 新增字段，表示用户ID
}

type CreateOrderResp struct {
	Success bool       `json:"success"`
	Items   []CartItem `json:"items"` // 返回被选中的商品
}

type ListCartsReq struct {
	UserId uuid.UUID `json:"user_id"` // 新增字段，表示用户ID
}
