.PHONY: help build test test-unit test-integration coverage lint fmt clean install deps run

# 変数定義
BINARY_NAME=jquants
BINARY_PATH=./cmd/jquants
COVERAGE_FILE=coverage.out
GO=go
GOFLAGS=-v

# デフォルトターゲット
.DEFAULT_GOAL := help

## help: Makefileのヘルプを表示
help:
	@echo "使用可能なターゲット:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
	@echo ""

## deps: 依存パッケージをインストール
deps:
	@echo "依存パッケージをインストール中..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✓ 依存パッケージのインストール完了"

## build: バイナリをビルド
build:
	@echo "バイナリをビルド中..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(BINARY_PATH)
	@echo "✓ ビルド完了: $(BINARY_NAME)"

## install: バイナリをインストール
install:
	@echo "バイナリをインストール中..."
	$(GO) install $(GOFLAGS) $(BINARY_PATH)
	@echo "✓ インストール完了"

## run: アプリケーションを実行
run: build
	@echo "アプリケーションを実行中..."
	./$(BINARY_NAME)

## test: 全テストを実行
test:
	@echo "全テストを実行中..."
	@./scripts/test.sh

## test-unit: 単体テストのみ実行
test-unit:
	@echo "単体テストを実行中..."
	$(GO) test -short -race $(GOFLAGS) ./...
	@echo "✓ 単体テスト完了"

## test-integration: 結合テストのみ実行
test-integration:
	@echo "結合テストを実行中..."
	@if [ -z "$$JQUANTS_API_KEY" ]; then \
		echo "警告: JQUANTS_API_KEY が設定されていません"; \
		echo "結合テストをスキップします"; \
	else \
		$(GO) test $(GOFLAGS) ./tests/integration/...; \
		echo "✓ 結合テスト完了"; \
	fi

## coverage: テストカバレッジを測定
coverage:
	@echo "テストカバレッジを測定中..."
	@./scripts/coverage.sh
	@echo ""
	@echo "詳細なカバレッジレポートを表示する場合:"
	@echo "  make coverage-html"

## coverage-html: カバレッジレポートをHTMLで表示
coverage-html:
	@if [ -f $(COVERAGE_FILE) ]; then \
		echo "ブラウザでカバレッジレポートを開いています..."; \
		$(GO) tool cover -html=$(COVERAGE_FILE); \
	else \
		echo "エラー: $(COVERAGE_FILE) が見つかりません"; \
		echo "先に 'make coverage' を実行してください"; \
		exit 1; \
	fi

## lint: コードをlintチェック
lint:
	@echo "コードをlintチェック中..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
		echo "✓ lintチェック完了"; \
	else \
		echo "エラー: golangci-lint がインストールされていません"; \
		echo "インストール方法: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

## fmt: コードをフォーマット
fmt:
	@echo "コードをフォーマット中..."
	$(GO) fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi
	@echo "✓ フォーマット完了"

## vet: go vetを実行
vet:
	@echo "go vetを実行中..."
	$(GO) vet ./...
	@echo "✓ go vet完了"

## clean: ビルド成果物とキャッシュを削除
clean:
	@echo "クリーンアップ中..."
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE)
	@rm -rf dist/
	@$(GO) clean -cache -testcache -modcache
	@echo "✓ クリーンアップ完了"

## clean-build: ビルド成果物のみ削除
clean-build:
	@echo "ビルド成果物を削除中..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@echo "✓ ビルド成果物の削除完了"

## check: フォーマット、vet、lintを実行
check: fmt vet lint
	@echo "✓ すべてのチェック完了"

## ci: CI環境で実行するタスク
ci: deps check test coverage
	@echo "✓ CIタスク完了"

## pre-commit: コミット前に実行すべきタスク
pre-commit: fmt vet test-unit
	@echo "✓ コミット前チェック完了"

## release-dry-run: リリースのドライラン（実際にはリリースしない）
release-dry-run:
	@echo "リリースのドライランを実行中..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean --skip=publish; \
		echo "✓ リリースドライラン完了"; \
	else \
		echo "エラー: goreleaser がインストールされていません"; \
		echo "インストール: go install github.com/goreleaser/goreleaser@latest"; \
		exit 1; \
	fi

## mod-tidy: go.modをクリーンアップ
mod-tidy:
	@echo "go.modをクリーンアップ中..."
	$(GO) mod tidy
	@echo "✓ go.modクリーンアップ完了"

## mod-verify: 依存関係を検証
mod-verify:
	@echo "依存関係を検証中..."
	$(GO) mod verify
	@echo "✓ 依存関係の検証完了"

## version: バージョン情報を表示
version:
	@if [ -f $(BINARY_NAME) ]; then \
		./$(BINARY_NAME) version 2>/dev/null || echo "version コマンド未実装"; \
	else \
		echo "エラー: $(BINARY_NAME) が見つかりません"; \
		echo "先に 'make build' を実行してください"; \
	fi

## dev: 開発環境のセットアップ
dev: deps
	@echo "開発ツールをインストール中..."
	@echo "golangci-lintのインストール..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "goimportsのインストール..."
	@if ! command -v goimports >/dev/null 2>&1; then \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@echo "✓ 開発環境のセットアップ完了"

## watch: ファイル変更を監視してテストを自動実行（要: entr）
watch:
	@if command -v entr >/dev/null 2>&1; then \
		echo "ファイル変更を監視中... (Ctrl+Cで停止)"; \
		find . -name '*.go' | entr -c make test-unit; \
	else \
		echo "エラー: entr がインストールされていません"; \
		echo "インストール (macOS): brew install entr"; \
		echo "インストール (Linux): apt-get install entr"; \
		exit 1; \
	fi
