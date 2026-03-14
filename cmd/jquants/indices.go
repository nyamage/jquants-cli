package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var indicesCmd = &cobra.Command{
	Use:   "indices",
	Short: "Get index data",
	Long:  "Get index data including TOPIX and other indices",
}

var indicesTopixCmd = &cobra.Command{
	Use:   "topix",
	Short: "Get TOPIX index data",
	Long: `Get TOPIX (Tokyo Stock Price Index) daily bars.

Examples:
  # Get all TOPIX data
  jquants indices topix

  # Get TOPIX data within a date range
  jquants indices topix --from 2025-03-01 --to 2025-03-14

  # Output as JSON
  jquants indices topix --output json

  # Output as CSV
  jquants indices topix --output csv > topix.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		// TOPIXデータを取得
		bars, err := client.GetTopixBars(from, to)
		if err != nil {
			return fmt.Errorf("failed to get TOPIX bars: %w", err)
		}

		if len(bars) == 0 {
			fmt.Println("No TOPIX data found")
			return nil
		}

		// 結果を出力
		if err := formatter.Format(os.Stdout, bars); err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		return nil
	},
}

var indicesDailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Get daily index bars",
	Long: `Get daily index bars for various indices.

Examples:
  # Get all index data for a specific index
  jquants indices daily --code 0028

  # Get index data for a specific index on a specific date
  jquants indices daily --code 0028 --date 2025-03-14

  # Get index data for a specific index within a date range
  jquants indices daily --code 0028 --from 2025-03-01 --to 2025-03-14

  # Get all indices data on a specific date
  jquants indices daily --date 2025-03-14

  # Output as JSON
  jquants indices daily --code 0028 --output json

  # Output as CSV
  jquants indices daily --code 0028 --output csv > indices.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		code, _ := cmd.Flags().GetString("code")
		date, _ := cmd.Flags().GetString("date")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		// codeまたはdateのいずれかが必須
		if code == "" && date == "" {
			return fmt.Errorf("either --code or --date is required")
		}

		// 指数データを取得
		bars, err := client.GetIndexBars(code, date, from, to)
		if err != nil {
			return fmt.Errorf("failed to get index bars: %w", err)
		}

		if len(bars) == 0 {
			fmt.Println("No index data found")
			return nil
		}

		// 結果を出力
		if err := formatter.Format(os.Stdout, bars); err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		return nil
	},
}

func init() {
	// indicesサブコマンドにフラグを追加
	indicesTopixCmd.Flags().String("from", "", "Start date for range query (YYYYMMDD or YYYY-MM-DD)")
	indicesTopixCmd.Flags().String("to", "", "End date for range query (YYYYMMDD or YYYY-MM-DD)")

	indicesDailyCmd.Flags().String("code", "", "Index code (e.g., 0028)")
	indicesDailyCmd.Flags().String("date", "", "Specific date (YYYYMMDD or YYYY-MM-DD)")
	indicesDailyCmd.Flags().String("from", "", "Start date for range query (YYYYMMDD or YYYY-MM-DD)")
	indicesDailyCmd.Flags().String("to", "", "End date for range query (YYYYMMDD or YYYY-MM-DD)")

	// サブコマンドを追加
	indicesCmd.AddCommand(indicesTopixCmd)
	indicesCmd.AddCommand(indicesDailyCmd)
}
