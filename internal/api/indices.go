package api

import (
	"encoding/json"
	"net/url"
)

// IndexBar は指数四本値を表します
type IndexBar struct {
	Date string  `json:"Date"` // 日付 (YYYY-MM-DD)
	Code string  `json:"Code"` // 指数コード（idx-bars-dailyのみ）
	O    float64 `json:"O"`    // 始値
	H    float64 `json:"H"`    // 高値
	L    float64 `json:"L"`    // 安値
	C    float64 `json:"C"`    // 終値
}

// IndexBarResponse はGET /v2/indices/bars/daily, /v2/indices/bars/daily/topixのレスポンスです
type IndexBarResponse struct {
	Data          []IndexBar `json:"data"`
	PaginationKey string     `json:"pagination_key,omitempty"`
}

// GetTopixBars はTOPIX指数四本値を取得します
//
// パラメータ:
//   - from: 期間の開始日 (optional, YYYYMMDD or YYYY-MM-DD)
//   - to: 期間の終了日 (optional, YYYYMMDD or YYYY-MM-DD)
//
// 返り値:
//   - TOPIX指数四本値のリスト
//   - エラー
func (c *Client) GetTopixBars(from, to string) ([]IndexBar, error) {
	// クエリパラメータを構築
	params := url.Values{}
	if from != "" {
		params.Set("from", from)
	}
	if to != "" {
		params.Set("to", to)
	}

	// APIリクエストを実行
	body, err := c.Get("/indices/bars/daily/topix", params)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp IndexBarResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GetIndexBars は指数四本値を取得します
//
// パラメータ:
//   - code: 指数コード (optional, ただしdateと排他的に必須)
//   - date: 日付 (optional, YYYYMMDD or YYYY-MM-DD)
//   - from: 期間の開始日 (optional, YYYYMMDD or YYYY-MM-DD)
//   - to: 期間の終了日 (optional, YYYYMMDD or YYYY-MM-DD)
//
// 返り値:
//   - 指数四本値のリスト
//   - エラー
func (c *Client) GetIndexBars(code, date, from, to string) ([]IndexBar, error) {
	// クエリパラメータを構築
	params := url.Values{}
	if code != "" {
		params.Set("code", code)
	}
	if date != "" {
		params.Set("date", date)
	}
	if from != "" {
		params.Set("from", from)
	}
	if to != "" {
		params.Set("to", to)
	}

	// APIリクエストを実行
	body, err := c.Get("/indices/bars/daily", params)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp IndexBarResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
