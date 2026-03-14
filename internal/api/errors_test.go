package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAPIErrorError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		expectedString string
	}{
		{
			name:           "401 Unauthorized",
			statusCode:     http.StatusUnauthorized,
			message:        "Invalid API key",
			expectedString: "API error (401): Invalid API key",
		},
		{
			name:           "429 Too Many Requests",
			statusCode:     http.StatusTooManyRequests,
			message:        "Rate limit exceeded",
			expectedString: "API error (429): Rate limit exceeded",
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			message:        "Server error",
			expectedString: "API error (500): Server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{
				StatusCode: tt.statusCode,
				Message:    tt.message,
			}

			if err.Error() != tt.expectedString {
				t.Errorf("Expected error string %s, got %s", tt.expectedString, err.Error())
			}
		})
	}
}

func TestHandleAPIError(t *testing.T) {
	tests := []struct {
		name               string
		statusCode         int
		responseBody       string
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:               "401 Unauthorized",
			statusCode:         http.StatusUnauthorized,
			responseBody:       "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    "Invalid API key. Please check your API key configuration.",
		},
		{
			name:               "429 Too Many Requests",
			statusCode:         http.StatusTooManyRequests,
			responseBody:       "",
			expectedStatusCode: http.StatusTooManyRequests,
			expectedMessage:    "Rate limit exceeded. Please wait and try again.",
		},
		{
			name:               "500 Internal Server Error",
			statusCode:         http.StatusInternalServerError,
			responseBody:       "",
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "J-Quants API server error. Please try again later.",
		},
		{
			name:               "502 Bad Gateway",
			statusCode:         http.StatusBadGateway,
			responseBody:       "",
			expectedStatusCode: http.StatusBadGateway,
			expectedMessage:    "J-Quants API server error. Please try again later.",
		},
		{
			name:               "503 Service Unavailable",
			statusCode:         http.StatusServiceUnavailable,
			responseBody:       "",
			expectedStatusCode: http.StatusServiceUnavailable,
			expectedMessage:    "J-Quants API server error. Please try again later.",
		},
		{
			name:               "400 Bad Request with custom message",
			statusCode:         http.StatusBadRequest,
			responseBody:       "Invalid parameter",
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid parameter",
		},
		{
			name:               "404 Not Found with custom message",
			statusCode:         http.StatusNotFound,
			responseBody:       "Resource not found",
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックレスポンスを作成
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
			}

			// handleAPIErrorを実行
			err := handleAPIError(resp)

			// エラーがAPIErrorであることを確認
			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("Expected APIError, got %T", err)
			}

			// ステータスコードを検証
			if apiErr.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, apiErr.StatusCode)
			}

			// メッセージを検証
			if apiErr.Message != tt.expectedMessage {
				t.Errorf("Expected message %s, got %s", tt.expectedMessage, apiErr.Message)
			}
		})
	}
}
