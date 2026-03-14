package main

import (
	"fmt"
	"os"

	"github.com/nyamage/jquants-cli/internal/api"
	"github.com/nyamage/jquants-cli/internal/config"
	"github.com/nyamage/jquants-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	// グローバルフラグ
	apiKey       string
	outputFormat string
	verbose      bool

	// グローバル変数
	client    *api.Client
	formatter output.Formatter
)

var rootCmd = &cobra.Command{
	Use:   "jquants",
	Short: "J-Quants API CLI tool",
	Long: `jquants-cli is a command-line tool for accessing J-Quants API.
It provides easy access to Japanese stock market data including equities,
prices, and financial information.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// configコマンドはAPIキー不要なのでスキップ
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return nil
		}
		if cmd.Name() == "config" {
			return nil
		}

		// APIキーの取得（優先順位: フラグ > 環境変数 > 設定ファイル）
		if apiKey == "" {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if cfg.APIKey != "" {
				apiKey = cfg.APIKey
			}
		}

		// APIキーが設定されていない場合はエラー
		if apiKey == "" {
			return fmt.Errorf("API key is required. Set it via --api-key flag, JQUANTS_API_KEY environment variable, or 'jquants config set-api-key' command")
		}

		// APIクライアントを作成
		client = api.NewClient(apiKey)

		// 出力フォーマッターを作成
		var err error
		formatter, err = output.NewFormatter(output.Format(outputFormat))
		if err != nil {
			return err
		}

		// JSONフォーマットの場合はインデントを有効化
		if jsonFormatter, ok := formatter.(*output.JSONFormatter); ok {
			jsonFormatter.Indent = true
		}

		return nil
	},
}

func init() {
	// グローバルフラグを定義
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "J-Quants API key")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, csv)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// サブコマンドを追加
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(equitiesCmd)
	rootCmd.AddCommand(pricesCmd)
	rootCmd.AddCommand(indicesCmd)
}

// configCmd は設定管理コマンド
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Manage J-Quants CLI configuration including API key",
}

var configSetAPIKeyCmd = &cobra.Command{
	Use:   "set-api-key [API_KEY]",
	Short: "Set J-Quants API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{
			APIKey: args[0],
		}
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("API key has been saved successfully")
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.APIKey != "" {
			fmt.Printf("API Key: %s...%s\n", cfg.APIKey[:8], cfg.APIKey[len(cfg.APIKey)-4:])
		} else {
			fmt.Println("API Key: (not set)")
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetAPIKeyCmd)
	configCmd.AddCommand(configShowCmd)
}

// エラーハンドリングヘルパー
func handleError(err error) {
	if verbose {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(1)
}
