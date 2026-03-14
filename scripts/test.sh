#!/bin/bash
# テスト実行スクリプト

set -e

echo "========================================="
echo " jquants-cli Test Suite"
echo "========================================="
echo ""

# 色付き出力用
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 単体テスト実行
echo -e "${YELLOW}=== Running Unit Tests ===${NC}"
go test -v -race -short ./... || {
    echo -e "${RED}✗ Unit tests failed${NC}"
    exit 1
}
echo -e "${GREEN}✓ Unit tests passed${NC}"
echo ""

# 結合テスト実行（APIキーが設定されている場合のみ）
if [ -n "$JQUANTS_API_KEY" ]; then
    echo -e "${YELLOW}=== Running Integration Tests ===${NC}"
    go test -v ./tests/integration/... || {
        echo -e "${RED}✗ Integration tests failed${NC}"
        exit 1
    }
    echo -e "${GREEN}✓ Integration tests passed${NC}"
    echo ""
else
    echo -e "${YELLOW}ℹ Skipping integration tests (JQUANTS_API_KEY not set)${NC}"
    echo ""
fi

# カバレッジチェック
echo -e "${YELLOW}=== Checking Test Coverage ===${NC}"
./scripts/coverage.sh || {
    echo -e "${RED}✗ Coverage check failed${NC}"
    exit 1
}

echo ""
echo -e "${GREEN}========================================="
echo -e " ✓ All tests passed!"
echo -e "=========================================${NC}"
