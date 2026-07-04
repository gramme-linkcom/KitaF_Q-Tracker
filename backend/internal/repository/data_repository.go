package repository

import (
	"database/sql"
)

type Ticket struct {
	Number 		int    `json:"number"`
	DeviceID	string `json:"device_id"`
	Status 		string `json:"status"` // "waiting", "called", "canceled", "absent"
}

// RoomStatus は現在の部屋全体の状況を表す構造体
type RoomStatus struct {
	CurrentNumber int  `json:"currentNumber"`
	IsActive      bool `json:"isActive"`
}

func GetRoomStatus(db *sql.DB) (RoomStatus, error) {
	var room RoomStatus
	err := db.QueryRow("SELECT current_number, is_active FROM room_status WHERE id = 1").Scan(&room.CurrentNumber, &room.IsActive)
	return room, err
}

func GetAheadGroups(db *sql.DB, myNumberStr string) (int) {
	myAheadGroups := 0
	query := "SELECT COUNT(*) FROM tickets WHERE number < ? AND status = 'waiting'"
	err := db.QueryRow(query, myNumberStr).Scan(&myAheadGroups)
	if err != nil {
		// エラー時は安全に0にしておく
		return 0
	}
	return myAheadGroups 
}

// GetActiveTickets は待機中("waiting")のチケットの一覧を番号順にそのまま取得する
func GetActiveTickets(db *sql.DB) ([]Ticket, error) {
	rows, err := db.Query("SELECT number, device_id, status FROM tickets WHERE status = 'waiting' ORDER BY number ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.Number, &t.DeviceID, &t.Status); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	if tickets == nil {
		tickets = []Ticket{}
	}
	return tickets, nil
}
