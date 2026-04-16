package pkg

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/server-selfish/backend/internal/constant"
)

func StringToPgUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	return uuid, err
}

func GetFileNameByTechstack(name string) string {
	switch strings.ToLower(name) {
	case "node":
		return constant.NODE_DOCKERFILE_TEMPLATE
	case "go":
		return constant.GO_DOCKERFILE_TEMPLATE
	case "python":
		return constant.PYTHON_DOCKERFILE_TEMPLATE
	default:
		return constant.NODE_DOCKERFILE_TEMPLATE
	}
}

func ShellToExecForm(cmd string) string {
	args := strings.Fields(cmd)
	b, _ := json.Marshal(args)
	return string(b)
}

func PgUUIDFromUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func PgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t.UTC(),
		Valid: true,
	}
}

func ToPgText(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{Valid: false}
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func StrPtr(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return &v
}
