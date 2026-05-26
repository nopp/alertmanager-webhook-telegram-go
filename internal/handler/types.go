package handler

import "time"

// alertmanagerPayload is the top-level Alertmanager webhook payload.
type alertmanagerPayload struct {
	Receiver          string       `json:"receiver"`
	Status            string       `json:"status"`
	Alerts            []alert      `json:"alerts"`
	GroupLabels       groupLabels  `json:"groupLabels"`
	CommonLabels      commonLabels `json:"commonLabels"`
	CommonAnnotations annotations  `json:"commonAnnotations"`
	ExternalURL       string       `json:"externalURL"`
	Version           string       `json:"version"`
	GroupKey          string       `json:"groupKey"`
}

type alert struct {
	Status       string      `json:"status"`
	Labels       alertLabels `json:"labels"`
	Annotations  annotations `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
	Fingerprint  string      `json:"fingerprint"`
}

type alertLabels struct {
	Name      string `json:"name"`
	Instance  string `json:"instance"`
	Alertname string `json:"alertname"`
	Service   string `json:"service"`
	Severity  string `json:"severity"`
}

type annotations struct {
	Info        string `json:"info"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type groupLabels struct {
	Alertname string `json:"alertname"`
}

type commonLabels struct {
	Alertname string `json:"alertname"`
	Service   string `json:"service"`
	Severity  string `json:"severity"`
}
