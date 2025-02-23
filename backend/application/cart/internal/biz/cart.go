package biz

type CartItem struct {
	MerchantId string `json:"merchant_id"` // 新增字段，表示商家ID
	ProductId  uint32 `json:"product_id"`  // 商品ID
	Quantity   int32  `json:"quantity"`    // 商品数量
	Selected   bool   `json:"selected"`    // 是否选中
}

type UpsertItemReq struct {
	UserId string   `json:"user_id"` // 新增字段，表示用户ID
	Item   CartItem `json:"item"`
}

type UpsertItemResp struct {
	Success bool `json:"success"`
}

type EmptyCartReq struct {
	UserId string `json:"user_id"` // 新增字段，表示用户ID
}

type EmptyCartResp struct {
	Success bool `json:"success"`
}

type GetCartReq struct {
	UserId string `json:"user_id"` // 新增字段，表示用户ID
}

type GetCartResp struct {
	Cart Cart `json:"cart"`
}

type Cart struct {
	UserId string     `json:"user_id"` // 新增字段，表示用户ID
	Items  []CartItem `json:"items"`   // 购物车商品列表
}

type RemoveCartItemReq struct {
	UserId     string `json:"user_id"`     // 新增字段，表示用户ID
	MerchantId string `json:"merchant_id"` // 新增字段，表示商家ID
	ProductId  uint32 `json:"product_id"`
}

type RemoveCartItemResp struct {
	Success bool `json:"success"`
}

type CheckCartItemReq struct {
	UserId     string `json:"user_id"`     // 新增字段，表示用户ID
	MerchantId string `json:"merchant_id"` // 新增字段，表示商家ID
	ProductId  uint32 `json:"product_id"`
}

type CheckCartItemResp struct {
	Success bool `json:"success"`
}

type UncheckCartItemReq struct {
	UserId     string `json:"user_id"`     // 新增字段，表示用户ID
	MerchantId string `json:"merchant_id"` // 新增字段，表示商家ID
	ProductId  uint32 `json:"product_id"`
}

type UncheckCartItemResp struct {
	Success bool `json:"success"`
}

type CreateOrderReq struct {
	UserId string `json:"user_id"` // 新增字段，表示用户ID
}

type CreateOrderResp struct {
	Success bool       `json:"success"`
	Items   []CartItem `json:"items"` // 返回被选中的商品
}

type CreateCartReq struct {
	UserId   string `json:"user_id"`   // 新增字段，表示用户ID
	CartName string `json:"cart_name"` // 购物车名称
}

type CreateCartResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"` // 购物车创建反馈信息
}

type ListCartsReq struct {
	UserId string `json:"user_id"` // 新增字段，表示用户ID
}

type CartSummary struct {
	CartId   uint32 `json:"cart_id"`   // 购物车ID
	CartName string `json:"cart_name"` // 购物车名称
}

type ListCartsResp struct {
	Carts []CartSummary `json:"carts"` // 返回购物车列表
}
