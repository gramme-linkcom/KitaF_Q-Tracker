package service

import (
	"fmt"
	"kfqt_backend/internal/repository"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (env *APIEnv) GetExistsTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	myNumberStr := r.URL.Query().Get("myNumber")
	ticketUUIDStr  := r.URL.Query().Get("uuid")

	myNumber, err := strconv.Atoi(myNumberStr)
	if err != nil {
		// 数字じゃない文字（"abc"など）が入っていたらエラー
		http.Error(w, `{"error": "myNumberは正しい数値で指定してください"}`, http.StatusBadRequest)
		return
	}

	ticketUUID, err := uuid.Parse(ticketUUIDStr)
	if err != nil {
		http.Error(w, `{"error": "不正なUUIDフォーマットです"}`, http.StatusBadRequest)
		return
	}
	
	isTicketAvailable, err := repository.ExistsTicket(env.DB, ticketUUID, myNumber)
	if err != nil {
		http.Error(w, `{"error": "データベースからデータを取得できませんでした。"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := fmt.Sprintf(`{"isTicketAvailable": %t}`, isTicketAvailable)
	w.Write([]byte(resp))
}
