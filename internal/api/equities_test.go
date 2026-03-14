package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetEquityMaster(t *testing.T) {
	// テストデータを読み込み
	testData, err := os.ReadFile("testdata/equity_master.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	tests := []struct {
		name           string
		code           string
		date           string
		expectedCount  int
		expectedStatus int
		responseBody   string
	}{
		{
			name:           "全銘柄取得",
			code:           "",
			date:           "",
			expectedCount:  3,
			expectedStatus: http.StatusOK,
			responseBody:   string(testData),
		},
		{
			name:           "特定銘柄取得",
			code:           "86970",
			date:           "",
			expectedCount:  1,
			expectedStatus: http.StatusOK,
			responseBody:   `{"data":[{"Date":"2025-03-14","Code":"86970","CoName":"日本取引所グループ","CoNameEn":"Japan Exchange Group,Inc.","S17":"16","S17Nm":"金融（除く銀行）","S33":"7200","S33Nm":"その他金融業","ScaleCat":"TOPIX Large70","Mkt":"0111","MktNm":"プライム","Mrgn":"1","MrgnNm":"信用"}]}`,
		},
		{
			name:           "日付指定取得",
			code:           "",
			date:           "2025-03-14",
			expectedCount:  3,
			expectedStatus: http.StatusOK,
			responseBody:   string(testData),
		},
		{
			name:           "銘柄コードと日付指定",
			code:           "86970",
			date:           "2025-03-14",
			expectedCount:  1,
			expectedStatus: http.StatusOK,
			responseBody:   `{"data":[{"Date":"2025-03-14","Code":"86970","CoName":"日本取引所グループ","CoNameEn":"Japan Exchange Group,Inc.","S17":"16","S17Nm":"金融（除く銀行）","S33":"7200","S33Nm":"その他金融業","ScaleCat":"TOPIX Large70","Mkt":"0111","MktNm":"プライム","Mrgn":"1","MrgnNm":"信用"}]}`,
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

				// レスポンスを返す
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.expectedStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// テスト用のクライアントを作成
			client := NewClient("test-api-key")
			client.BaseURL = server.URL

			// GetEquityMasterを実行
			equities, err := client.GetEquityMaster(tt.code, tt.date)
			if err != nil {
				t.Fatalf("GetEquityMaster failed: %v", err)
			}

			// 結果を検証
			if len(equities) != tt.expectedCount {
				t.Errorf("Expected %d equities, got %d", tt.expectedCount, len(equities))
			}

			// 最初のエントリのフィールドを検証
			if len(equities) > 0 {
				first := equities[0]
				if first.Code == "" {
					t.Error("Code should not be empty")
				}
				if first.CoName == "" {
					t.Error("CoName should not be empty")
				}
				if first.Date == "" {
					t.Error("Date should not be empty")
				}
			}
		})
	}
}

func TestGetEquityMasterError(t *testing.T) {
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

			// GetEquityMasterを実行
			_, err := client.GetEquityMaster("", "")
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
