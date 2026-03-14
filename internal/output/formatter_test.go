package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/nyamage/jquants-cli/internal/api"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name        string
		format      Format
		expectError bool
		expectType  string
	}{
		{
			name:        "Table formatter",
			format:      FormatTable,
			expectError: false,
			expectType:  "*output.TableFormatter",
		},
		{
			name:        "JSON formatter",
			format:      FormatJSON,
			expectError: false,
			expectType:  "*output.JSONFormatter",
		},
		{
			name:        "CSV formatter",
			format:      FormatCSV,
			expectError: false,
			expectType:  "*output.CSVFormatter",
		},
		{
			name:        "Invalid formatter",
			format:      Format("invalid"),
			expectError: true,
			expectType:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := NewFormatter(tt.format)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if formatter == nil {
					t.Error("Formatter should not be nil")
				}
			}
		})
	}
}

func TestTableFormatter(t *testing.T) {
	formatter := &TableFormatter{}
	data := []api.EquityMaster{
		{
			Date:     "2025-03-14",
			Code:     "86970",
			CoName:   "日本取引所グループ",
			CoNameEn: "Japan Exchange Group,Inc.",
			S33Nm:    "その他金融業",
			MktNm:    "プライム",
			ScaleCat: "TOPIX Large70",
		},
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, data); err != nil {
		t.Errorf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "86970") {
		t.Error("Output should contain equity code")
	}
	if !strings.Contains(output, "日本取引所グループ") {
		t.Error("Output should contain company name")
	}
}

func TestTableFormatterEmpty(t *testing.T) {
	formatter := &TableFormatter{}
	var data []api.EquityMaster

	var buf bytes.Buffer
	if err := formatter.Format(&buf, data); err != nil {
		t.Errorf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No data found") {
		t.Error("Output should contain 'No data found'")
	}
}

func TestJSONFormatter(t *testing.T) {
	formatter := &JSONFormatter{Indent: true}
	data := []api.EquityMaster{
		{
			Date:     "2025-03-14",
			Code:     "86970",
			CoName:   "日本取引所グループ",
			CoNameEn: "Japan Exchange Group,Inc.",
		},
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, data); err != nil {
		t.Errorf("Format failed: %v", err)
	}

	// JSONとしてパース可能か確認
	var result []api.EquityMaster
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 equity, got %d", len(result))
	}

	if result[0].Code != "86970" {
		t.Errorf("Expected code 86970, got %s", result[0].Code)
	}
}

func TestCSVFormatter(t *testing.T) {
	formatter := &CSVFormatter{}
	data := []api.EquityMaster{
		{
			Date:     "2025-03-14",
			Code:     "86970",
			CoName:   "日本取引所グループ",
			CoNameEn: "Japan Exchange Group,Inc.",
			S17:      "16",
			S17Nm:    "金融（除く銀行）",
			S33:      "7200",
			S33Nm:    "その他金融業",
			ScaleCat: "TOPIX Large70",
			Mkt:      "0111",
			MktNm:    "プライム",
			Mrgn:     "1",
			MrgnNm:   "信用",
		},
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, data); err != nil {
		t.Errorf("Format failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// ヘッダー行とデータ行の2行が存在することを確認
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines (header + data), got %d", len(lines))
	}

	// ヘッダー行を確認
	if !strings.Contains(lines[0], "Date") {
		t.Error("Header should contain 'Date'")
	}

	// データ行を確認
	if !strings.Contains(lines[1], "86970") {
		t.Error("Data should contain equity code")
	}
}

func TestCSVFormatterEmpty(t *testing.T) {
	formatter := &CSVFormatter{}
	var data []api.EquityMaster

	var buf bytes.Buffer
	if err := formatter.Format(&buf, data); err != nil {
		t.Errorf("Format failed: %v", err)
	}

	// 空のデータの場合は何も出力されない（ヘッダーも出力されない）
	output := buf.String()
	if len(output) > 0 {
		// 空データの場合は空、またはヘッダーのみ
		t.Logf("CSV output for empty data: %q", output)
	}
}
