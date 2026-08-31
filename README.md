# gopoke

URLの状態確認やレスポンス時間の計測、ICMP Ping、リダイレクト追跡ができるシンプルなCLIツールです。

## 主な機能

- **HTTPステータス・レスポンス計測**: 指定URLへのGETリクエスト、Content-Type、データサイズ、応答時間の確認
- **JSON出力**: 出力結果をパースしやすいJSON形式で取得
- **リダイレクト追跡**: リダイレクトの遷移チェーンと最終ステータスコードの表示
- **Ping**: ICMP Echo Requestによる疎通確認（※要管理者/root権限）

## ビルド / 実行方法

```bash
# ビルド
go build -o gopoke ./cmd/gopoke
# 実行
./gopoke https://example.com
```

## 使い方 (Usage)

```bash
gopoke [オプション] <URLまたはホスト>
```

### オプション一覧

| オプション | デフォルト | 説明 |
| :--- | :--- | :--- |
| `-timeout <秒>` | `10` | リクエストのタイムアウト時間（秒） |
| `-json` | `false` | 結果をJSON形式で出力 |
| `-tracking` | `false` | リダイレクトの遷移を追跡 |
| `-max-redirects <回数>` | `10` | リダイレクト追跡回数の上限（`0`で無制限） |
| `-ping` | `false` | ICMP Pingを実行（※要管理者権限） |

---

## 実行例

### 1. 通常の確認 (Poke)
```bash
$ gopoke https://example.com
Response: 200 OK
Content type: text/html; charset=UTF-8
Content length: 1256 bytes
Body size: 1256 bytes
Time: 142.35ms
```

### 2. JSON形式で出力
```bash
$ gopoke -json https://google.com                                                                                 
{
  "status": "200 OK",
  "content_type": "text/html; charset=ISO-8859-1",
  "content_length": -1,
  "body_size": 84420,
  "elapsed": "301.4615ms"
}
```

### 3. リダイレクト追跡
```bash
$ gopoke -tracking http://google.com
リダイレクトはありませんでした
$ gopoke -tracking https://x.gd/GoFKp
Redirect chain:
1: https://www.google.com/
Final HTTP status code: 200
```

### 4. Ping送信
```bash
$ gopoke -ping example.com
Ping result: 8 bytes, Type: 0
```
