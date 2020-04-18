package alert

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	botToken = "xxbotTokenxx"
	chatID   = 666777666
)

type alertmanagerAlert struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Name      string `json:"name"`
			Instance  string `json:"instance"`
			Alertname string `json:"alertname"`
			Service   string `json:"service"`
			Severity  string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Info        string `json:"info"`
			Description string `json:"description"`
			Summary     string `json:"summary"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
		Service   string `json:"service"`
		Severity  string `json:"severity"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL string `json:"externalURL"`
	Version     string `json:"version"`
	GroupKey    string `json:"groupKey"`
}

func message(w http.ResponseWriter, message string) {
	log.Println(message)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

// ToTelegram function responsible to send msg to telegram
func ToTelegram(w http.ResponseWriter, r *http.Request) {

	var alerts alertmanagerAlert

	bot, err := botapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	_ = json.NewDecoder(r.Body).Decode(&alerts)

	for _, alert := range alerts.Alerts {
		log.Println(alert.Status)
		telegramMsg := "Status: " + alerts.Status + "\n"
		if alert.Labels.Name != "" {
			telegramMsg += "Instance: " + alert.Labels.Instance + "(" + alert.Labels.Name + ")\n"
		}
		if alert.Annotations.Info != "" {
			telegramMsg += "Info: " + alert.Annotations.Info + "\n"
		}
		if alert.Annotations.Summary != "" {
			telegramMsg += "Summary: " + alert.Annotations.Summary + "\n"
		}
		if alert.Annotations.Description != "" {
			telegramMsg += "Description: " + alert.Annotations.Description + "\n"
		}
		if alert.Status == "resolved" {
			telegramMsg += "Resolved: " + alert.EndsAt.Format("2006-01-02 15:04:05")
		} else if alert.Status == "firing" {
			telegramMsg += "Started: " + alert.StartsAt.Format("2006-01-02 15:04:05")
		}

		msg := botapi.NewMessage(-chatID, telegramMsg)
		bot.Send(msg)
	}

	log.Println(alerts)
	json.NewEncoder(w).Encode(alerts)

}
