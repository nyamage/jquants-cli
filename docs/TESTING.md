# テスト戦略ドキュメント

jquants-cliプロジェクトのテスト方針と実装ガイドラインを定義します。

## 目次

1. [テスト方針](#テスト方針)
2. [テストディレクトリ構成](#テストディレクトリ構成)
3. [単体テスト](#単体テスト)
4. [結合テスト](#結合テスト)
5. [テストデータ管理](#テストデータ管理)
6. [モック戦略](#モック戦略)
7. [テストカバレッジ](#テストカバレッジ)
8. [CI/CD統合](#cicd統合)

## テスト方針

### 基本原則

1. **テストファースト**: 新しい機能を実装する前にテストを書く（可能な限り）
2. **高いカバレッジ**: 80%以上のテストカバレッジを目指す
3. **高速な実行**: 単体テストは1秒以内、全テストは30秒以内に完了
4. **独立性**: テストは他のテストに依存しない
5. **再現性**: いつ実行しても同じ結果が得られる

### テストレベル

| レベル | 目的 | 実行頻度 | 実行時間 |
|--------|------|---------|---------|
| 単体テスト | 関数・メソッドの動作確認 | コミット前 | < 1秒 |
| 結合テスト | コンポーネント間の連携確認 | PR作成時 | < 30秒 |
| E2Eテスト | CLI全体の動作確認 | リリース前 | < 5分 |

### テスト戦略マトリクス

| コンポーネント | 単体テスト | 結合テスト | E2Eテスト |
|--------------|----------|----------|----------|
| APIクライアント | ✓ | ✓ | - |
| 設定管理 | ✓ | ✓ | - |
| 出力フォーマット | ✓ | - | - |
| CLIコマンド | - | ✓ | ✓ |
| エンドツーエンド | - | - | ✓ |

## テストディレクトリ構成

```
jquants-cli/
├── internal/
│   ├── api/
│   │   ├── client.go
│   │   ├── client_test.go          # 単体テスト
│   │   ├── equities.go
│   │   ├── equities_test.go        # 単体テスト
│   │   └── testdata/               # テストデータ
│   │       ├── equity_master.json
│   │       ├── prices_daily.json
│   │       └── fins_summary.json
│   ├── config/
│   │   ├── config.go
│   │   ├── config_test.go          # 単体テスト
│   │   └── testdata/
│   │       └── config.yaml
│   └── output/
│       ├── formatter.go
│       ├── formatter_test.go       # 単体テスト
│       └── testdata/
│           └── sample_data.json
├── tests/
│   ├── integration/                # 結合テスト
│   │   ├── api_test.go
│   │   ├── cli_test.go
│   │   └── config_test.go
│   └── fixtures/                   # テストフィクスチャ
│       ├── api_responses/
│       └── config_files/
└── scripts/
    ├── test.sh                     # テスト実行スクリプト
    └── coverage.sh                 # カバレッジ計測スクリプト
```

### ディレクトリ規約

- **`*_test.go`**: 単体テストファイル（対象ファイルと同じディレクトリ）
- **`testdata/`**: 各パッケージのテストデータ（Goのビルドから除外される）
- **`tests/integration/`**: 結合テスト（複数パッケージにまたがるテスト）
- **`tests/fixtures/`**: 共通テストデータ

## 単体テスト

### 基本構造

```go
// internal/api/client_test.go
package api

import (
    "testing"
)

func TestClient_Get(t *testing.T) {
    // Arrange (準備)
    client := NewClient("test-api-key")

    // Act (実行)
    result, err := client.Get("/test", nil)

    // Assert (検証)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result == nil {
        t.Error("result should not be nil")
    }
}
```

### テーブル駆動テスト

複数のテストケースを効率的に記述:

```go
func TestMaskAPIKey(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "normal key",
            input:    "abcdefghijklmnop",
            expected: "abcd****mnop",
        },
        {
            name:     "short key",
            input:    "abc",
            expected: "****",
        },
        {
            name:     "empty key",
            input:    "",
            expected: "****",
        },
        {
            name:     "exactly 8 chars",
            input:    "abcdefgh",
            expected: "****",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := maskAPIKey(tt.input)
            if got != tt.expected {
                t.Errorf("maskAPIKey(%q) = %q, want %q",
                    tt.input, got, tt.expected)
            }
        })
    }
}
```

### サブテスト

関連するテストをグループ化:

```go
func TestEquityMasterResponse(t *testing.T) {
    t.Run("valid response", func(t *testing.T) {
        // テスト実装
    })

    t.Run("empty data", func(t *testing.T) {
        // テスト実装
    })

    t.Run("malformed json", func(t *testing.T) {
        // テスト実装
    })
}
```

### エラーケースのテスト

```go
func TestClient_Get_Errors(t *testing.T) {
    tests := []struct {
        name        string
        statusCode  int
        expectError bool
        errorMsg    string
    }{
        {
            name:        "401 Unauthorized",
            statusCode:  401,
            expectError: true,
            errorMsg:    "Invalid API key",
        },
        {
            name:        "429 Rate Limit",
            statusCode:  429,
            expectError: true,
            errorMsg:    "Rate limit exceeded",
        },
        {
            name:        "500 Server Error",
            statusCode:  500,
            expectError: true,
            errorMsg:    "server error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // モックサーバーで指定のステータスコードを返す
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.statusCode)
            }))
            defer server.Close()

            client := NewClient("test-key")
            client.BaseURL = server.URL

            _, err := client.Get("/test", nil)

            if tt.expectError {
                if err == nil {
                    t.Fatal("expected error, got nil")
                }
                if !strings.Contains(err.Error(), tt.errorMsg) {
                    t.Errorf("error message %q should contain %q", err.Error(), tt.errorMsg)
                }
            } else {
                if err != nil {
                    t.Fatalf("unexpected error: %v", err)
                }
            }
        })
    }
}
```

### テストデータの使用

```go
func TestParseEquityMaster(t *testing.T) {
    // testdataディレクトリからJSONを読み込み
    data, err := os.ReadFile("testdata/equity_master.json")
    if err != nil {
        t.Fatalf("failed to read test data: %v", err)
    }

    var resp EquityMasterResponse
    if err := json.Unmarshal(data, &resp); err != nil {
        t.Fatalf("failed to unmarshal: %v", err)
    }

    if len(resp.Data) == 0 {
        t.Error("expected data, got empty")
    }

    // 最初のレコードを検証
    first := resp.Data[0]
    if first.Code != "86970" {
        t.Errorf("expected code 86970, got %s", first.Code)
    }
}
```

## 結合テスト

### APIクライアント結合テスト

複数のコンポーネントを組み合わせてテスト:

```go
// tests/integration/api_test.go
package integration

import (
    "testing"
    "github.com/nyamage/jquants-cli/internal/api"
    "github.com/nyamage/jquants-cli/internal/config"
)

func TestAPIClient_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // 設定を読み込み
    cfg, err := config.Load()
    if err != nil {
        t.Fatalf("failed to load config: %v", err)
    }

    if cfg.APIKey == "" {
        t.Skip("API key not set, skipping integration test")
    }

    // 実際のAPIクライアントを作成
    client := api.NewClient(cfg.APIKey)

    t.Run("GetEquityMaster", func(t *testing.T) {
        data, err := client.GetEquityMaster("86970", "")
        if err != nil {
            t.Fatalf("failed to get equity master: %v", err)
        }

        if len(data) == 0 {
            t.Error("expected data, got empty")
        }

        // データの妥当性を検証
        if data[0].Code != "86970" {
            t.Errorf("expected code 86970, got %s", data[0].Code)
        }
    })
}
```

### CLIコマンド結合テスト

```go
// tests/integration/cli_test.go
package integration

import (
    "bytes"
    "os/exec"
    "strings"
    "testing"
)

func TestCLI_EquitiesList(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // CLIをビルド
    buildCmd := exec.Command("go", "build", "-o", "jquants-test", "./cmd/jquants")
    if err := buildCmd.Run(); err != nil {
        t.Fatalf("failed to build CLI: %v", err)
    }
    defer os.Remove("jquants-test")

    tests := []struct {
        name    string
        args    []string
        wantErr bool
        contain string
    }{
        {
            name:    "list all equities",
            args:    []string{"equities", "list"},
            wantErr: false,
            contain: "CODE",
        },
        {
            name:    "get specific equity",
            args:    []string{"equities", "get", "86970"},
            wantErr: false,
            contain: "86970",
        },
        {
            name:    "json output",
            args:    []string{"equities", "list", "--output", "json"},
            wantErr: false,
            contain: `"data"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var stdout, stderr bytes.Buffer
            cmd := exec.Command("./jquants-test", tt.args...)
            cmd.Stdout = &stdout
            cmd.Stderr = &stderr

            err := cmd.Run()

            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
            } else {
                if err != nil {
                    t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr.String())
                }
            }

            if tt.contain != "" && !strings.Contains(stdout.String(), tt.contain) {
                t.Errorf("output should contain %q, got:\n%s", tt.contain, stdout.String())
            }
        })
    }
}
```

### ページネーション結合テスト

```go
func TestPagination_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    client := api.NewClient(getAPIKey(t))

    t.Run("fetch all pages", func(t *testing.T) {
        allData := []api.PriceDaily{}
        params := url.Values{}
        params.Set("date", "2025-03-01")

        err := client.FetchAllPages("/equities/bars/daily", params, func(data []byte) error {
            var resp api.PriceDailyResponse
            if err := json.Unmarshal(data, &resp); err != nil {
                return err
            }
            allData = append(allData, resp.Data...)
            return nil
        })

        if err != nil {
            t.Fatalf("failed to fetch all pages: %v", err)
        }

        if len(allData) == 0 {
            t.Error("expected data, got empty")
        }

        t.Logf("fetched %d records", len(allData))
    })
}
```

## テストデータ管理

### testdataディレクトリ

各パッケージの`testdata/`ディレクトリにサンプルデータを配置:

```
internal/api/testdata/
├── equity_master.json          # 銘柄マスターのサンプル
├── prices_daily.json           # 株価データのサンプル
├── prices_daily_paginated.json # ページネーション付き
├── fins_summary.json           # 財務情報のサンプル
├── bulk_list.json              # Bulkファイル一覧
└── error_responses.json        # エラーレスポンス集
```

### サンプルデータの作成

```json
// internal/api/testdata/equity_master.json
{
  "data": [
    {
      "Date": "2025-03-14",
      "Code": "86970",
      "CoName": "日本取引所グループ",
      "CoNameEn": "Japan Exchange Group,Inc.",
      "S17": "16",
      "S17Nm": "金融（除く銀行）",
      "S33": "7200",
      "S33Nm": "その他金融業",
      "ScaleCat": "TOPIX Large70",
      "Mkt": "0111",
      "MktNm": "プライム",
      "Mrgn": "1",
      "MrgnNm": "信用"
    }
  ]
}
```

### テストヘルパー関数

共通処理をヘルパー関数化:

```go
// internal/api/testhelper_test.go
package api

import (
    "os"
    "path/filepath"
    "testing"
)

// loadTestData はtestdataディレクトリからファイルを読み込む
func loadTestData(t *testing.T, filename string) []byte {
    t.Helper()

    path := filepath.Join("testdata", filename)
    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("failed to read test data %s: %v", filename, err)
    }

    return data
}

