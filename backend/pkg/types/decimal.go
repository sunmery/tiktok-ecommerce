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

// NumericToFloat 将 pgtype.Numeric 转换为 float64（兼容新版 pgx）
func NumericToFloat(n pgtype.Numeric) (float64, error) {
	// 1. 检查 Numeric 是否有效（非 NULL）
	if !n.Valid {
		return 0, fmt.Errorf("numeric value is NULL")
	}

	// 2. 使用内置方法转换为 float64
	floatValue, err := n.Float64Value()
	if err != nil {
		return 0, fmt.Errorf("convert numeric to float64 failed: %w", err)
	}

	// 3. 检查转换结果是否有效
	if !floatValue.Valid {
		return 0, fmt.Errorf("numeric cannot be represented as float64")
	}

	return floatValue.Float64, nil
}
