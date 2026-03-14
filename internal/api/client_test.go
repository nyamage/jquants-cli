package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client.APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.APIKey)
	}

	if client.BaseURL != DefaultBaseURL {
		t.Errorf("Expected base URL %s, got %s", DefaultBaseURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}

	if client.HTTPClient.Timeout != DefaultTimeout {
		t.Errorf("Expected timeout %v, got %v", DefaultTimeout, client.HTTPClient.Timeout)
	}
}

func TestClientGet(t *testing.T) {
	tests := []struct {
		name           string
		endpoint       string
		params         url.Values
		expectedPath   string
		expectedQuery  string
		responseStatus int
		responseBody   string
		expectError    bool
	}{
		{
			name:           "パラメータなしのGETリクエスト",
			endpoint:       "/test",
			params:         nil,
			expectedPath:   "/test",
			expectedQuery:  "",
			responseStatus: http.StatusOK,
			responseBody:   `{"result": "ok"}`,
			expectError:    false,
		},
		{
			name:     "パラメータ付きのGETリクエスト",
			endpoint: "/test",
			params: url.Values{
				"code": []string{"86970"},
				"date": []string{"2025-03-14"},
			},
			expectedPath:   "/test",
			expectedQuery:  "code=86970&date=2025-03-14",
			responseStatus: http.StatusOK,
			responseBody:   `{"result": "ok"}`,
			expectError:    false,
		},
		{
			name:           "200以外のステータスコード",
			endpoint:       "/test",
			params:         nil,
			expectedPath:   "/test",
			expectedQuery:  "",
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"error": "bad request"}`,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサーバーを作成
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// HTTPメソッドを検証
				if r.Method != "GET" {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				// パスを検証
				if r.URL.Path != tt.expectedPath {
					t.Errorf("Expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				// クエリパラメータを検証
				if tt.expectedQuery != "" {
					if r.URL.RawQuery != tt.expectedQuery {
						t.Errorf("Expected query %s, got %s", tt.expectedQuery, r.URL.RawQuery)
					}
				}

				// x-api-keyヘッダーを検証
				apiKey := r.Header.Get("x-api-key")
				if apiKey == "" {
					t.Error("x-api-key header is missing")
				}

				// レスポンスを返す
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// テスト用のクライアントを作成
			client := NewClient("test-api-key")
			client.BaseURL = server.URL

			// Getメソッドを実行
			body, err := client.Get(tt.endpoint, tt.params)

			// エラーハンドリングを検証
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if string(body) != tt.responseBody {
					t.Errorf("Expected body %s, got %s", tt.responseBody, string(body))
				}
			}
		})
	}
}

func TestClientGetTimeout(t *testing.T) {
	// タイムアウトするサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 短いタイムアウトでクライアントを作成
	client := NewClient("test-api-key")
	client.BaseURL = server.URL
	client.HTTPClient.Timeout = 100 * time.Millisecond

	// タイムアウトエラーが発生することを確認
	_, err := client.Get("/test", nil)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
