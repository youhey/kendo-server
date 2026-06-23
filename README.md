# kendo-server

ESP32-C3 の `kendo-node` から送信される 1 秒単位の振動センサー集約データを受信し、SQLite に保存する軽量 API サーバーです。

このフェーズでは保存機能だけを実装しています。衝撃判定、イベント生成、通知、Web UI、MQTT、複雑な集計は含みません。

## 環境変数

| 変数 | デフォルト | 説明 |
| --- | --- | --- |
| `KENDO_ADDR` | `:8080` | HTTP listen address |
| `KENDO_DB_PATH` | `/data/kendo.sqlite3` | SQLite database path |
| `KENDO_API_TOKEN` | なし | API Bearer Token。未設定なら起動エラー |

## 起動

```bash
docker compose up -d --build
```

停止:

```bash
docker compose down
```

SQLite データベースは `./data/kendo.sqlite3` に保存されます。

## API

`/healthz` 以外は Bearer Token が必要です。

```http
Authorization: Bearer <KENDO_API_TOKEN>
```

### GET /healthz

```bash
curl -sS http://localhost:9401/healthz
```

レスポンス:

```json
{"ok":true}
```

### POST /api/v1/samples

```bash
curl -sS -X POST http://localhost:9401/api/v1/samples \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer change-me' \
  -d '{
    "node_id": "ceiling-01",
    "seq": 1,
    "measured_at": "2026-06-23T15:00:00+09:00",
    "adxl": {
      "x": 0.012,
      "y": -0.035,
      "z": 0.991,
      "mag": 0.992,
      "rms": 0.018,
      "peak": 1.087
    },
    "piezo": {
      "raw": 1810,
      "min": 1788,
      "max": 2730,
      "peak": 920
    }
  }'
```

正常レスポンス:

```json
{"ok":true}
```

### GET /api/v1/samples/recent

```bash
curl -sS 'http://localhost:9401/api/v1/samples/recent?node_id=ceiling-01&limit=10' \
  -H 'Authorization: Bearer change-me'
```

`limit` のデフォルトは `100`、最大値は `1000` です。`node_id` を省略した場合は全 node が対象です。

レスポンス例:

```json
{
  "ok": true,
  "samples": [
    {
      "id": 1,
      "node_id": "ceiling-01",
      "seq": 1,
      "measured_at": "2026-06-23T15:00:00+09:00",
      "received_at": "2026-06-23T06:00:01Z",
      "adxl": {
        "x": 0.012,
        "y": -0.035,
        "z": 0.991,
        "mag": 0.992,
        "rms": 0.018,
        "peak": 1.087
      },
      "piezo": {
        "raw": 1810,
        "min": 1788,
        "max": 2730,
        "peak": 920
      },
      "created_at": "2026-06-23T06:00:01Z"
    }
  ]
}
```
