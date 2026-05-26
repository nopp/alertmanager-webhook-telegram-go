package main

import (
	"log"
	"net/http"

	"alertmanager-webhook-telegram-go/internal/config"
	"alertmanager-webhook-telegram-go/internal/handler"
	"alertmanager-webhook-telegram-go/internal/telegram"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	bot, err := telegram.NewClient(cfg.BotToken)
	if err != nil {
		log.Fatalf("failed to create telegram bot: %v", err)
	}

	h := handler.New(bot, cfg.ChatID)

	router := mux.NewRouter()
	router.HandleFunc("/alert", h.Alert).Methods("POST")

	log.Printf("listening on %s", cfg.ListenAddr)
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, router))
}
