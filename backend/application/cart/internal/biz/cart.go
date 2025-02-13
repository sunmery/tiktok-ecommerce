package biz

type CartItem struct {
	ProductId uint32 `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type UpsertItemReq struct {
	Owner string   `json:"owner"`
	Name  string   `json:"name"`
	Item  CartItem `json:"item"`
}

type UpsertItemResp struct {
	Success bool `json:"success"`
}

type EmptyCartReq struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type EmptyCartResp struct {
	Success bool `json:"success"`
}

type GetCartReq struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type GetCartResp struct {
	Cart Cart `json:"cart"`
}

type Cart struct {
	Owner string     `json:"owner"`
	Name  string     `json:"name"`
	Items []CartItem `json:"items"`
}

type RemoveCartItemReq struct {
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	ProductId uint32 `json:"product_id"`
}

type RemoveCartItemResp struct {
	Success bool `json:"success"`
}
