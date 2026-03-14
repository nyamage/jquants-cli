package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	"github.com/nyamage/jquants-cli/internal/api"
)

// CSVFormatter はCSV形式でデータを出力します
type CSVFormatter struct{}

// Format はデータをCSV形式で出力します
func (f *CSVFormatter) Format(w io.Writer, data interface{}) error {
	switch v := data.(type) {
	case []api.EquityMaster:
		return f.formatEquityMaster(w, v)
	case []api.DailyBar:
		return f.formatDailyBar(w, v)
	case []api.IndexBar:
		return f.formatIndexBar(w, v)
	default:
		return f.formatGeneric(w, data)
	}
}

// formatEquityMaster は銘柄情報をCSV形式で出力します
func (f *CSVFormatter) formatEquityMaster(w io.Writer, equities []api.EquityMaster) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// ヘッダー行を出力
	header := []string{
		"Date",
		"Code",
		"CoName",
		"CoNameEn",
		"S17",
		"S17Nm",
		"S33",
		"S33Nm",
		"ScaleCat",
		"Mkt",
		"MktNm",
		"Mrgn",
		"MrgnNm",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// データ行を出力
	for _, eq := range equities {
		record := []string{
			eq.Date,
			eq.Code,
			eq.CoName,
			eq.CoNameEn,
			eq.S17,
			eq.S17Nm,
			eq.S33,
			eq.S33Nm,
			eq.ScaleCat,
			eq.Mkt,
			eq.MktNm,
			eq.Mrgn,
			eq.MrgnNm,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// formatDailyBar は株価四本値をCSV形式で出力します
func (f *CSVFormatter) formatDailyBar(w io.Writer, bars []api.DailyBar) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// ヘッダー行を出力
	header := []string{
		"Date",
		"Code",
		"O",
		"H",
		"L",
		"C",
		"UL",
		"LL",
		"Vo",
		"Va",
		"AdjFactor",
		"AdjO",
		"AdjH",
		"AdjL",
		"AdjC",
		"AdjVo",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// データ行を出力
	for _, bar := range bars {
		record := []string{
			bar.Date,
			bar.Code,
			fmt.Sprintf("%f", bar.O),
			fmt.Sprintf("%f", bar.H),
			fmt.Sprintf("%f", bar.L),
			fmt.Sprintf("%f", bar.C),
			bar.UL,
			bar.LL,
			fmt.Sprintf("%f", bar.Vo),
			fmt.Sprintf("%f", bar.Va),
			fmt.Sprintf("%f", bar.AdjFactor),
			fmt.Sprintf("%f", bar.AdjO),
			fmt.Sprintf("%f", bar.AdjH),
			fmt.Sprintf("%f", bar.AdjL),
			fmt.Sprintf("%f", bar.AdjC),
			fmt.Sprintf("%f", bar.AdjVo),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// formatIndexBar は指数四本値をCSV形式で出力します
func (f *CSVFormatter) formatIndexBar(w io.Writer, bars []api.IndexBar) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Codeフィールドがある場合とない場合で分ける
	hasCode := len(bars) > 0 && bars[0].Code != ""

	// ヘッダー行を出力
	var header []string
	if hasCode {
		header = []string{"Date", "Code", "O", "H", "L", "C"}
	} else {
		header = []string{"Date", "O", "H", "L", "C"}
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// データ行を出力
	for _, bar := range bars {
		var record []string
		if hasCode {
			record = []string{
				bar.Date,
				bar.Code,
				fmt.Sprintf("%f", bar.O),
				fmt.Sprintf("%f", bar.H),
				fmt.Sprintf("%f", bar.L),
				fmt.Sprintf("%f", bar.C),
			}
		} else {
			record = []string{
				bar.Date,
				fmt.Sprintf("%f", bar.O),
				fmt.Sprintf("%f", bar.H),
				fmt.Sprintf("%f", bar.L),
				fmt.Sprintf("%f", bar.C),
			}
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// formatGeneric は汎用的なCSV出力を行います
func (f *CSVFormatter) formatGeneric(w io.Writer, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	if v.Len() == 0 {
		return nil
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 最初の要素から列名を取得
	first := v.Index(0)
	if first.Kind() != reflect.Struct {
		return fmt.Errorf("slice elements must be structs")
	}

	// ヘッダー行を出力
	t := first.Type()
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		headers = append(headers, t.Field(i).Name)
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// データ行を出力
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		var row []string
		for j := 0; j < elem.NumField(); j++ {
			row = append(row, fmt.Sprintf("%v", elem.Field(j).Interface()))
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
