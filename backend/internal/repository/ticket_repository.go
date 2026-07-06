package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

func ExistsTicket(db *sql.DB, id uuid.UUID, number int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tickets WHERE number = ? AND uuid = ?);"
	var exists bool
	err := db.QueryRow(query, number, id.String()).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
