package convert

import (
	"backend/application/order/internal/biz"

	pb "backend/api/order/v1"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPbOrderItems(items []biz.OrderItem) []*pb.OrderItem {
	var pbItems []*pb.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			Item: &pb.CartItem{
				ProductId: uint32(item.Id),       // 传递商品ID
				Quantity:  uint32(item.Quantity), // 传递商品数量
			},
		})
	}
	return pbItems
}

func Float32ToNumeric(value float32) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	// 使用 Scan 方法将 float64 转换为 pgtype.Numeric
	err := numeric.Scan(float64(value))
	if err != nil {
		return pgtype.Numeric{}, err
	}
	return numeric, nil
}
