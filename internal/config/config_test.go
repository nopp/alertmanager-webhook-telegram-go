package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, kv map[string]string) {
	t.Helper()
	for k, v := range kv {
		t.Setenv(k, v)
	}
}

func TestLoad_AllEnvSet(t *testing.T) {
	setEnv(t, map[string]string{
		"BOT_TOKEN":   "token123",
		"CHAT_ID":     "-987654",
		"LISTEN_ADDR": "127.0.0.1:8080",
	})

	cfg := Load()

	if cfg.BotToken != "token123" {
		t.Errorf("BotToken: got %q, want %q", cfg.BotToken, "token123")
	}
	if cfg.ChatID != -987654 {
		t.Errorf("ChatID: got %d, want -987654", cfg.ChatID)
	}
	if cfg.ListenAddr != "127.0.0.1:8080" {
		t.Errorf("ListenAddr: got %q, want %q", cfg.ListenAddr, "127.0.0.1:8080")
	}
}

func TestLoad_DefaultListenAddr(t *testing.T) {
	setEnv(t, map[string]string{
		"BOT_TOKEN": "tok",
		"CHAT_ID":   "42",
	})
	os.Unsetenv("LISTEN_ADDR")

	cfg := Load()

	if cfg.ListenAddr != "0.0.0.0:9229" {
		t.Errorf("ListenAddr: got %q, want default 0.0.0.0:9229", cfg.ListenAddr)
	}
}

func TestLoad_PositiveChatID(t *testing.T) {
	setEnv(t, map[string]string{
		"BOT_TOKEN": "tok",
		"CHAT_ID":   "99",
	})

	cfg := Load()

	if cfg.ChatID != 99 {
		t.Errorf("ChatID: got %d, want 99", cfg.ChatID)
	}
}
