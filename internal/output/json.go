package output

import (
	"encoding/json"
	"io"
)

// JSONFormatter はJSON形式でデータを出力します
type JSONFormatter struct {
	// Indent は出力をインデントするかどうか
	Indent bool
}

// Format はデータをJSON形式で出力します
func (f *JSONFormatter) Format(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	if f.Indent {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(data)
}
