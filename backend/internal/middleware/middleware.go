package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SameSiteOnlyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. リバースプロキシ（Cloudflareやngrok）を考慮して正しいホスト名を取得
		currentHost := r.Header.Get("X-Forwarded-Host")
		if currentHost == "" {
			currentHost = r.Host
		}

		// 2. 正しいプロトコル（http / https）を取得
		scheme := "http://"
		if r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https://"
		}

		// これが「このサーバー自身の現在のオリジン」（例: https://kfs-wp.9ramme.net や http://localhost:8080）
		myOrigin := fmt.Sprintf("%s%s", scheme, currentHost)

		// 3. CORSヘッダーを動的に設定
		w.Header().Set("Access-Control-Allow-Origin", myOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 4. Origin と Referer の厳密かつ汎用的な検証
		origin := r.Header.Get("Origin")
		referer := r.Header.Get("Referer")

		isAllowed := false

		// 💡 ポート番号の有無による文字列ミスマッチを防ぐため、ドメイン部分（Host）を取り出して比較する
		myHost := currentHost
		if idx := strings.Index(myHost, ":"); idx != -1 {
			myHost = myHost[:idx] // ポート番号（:8080など）をカットして純粋なドメイン名にする
		}

		if origin == "" && referer == "" {
			// 直接アクセス、またはモバイルアプリ等からの通信は許可
			isAllowed = true
		}

		if origin != "" {
			// OriginのURLをパースしてホスト名を取り出す
			if u, err := url.Parse(origin); err == nil {
				originHost := u.Host
				if idx := strings.Index(originHost, ":"); idx != -1 {
					originHost = originHost[:idx]
				}
				// 🚀 自分のホスト名と完全に一致するか検証（これならドメインが変わっても自動追従！）
				if originHost == myHost {
					isAllowed = true
				}
			}
		}

		if referer != "" && !isAllowed {
			// RefererのURLをパースしてホスト名を取り出す
			if u, err := url.Parse(referer); err == nil {
				refererHost := u.Host
				if idx := strings.Index(refererHost, ":"); idx != -1 {
					refererHost = refererHost[:idx]
				}
				// 🚀 リファラーのホスト名が自分と一致するか検証
				if refererHost == myHost {
					isAllowed = true
				}
			}
		}

		// 5. 不正なアクセスは一斉遮断
		if !isAllowed {
			http.Error(w, "Access Denied: 不正なオリジンからのアクセスです", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
