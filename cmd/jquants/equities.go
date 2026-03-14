package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var equitiesCmd = &cobra.Command{
	Use:   "equities",
	Short: "Manage equity information",
	Long:  "Get equity master data including company names, sectors, and market information",
}

var equitiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all equities",
	Long: `List all equities or filter by date.

Examples:
  # List all equities as of today
  jquants equities list

  # List all equities as of specific date
  jquants equities list --date 2025-03-14

  # Output as JSON
  jquants equities list --output json

  # Output as CSV
  jquants equities list --output csv > equities.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")

		// 銘柄情報を取得
		equities, err := client.GetEquityMaster("", date)
		if err != nil {
			return fmt.Errorf("failed to get equity master: %w", err)
		}

		// 結果を出力
		if err := formatter.Format(os.Stdout, equities); err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		return nil
	},
}

var equitiesGetCmd = &cobra.Command{
	Use:   "get [CODE]",
	Short: "Get specific equity information",
	Long: `Get information for a specific equity by code.

Examples:
  # Get equity information for code 86970
  jquants equities get 86970

  # Get equity information as of specific date
  jquants equities get 86970 --date 2025-03-14

  # Output as JSON
  jquants equities get 86970 --output json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]
		date, _ := cmd.Flags().GetString("date")

		// 銘柄情報を取得
		equities, err := client.GetEquityMaster(code, date)
		if err != nil {
			return fmt.Errorf("failed to get equity master: %w", err)
		}

		if len(equities) == 0 {
			fmt.Printf("No equity found with code: %s\n", code)
			return nil
		}

		// 結果を出力
		if err := formatter.Format(os.Stdout, equities); err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		return nil
	},
}

func init() {
	// equitiesサブコマンドにフラグを追加
	equitiesListCmd.Flags().String("date", "", "Reference date (YYYYMMDD or YYYY-MM-DD)")
	equitiesGetCmd.Flags().String("date", "", "Reference date (YYYYMMDD or YYYY-MM-DD)")

	// サブコマンドを追加
	equitiesCmd.AddCommand(equitiesListCmd)
	equitiesCmd.AddCommand(equitiesGetCmd)
}
