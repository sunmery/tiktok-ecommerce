package data

import "testing"

func TestParseSubOrders(t *testing.T) {
	// 模拟数据库返回的两种 items 格式
	testCases := []struct {
		name     string
		input    string
		expected int // 期望解析出的订单项数量
	}{
		{
			name: "数组格式",
			input: `[{
                "id": "sub1",
                "merchant_id": "a1b2c3d4-...",
                "total_amount": 100,
                "items": [{"item": {"merchant_id": "m1", "product_id": "p1", "quantity": 2}, "cost": 50}]
            }]`,
			expected: 1,
		},
		{
			name: "对象格式（兼容旧数据）",
			input: `[{
                "id": "sub1",
                "merchant_id": "a1b2c3d4-...",
                "total_amount": 100,
                "items": {"item": {"merchant_id": "m1", "product_id": "p1", "quantity": 2}, "cost": 50}
            }]`,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			subOrders, err := parseSubOrders([]byte(tc.input))
			if err != nil {
				t.Fatalf("解析失败: %v", err)
			}
			if len(subOrders[0].Items) != tc.expected {
				t.Errorf("期望 %d 个订单项，实际 %d", tc.expected, len(subOrders[0].Items))
			}
		})
	}
}
