package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetDailyBars(t *testing.T) {
	// テストデータを読み込み
	testData, err := os.ReadFile("testdata/daily_bars.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	tests := []struct {
		name           string
		code           string
		date           string
		from           string
		to             string
		expectedCount  int
		expectedStatus int
		responseBody   string
	}{
		{
			name:           "特定銘柄の全期間データ取得",
			code:           "86970",
			date:           "",
			from:           "",
			to:             "",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			responseBody:   string(testData),
		},
		{
			name:           "特定銘柄の特定日データ取得",
			code:           "86970",
			date:           "2025-03-14",
			from:           "",
			to:             "",
			expectedCount:  1,
			expectedStatus: http.StatusOK,
			responseBody:   `{"data":[{"Date":"2025-03-14","Code":"86970","O":2050.0,"H":2075.0,"L":2040.0,"C":2060.0,"UL":"0","LL":"0","Vo":2500000.0,"Va":5150000000.0,"AdjFactor":1.0,"AdjO":2050.0,"AdjH":2075.0,"AdjL":2040.0,"AdjC":2060.0,"AdjVo":2500000.0,"MO":2050.0,"MH":2075.0,"ML":2045.0,"MC":2055.0,"MUL":"0","MLL":"0","MVo":1300000.0,"MVa":2671500000.0,"MAdjO":2050.0,"MAdjH":2075.0,"MAdjL":2045.0,"MAdjC":2055.0,"MAdjVo":1300000.0,"AO":2055.0,"AH":2070.0,"AL":2040.0,"AC":2060.0,"AUL":"0","ALL":"0","AVo":1200000.0,"AVa":2478500000.0,"AAdjO":2055.0,"AAdjH":2070.0,"AAdjL":2040.0,"AAdjC":2060.0,"AAdjVo":1200000.0}]}`,
		},
		{
			name:           "特定銘柄の期間指定データ取得",
			code:           "86970",
			date:           "",
			from:           "2025-03-13",
			to:             "2025-03-14",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			responseBody:   string(testData),
		},
		{
			name:           "特定日の全銘柄データ取得",
			code:           "",
			date:           "2025-03-14",
			from:           "",
			to:             "",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			responseBody:   string(testData),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサーバーを作成
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// リクエストヘッダーを検証
				if r.Header.Get("x-api-key") == "" {
					t.Error("x-api-key header is missing")
				}

				// クエリパラメータを検証
				query := r.URL.Query()
				if tt.code != "" && query.Get("code") != tt.code {
					t.Errorf("Expected code %s, got %s", tt.code, query.Get("code"))
				}
				if tt.date != "" && query.Get("date") != tt.date {
					t.Errorf("Expected date %s, got %s", tt.date, query.Get("date"))
				}
				if tt.from != "" && query.Get("from") != tt.from {
					t.Errorf("Expected from %s, got %s", tt.from, query.Get("from"))
				}
				if tt.to != "" && query.Get("to") != tt.to {
					t.Errorf("Expected to %s, got %s", tt.to, query.Get("to"))
				}

				// レスポンスを返す
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.expectedStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// テスト用のクライアントを作成
			client := NewClient("test-api-key")
			client.BaseURL = server.URL

			// GetDailyBarsを実行
			bars, err := client.GetDailyBars(tt.code, tt.date, tt.from, tt.to)
			if err != nil {
				t.Fatalf("GetDailyBars failed: %v", err)
			}

			// 結果を検証
			if len(bars) != tt.expectedCount {
				t.Errorf("Expected %d bars, got %d", tt.expectedCount, len(bars))
			}

			// 最初のエントリのフィールドを検証
			if len(bars) > 0 {
				first := bars[0]
				if first.Code == "" {
					t.Error("Code should not be empty")
				}
				if first.Date == "" {
					t.Error("Date should not be empty")
				}
				if first.C == 0 {
					t.Error("Close price should not be zero")
				}
				if first.AdjC == 0 {
					t.Error("Adjusted close price should not be zero")
				}
			}
		})
	}
}

func TestGetDailyBarsError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedErrMsg string
	}{
		{
			name:           "401 Unauthorized",
			statusCode:     http.StatusUnauthorized,
			responseBody:   "",
			expectedErrMsg: "Invalid API key",
		},
		{
			name:           "429 Rate Limit",
			statusCode:     http.StatusTooManyRequests,
			responseBody:   "",
			expectedErrMsg: "Rate limit exceeded",
		},
		{
			name:           "500 Server Error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   "",
			expectedErrMsg: "server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサーバーを作成
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// テスト用のクライアントを作成
			client := NewClient("test-api-key")
			client.BaseURL = server.URL

			// GetDailyBarsを実行
			_, err := client.GetDailyBars("86970", "", "", "")
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			// エラーメッセージを検証
			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("Expected APIError, got %T", err)
			}

			if apiErr.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, apiErr.StatusCode)
			}
		})
	}
}