// setupTestClient はテスト用のクライアントをセットアップ
func setupTestClient(t *testing.T) *Client {
    t.Helper()
    return NewClient("test-api-key")
}
```

## モック戦略

### httptestによるモックサーバー

```go
func TestClient_GetEquityMaster(t *testing.T) {
    // モックサーバーを起動
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // リクエストを検証
        if r.Header.Get("x-api-key") != "test-key" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // テストデータを返す
        w.Header().Set("Content-Type", "application/json")
        data := loadTestData(t, "equity_master.json")
        w.Write(data)
    }))
    defer server.Close()

    // クライアントにモックサーバーのURLを設定
    client := NewClient("test-key")
    client.BaseURL = server.URL

    // テスト実行
    result, err := client.GetEquityMaster("86970", "")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(result) == 0 {
        t.Error("expected data, got empty")
    }
}
```

### インターフェースによるモック

依存関係の注入を使用:

```go
// internal/api/client.go
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type Client struct {
    BaseURL    string
    APIKey     string
    HTTPClient HTTPClient // インターフェース
}

// internal/api/client_test.go
type mockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}

func TestClient_WithMock(t *testing.T) {
    mockClient := &mockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            // カスタムレスポンスを返す
            return &http.Response{
                StatusCode: 200,
                Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
            }, nil
        },
    }

    client := &Client{
        BaseURL:    "https://api.jquants.com/v2",
        APIKey:     "test-key",
        HTTPClient: mockClient,
    }

    // テスト実行
}
```

### モックファイルシステム

設定ファイルテスト用:

```go
func TestConfig_Load(t *testing.T) {
    // 一時ディレクトリを作成
    tmpDir := t.TempDir()

    // テスト用設定ファイルを作成
    configPath := filepath.Join(tmpDir, "config.yaml")
    configData := []byte("api_key: test-key-12345\n")
    if err := os.WriteFile(configPath, configData, 0600); err != nil {
        t.Fatalf("failed to write config: %v", err)
    }

    // 環境変数を一時的に変更
    t.Setenv("HOME", tmpDir)

    // 設定を読み込み
    cfg, err := config.Load()
    if err != nil {
        t.Fatalf("failed to load config: %v", err)
    }

    if cfg.APIKey != "test-key-12345" {
        t.Errorf("expected api_key test-key-12345, got %s", cfg.APIKey)
    }
}
```

## テストカバレッジ

### カバレッジ目標

| コンポーネント | 目標カバレッジ |
|--------------|--------------|
| APIクライアント | 85%以上 |
| 設定管理 | 80%以上 |
| 出力フォーマット | 90%以上 |
| ユーティリティ | 80%以上 |
| 全体 | 80%以上 |

### カバレッジ測定

```bash
# 全パッケージのカバレッジを測定
go test -coverprofile=coverage.out ./...

