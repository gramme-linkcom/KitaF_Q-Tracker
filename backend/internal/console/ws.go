package console

import (
	"encoding/json"
	"log"
	"net/http"

	"kfqt_backend/internal/api"
	"kfqt_backend/internal/model"
	"kfqt_backend/internal/system"

	"github.com/gorilla/websocket"
)

type Response struct {
	Action    string `json:"action"`
	RequestID string `json:"request_id"` // どのリクエストへの返子か特定するため
	Status    string `json:"status"`     // "success" や "error"
	Message   string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler は最小限の接続維持と初期データ送信を行います
func WebSocketHandler(env *api.APIEnv, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket connection failed:", err)
		return
	}
	defer conn.Close()

	sendInitialState(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var newConfigData model.Config
		err = json.Unmarshal(msg, &newConfigData)
		if err != nil {
			log.Printf("[ERROR] JSONの変換に失敗しました: %v", err)
			return
		}
		system.SaveConfig(newConfigData)
		
	}
}

func sendInitialState(conn *websocket.Conn) {
	cfg := system.ReadConfig()

	// admin.html の updateDOM がそのままパースできる器
	initialData := map[string]interface{}{
		"nextNumber":    1,
		"currentNumber": 0,
		"waitingGroups": 0,
		"tickets":       []interface{}{}, // 待ち列一覧（空）
		"config": map[string]interface{}{
			"page_title":              cfg.PageTitle,
			"room_name":               cfg.RoomName,
			"time_required":           cfg.TimeRequired,
			"time_required_range_min": cfg.TimeRequiredRangeMin,
			"time_required_range_max": cfg.TimeRequiredRangeMax,
			"serve_start_time":        cfg.ServeStartTime,
			"serve_end_time":          cfg.ServeEndTime,
			"infomation":              cfg.Infomation,
			"is_booking_available":   cfg.IsBookingAvailable,
			"admin_console_address":  cfg.AdminConsoleAddress,
		},
	}

	_ = conn.WriteJSON(initialData)
}
