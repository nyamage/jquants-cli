package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config はアプリケーションの設定を表します
type Config struct {
	APIKey string `mapstructure:"api_key"`
}

// Load は設定を読み込みます
// 優先順位: 1. 環境変数 2. 設定ファイル
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// ホームディレクトリから.jquantsディレクトリを取得
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".jquants")
		viper.AddConfigPath(configPath)
	}

	// 環境変数を読み込み
	viper.SetEnvPrefix("JQUANTS")
	viper.AutomaticEnv()
	// APIキーを環境変数にバインド
	viper.BindEnv("api_key")

	// 設定ファイルを読み込み（存在しない場合はスキップ）
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save は設定をファイルに保存します
func Save(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".jquants")
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return err
	}

	viper.Set("api_key", cfg.APIKey)

	configFile := filepath.Join(configPath, "config.yaml")
	return viper.WriteConfigAs(configFile)
}
