# Claude向け開発ガイドライン

このドキュメントは、Claude（AI）がjquants-cliプロジェクトの開発を支援する際の指針を示します。

## 基本方針

### 必読ドキュメント

このプロジェクトの開発を行う際は、**必ず**以下のドキュメントに従ってください:

1. **[docs/DESIGN.md](docs/DESIGN.md)** - 詳細設計ドキュメント
   - アーキテクチャ設計
   - ディレクトリ構成と各層の責務
   - 実装パターンとサンプルコード
   - 実装フェーズ

2. **[docs/API.md](docs/API.md)** - J-Quants API仕様
   - エンドポイント仕様
   - リクエスト/レスポンス形式
   - エラーハンドリング
   - ページネーション処理

3. **[README.md](README.md)** - プロジェクト概要
   - ユーザー向けドキュメント
   - CLIの使用方法
   - コマンド仕様

### ドキュメント優先原則

- コードを書く前に、該当するドキュメントを参照すること
- ドキュメントと実装が矛盾する場合は、ユーザーに確認すること
- 新しい機能を追加する場合は、ドキュメントも更新すること

## 実装ガイドライン

### 1. 実装フェーズの遵守

**DESIGN.md**で定義された実装フェーズに従って開発を進めてください:

- **Phase 1**: 基本機能（MVP） - 銘柄情報取得
- **Phase 2**: データ取得機能拡充 - 株価・財務データ
- **Phase 3**: 高度な機能 - Bulk API、UX向上
- **Phase 4**: 最適化・拡張

**重要**: 前のフェーズが完了するまで、次のフェーズに進まないこと。

### 2. ディレクトリ構成の厳守

```
cmd/jquants/     - CLIコマンド定義のみ
internal/api/    - J-Quants API通信ロジック
internal/config/ - 設定ファイル管理
internal/output/ - 出力フォーマット処理
```

各ディレクトリの責務を超えたコードは書かないこと。

### 3. Makefileの使用（必須）

このプロジェクトでは、開発タスクを標準化するためにMakefileを提供しています。**必ずMakefileを使用してください**。

#### 主要コマンド

```bash
# ヘルプ表示
make help

# ビルドとテスト
make build             # バイナリをビルド
make test              # 全テスト実行
make test-unit         # 単体テストのみ
make coverage          # カバレッジ測定

# コード品質
make fmt               # コードフォーマット
make vet               # go vet実行
make lint              # golangci-lint実行
make check             # fmt + vet + lint

# 開発ワークフロー
make pre-commit        # コミット前チェック
make ci                # CI環境タスク

# その他
make clean             # クリーンアップ
make dev               # 開発環境セットアップ
```

#### 開発フローでの使用

**実装開始時**:
```bash
# 開発ツールのインストール（初回のみ）
make dev
```

**コード実装中**:
```bash
# フォーマット適用
make fmt

# 単体テスト実行
make test-unit
```

**コミット前**:
```bash
# 必須: コミット前チェック
make pre-commit
```

**PR作成前**:
```bash
# 必須: すべてのチェックを実行
make ci
```

#### Makefileを使うべき理由

1. **一貫性**: すべての開発者が同じコマンドを使用
2. **効率性**: 複雑なコマンドを短いエイリアスで実行
3. **品質保証**: 必要なチェックを自動実行
4. **CI/CDとの整合性**: CIで実行されるのと同じタスクをローカルで実行

#### 直接goコマンドを使わない

```bash
# Bad: 直接goコマンドを使う
go test ./...
go build ./cmd/jquants

# Good: Makefileを使う
make test
make build
```

例外: デバッグや特殊なフラグが必要な場合のみ、直接goコマンドを使用可能。

### 4. コーディング規約

#### Go言語標準に従う

```go
// Good: Goの命名規則に従う
type APIClient struct {
    BaseURL string
    APIKey  string
}

func (c *APIClient) GetEquityMaster(code string) (*EquityMaster, error) {
    // ...
}

// Bad: 他言語の命名規則を持ち込まない
type apiClient struct {
    base_url string
    api_key  string
}

func (c *apiClient) get_equity_master(code string) (*EquityMaster, error) {
    // ...
}
```

