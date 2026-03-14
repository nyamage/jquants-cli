package api

import (
	"encoding/json"
	"net/url"
)

// DailyBar は株価四本値を表します
type DailyBar struct {
	Date      string  `json:"Date"`      // 日付 (YYYY-MM-DD)
	Code      string  `json:"Code"`      // 銘柄コード
	O         float64 `json:"O"`         // 始値（調整前）
	H         float64 `json:"H"`         // 高値（調整前）
	L         float64 `json:"L"`         // 安値（調整前）
	C         float64 `json:"C"`         // 終値（調整前）
	UL        string  `json:"UL"`        // 日通ストップ高フラグ
	LL        string  `json:"LL"`        // 日通ストップ安フラグ
	Vo        float64 `json:"Vo"`        // 取引高（調整前）
	Va        float64 `json:"Va"`        // 取引代金
	AdjFactor float64 `json:"AdjFactor"` // 調整係数
	AdjO      float64 `json:"AdjO"`      // 調整済み始値
	AdjH      float64 `json:"AdjH"`      // 調整済み高値
	AdjL      float64 `json:"AdjL"`      // 調整済み安値
	AdjC      float64 `json:"AdjC"`      // 調整済み終値
	AdjVo     float64 `json:"AdjVo"`     // 調整済み取引高
	// 以下、Premiumプランのみ取得可能なフィールド
	MO     float64 `json:"MO"`     // 前場始値
	MH     float64 `json:"MH"`     // 前場高値
	ML     float64 `json:"ML"`     // 前場安値
	MC     float64 `json:"MC"`     // 前場終値
	MUL    string  `json:"MUL"`    // 前場ストップ高フラグ
	MLL    string  `json:"MLL"`    // 前場ストップ安フラグ
	MVo    float64 `json:"MVo"`    // 前場売買高
	MVa    float64 `json:"MVa"`    // 前場取引代金
	MAdjO  float64 `json:"MAdjO"`  // 調整済み前場始値
	MAdjH  float64 `json:"MAdjH"`  // 調整済み前場高値
	MAdjL  float64 `json:"MAdjL"`  // 調整済み前場安値
	MAdjC  float64 `json:"MAdjC"`  // 調整済み前場終値
	MAdjVo float64 `json:"MAdjVo"` // 調整済み前場売買高
	AO     float64 `json:"AO"`     // 後場始値
	AH     float64 `json:"AH"`     // 後場高値
	AL     float64 `json:"AL"`     // 後場安値
	AC     float64 `json:"AC"`     // 後場終値
	AUL    string  `json:"AUL"`    // 後場ストップ高フラグ
	ALL    string  `json:"ALL"`    // 後場ストップ安フラグ
	AVo    float64 `json:"AVo"`    // 後場売買高
	AVa    float64 `json:"AVa"`    // 後場取引代金
	AAdjO  float64 `json:"AAdjO"`  // 調整済み後場始値
	AAdjH  float64 `json:"AAdjH"`  // 調整済み後場高値
	AAdjL  float64 `json:"AAdjL"`  // 調整済み後場安値
	AAdjC  float64 `json:"AAdjC"`  // 調整済み後場終値
	AAdjVo float64 `json:"AAdjVo"` // 調整済み後場売買高
}

// DailyBarResponse はGET /v2/equities/bars/dailyのレスポンスです
type DailyBarResponse struct {
	Data          []DailyBar `json:"data"`
	PaginationKey string     `json:"pagination_key,omitempty"`
}

// GetDailyBars は株価四本値を取得します
//
// パラメータ:
//   - code: 銘柄コード (optional, ただしdateと排他的に必須)
//   - date: 日付 (optional, YYYYMMDD or YYYY-MM-DD)
//   - from: 期間の開始日 (optional, YYYYMMDD or YYYY-MM-DD)
//   - to: 期間の終了日 (optional, YYYYMMDD or YYYY-MM-DD)
//
// 返り値:
//   - 株価四本値のリスト
//   - エラー
func (c *Client) GetDailyBars(code, date, from, to string) ([]DailyBar, error) {
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
	body, err := c.Get("/equities/bars/daily", params)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp DailyBarResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
