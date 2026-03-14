package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetTopixBars(t *testing.T) {
	// テストデータを読み込み
	testData, err := os.ReadFile("testdata/topix_bars.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	tests := []struct {
		name           string
		from           string
		to             string
		expectedCount  int
		expectedStatus int
		responseBody   string
	}{
		{
			name:           "全期間データ取得",
			from:           "",
			to:             "",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			responseBody:   string(testData),
		},
		{
			name:           "期間指定データ取得",
			from:           "2025-03-13",
			to:             "2025-03-14",
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

			// GetTopixBarsを実行
			bars, err := client.GetTopixBars(tt.from, tt.to)
			if err != nil {
				t.Fatalf("GetTopixBars failed: %v", err)
			}

			// 結果を検証
			if len(bars) != tt.expectedCount {
				t.Errorf("Expected %d bars, got %d", tt.expectedCount, len(bars))
			}

			// 最初のエントリのフィールドを検証
			if len(bars) > 0 {
				first := bars[0]
				if first.Date == "" {
					t.Error("Date should not be empty")
				}
				if first.C == 0 {
					t.Error("Close price should not be zero")
				}
			}
		})
	}
}

func TestGetIndexBars(t *testing.T) {
	// テストデータを読み込み
	testData, err := os.ReadFile("testdata/index_bars.json")
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
			name:           "特定指数の全期間データ取得",
			code:           "0028",
			date:           "",
			from:           "",
			to:             "",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			responseBody:   string(testData),
		},
		{
			name:           "特定指数の特定日データ取得",
			code:           "0028",
			date:           "2025-03-14",
			from:           "",
			to:             "",
			expectedCount:  1,
			expectedStatus: http.StatusOK,
			responseBody:   `{"data":[{"Date":"2025-03-14","Code":"0028","O":1205.30,"H":1210.45,"L":1202.00,"C":1208.25}]}`,
		},
		{
			name:           "特定指数の期間指定データ取得",
			code:           "0028",
			date:           "",
			from:           "2025-03-13",
			to:             "2025-03-14",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			responseBody:   string(testData),
		},
		{
			name:           "特定日の全指数データ取得",
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

			// GetIndexBarsを実行
			bars, err := client.GetIndexBars(tt.code, tt.date, tt.from, tt.to)
			if err != nil {
				t.Fatalf("GetIndexBars failed: %v", err)
			}

			// 結果を検証
			if len(bars) != tt.expectedCount {
				t.Errorf("Expected %d bars, got %d", tt.expectedCount, len(bars))
			}

			// 最初のエントリのフィールドを検証
			if len(bars) > 0 {
				first := bars[0]
				if first.Date == "" {
					t.Error("Date should not be empty")
				}
				if first.C == 0 {
					t.Error("Close price should not be zero")
				}
			}
		})
	}
}

func TestGetIndexBarsError(t *testing.T) {
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

			// GetIndexBarsを実行
			_, err := client.GetIndexBars("0028", "", "", "")
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
