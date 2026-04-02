package pkg

import "github.com/jackc/pgx/v5/pgtype"

func StringToPgUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	return uuid, err
}
