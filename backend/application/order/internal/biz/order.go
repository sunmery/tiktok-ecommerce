package biz

type PlaceOrderReq struct {
	OrderId      int32   // 订单ID
	Name         string  // 订单名称
	UserId       uint32  // 用户ID
	UserCurrency string  // 货币类型
	Address      Address // 地址信息
	Items        []Item  // 商品列表
	Email        string  // 用户邮箱
}

type Address struct {
	StreetAddress string // 街道地址
	City          string // 城市
	State         string // 州/省
	Country       string // 国家
	ZipCode       int32  // 邮政编码
}
type Item struct {
	Id          int32   // 商品ID
	Name        string  // 商品名称
	Description string  // 商品描述
	Price       float32 // 商品单价
	Quantity    int32   // 商品数量
}

type PlaceOrderResp struct {
	Order OrderResult // 订单结果
}

type OrderResult struct {
	OrderId int32 // 订单ID
}
type ListOrderReq struct {
	UserId uint32 `json:"user_id"` // 用户ID
}

type ListOrderResp struct {
	Orders []OrderSummary // 订单列表
}

// 修改后的 biz.OrderSummary 结构体
type OrderSummary struct {
	OrderId      string      // 订单 ID
	Status       string      // 订单状态
	CreatedAt    int32       // 创建时间（Unix 时间戳）
	Address      Address     // 地址
	Email        string      // 用户邮箱
	UserId       uint32      // 用户 ID
	UserCurrency string      // 用户货币类型
	OrderItems   []OrderItem // 订单商品列表
}

type OrderItem struct {
	Id        int32 // 商品ID
	OrderId   int32
	ProductId int32   // 商品ID
	Name      string  // 商品名称
	Price     float32 // 商品单价
	Quantity  int32   // 商品数量
}

type MarkOrderPaidReq struct {
	UserId  uint32 // 用户ID
	OrderId string // 订单ID
}

type MarkOrderPaidResp struct {
	Success bool `json:"success"`
}
