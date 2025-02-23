package types

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// FromPgUUID 将 pgtype.UUID 转换为 google/uuid
func FromPgUUID(pgU pgtype.UUID) (uuid.UUID, error) {
	if !pgU.Valid {
		return uuid.Nil, fmt.Errorf("invalid UUID")
	}
	return uuid.FromBytes(pgU.Bytes[:])
}

// ToPgUUID 将 google/uuid 转换为 pgtype.UUID
func ToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u, // 假设 pgtype.UUID 使用 [16]byte 存储
		Valid: true,
	}
}
