package api

import (
	"fmt"
	"io"
	"net/http"
)

// APIError はJ-Quants APIのエラーを表します
type APIError struct {
	StatusCode int
	Message    string
}

// Error はerrorインターフェースを実装します
func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// handleAPIError はHTTPレスポンスからAPIErrorを生成します
func handleAPIError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "Invalid API key. Please check your API key configuration.",
		}
	case http.StatusTooManyRequests:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "Rate limit exceeded. Please wait and try again.",
		}
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "J-Quants API server error. Please try again later.",
		}
	default:
		body, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}
}
