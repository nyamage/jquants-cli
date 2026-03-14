package api

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL はJ-Quants API v2のベースURL
	DefaultBaseURL = "https://api.jquants.com/v2"
	// DefaultTimeout はHTTPリクエストのデフォルトタイムアウト
	DefaultTimeout = 30 * time.Second
)

// Client はJ-Quants APIクライアントです
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient は新しいAPIクライアントを作成します
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// Get はGETリクエストを実行します
func (c *Client) Get(endpoint string, params url.Values) ([]byte, error) {
	// URLを構築
	reqURL := c.BaseURL + endpoint
	if params != nil {
		reqURL += "?" + params.Encode()
	}

	// リクエストを作成
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	// 認証ヘッダーを設定
	req.Header.Set("x-api-key", c.APIKey)

	// リクエストを実行
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// エラーハンドリング
	if resp.StatusCode != http.StatusOK {
		return nil, handleAPIError(resp)
	}

	// レスポンスボディを読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
