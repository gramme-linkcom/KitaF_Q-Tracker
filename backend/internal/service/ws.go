package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"kfqt_backend/internal/model"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"

	"github.com/gorilla/websocket"
)

type BroadcastDatas struct {
	PushType	string
	Queue		[]interface{}
}

var ActiveAdminConn *websocket.Conn
var ConnMu sync.Mutex

func BroadcastQueue(data BroadcastDatas) {
	ConnMu.Lock()
	defer ConnMu.Unlock()

	if ActiveAdminConn == nil {
		log.Println("[LOG] 接続中の管理画面が存在しませんでした。")
		return
	}

	payload := map[string]interface{}{
		"type":  data.PushType,
		"queue": data.Queue,
	}

	// ログイン中の「唯一の1台」に直接送信！
	err := ActiveAdminConn.WriteJSON(payload)
	if err != nil {
		log.Println("[WS_WRITE_ERROR] 送信失敗、接続を切断します:", err)
		ActiveAdminConn.Close()
		ActiveAdminConn = nil
	}

	log.Println("BroadcastQueue を実行しました。")
}

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
func (env *APIEnv) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAdminAuth(r) {
		http.Error(w, "Unauthorized Session", http.StatusUnauthorized)
		log.Println("[WS_AUTH_ERROR] 無効なセッションからのWebSocket接続を拒否しました")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket connection failed:", err)
		return
	}
	
	ConnMu.Lock()
	ActiveAdminConn = conn
	ConnMu.Unlock()
	log.Println("[WS_CONNECT] 管理画面のWebSocketが正常に接続・登録されました。")

	defer func() {
		conn.Close()

		ConnMu.Lock()
		if ActiveAdminConn == conn { // 自分自身の接続ならnilにする
			ActiveAdminConn = nil
		}
		ConnMu.Unlock()
		log.Println("[WS_DISCONNECT] 管理画面のWebSocketが切断されました。")

		cookie, err := r.Cookie("admin_session")
		if err == nil {
			sessionMu.Lock()
			delete(adminSessions, cookie.Value)
			log.Printf("[WS_DISCONNECT] セッション %s を名簿から削除し、ロックを解放しました\n", cookie.Value[:8])
			sessionMu.Unlock()
		}
	}()

	env.sendInitialState(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if !checkAdminAuth(r) {
			log.Println("[WS_AUTH_ERROR] 操作中にセッションが無効化されたため、パケットを破棄しました")
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Session Expired"))
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

func (env *APIEnv) sendInitialState(conn *websocket.Conn) {
	tickets, err := repository.GetActiveTickets(env.DB)
	if err != nil {
		log.Fatalln("DBからのデータ取得に失敗しました。DBが破損している可能性があります。")
		return
	}

	// 取得したデータをもとに、API側でロジック計算を行う
	waitingGroups := len(tickets)
	currentNumber := 0
	nextNumber	  := 0
	if len(tickets) > 0 {
		currentNumber = tickets[0].Number
	}
	if len(tickets) > 1 {
		nextNumber = tickets[1].Number
	}

	cfg := system.ReadConfig()

	// admin.html の updateDOM がそのままパースできる器
	initialData := map[string]interface{}{
		"nextNumber":    nextNumber,
		"currentNumber": currentNumber,
		"waitingGroups": waitingGroups,
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
