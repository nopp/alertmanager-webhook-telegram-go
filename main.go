package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type alertmanagerAlert struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname string `json:"alertname"`
			Service   string `json:"service"`
			Severity  string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Summary string `json:"summary"`
		} `json:"annotations"`
		StartsAt     string    `json:"startsAt"`
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

func alertToTelegram(w http.ResponseWriter, r *http.Request) {

	var alert alertmanagerAlert

	_ = json.NewDecoder(r.Body).Decode(&alert)

	log.Println(alert)
	json.NewEncoder(w).Encode(alert)

}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/alert", alertToTelegram).Methods("POST")

	log.Fatal(http.ListenAndServe("0.0.0.0:8686", router))
}
