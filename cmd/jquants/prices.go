package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var pricesCmd = &cobra.Command{
	Use:   "prices",
	Short: "Get stock price data",
	Long:  "Get stock price data including daily bars (OHLCV)",
}

var pricesDailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Get daily price bars (OHLCV)",
	Long: `Get daily price bars (OHLCV) for stocks.

Examples:
  # Get all price data for a specific stock
  jquants prices daily --code 86970

  # Get price data for a specific stock on a specific date
  jquants prices daily --code 86970 --date 2025-03-14

  # Get price data for a specific stock within a date range
  jquants prices daily --code 86970 --from 2025-03-01 --to 2025-03-14

  # Get all stocks' price data on a specific date
  jquants prices daily --date 2025-03-14

  # Output as JSON
  jquants prices daily --code 86970 --output json

  # Output as CSV
  jquants prices daily --code 86970 --output csv > prices.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		code, _ := cmd.Flags().GetString("code")
		date, _ := cmd.Flags().GetString("date")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		// codeまたはdateのいずれかが必須
		if code == "" && date == "" {
			return fmt.Errorf("either --code or --date is required")
		}

		// 株価データを取得
		bars, err := client.GetDailyBars(code, date, from, to)
		if err != nil {
			return fmt.Errorf("failed to get daily bars: %w", err)
		}

		if len(bars) == 0 {
			fmt.Println("No price data found")
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
	// pricesサブコマンドにフラグを追加
	pricesDailyCmd.Flags().String("code", "", "Stock code (e.g., 86970)")
	pricesDailyCmd.Flags().String("date", "", "Specific date (YYYYMMDD or YYYY-MM-DD)")
	pricesDailyCmd.Flags().String("from", "", "Start date for range query (YYYYMMDD or YYYY-MM-DD)")
	pricesDailyCmd.Flags().String("to", "", "End date for range query (YYYYMMDD or YYYY-MM-DD)")

	// サブコマンドを追加
	pricesCmd.AddCommand(pricesDailyCmd)
}
