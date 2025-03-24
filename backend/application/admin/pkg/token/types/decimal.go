package types

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

// Float64ToNumeric 将 float64 转换为 pgtype.Numeric
func Float64ToNumeric(value float64) (pgtype.Numeric, error) {
	// 将 float64 转换为字符串
	str := fmt.Sprintf("%v", value)

	// 创建 pgtype.Numeric 实例
	var numeric pgtype.Numeric
	err := numeric.Scan(str)
	if err != nil {
		return pgtype.Numeric{}, fmt.Errorf("failed to convert float64 to pgtype.Numeric: %w", err)
	}

	return numeric, nil
}
