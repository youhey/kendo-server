# kendo-server

ESP32-C3 の `kendo-node` から送信される 1 秒単位の振動センサー集約データを受信し、SQLite に保存する軽量 API サーバーです。

このフェーズでは保存機能だけを実装しています。衝撃判定、イベント生成、通知、Web UI、MQTT、複雑な集計は含みません。

## 環境変数

| 変数 | デフォルト | 説明 |
| --- | --- | --- |
| `KENDO_HOST_PORT` | `9401` | Docker ホスト側の公開ポート |
| `KENDO_HTTP_PORT` | `8080` | コンテナ内の HTTP ポート |
| `KENDO_DB_PATH` | `/data/kendo.sqlite3` | SQLite database path |
| `KENDO_DATA_DIR` | `./data` | ホスト側の SQLite 保存ディレクトリ |

Docker Compose は `.env` を読み込んで `docker-compose.yml` のポート、ボリューム、コンテナ環境変数に反映します。初回は `.env.example` をコピーして必要な値を変更してください。

```bash
cp .env.example .env
```

## 起動

```bash
docker compose up -d --build
```

停止:

```bash
docker compose down
```

デフォルトでは `http://localhost:9401` で起動し、SQLite データベースは `./data/kendo.sqlite3` に保存されます。

## API

この MVP では Bearer Token 認証は実装していません。`/healthz`、`/api/v1/samples`、`/api/v1/samples/recent` はすべて認証なしで利用します。

家庭内 LAN / Scum サーバー上の Docker Compose 運用を前提に、ESP32-C3 側のトークン管理・設定更新コストを避け、データ収集の安定性を優先しています。外部公開、VPN 越し利用、リバースプロキシ配下での利用を行う場合は、その時点で Bearer Token またはリバースプロキシ側の認証を追加します。

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
curl -sS 'http://localhost:9401/api/v1/samples/recent?node_id=ceiling-01&limit=10'
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
