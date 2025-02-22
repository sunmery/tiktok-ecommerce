package biz

type CartItem struct {
	ProductId uint32 `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Selected  bool   `json:"selected"` // 新增字段，表示商品是否被选中
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

type CheckCartItemReq struct {
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	ProductId uint32 `json:"product_id"`
}

type CheckCartItemResp struct {
	Success bool `json:"success"`
}

type UncheckCartItemReq struct {
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	ProductId uint32 `json:"product_id"`
}

type UncheckCartItemResp struct {
	Success bool `json:"success"`
}

type CreateOrderReq struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CreateOrderResp struct {
	Success bool       `json:"success"`
	Items   []CartItem `json:"items"`
}

type CreateCartReq struct {
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	CartName string `json:"cart_name"` // 新增字段，表示购物车名称
}

type CreateCartResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"` // 购物车创建反馈信息
}

type ListCartsReq struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CartSummary struct {
	CartId   uint32 `json:"cart_id"`   // 购物车ID
	CartName string `json:"cart_name"` // 购物车名称
}

type ListCartsResp struct {
	Carts []CartSummary `json:"carts"` // 返回购物车列表
}
