package biz

type CartItem struct {
	ProductId uint32 `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type UpsertItemReq struct {
	UserId uint32   `json:"user_id"`
	Item   CartItem `json:"item"`
}

type UpsertItemResp struct {
	Success bool `json:"success"`
}

type EmptyCartReq struct {
	UserId uint32 `json:"user_id"`
}

type EmptyCartResp struct {
	Success bool `json:"success"`
}

type GetCartReq struct {
	UserId uint32 `json:"user_id"`
}

type GetCartResp struct {
	Cart Cart `json:"cart"`
}

type Cart struct {
	UserId uint32     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type RemoveCartItemReq struct {
	UserId    uint32 `json:"user_id"`
	ProductId uint32 `json:"product_id"`
}

type RemoveCartItemResp struct {
	Success bool `json:"success"`
}
