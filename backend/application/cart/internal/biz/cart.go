package biz

type CartItem struct {
	ProductId uint32 `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type AddItemReq struct {
	UserId uint32   `json:"user_id"`
	Item   CartItem `json:"item"`
}

type AddItemResp struct{}

type EmptyCartReq struct {
	UserId uint32 `json:"user_id"`
}

type EmptyCartResp struct{}

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

type UpdateItemReq struct {
	UserId uint32   `json:"user_id"`
	Item   CartItem `json:"item"`
}

type UpdateItemResp struct{}

type RemoveItemReq struct {
	UserId    uint32 `json:"user_id"`
	ProductId uint32 `json:"product_id"`
}

type RemoveItemResp struct{}
