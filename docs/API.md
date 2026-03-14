# J-Quants API 仕様メモ

このドキュメントは、jquants-cli実装に必要なJ-Quants API v2の仕様をまとめたものです。

## 基本情報

- **ベースURL**: `https://api.jquants.com/v2`
- **認証**: `x-api-key` ヘッダーでAPIキーを送信
- **レスポンス形式**: JSON

## 認証

### APIキーの取得

1. [J-Quants](https://jpx-jquants.com/)でアカウント作成
2. ダッシュボードからAPIキーを取得

### リクエストヘッダー

```
x-api-key: YOUR_API_KEY
```

## エンドポイント一覧

### 1. 上場銘柄一覧 (`/equities/master`)

**エンドポイント**: `GET /v2/equities/master`

**パラメータ**:
- `code` (optional): 銘柄コード (例: 86970)
- `date` (optional): 日付 (YYYY-MM-DD または YYYYMMDD)

**レスポンス例**:
```json
{
  "data": [
    {
      "Date": "2022-11-11",
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

**主要フィールド**:
- `Date`: 情報適用年月日
- `Code`: 銘柄コード
- `CoName`: 会社名（日本語）
- `CoNameEn`: 会社名（英語）
- `S17`: 17業種コード
- `S17Nm`: 17業種名
- `S33`: 33業種コード
- `S33Nm`: 33業種名
- `ScaleCat`: 規模コード
- `Mkt`: 市場区分コード
- `MktNm`: 市場区分名
- `Mrgn`: 貸借信用区分 (Standard/Premium プランのみ)
- `MrgnNm`: 貸借信用区分名 (Standard/Premium プランのみ)

### 2. 株価四本値 (`/equities/bars/daily`)

**エンドポイント**: `GET /v2/equities/bars/daily`

**パラメータ**:
- `code` (optional): 銘柄コード
- `date` (optional): 日付
- `from` (optional): 開始日
- `to` (optional): 終了日
- `pagination_key` (optional): ページネーションキー

**注意**: `code` または `date` のいずれか必須

**レスポンス例**:
```json
{
  "data": [
    {
      "Date": "2023-03-24",
      "Code": "86970",
      "O": 2047.0,
      "H": 2069.0,
      "L": 2035.0,
      "C": 2045.0,
      "UL": "0",
      "LL": "0",
      "Vo": 2202500.0,
      "Va": 4507051850.0,
      "AdjFactor": 1.0,
      "AdjO": 2047.0,
      "AdjH": 2069.0,
      "AdjL": 2035.0,
      "AdjC": 2045.0,
      "AdjVo": 2202500.0
    }
  ],
  "pagination_key": "value1.value2."
}
```

**主要フィールド**:
- `Date`: 日付
- `Code`: 銘柄コード
- `O`: 始値（調整前）
- `H`: 高値（調整前）
- `L`: 安値（調整前）
- `C`: 終値（調整前）
- `UL`: ストップ高フラグ
- `LL`: ストップ安フラグ
- `Vo`: 取引高（調整前）
- `Va`: 取引代金
- `AdjFactor`: 調整係数
- `AdjO`: 調整済み始値
- `AdjH`: 調整済み高値
- `AdjL`: 調整済み安値
- `AdjC`: 調整済み終値
- `AdjVo`: 調整済み取引高
- `pagination_key`: 次のページのキー（データが続く場合のみ）

**前場・後場データ (Premiumプランのみ)**:
- `MO`, `MH`, `ML`, `MC`: 前場四本値
- `AO`, `AH`, `AL`, `AC`: 後場四本値

### 3. 財務情報 (`/fins/summary`)

**エンドポイント**: `GET /v2/fins/summary`

**パラメータ**:
- `code` (optional): 銘柄コード
- `date` (optional): 開示日付
- `pagination_key` (optional): ページネーションキー

**注意**: `code` または `date` のいずれか必須

**レスポンス例**:
```json
{
  "data": [
    {
      "DiscDate": "2023-01-30",
      "DiscTime": "12:00:00",
      "Code": "86970",
      "DiscNo": "20230127594871",
      "DocType": "3QFinancialStatements_Consolidated_IFRS",
      "CurPerType": "3Q",
      "CurPerSt": "2022-04-01",
      "CurPerEn": "2022-12-31",
      "Sales": "100529000000",
      "OP": "51765000000",
      "NP": "35175000000",
      "EPS": "66.76",
      "TA": "79205861000000",
      "Eq": "320021000000",
      "DivAnn": "62.0",
      "FSales": "132500000000",
      "FNP": "45000000000"
    }
  ],
  "pagination_key": "value1.value2."
}
```

**主要フィールド**:
- `DiscDate`: 開示日
- `DiscTime`: 開示時刻
- `Code`: 銘柄コード
- `DiscNo`: 開示番号
- `DocType`: 開示書類種別
- `CurPerType`: 会計期間種類 (1Q, 2Q, 3Q, 4Q, FY)
- `Sales`: 売上高
- `OP`: 営業利益
- `OdP`: 経常利益
- `NP`: 当期純利益
- `EPS`: 一株あたり当期純利益
- `TA`: 総資産
- `Eq`: 純資産
- `DivAnn`: 一株あたり配当実績（年間）
- `FSales`: 売上高予想
- `FNP`: 当期純利益予想

### 4. Bulk API - ファイル一覧 (`/bulk/list`)

**エンドポイント**: `GET /v2/bulk/list`

**パラメータ**:
- `endpoint` (optional): エンドポイント名 (例: /equities/bars/daily)
- `date` (optional): 日付 (YYYY-MM または YYYY-MM-DD)
- `from` (optional): 開始日 (YYYY-MM または YYYY-MM-DD)
- `to` (optional): 終了日 (YYYY-MM または YYYY-MM-DD)

**注意**: `endpoint` または `date` のいずれか必須

**レスポンス例**:
```json
{
  "data": [
    {
      "Key": "equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz",
      "LastModified": "2025-11-07T20:48:51.295000+00:00",
      "Size": 6933528
    }
  ]
}
```

**フィールド**:
- `Key`: ファイルキー（ダウンロード時に使用）
- `LastModified`: 最終更新日時
- `Size`: ファイルサイズ（バイト）

### 5. Bulk API - ファイルダウンロード (`/bulk/get`)

**エンドポイント**: `GET /v2/bulk/get`

**パラメータ**:
- `key` (optional): ファイルキー
- `endpoint` + `date` (optional): エンドポイントと日付の組み合わせ

**注意**: `key` または `endpoint`+`date` のいずれか必須

**レスポンス例**:
```json
{
  "url": "https://example.presigned-url.com/..."
}
```

**フィールド**:
- `url`: ダウンロード用署名付きURL（有効期限5分）

**ダウンロード手順**:
1. `/bulk/list`でファイル一覧取得
2. `/bulk/get`で署名付きURL取得
3. URLからgzip圧縮CSVをダウンロード
4. gzip解凍して使用

## ページネーション処理

大量データを取得する一部のエンドポイントはページネーションに対応しています。

**対応エンドポイント**:
- `/equities/bars/daily`
- `/fins/summary`

**処理フロー**:
```
1. 初回リクエスト（pagination_keyなし）
2. レスポンスのpagination_keyをチェック
3. pagination_keyがある場合、次のリクエストに含める
4. pagination_keyがnullまたは空になるまで繰り返す
```

**実装例**:
```go
allData := []DataType{}
paginationKey := ""

for {
    params := url.Values{}
    params.Set("date", "2025-03-01")
    if paginationKey != "" {
        params.Set("pagination_key", paginationKey)
    }

    resp := GetData(params)
    allData = append(allData, resp.Data...)

    if resp.PaginationKey == "" {
        break
    }
    paginationKey = resp.PaginationKey

    // レート制限対策
    time.Sleep(500 * time.Millisecond)
}
```

## レート制限

J-Quants APIにはレート制限があります。

**対策**:
- リクエスト間に500ms程度の待機時間を設ける
- 429エラー時は適切な待機時間後にリトライ
- Bulk APIを活用して一括取得する

## エラーハンドリング

### HTTPステータスコード

| コード | 意味 | 対処方法 |
|-------|------|---------|
| 200 | 成功 | - |
| 400 | リクエストエラー | パラメータを確認 |
| 401 | 認証エラー | APIキーを確認 |
| 404 | データなし | パラメータを確認 |
| 429 | レート制限 | 待機後リトライ |
| 500 | サーバーエラー | 待機後リトライ |
| 502 | Bad Gateway | 待機後リトライ |
| 503 | Service Unavailable | 待機後リトライ |

### リトライ戦略

5xxエラーの場合のみリトライを実施:
- 最大リトライ回数: 3回
- 待機時間: 指数バックオフ (1秒 → 2秒 → 4秒)

## プラン別データ取得範囲

| プラン | 取得可能期間 |
|--------|------------|
| Free | 12週間前〜2年12週間前まで |
| Light | 5年前まで |
| Standard | 10年前まで |
| Premium | 20年前まで |

## Bulk API対応エンドポイント

以下のエンドポイントはBulk APIでCSV一括取得が可能:

- `/equities/master` - 上場銘柄一覧
- `/equities/bars/daily` - 株価四本値
- `/fins/summary` - 財務情報
- `/equities/investor-types` - 投資部門別情報
- `/indices/bars/daily/topix` - TOPIX四本値
- `/indices/bars/daily` - 指数四本値
- その他多数

詳細は[公式ドキュメント](https://jpx-jquants.com/ja/spec/data-spec)を参照。

## 参考リンク

- [J-Quants API仕様書](https://jpx-jquants.com/ja/spec)
- [J-Quantsヘルプページ](https://jpx-jquants.com/ja/help)
- [V1からV2への移行ガイド](https://jpx-jquants.com/ja/spec/migration-v1-v2)
