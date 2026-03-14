# jquants-cli

J-Quants API用のコマンドラインツール（Go実装）

## 概要

jquants-cliは、[J-Quants API](https://jpx-jquants.com/)を利用して日本株の株価データ、財務情報、銘柄情報を取得するためのコマンドラインツールです。

## 特徴

- シンプルで直感的なCLIインターフェース
- 複数の出力フォーマット対応（テーブル、JSON、CSV）
- ページネーション自動処理
- Bulk APIによる大量データ一括取得
- 設定ファイル/環境変数によるAPIキー管理

## インストール

```bash
go install github.com/nyamage/jquants-cli/cmd/jquants@latest
```

または、ソースからビルド:

```bash
git clone https://github.com/nyamage/jquants-cli.git
cd jquants-cli
go build -o jquants ./cmd/jquants
```

## 事前準備

### APIキーの取得

1. [J-Quants](https://jpx-jquants.com/)でアカウント作成
2. ダッシュボードからAPIキーを取得

### APIキーの設定

方法1: 設定ファイル

```bash
jquants config set-api-key YOUR_API_KEY
```

方法2: 環境変数

```bash
export JQUANTS_API_KEY=YOUR_API_KEY
```

方法3: コマンドラインフラグ

```bash
jquants --api-key YOUR_API_KEY <command>
```

## 使い方

### 銘柄情報の取得

```bash
# 全銘柄一覧を取得
jquants equities list

# 特定日付の銘柄情報を取得
jquants equities list --date 2025-03-01

# 特定銘柄の情報を取得
jquants equities get 86970
```

### 株価データの取得

```bash
# 特定銘柄の株価を全期間取得
jquants prices daily --code 86970

# 期間を指定して取得
jquants prices daily --code 86970 --from 2025-01-01 --to 2025-03-01

# 特定日付の全銘柄株価を取得
jquants prices daily --date 2025-03-14
```

### 財務情報の取得

```bash
# 特定銘柄の財務情報を取得
jquants fins get --code 86970

# 特定日付に開示された財務情報を取得
jquants fins list --date 2025-03-01
```

### Bulk APIでデータ一括取得

```bash
# 利用可能なファイル一覧を取得
jquants bulk list --endpoint /equities/bars/daily

# 期間を指定してファイル一覧を取得
jquants bulk list --endpoint /equities/bars/daily --from 2025-01 --to 2025-03

# ファイルをダウンロード
jquants bulk download --endpoint /equities/bars/daily --date 2025-01 --output-dir ./data
```

### 出力フォーマットの指定

```bash
# テーブル形式（デフォルト）
jquants equities list

# JSON形式
jquants equities list --output json

# CSV形式
jquants equities list --output csv > equities.csv
```

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| `jquants config set-api-key` | APIキーを設定 |
| `jquants config show` | 現在の設定を表示 |
| `jquants equities list` | 銘柄一覧を取得 |
| `jquants equities get` | 特定銘柄の情報を取得 |
| `jquants prices daily` | 株価四本値（日次）を取得 |
| `jquants fins get` | 財務情報を取得 |
| `jquants fins list` | 財務情報一覧を取得 |
| `jquants bulk list` | Bulk APIファイル一覧を取得 |
| `jquants bulk download` | Bulk APIファイルをダウンロード |

## グローバルフラグ

| フラグ | 短縮形 | 説明 | デフォルト |
|-------|-------|------|-----------|
| `--api-key` | | APIキー | 環境変数 `JQUANTS_API_KEY` |
| `--output` | `-o` | 出力形式 (table/json/csv) | table |
| `--verbose` | `-v` | 詳細ログ出力 | false |

## プロジェクト構成

```
jquants-cli/
├── cmd/
│   └── jquants/          # メインCLI
├── internal/
│   ├── api/              # APIクライアント
│   ├── config/           # 設定管理
│   └── output/           # 出力フォーマット
├── docs/                 # ドキュメント
├── go.mod
├── go.sum
└── README.md
```

詳細な設計方針については [docs/DESIGN.md](docs/DESIGN.md) を参照してください。

## 開発

### 開発環境のセットアップ

```bash
# 依存パッケージと開発ツールをインストール
make dev

# または手動でインストール
go mod download
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
```

### Makefileコマンド

このプロジェクトではMakefileを提供しています。利用可能なコマンド:

```bash
# ヘルプを表示
make help

# ビルド
make build

# テスト実行
make test              # 全テスト実行
make test-unit         # 単体テストのみ
make test-integration  # 結合テストのみ

# カバレッジ測定
make coverage          # カバレッジ測定
make coverage-html     # HTMLレポート表示

# コード品質チェック
make fmt               # コードフォーマット
make vet               # go vet実行
make lint              # golangci-lint実行
make check             # fmt + vet + lint

# クリーンアップ
make clean             # 全クリーンアップ
make clean-build       # ビルド成果物のみ削除

# その他
make install           # バイナリをインストール
make run               # ビルド後に実行
make pre-commit        # コミット前チェック
make ci                # CI環境用タスク
```

### 開発ワークフロー

#### 新機能開発時

```bash
# 1. 開発環境セットアップ（初回のみ）
make dev

# 2. コード実装

# 3. フォーマットとテスト
make fmt
make test-unit

# 4. コミット前の最終チェック
make pre-commit
```

#### Pull Request作成前

```bash
# すべてのチェックを実行
make ci
```

### 依存ライブラリ

- [spf13/cobra](https://github.com/spf13/cobra) - CLIフレームワーク
- [spf13/viper](https://github.com/spf13/viper) - 設定管理
- [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) - テーブル表示

### 手動ビルド・テスト

Makefileを使わない場合:

```bash
# ビルド
go build -o jquants ./cmd/jquants

# テスト
go test ./...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## ライセンス

MIT License

## 関連リンク

- [J-Quants API仕様書](https://jpx-jquants.com/ja/spec)
- [J-Quantsヘルプページ](https://jpx-jquants.com/ja/help)

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。
