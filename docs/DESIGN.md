# jquants-cli 設計ドキュメント

## 目次

1. [概要](#概要)
2. [アーキテクチャ](#アーキテクチャ)
3. [APIクライアント設計](#apiクライアント設計)
4. [コマンド設計](#コマンド設計)
5. [設定管理](#設定管理)
6. [出力フォーマット](#出力フォーマット)
7. [エラーハンドリング](#エラーハンドリング)
8. [テスト戦略](#テスト戦略)
9. [実装フェーズ](#実装フェーズ)

## 概要

jquants-cliは、J-Quants API v2を利用した日本株データ取得CLIツールです。

### 設計原則

- **シンプル**: 直感的で学習コストの低いCLIインターフェース
- **拡張性**: 新しいエンドポイントの追加が容易
- **堅牢性**: 適切なエラーハンドリングとリトライ機能
- **効率性**: ページネーション自動処理とBulk API対応

## アーキテクチャ

### ディレクトリ構成

```
jquants-cli/
├── cmd/
│   └── jquants/              # メインCLIエントリポイント
│       ├── main.go           # アプリケーションエントリポイント
│       ├── root.go           # ルートコマンド
│       ├── config.go         # 設定コマンド
│       ├── equities.go       # 銘柄情報コマンド
│       ├── prices.go         # 株価データコマンド
│       ├── fins.go           # 財務情報コマンド
│       └── bulk.go           # Bulk APIコマンド
├── internal/
│   ├── api/                  # APIクライアント
│   │   ├── client.go         # ベースクライアント
│   │   ├── equities.go       # 銘柄情報API
│   │   ├── prices.go         # 株価データAPI
│   │   ├── financials.go     # 財務情報API
│   │   ├── bulk.go           # Bulk API
│   │   └── pagination.go     # ページネーション処理
│   ├── config/               # 設定管理
│   │   └── config.go         # 設定ファイル読み書き
│   └── output/               # 出力フォーマット
│       ├── formatter.go      # フォーマッターインターフェース
│       ├── table.go          # テーブル形式出力
│       ├── json.go           # JSON形式出力
│       └── csv.go            # CSV形式出力
├── docs/                     # ドキュメント
│   ├── DESIGN.md             # 設計ドキュメント（本ファイル）
│   └── API.md                # API仕様メモ
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

### レイヤー構成

```
┌─────────────────────────────┐
│     CLI Layer (cmd/)        │  ← ユーザーインターフェース
├─────────────────────────────┤
│   Business Logic Layer      │  ← コマンド処理ロジック
├─────────────────────────────┤
│  API Client Layer           │  ← J-Quants API通信
│  (internal/api/)            │
├─────────────────────────────┤
│  Output Layer               │  ← データフォーマット・表示
│  (internal/output/)         │
├─────────────────────────────┤
│  Config Layer               │  ← 設定管理
│  (internal/config/)         │
└─────────────────────────────┘
```

## APIクライアント設計

### Client構造体

```go
// internal/api/client.go
package api

import (
    "net/http"
    "time"
)

type Client struct {
    BaseURL    string
    APIKey     string
    HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
    return &Client{
        BaseURL: "https://api.jquants.com/v2",
        APIKey:  apiKey,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *Client) Get(endpoint string, params url.Values) ([]byte, error) {
    req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("x-api-key", c.APIKey)
    req.URL.RawQuery = params.Encode()

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // エラーハンドリング
    if resp.StatusCode != http.StatusOK {
        return nil, handleAPIError(resp)
    }

    return io.ReadAll(resp.Body)
}
```

### エンドポイント別クライアント

各エンドポイントごとに専用の関数を実装:

```go
// internal/api/equities.go
package api

type EquityMaster struct {
    Date     string `json:"Date"`
    Code     string `json:"Code"`
    CoName   string `json:"CoName"`
    CoNameEn string `json:"CoNameEn"`
    S17      string `json:"S17"`
    S17Nm    string `json:"S17Nm"`
    S33      string `json:"S33"`
    S33Nm    string `json:"S33Nm"`
    ScaleCat string `json:"ScaleCat"`
    Mkt      string `json:"Mkt"`
    MktNm    string `json:"MktNm"`
    Mrgn     string `json:"Mrgn"`
    MrgnNm   string `json:"MrgnNm"`
}

type EquityMasterResponse struct {
    Data []EquityMaster `json:"data"`
}

func (c *Client) GetEquityMaster(code, date string) ([]EquityMaster, error) {
    params := url.Values{}
    if code != "" {
        params.Set("code", code)
    }
    if date != "" {
        params.Set("date", date)
    }

    body, err := c.Get("/equities/master", params)
    if err != nil {
        return nil, err
    }

    var resp EquityMasterResponse
    if err := json.Unmarshal(body, &resp); err != nil {
        return nil, err
    }

    return resp.Data, nil
}
```

### ページネーション処理

```go
// internal/api/pagination.go
package api

import (
    "time"
)

type PaginatedResponse struct {
    Data          json.RawMessage `json:"data"`
    PaginationKey string          `json:"pagination_key,omitempty"`
}

// FetchAllPages はページネーションを自動処理して全データを取得
func (c *Client) FetchAllPages(endpoint string, params url.Values, processPage func([]byte) error) error {
    paginationKey := ""

    for {
        if paginationKey != "" {
            params.Set("pagination_key", paginationKey)
        }

        body, err := c.Get(endpoint, params)
        if err != nil {
            return err
        }

        var resp PaginatedResponse
        if err := json.Unmarshal(body, &resp); err != nil {
            return err
        }

        // ページデータを処理
        if err := processPage(resp.Data); err != nil {
            return err
        }

        // 次のページがあるかチェック
        if resp.PaginationKey == "" {
            break
        }

        paginationKey = resp.PaginationKey

        // レート制限対策: 500ms待機
        time.Sleep(500 * time.Millisecond)
    }

    return nil
}
```

### Bulk API処理

```go
// internal/api/bulk.go
package api

import (
    "compress/gzip"
    "io"
    "os"
)

type BulkFile struct {
    Key          string `json:"Key"`
    LastModified string `json:"LastModified"`
    Size         int64  `json:"Size"`
}

type BulkListResponse struct {
    Data []BulkFile `json:"data"`
}

type BulkGetResponse struct {
    URL string `json:"url"`
}

func (c *Client) ListBulkFiles(endpoint, date, from, to string) ([]BulkFile, error) {
    params := url.Values{}
    if endpoint != "" {
        params.Set("endpoint", endpoint)
    }
    if date != "" {
        params.Set("date", date)
    }
    if from != "" {
        params.Set("from", from)
    }
    if to != "" {
        params.Set("to", to)
    }

    body, err := c.Get("/bulk/list", params)
    if err != nil {
        return nil, err
    }

    var resp BulkListResponse
    if err := json.Unmarshal(body, &resp); err != nil {
        return nil, err
    }

    return resp.Data, nil
}

func (c *Client) GetBulkURL(key string) (string, error) {
    params := url.Values{}
    params.Set("key", key)

    body, err := c.Get("/bulk/get", params)
    if err != nil {
        return "", err
    }

    var resp BulkGetResponse
    if err := json.Unmarshal(body, &resp); err != nil {
        return "", err
    }

    return resp.URL, nil
}

func (c *Client) DownloadBulkFile(url, outputPath string) error {
    resp, err := c.HTTPClient.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // gzip解凍
    gzReader, err := gzip.NewReader(resp.Body)
    if err != nil {
        return err
    }
    defer gzReader.Close()

    // ファイル作成
    outFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outFile.Close()

    // データ書き込み
    _, err = io.Copy(outFile, gzReader)
    return err
}
```

## コマンド設計

### コマンド階層

```
jquants
├── config
│   ├── set-api-key
│   └── show
├── equities
│   ├── list
│   └── get
├── prices
│   └── daily
├── fins
│   ├── get
│   └── list
└── bulk
    ├── list
    └── download
```

### Cobraによる実装例

```go
// cmd/jquants/root.go
package main

import (
    "github.com/spf13/cobra"
)

var (
    apiKey  string
    output  string
    verbose bool
)

var rootCmd = &cobra.Command{
    Use:   "jquants",
    Short: "J-Quants API CLI tool",
    Long:  "Command-line tool to interact with J-Quants API for Japanese stock data",
}

func init() {
    rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "J-Quants API key")
    rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format (table|json|csv)")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}
```

```go
// cmd/jquants/equities.go
package main

import (
    "github.com/spf13/cobra"
)

var equitiesCmd = &cobra.Command{
    Use:   "equities",
    Short: "Equity data commands",
}

var equitiesListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all equities",
    RunE: func(cmd *cobra.Command, args []string) error {
        // 実装
        return nil
    },
}

var equitiesGetCmd = &cobra.Command{
    Use:   "get <code>",
    Short: "Get equity information by code",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // 実装
        return nil
    },
}

func init() {
    equitiesCmd.AddCommand(equitiesListCmd)
    equitiesCmd.AddCommand(equitiesGetCmd)
    rootCmd.AddCommand(equitiesCmd)

    // フラグ定義
    equitiesListCmd.Flags().String("date", "", "Date (YYYY-MM-DD)")
}
```

## 設定管理

### 設定ファイル

設定ファイルパス: `~/.jquants/config.yaml`

```yaml
api_key: YOUR_API_KEY_HERE
```

### 設定優先順位

1. コマンドラインフラグ (`--api-key`)
2. 環境変数 (`JQUANTS_API_KEY`)
3. 設定ファイル (`~/.jquants/config.yaml`)

### 実装

```go
// internal/config/config.go
package config

import (
    "os"
    "path/filepath"

    "github.com/spf13/viper"
)

type Config struct {
    APIKey string `mapstructure:"api_key"`
}

func Load() (*Config, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    configPath := filepath.Join(home, ".jquants")

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(configPath)

    viper.SetEnvPrefix("JQUANTS")
    viper.AutomaticEnv()

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
```

## 出力フォーマット

### Formatterインターフェース

```go
// internal/output/formatter.go
package output

type Formatter interface {
    Format(data interface{}) (string, error)
}

func NewFormatter(format string) Formatter {
    switch format {
    case "json":
        return &JSONFormatter{}
    case "csv":
        return &CSVFormatter{}
    default:
        return &TableFormatter{}
    }
}
```

### テーブル形式

```go
// internal/output/table.go
package output

import (
    "bytes"
    "github.com/olekukonko/tablewriter"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(data interface{}) (string, error) {
    buf := new(bytes.Buffer)
    table := tablewriter.NewWriter(buf)

    // データに応じてヘッダーとデータを設定
    // 実装詳細は省略

    table.Render()
    return buf.String(), nil
}
```

### JSON形式

```go
// internal/output/json.go
package output

import (
    "encoding/json"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(data interface{}) (string, error) {
    bytes, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}
```

### CSV形式

```go
// internal/output/csv.go
package output

import (
    "encoding/csv"
    "bytes"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(data interface{}) (string, error) {
    buf := new(bytes.Buffer)
    writer := csv.NewWriter(buf)

    // データに応じてCSV出力
    // 実装詳細は省略

    writer.Flush()
    return buf.String(), writer.Error()
}
```

## エラーハンドリング

### エラー種別

```go
// internal/api/errors.go
package api

import (
    "fmt"
    "net/http"
)

type APIError struct {
    StatusCode int
    Message    string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

func handleAPIError(resp *http.Response) error {
    switch resp.StatusCode {
    case http.StatusUnauthorized:
        return &APIError{
            StatusCode: resp.StatusCode,
            Message:    "Invalid API key. Please check your API key configuration.",
        }
    case http.StatusTooManyRequests:
        return &APIError{
            StatusCode: resp.StatusCode,
            Message:    "Rate limit exceeded. Please wait and try again.",
        }
    case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
        return &APIError{
            StatusCode: resp.StatusCode,
            Message:    "J-Quants API server error. Please try again later.",
        }
    default:
        body, _ := io.ReadAll(resp.Body)
        return &APIError{
            StatusCode: resp.StatusCode,
            Message:    string(body),
        }
    }
}
```

### リトライ機能

```go
// internal/api/retry.go
package api

import (
    "time"
)

func (c *Client) GetWithRetry(endpoint string, params url.Values, maxRetries int) ([]byte, error) {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        body, err := c.Get(endpoint, params)
        if err == nil {
            return body, nil
        }

        lastErr = err

        // 5xxエラーの場合のみリトライ
        if apiErr, ok := err.(*APIError); ok {
            if apiErr.StatusCode >= 500 {
                waitTime := time.Duration(i+1) * time.Second
                time.Sleep(waitTime)
                continue
            }
        }

        // その他のエラーは即座に返す
        return nil, err
    }

    return nil, lastErr
}
```

## テスト戦略

### 概要

このプロジェクトでは、高品質なコードを維持するために包括的なテスト戦略を採用します。詳細は [docs/TESTING.md](TESTING.md) を参照してください。

### テストレベル

| レベル | 目的 | カバレッジ目標 |
|--------|------|---------------|
| 単体テスト | 関数・メソッドの動作確認 | 80%以上 |
| 結合テスト | コンポーネント間連携確認 | 主要フロー100% |
| E2Eテスト | CLI全体の動作確認 | 主要コマンド100% |

### ディレクトリ構成（テスト追加版）

```
jquants-cli/
├── internal/
│   ├── api/
│   │   ├── client.go
│   │   ├── client_test.go          # 単体テスト
│   │   └── testdata/               # テストデータ
│   ├── config/
│   │   ├── config.go
│   │   ├── config_test.go
│   │   └── testdata/
│   └── output/
│       ├── formatter.go
│       ├── formatter_test.go
│       └── testdata/
├── tests/
│   ├── integration/                # 結合テスト
│   └── fixtures/                   # テストフィクスチャ
└── scripts/
    ├── test.sh                     # テスト実行
    └── coverage.sh                 # カバレッジ計測
```

### テスト実装方針

#### 単体テスト

各パッケージの`*_test.go`ファイルに実装:

```go
func TestClient_Get(t *testing.T) {
    // httptestでモックサーバーを作成
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"data":[]}`))
    }))
    defer server.Close()

    client := NewClient("test-key")
    client.BaseURL = server.URL

    result, err := client.Get("/test", nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // アサーション...
}
```

#### 結合テスト

`tests/integration/`ディレクトリに実装:

```go
func TestAPIClient_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // 実際のAPIキーを使ってテスト
    cfg, err := config.Load()
    if err != nil || cfg.APIKey == "" {
        t.Skip("API key not configured")
    }

    client := api.NewClient(cfg.APIKey)
    // 実際のAPIを呼び出してテスト...
}
```

### テスト実行コマンド

```bash
# 単体テストのみ実行（高速）
go test -short ./...

# 全テスト実行
go test ./...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./...

# カバレッジレポート表示
go tool cover -html=coverage.out

# 結合テストのみ実行
go test ./tests/integration/...
```

### モック戦略

- **httptestパッケージ**: APIクライアントのモックサーバー
- **インターフェース**: 依存性の注入でテスト容易性向上
- **testdataディレクトリ**: サンプルJSONレスポンスを配置

### CI/CD統合

GitHub Actionsで自動テスト:

```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: ./scripts/coverage.sh  # 80%以上チェック
```

## 実装フェーズ

### Phase 1: 基本機能（MVP）

**目標**: 基本的なデータ取得機能を実装

1. プロジェクト構成とCLI基盤
   - Cobra導入
   - ルートコマンド実装
   - グローバルフラグ設定

2. 認証・設定管理
   - 設定ファイル読み書き
   - 環境変数対応
   - `config set-api-key` / `config show` コマンド

3. 基本的なAPIクライアント
   - Client構造体実装
   - GET リクエスト基本機能
   - エラーハンドリング

4. 銘柄一覧取得
   - `/equities/master` エンドポイント対応
   - `equities list` / `equities get` コマンド
   - テーブル形式出力

**成果物**:
- 銘柄情報が取得できるCLI

### Phase 2: データ取得機能拡充

**目標**: 株価・財務データ取得とページネーション対応

5. 株価データ取得
   - `/equities/bars/daily` エンドポイント対応
   - `prices daily` コマンド
   - ページネーション自動処理

6. 財務情報取得
   - `/fins/summary` エンドポイント対応
   - `fins get` / `fins list` コマンド

7. 出力フォーマット対応
   - JSON形式出力
   - CSV形式出力
   - フォーマット切り替え機能

**成果物**:
- 株価・財務データが取得できるCLI
- 複数フォーマット対応

### Phase 3: 高度な機能

**目標**: Bulk API対応と使い勝手向上

8. Bulk API実装
   - `/bulk/list` / `/bulk/get` エンドポイント対応
   - `bulk list` / `bulk download` コマンド
   - gzip解凍処理

9. 進捗表示・UX向上
   - プログレスバー表示
   - 詳細ログ出力（verboseモード）
   - エラーメッセージ改善

10. リトライ機能
    - 5xxエラー時の自動リトライ
    - 指数バックオフ

**成果物**:
- 本番利用可能なフル機能CLI

### Phase 4: 最適化・拡張（オプション）

11. キャッシュ機能
    - ローカルキャッシュ実装
    - キャッシュ有効期限管理

12. その他エンドポイント対応
    - `/indices/bars/daily` (指数四本値)
    - `/equities/bars/minute` (分足データ)
    - `/markets/*` (市場データ)

13. テスト充実
    - ユニットテスト
    - インテグレーションテスト
    - モックAPI対応

## セキュリティ考慮事項

### APIキー管理

- 設定ファイルのパーミッション: `0600`
- ログ出力時のAPIキーマスキング
- `.gitignore`に設定ファイルを追加

### 実装例

```go
// APIキーのマスキング
func maskAPIKey(key string) string {
    if len(key) <= 8 {
        return "****"
    }
    return key[:4] + "****" + key[len(key)-4:]
}

// ログ出力時
log.Printf("Using API Key: %s", maskAPIKey(apiKey))
```

## 参考資料

- [J-Quants API仕様書](https://jpx-jquants.com/ja/spec)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Viper Documentation](https://github.com/spf13/viper)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