#### エラーハンドリング

```go
// Good: エラーは必ず処理する
data, err := client.GetEquityMaster("86970")
if err != nil {
    return fmt.Errorf("failed to get equity master: %w", err)
}

// Bad: エラーを無視しない
data, _ := client.GetEquityMaster("86970")
```

#### コンテキストの活用

```go
// Good: 長時間処理にはcontextを渡す
func (c *Client) GetWithContext(ctx context.Context, endpoint string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+endpoint, nil)
    // ...
}

// Bad: contextなしで長時間処理
func (c *Client) Get(endpoint string) ([]byte, error) {
    req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
    // ...
}
```

### 5. セキュリティチェックリスト

実装時に必ず以下を確認:

- [ ] APIキーをログに出力していないか
- [ ] APIキーをエラーメッセージに含めていないか
- [ ] 設定ファイルのパーミッションは適切か（0600）
- [ ] 環境変数から機密情報を取得する際、デフォルト値を設定していないか
- [ ] ファイルパスの検証を行っているか（パストラバーサル対策）
- [ ] 外部入力のバリデーションを行っているか

```go
// Good: APIキーをマスクする
func maskAPIKey(key string) string {
    if len(key) <= 8 {
        return "****"
    }
    return key[:4] + "****" + key[len(key)-4:]
}

log.Printf("Using API Key: %s", maskAPIKey(apiKey))

// Bad: APIキーをそのまま出力
log.Printf("Using API Key: %s", apiKey)
```

### 6. エラーメッセージの方針

#### ユーザーフレンドリーなメッセージ

```go
// Good: 具体的で解決方法を示す
if apiKey == "" {
    return errors.New("API key is not set. Please set it via 'jquants config set-api-key' or environment variable JQUANTS_API_KEY")
}

// Bad: 曖昧で不親切
if apiKey == "" {
    return errors.New("invalid configuration")
}
```

#### エラーのラップ

```go
// Good: エラーチェーンを保持
resp, err := client.Get("/equities/master", params)
if err != nil {
    return nil, fmt.Errorf("failed to fetch equity master data: %w", err)
}

// Bad: エラー情報を失う
resp, err := client.Get("/equities/master", params)
if err != nil {
    return nil, errors.New("API error")
}
```

### 7. テストの方針

**重要**: テスト戦略の詳細は [docs/TESTING.md](docs/TESTING.md) を必ず参照してください。

#### Makefileでのテスト実行

テストは必ずMakefileを経由して実行:

```bash
# 推奨: Makefileを使用
make test-unit         # 単体テスト
make coverage          # カバレッジ測定
make pre-commit        # コミット前チェック

# 非推奨: 直接goコマンド（特殊な場合のみ）
go test -short ./...
```

#### テスト必須化ポリシー

すべての新機能・修正には以下のテストが必要:

- [ ] 単体テスト（カバレッジ80%以上）
- [ ] エラーケースのテスト
- [ ] 必要に応じて結合テスト

#### テストレベル別の実装義務

| レベル | 実装タイミング | 必須度 |
|--------|--------------|-------|
| 単体テスト | コード実装と同時 | 必須 |
| 結合テスト | コンポーネント統合時 | 推奨 |
| E2Eテスト | Phase 3以降 | 推奨 |

#### 単体テストの実装パターン

```go
// テーブル駆動テストを推奨
func TestMaskAPIKey(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"normal key", "abcdefghijklmnop", "abcd****mnop"},
        {"short key", "abc", "****"},
        {"empty key", "", "****"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := maskAPIKey(tt.input)
            if got != tt.expected {
                t.Errorf("maskAPIKey(%q) = %q, want %q", tt.input, got, tt.expected)
            }
        })
    }
}
```

#### モック戦略

httptestを使ったモックサーバー:

