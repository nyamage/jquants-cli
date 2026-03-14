#!/bin/bash
# カバレッジ測定スクリプト

set -e

# 色付き出力用
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# カバレッジ目標
TARGET_COVERAGE=80

echo "Generating coverage report..."
go test -coverprofile=coverage.out -covermode=atomic ./...

echo ""
echo "========================================="
echo " Coverage Summary"
echo "========================================="

# 全体のカバレッジを表示
go tool cover -func=coverage.out | tail -1

echo ""
echo "Packages with coverage < ${TARGET_COVERAGE}%:"
echo "========================================="

# 80%未満のパッケージを表示
LOW_COVERAGE=$(go tool cover -func=coverage.out | awk -v target="$TARGET_COVERAGE" '
    NR > 1 && $3 != "total:" {
        cov = substr($3, 1, length($3)-1)
        if (cov < target) {
            print $1, $3
        }
    }
')

if [ -z "$LOW_COVERAGE" ]; then
    echo -e "${GREEN}None - All packages meet the coverage target!${NC}"
else
    echo -e "${RED}$LOW_COVERAGE${NC}"
fi

echo ""

# 全体のカバレッジ率を取得
TOTAL_COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo "Total coverage: ${TOTAL_COVERAGE}%"
echo "Target coverage: ${TARGET_COVERAGE}%"

# カバレッジが目標未満の場合はエラー
if (( $(echo "$TOTAL_COVERAGE < $TARGET_COVERAGE" | bc -l) )); then
    echo -e "${RED}ERROR: Total coverage ${TOTAL_COVERAGE}% is below target ${TARGET_COVERAGE}%${NC}"
    echo ""
    echo "To view detailed coverage report, run:"
    echo "  go tool cover -html=coverage.out"
    exit 1
fi

echo -e "${GREEN}✓ Coverage check passed: ${TOTAL_COVERAGE}%${NC}"
echo ""
echo "To view detailed HTML coverage report, run:"
echo "  go tool cover -html=coverage.out"
