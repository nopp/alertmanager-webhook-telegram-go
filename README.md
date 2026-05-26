# Alertmanager webhook for Telegram (GO Version)

Alertmanager webhook receiver that forwards alerts to a Telegram chat.

Python version: https://github.com/nopp/alertmanager-webhook-telegram-python

## Requirements

- Go 1.22+
- Telegram bot token ([BotFather](https://t.me/BotFather))
- Telegram chat ID (negative for groups, e.g. `-123456789`)

## Configuration

All configuration is via environment variables — no hardcoded values.

| Variable      | Required | Default          | Description                  |
|---------------|----------|------------------|------------------------------|
| `BOT_TOKEN`   | ✅       | —                | Telegram bot token           |
| `CHAT_ID`     | ✅       | —                | Telegram chat ID (numeric)   |
| `LISTEN_ADDR` | ❌       | `0.0.0.0:9229`   | Address and port to bind     |

## Running

### Docker (recommended)

```bash
docker build -f docker/Dockerfile -t alertmanager-webhook-telegram-go .

docker run -d \
  -e BOT_TOKEN=your_bot_token \
  -e CHAT_ID=-123456789 \
  -p 9229:9229 \
  alertmanager-webhook-telegram-go
```

### From source

```bash
go build -o webhook ./cmd/webhook
BOT_TOKEN=xxx CHAT_ID=-123456789 ./webhook
```

### Tests

```bash
go test ./...
```

## Alertmanager configuration

```yaml
receivers:
  - name: telegram-webhook
    webhook_configs:
      - url: http://<host>:9229/alert
        send_resolved: true
```

## Project structure

```
cmd/webhook/          entry point
internal/config/      env var loading
internal/handler/     HTTP handler + Alertmanager types
internal/telegram/    Telegram bot client
docker/               Dockerfile
```

## Endpoint

| Method | Path     | Description                        |
|--------|----------|------------------------------------|
| POST   | `/alert` | Receives Alertmanager webhook payload |
