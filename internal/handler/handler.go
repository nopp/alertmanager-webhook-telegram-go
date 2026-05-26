package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"alertmanager-webhook-telegram-go/internal/telegram"
)

const timeDateFormat = "2006-01-02 15:04:05"

// Handler handles incoming Alertmanager webhook requests.
type Handler struct {
	telegram telegram.Sender
	chatID   int64
}

// New creates a Handler with the given Telegram sender and target chat ID.
func New(tg telegram.Sender, chatID int64) *Handler {
	return &Handler{telegram: tg, chatID: chatID}
}

// Alert is the HTTP handler for POST /alert.
func (h *Handler) Alert(w http.ResponseWriter, r *http.Request) {
	var payload alertmanagerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for _, a := range payload.Alerts {
		msg := buildMessage(payload.Status, a)
		if err := h.telegram.Send(-h.chatID, msg); err != nil {
			log.Printf("handler: failed to send telegram message: %v", err)
		}
	}

	log.Printf("handler: processed %d alert(s), group_status=%s", len(payload.Alerts), payload.Status)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("handler: failed to encode response: %v", err)
	}
}

// buildMessage formats a single alert into a human-readable Telegram message.
func buildMessage(groupStatus string, a alert) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Status: %s\n", groupStatus)

	if a.Labels.Instance != "" {
		if a.Labels.Name != "" {
			fmt.Fprintf(&sb, "Instance: %s (%s)\n", a.Labels.Instance, a.Labels.Name)
		} else {
			fmt.Fprintf(&sb, "Instance: %s\n", a.Labels.Instance)
		}
	}
	if a.Annotations.Info != "" {
		fmt.Fprintf(&sb, "Info: %s\n", a.Annotations.Info)
	}
	if a.Annotations.Summary != "" {
		fmt.Fprintf(&sb, "Summary: %s\n", a.Annotations.Summary)
	}
	if a.Annotations.Description != "" {
		fmt.Fprintf(&sb, "Description: %s\n", a.Annotations.Description)
	}

	switch a.Status {
	case "resolved":
		fmt.Fprintf(&sb, "Resolved: %s", a.EndsAt.Format(timeDateFormat))
	case "firing":
		fmt.Fprintf(&sb, "Started: %s", a.StartsAt.Format(timeDateFormat))
	}

	return sb.String()
}
