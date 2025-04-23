package types

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ToPgUUID 将 google/uuid 转换为 pgtype.UUID
func ToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u, // 假设 pgtype.UUID 使用 [16]byte 存储
		Valid: true,
	}
}