```go
func TestClient_GetEquityMaster(t *testing.T) {
    // モックサーバーを起動
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // リクエスト検証
        if r.Header.Get("x-api-key") != "test-key" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // テストデータを返す
        w.Header().Set("Content-Type", "application/json")
        data, _ := os.ReadFile("testdata/equity_master.json")
        w.Write(data)
    }))
    defer server.Close()

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

#### テストデータ管理

- テストデータは`testdata/`ディレクトリに配置
- JSONファイルは実際のAPIレスポンスを参考に作成
- 個人情報や機密情報は含めない

```
internal/api/testdata/
├── equity_master.json
├── prices_daily.json
├── fins_summary.json
└── error_responses.json
```

#### 結合テストのスキップ制御

```go
func TestAPIClient_Integration(t *testing.T) {
    // -short フラグで結合テストをスキップ
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // APIキーが設定されていない場合もスキップ
    apiKey := os.Getenv("JQUANTS_API_KEY")
    if apiKey == "" {
        t.Skip("JQUANTS_API_KEY not set")
    }

    // 実際のAPIを使ったテスト
    client := api.NewClient(apiKey)
    // ...
}
```

#### テスト実行コマンド

```bash
# 単体テストのみ（高速）
go test -short ./...

# 全テスト実行
go test ./...

# カバレッジ測定
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 特定パッケージのテスト
go test ./internal/api/...

# verbose出力
go test -v ./...

# 並列実行
go test -parallel 4 ./...
```

#### テスト作成チェックリスト

新しい関数・メソッドを実装したら:

- [ ] 正常系のテストを作成
- [ ] 異常系（エラーケース）のテストを作成
- [ ] 境界値のテストを作成
- [ ] テストデータを`testdata/`に配置
- [ ] `go test ./...`が通ることを確認
- [ ] カバレッジが80%以上であることを確認
- [ ] モックサーバーを適切に使用

#### テスト禁止事項

- ✗ 実際のAPIに直接依存するテスト（結合テスト除く）
- ✗ テスト間で状態を共有
- ✗ sleepを使った時間依存のテスト
- ✗ グローバル変数の変更
- ✗ テストの実行順序に依存
- ✗ ハードコードされたAPIキーやシークレット

### 8. パフォーマンス考慮事項

#### レート制限対策

```go
// Good: リクエスト間に待機時間を設ける
for _, code := range codes {
    data, err := client.GetEquityMaster(code)
    // 処理...
    time.Sleep(500 * time.Millisecond) // レート制限対策
}

// Bad: 連続リクエスト
for _, code := range codes {
    data, err := client.GetEquityMaster(code)
    // 処理...
}
```

#### 大量データ処理

```go
// Good: ストリーミング処理
func processLargeDataset(data []Record) error {
    for i := 0; i < len(data); i += 1000 {
        end := i + 1000
        if end > len(data) {
            end = len(data)
        }
        batch := data[i:end]
        if err := processBatch(batch); err != nil {
            return err
        }
    }
    return nil
}

// Bad: 全データをメモリに展開
func processLargeDataset(data []Record) error {
    allResults := make([]Result, len(data))
    // メモリ使用量が膨大に...
}
```

### 9. ログ出力の方針

#### レベル別のログ出力

```go
// verboseフラグがtrueの時のみ詳細ログを出力
if verbose {
    log.Printf("Fetching data from %s with params %v", endpoint, params)
}

// エラーは常に出力
if err != nil {
    log.Printf("Error: %v", err)
    return err
}
```

#### 構造化ログ（Phase 4以降で検討）

```go
// 将来的には構造化ログライブラリの導入を検討
log.Info("API request completed",
    "endpoint", endpoint,
    "duration", time.Since(start),
    "status", resp.StatusCode,
)
```

### 10. ドキュメント更新の義務

コードを変更した場合、以下を必ず更新:

1. **新しいコマンド追加時**:
   - `README.md` のコマンド一覧
   - `docs/DESIGN.md` のコマンド設計セクション

2. **新しいエンドポイント対応時**:
   - `docs/API.md` にエンドポイント仕様を追加

3. **設定項目追加時**:
   - `README.md` の設定方法セクション
   - `docs/DESIGN.md` の設定管理セクション

4. **出力フォーマット変更時**:
   - `README.md` の使用例
   - `docs/DESIGN.md` の出力フォーマットセクション

### 11. コミットメッセージ規約

```
<type>: <subject>

