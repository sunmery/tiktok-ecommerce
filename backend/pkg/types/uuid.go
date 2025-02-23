package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

// UUID 统一使用指针接收器
type UUID struct {
	uuid.UUID
}

// ParseUUID 从字符串转换（返回指针）
func ParseUUID(s string) (*UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("无效的UUID格式: %s", s)
	}
	return &UUID{UUID: u}, nil
}

// NewUUID 生成新UUID（返回指针）
func NewUUID() *UUID {
	return &UUID{UUID: uuid.New()}
}

// IsNil 空UUID判断（指针接收器）
func (u *UUID) IsNil() bool {
	return u == nil || u.UUID == uuid.Nil
}

// Scan 实现数据库扫描接口（指针接收器）
func (u *UUID) Scan(value interface{}) error {
	if value == nil {
		u.UUID = uuid.Nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return u.UnmarshalBinary(v)
	case string:
		return u.UnmarshalText([]byte(v))
	default:
		return fmt.Errorf("无法扫描UUID类型: %T", value)
	}
}

// Value 实现数据库值转换接口（指针接收器）
func (u *UUID) Value() (driver.Value, error) {
	if u.IsNil() {
		return nil, nil
	}
	return u.UUID.String(), nil
}

// MarshalJSON 序列化（指针接收器）
func (u *UUID) MarshalJSON() ([]byte, error) {
	if u.IsNil() {
		return json.Marshal(nil)
	}
	return json.Marshal(u.UUID.String())
}

// UnmarshalJSON 反序列化（指针接收器）
func (u *UUID) UnmarshalJSON(data []byte) error {
	if u == nil {
		return fmt.Errorf("尝试反序列化到nil指针")
	}

	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		u.UUID = uuid.Nil
		return nil
	}

	parsed, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	u.UUID = parsed
	return nil
}

// NullUUID 处理可空UUID（统一指针接收器）
type NullUUID struct {
	*UUID
	Valid bool
}

func (nu *NullUUID) Scan(value interface{}) error {
	if value == nil {
		nu.UUID = nil
		nu.Valid = false
		return nil
	}

	nu.Valid = true
	if nu.UUID == nil {
		nu.UUID = NewUUID()
	}
	return nu.UUID.Scan(value)
}

func (nu *NullUUID) Value() (driver.Value, error) {
	if !nu.Valid || nu.UUID == nil {
		return nil, nil
	}
	return nu.UUID.Value()
}

// 实现字符串转换接口
func (u *UUID) String() string {
	if u == nil {
		return uuid.Nil.String()
	}
	return u.UUID.String()
}