# カバレッジをブラウザで表示
go tool cover -html=coverage.out

# カバレッジ率を表示
go tool cover -func=coverage.out

# 80%未満のパッケージを検出
go tool cover -func=coverage.out | grep -v "100.0%" | grep -v "^total"
```

### カバレッジレポートスクリプト

```bash
#!/bin/bash
# scripts/coverage.sh

set -e

echo "Running tests with coverage..."
go test -coverprofile=coverage.out -covermode=atomic ./...

echo ""
echo "Coverage summary:"
go tool cover -func=coverage.out | tail -1

echo ""
echo "Packages with coverage < 80%:"
go tool cover -func=coverage.out | awk '$3 < 80.0 {print $1, $3}'

# カバレッジが80%未満の場合は失敗
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "ERROR: Total coverage $COVERAGE% is below 80%"
    exit 1
fi

echo "Coverage check passed: $COVERAGE%"
```

## CI/CD統合

### GitHub Actions設定例

```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'

    - name: Run unit tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Check coverage
      run: |
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        echo "Total coverage: $COVERAGE%"
        if (( $(echo "$COVERAGE < 80" | bc -l) )); then
          echo "ERROR: Coverage $COVERAGE% is below 80%"
          exit 1
        fi

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v4
      with:
        file: ./coverage.out
        flags: unittests
        name: codecov-umbrella

    - name: Run integration tests
      run: go test -v ./tests/integration/...
      env:
        JQUANTS_API_KEY: ${{ secrets.JQUANTS_API_KEY }}

  lint:
    name: Lint
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'

    - name: golangci-lint
      uses: golangci/golangci-lint-action@v4
      with:
        version: latest