<body>
```

**Type**:
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `refactor`: リファクタリング
- `test`: テスト追加・修正
- `chore`: ビルドプロセスやツールの変更

**例**:
```
feat: add equities list command

- Implement /equities/master endpoint client
- Add table format output
- Add basic error handling

Refs: Phase 1 - Step 4
```

### 12. 実装時の質問・確認事項

実装中に以下の状況に遭遇した場合、ユーザーに確認すること:

1. **ドキュメントに記載のない機能を追加する場合**
   - 「この機能はDESIGN.mdに記載されていませんが、追加しますか？」

2. **ドキュメントと実装が矛盾する場合**
   - 「DESIGN.mdでは〇〇となっていますが、実装上は△△が必要です。どちらに従いますか？」

3. **複数の実装方法がある場合**
   - 「AとBの実装方法があります。どちらを採用しますか？（理由：...）」

4. **セキュリティリスクがある場合**
   - 「この実装はセキュリティリスクがあります: ... 代替案: ...」

### 13. デバッグ・トラブルシューティング

#### デバッグ情報の出力

```go
// verboseモードでのデバッグ情報
if verbose {
    log.Printf("DEBUG: Request URL: %s", req.URL.String())
    log.Printf("DEBUG: Request Headers: %v", req.Header)
    log.Printf("DEBUG: Response Status: %d", resp.StatusCode)
}
```

#### エラー再現手順の記録

エラーを修正する際は、再現手順をコメントに記載:

```go
// Fix: Panic when pagination_key is null
// Reproduce: jquants prices daily --date 2025-03-14
// Root cause: json.Unmarshal sets empty string instead of nil
if resp.PaginationKey == "" || resp.PaginationKey == "null" {
    break
}
```

### 14. 依存ライブラリの管理

#### 新しいライブラリを追加する前に

1. 標準ライブラリで実現できないか検討
2. ライブラリのメンテナンス状況を確認
3. ユーザーに追加の許可を得る

#### 推奨ライブラリ（DESIGN.mdに記載）

- CLI: `spf13/cobra`
- 設定: `spf13/viper`
- テーブル表示: `olekukonko/tablewriter`

これら以外のライブラリを追加する場合は、ユーザーに確認すること。

### 15. パフォーマンステスト

Phase 3以降で以下を確認:

- [ ] 1000件以上のデータ取得が正常に動作するか
- [ ] メモリリークがないか
- [ ] ページネーション処理が効率的か
- [ ] Bulk API処理が高速か

### 16. リリース前チェックリスト

リリース前に以下を確認:

- [ ] すべてのテストが通過
- [ ] ドキュメントが最新
- [ ] README.mdの使用例が動作する
- [ ] `.gitignore`に機密情報が含まれていない
- [ ] `go mod tidy`実行済み
- [ ] ビルドが成功する（`go build ./cmd/jquants`）
- [ ] 異なるOS（Linux/macOS/Windows）でビルド可能か確認

## プロジェクト固有の重要事項

### J-Quants API の特性

1. **ページネーション**:
   - 大量データ取得時は必ずページネーション処理を実装
   - `pagination_key`が空になるまでループ

2. **レート制限**:
   - リクエスト間に500ms待機
   - 429エラー時は適切にリトライ

3. **Bulk API**:
   - 大量データ取得にはBulk APIを推奨
   - gzip解凍処理を忘れずに

4. **調整済み株価**:
   - `Adj*`フィールドが調整済みデータ
   - ユーザーには調整済みを推奨

### プロジェクトの哲学

1. **シンプル第一**: 複雑な実装より、シンプルで理解しやすい実装を優先
2. **ユーザーフレンドリー**: エラーメッセージは親切に、具体的に
3. **拡張性**: 新しいエンドポイント追加が容易な設計
4. **堅牢性**: エラーハンドリングとリトライ機能を重視

## まとめ

このプロジェクトに貢献する際は:

1. **DESIGN.md、API.md、README.mdを熟読**すること
2. **実装フェーズを守る**こと
3. **セキュリティを最優先**すること
4. **ドキュメントを常に更新**すること
5. **不明点はユーザーに確認**すること

Happy coding!
