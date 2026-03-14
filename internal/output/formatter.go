package output

import (
	"fmt"
	"io"
)

// Format は出力フォーマットの種類を表します
type Format string

const (
	// FormatTable はテーブル形式
	FormatTable Format = "table"
	// FormatJSON はJSON形式
	FormatJSON Format = "json"
	// FormatCSV はCSV形式
	FormatCSV Format = "csv"
)

// Formatter はデータを特定の形式で出力するインターフェースです
type Formatter interface {
	// Format はデータを指定された形式でフォーマットして出力します
	Format(w io.Writer, data interface{}) error
}

// NewFormatter は指定されたフォーマットに対応するFormatterを作成します
func NewFormatter(format Format) (Formatter, error) {
	switch format {
	case FormatTable:
		return &TableFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	case FormatCSV:
		return &CSVFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
