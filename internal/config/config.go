package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	BotToken   string
	ChatID     int64
	ListenAddr string
}

// Load reads configuration from environment variables and exits on missing required values.
//
// Required:
//   - BOT_TOKEN  — Telegram bot token
//   - CHAT_ID    — Telegram chat ID (numeric, negative for groups)
//
// Optional:
//   - LISTEN_ADDR — host:port to bind (default: 0.0.0.0:9229)
func Load() *Config {
	botToken := mustEnv("BOT_TOKEN")

	chatIDStr := mustEnv("CHAT_ID")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("config: invalid CHAT_ID %q: %v", chatIDStr, err)
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "0.0.0.0:9229"
	}

	return &Config{
		BotToken:   botToken,
		ChatID:     chatID,
		ListenAddr: listenAddr,
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("config: required env var %s is not set", key)
	}
	return v
}
