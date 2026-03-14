package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoad(t *testing.T) {
	// テスト後のクリーンアップ
	defer func() {
		viper.Reset()
		os.Unsetenv("JQUANTS_API_KEY")
	}()

	tests := []struct {
		name        string
		setupEnv    func()
		setupFile   func(t *testing.T) string
		expectedKey string
		expectError bool
	}{
		{
			name: "環境変数からAPIキーを読み込み",
			setupEnv: func() {
				os.Setenv("JQUANTS_API_KEY", "env-api-key")
				// viperが環境変数を確実に読み込むように設定
				viper.SetEnvPrefix("JQUANTS")
				viper.AutomaticEnv()
			},
			setupFile:   nil,
			expectedKey: "env-api-key",
			expectError: false,
		},
		{
			name:     "設定ファイルからAPIキーを読み込み",
			setupEnv: nil,
			setupFile: func(t *testing.T) string {
				// 一時ディレクトリを作成
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, "config.yaml")

				// 設定ファイルを作成
				content := []byte("api_key: file-api-key\n")
				if err := os.WriteFile(configPath, content, 0600); err != nil {
					t.Fatalf("Failed to create config file: %v", err)
				}

				// viperに設定ファイルのパスを設定
				viper.AddConfigPath(tmpDir)
				return tmpDir
			},
			expectedKey: "file-api-key",
			expectError: false,
		},
		{
			name:        "設定ファイルがない場合",
			setupEnv:    nil,
			setupFile:   nil,
			expectedKey: "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// viperをリセット
			viper.Reset()

			// 環境変数をクリア
			os.Unsetenv("JQUANTS_API_KEY")

			// テストのセットアップ
			if tt.setupEnv != nil {
				tt.setupEnv()
			}

			if tt.setupFile != nil {
				tt.setupFile(t)
			}

			// Loadを実行
			cfg, err := Load()

			// エラーハンドリングを検証
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// APIキーを検証
			if cfg.APIKey != tt.expectedKey {
				t.Errorf("Expected API key %s, got %s", tt.expectedKey, cfg.APIKey)
			}
		})
	}
}

func TestSave(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir := t.TempDir()

	// ホームディレクトリを一時ディレクトリに置き換え
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// テスト用の設定を作成
	cfg := &Config{
		APIKey: "test-api-key",
	}

	// Saveを実行
	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 設定ファイルが作成されたことを確認
	configPath := filepath.Join(tmpDir, ".jquants", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// ファイルの権限を確認
	info, err := os.Stat(filepath.Join(tmpDir, ".jquants"))
	if err != nil {
		t.Fatalf("Failed to stat config directory: %v", err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("Expected directory permissions 0700, got %o", info.Mode().Perm())
	}

	// 設定ファイルの内容を確認
	viper.Reset()
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	savedKey := viper.GetString("api_key")
	if savedKey != cfg.APIKey {
		t.Errorf("Expected saved API key %s, got %s", cfg.APIKey, savedKey)
	}
}

func TestEnvironmentVariablePriority(t *testing.T) {
	// テスト後のクリーンアップ
	defer func() {
		viper.Reset()
		os.Unsetenv("JQUANTS_API_KEY")
	}()

	// 一時ディレクトリを作成
	tmpDir := t.TempDir()

	// 設定ファイルを作成
	configPath := filepath.Join(tmpDir, "config.yaml")
	content := []byte("api_key: file-api-key\n")
	if err := os.WriteFile(configPath, content, 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// 環境変数を設定
	os.Setenv("JQUANTS_API_KEY", "env-api-key")

	// viperをリセットして再設定
	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.SetEnvPrefix("JQUANTS")
	viper.AutomaticEnv()

	// 設定をロード
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// 環境変数が優先されることを確認
	if cfg.APIKey != "env-api-key" {
		t.Errorf("Expected environment variable to take priority. Got %s", cfg.APIKey)
	}
}
