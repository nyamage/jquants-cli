package output

import (
	"fmt"
	"io"
	"reflect"

	"github.com/nyamage/jquants-cli/internal/api"
	"github.com/olekukonko/tablewriter"
)

// TableFormatter はテーブル形式でデータを出力します
type TableFormatter struct{}

// Format はデータをテーブル形式で出力します
func (f *TableFormatter) Format(w io.Writer, data interface{}) error {
	switch v := data.(type) {
	case []api.EquityMaster:
		return f.formatEquityMaster(w, v)
	case []api.DailyBar:
		return f.formatDailyBar(w, v)
	case []api.IndexBar:
		return f.formatIndexBar(w, v)
	default:
		return fmt.Errorf("unsupported data type: %T", data)
	}
}

// formatEquityMaster は銘柄情報をテーブル形式で出力します
func (f *TableFormatter) formatEquityMaster(w io.Writer, equities []api.EquityMaster) error {
	if len(equities) == 0 {
		fmt.Fprintln(w, "No data found")
		return nil
	}

	table := tablewriter.NewWriter(w)

	// ヘッダーを設定
	headers := []string{
		"Date",
		"Code",
		"Company Name",
		"Company Name (EN)",
		"Sector",
		"Market",
		"Scale",
	}

	// データ行を追加（ヘッダーを最初の行として追加）
	table.Append(headers)

	for _, eq := range equities {
		table.Append([]string{
			eq.Date,
			eq.Code,
			eq.CoName,
			eq.CoNameEn,
			eq.S33Nm,
			eq.MktNm,
			eq.ScaleCat,
		})
	}

	table.Render()
	return nil
}

// formatDailyBar は株価四本値をテーブル形式で出力します
func (f *TableFormatter) formatDailyBar(w io.Writer, bars []api.DailyBar) error {
	if len(bars) == 0 {
		fmt.Fprintln(w, "No data found")
		return nil
	}

	table := tablewriter.NewWriter(w)

	// ヘッダーを設定
	headers := []string{
		"Date",
		"Code",
		"Open",
		"High",
		"Low",
		"Close",
		"Volume",
		"Adj Close",
		"Adj Volume",
	}

	// データ行を追加（ヘッダーを最初の行として追加）
	table.Append(headers)

	for _, bar := range bars {
		table.Append([]string{
			bar.Date,
			bar.Code,
			fmt.Sprintf("%.2f", bar.O),
			fmt.Sprintf("%.2f", bar.H),
			fmt.Sprintf("%.2f", bar.L),
			fmt.Sprintf("%.2f", bar.C),
			fmt.Sprintf("%.0f", bar.Vo),
			fmt.Sprintf("%.2f", bar.AdjC),
			fmt.Sprintf("%.0f", bar.AdjVo),
		})
	}

	table.Render()
	return nil
}

// formatIndexBar は指数四本値をテーブル形式で出力します
func (f *TableFormatter) formatIndexBar(w io.Writer, bars []api.IndexBar) error {
	if len(bars) == 0 {
		fmt.Fprintln(w, "No data found")
		return nil
	}

	table := tablewriter.NewWriter(w)

	// ヘッダーを設定（Codeフィールドがある場合とない場合で分ける）
	var headers []string
	hasCode := len(bars) > 0 && bars[0].Code != ""

	if hasCode {
		headers = []string{
			"Date",
			"Code",
			"Open",
			"High",
			"Low",
			"Close",
		}
	} else {
		headers = []string{
			"Date",
			"Open",
			"High",
			"Low",
			"Close",
		}
	}

	// データ行を追加（ヘッダーを最初の行として追加）
	table.Append(headers)

	for _, bar := range bars {
		var row []string
		if hasCode {
			row = []string{
				bar.Date,
				bar.Code,
				fmt.Sprintf("%.2f", bar.O),
				fmt.Sprintf("%.2f", bar.H),
				fmt.Sprintf("%.2f", bar.L),
				fmt.Sprintf("%.2f", bar.C),
			}
		} else {
			row = []string{
				bar.Date,
				fmt.Sprintf("%.2f", bar.O),
				fmt.Sprintf("%.2f", bar.H),
				fmt.Sprintf("%.2f", bar.L),
				fmt.Sprintf("%.2f", bar.C),
			}
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

// formatGenericTable は汎用的なテーブル出力を行います
func (f *TableFormatter) formatGenericTable(w io.Writer, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	if v.Len() == 0 {
		fmt.Fprintln(w, "No data found")
		return nil
	}

	// 最初の要素から列名を取得
	first := v.Index(0)
	if first.Kind() != reflect.Struct {
		return fmt.Errorf("slice elements must be structs")
	}

	// ヘッダーを作成
	var headers []string
	t := first.Type()
	for i := 0; i < t.NumField(); i++ {
		headers = append(headers, t.Field(i).Name)
	}

	table := tablewriter.NewWriter(w)

	// ヘッダー行を追加
	table.Append(headers)

	// データ行を追加
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		var row []string
		for j := 0; j < elem.NumField(); j++ {
			row = append(row, fmt.Sprintf("%v", elem.Field(j).Interface()))
		}
		table.Append(row)
	}

	table.Render()
	return nil
}
