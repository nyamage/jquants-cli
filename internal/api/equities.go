package api

import (
	"encoding/json"
	"net/url"
)

// EquityMaster は銘柄情報を表します
type EquityMaster struct {
	Date      string `json:"Date"`      // 情報適用年月日 (YYYY-MM-DD)
	Code      string `json:"Code"`      // 銘柄コード
	CoName    string `json:"CoName"`    // 会社名
	CoNameEn  string `json:"CoNameEn"`  // 会社名（英語）
	S17       string `json:"S17"`       // 17業種コード
	S17Nm     string `json:"S17Nm"`     // 17業種コード名
	S33       string `json:"S33"`       // 33業種コード
	S33Nm     string `json:"S33Nm"`     // 33業種コード名
	ScaleCat  string `json:"ScaleCat"`  // 規模コード
	Mkt       string `json:"Mkt"`       // 市場区分コード
	MktNm     string `json:"MktNm"`     // 市場区分名
	Mrgn      string `json:"Mrgn"`      // 貸借信用区分 (Standard/Premiumプランのみ)
	MrgnNm    string `json:"MrgnNm"`    // 貸借信用区分名 (Standard/Premiumプランのみ)
}

// EquityMasterResponse はGET /v2/equities/masterのレスポンスです
type EquityMasterResponse struct {
	Data []EquityMaster `json:"data"`
}

// GetEquityMaster は銘柄情報を取得します
//
// パラメータ:
//   - code: 銘柄コード (optional)
//   - date: 基準となる日付 (optional, YYYYMMDD or YYYY-MM-DD)
//
// 返り値:
//   - 銘柄情報のリスト
//   - エラー
func (c *Client) GetEquityMaster(code, date string) ([]EquityMaster, error) {
	// クエリパラメータを構築
	params := url.Values{}
	if code != "" {
		params.Set("code", code)
	}
	if date != "" {
		params.Set("date", date)
	}

	// APIリクエストを実行
	body, err := c.Get("/equities/master", params)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp EquityMasterResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
