package repository

import (
	"database/sql"
)

func CreateUserTicket(db *sql.DB, pushToken string) (int, error) {
	query := "INSERT INTO tickets (status, device_id) VALUES ('waiting', ?)"
	
	result, err := db.Exec(query, pushToken)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

// CancelUserTicket はユーザーが自分のスマホから整理券をキャンセルした時にステータスを書き換える
func CancelUserTicket(db *sql.DB, bookingNumber int) error {
	query := "UPDATE tickets SET status = 'canceled' WHERE number = ? AND status = 'waiting'"
	
	_, err := db.Exec(query, bookingNumber)
	return err
}