```

### テスト実行スクリプト

```bash
#!/bin/bash
# scripts/test.sh

set -e

echo "=== Running Unit Tests ==="
go test -v -race -short ./...

echo ""
echo "=== Running Integration Tests ==="
go test -v ./tests/integration/...

echo ""
echo "=== Running Coverage Check ==="
./scripts/coverage.sh

echo ""
echo "✓ All tests passed!"
```

## ベストプラクティス

### DO（推奨）

- ✓ テストは独立させる（他のテストに依存しない）
- ✓ テストデータは`testdata/`ディレクトリに配置
- ✓ エラーケースを必ずテストする
- ✓ テーブル駆動テストを活用する
- ✓ `t.Helper()`を使ってヘルパー関数を作る
- ✓ `t.Cleanup()`でリソースを確実にクリーンアップ
- ✓ `testing.Short()`で時間のかかるテストをスキップ可能にする

### DON'T（非推奨）

- ✗ テスト間で状態を共有しない
- ✗ 外部APIに直接依存しない（モックを使う）
- ✗ テストでsleep()を使わない（タイムアウトを使う）
- ✗ グローバル変数を変更しない
- ✗ テストに副作用を残さない

### テスト命名規約

```go
// 良い例
func TestClient_GetEquityMaster_Success(t *testing.T) {}
func TestClient_GetEquityMaster_InvalidAPIKey(t *testing.T) {}
func TestConfig_Load_FileNotFound(t *testing.T) {}

// 悪い例
func TestClient(t *testing.T) {}
func Test1(t *testing.T) {}
func TestStuff(t *testing.T) {}
```

パターン: `Test<Type>_<Method>_<Scenario>`

## まとめ

### テスト実装チェックリスト

新機能を実装する際は以下を確認:

- [ ] 単体テストを作成（カバレッジ80%以上）
- [ ] エラーケースのテストを作成
- [ ] 必要に応じて結合テストを作成
- [ ] テストデータを`testdata/`に配置
- [ ] `go test ./...`が通ることを確認
- [ ] カバレッジが基準を満たすことを確認
- [ ] テストが独立して実行できることを確認
- [ ] CIで自動実行されることを確認

### 参考資料

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Testify - Testing toolkit](https://github.com/stretchr/testify)
- [httptest - HTTP testing utilities](https://golang.org/pkg/net/http/httptest/)
