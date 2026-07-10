package model

type BookingRequest struct {
	PushToken    string `json:"pushToken"`    // スマホのPush通知用識別子（通知拒否や非対応なら空文字 ""）
	ReservedTime string `json:"reservedTime"` // 指定された予約時間（例: "10:00 - 10:30"）
}

// 予約成功時にフロントへ返すJSONの形
type BookingResponse struct {
	BookingNumber int  `json:"bookingNumber"` // 発行された整理券番号
	Uuid		  string `json:"uuid"`
	Success       bool `json:"success"`
}

// フロントから送られてくるキャンセルのリクエスト構造体
type CancelRequest struct {
	BookingNumber int `json:"bookingNumber"`
}
