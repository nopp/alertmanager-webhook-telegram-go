package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mockSender records calls to Send for assertion in tests.
type mockSender struct {
	calls []mockCall
	err   error
}

type mockCall struct {
	chatID int64
	text   string
}

func (m *mockSender) Send(chatID int64, text string) error {
	m.calls = append(m.calls, mockCall{chatID: chatID, text: text})
	return m.err
}

// ── buildMessage ─────────────────────────────────────────────────────────────

func TestBuildMessage_Firing(t *testing.T) {
	a := alert{
		Status: "firing",
		Labels: alertLabels{
			Instance: "web-01",
			Name:     "frontend",
		},
		Annotations: annotations{
			Summary:     "High CPU",
			Description: "CPU above 90%",
		},
		StartsAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	got := buildMessage("firing", a)

	assertContains(t, got, "Status: firing")
	assertContains(t, got, "Instance: web-01 (frontend)")
	assertContains(t, got, "Summary: High CPU")
	assertContains(t, got, "Description: CPU above 90%")
	assertContains(t, got, "Started: 2024-01-15 10:00:00")
}

func TestBuildMessage_Resolved(t *testing.T) {
	a := alert{
		Status: "resolved",
		Labels: alertLabels{Instance: "db-01"},
		Annotations: annotations{
			Info: "Back to normal",
		},
		EndsAt: time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC),
	}

	got := buildMessage("resolved", a)

	assertContains(t, got, "Status: resolved")
	assertContains(t, got, "Instance: db-01")
	assertContains(t, got, "Info: Back to normal")
	assertContains(t, got, "Resolved: 2024-01-15 11:30:00")
}

func TestBuildMessage_InstanceWithoutName(t *testing.T) {
	a := alert{
		Status: "firing",
		Labels: alertLabels{Instance: "10.0.0.1:9100"},
	}

	got := buildMessage("firing", a)

	assertContains(t, got, "Instance: 10.0.0.1:9100")
	assertNotContains(t, got, "()")
}

func TestBuildMessage_NoInstanceNoAnnotations(t *testing.T) {
	a := alert{Status: "firing"}

	got := buildMessage("firing", a)

	assertContains(t, got, "Status: firing")
	assertNotContains(t, got, "Instance:")
	assertNotContains(t, got, "Info:")
	assertNotContains(t, got, "Summary:")
	assertNotContains(t, got, "Description:")
}

// ── Alert HTTP handler ────────────────────────────────────────────────────────

func TestAlertHandler_ValidPayload(t *testing.T) {
	mock := &mockSender{}
	h := New(mock, 123456)

	payload := alertmanagerPayload{
		Status: "firing",
		Alerts: []alert{
			{
				Status:  "firing",
				Labels:  alertLabels{Instance: "host-01"},
				StartsAt: time.Now(),
			},
			{
				Status:  "firing",
				Labels:  alertLabels{Instance: "host-02"},
				StartsAt: time.Now(),
			},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/alert", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Alert(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if len(mock.calls) != 2 {
		t.Fatalf("expected 2 telegram sends, got %d", len(mock.calls))
	}
	// chatID must be negated
	for _, c := range mock.calls {
		if c.chatID != -123456 {
			t.Errorf("expected chatID -123456, got %d", c.chatID)
		}
	}
}

func TestAlertHandler_InvalidJSON(t *testing.T) {
	mock := &mockSender{}
	h := New(mock, 123456)

	req := httptest.NewRequest(http.MethodPost, "/alert", bytes.NewReader([]byte("not-json")))
	rec := httptest.NewRecorder()

	h.Alert(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(mock.calls) != 0 {
		t.Errorf("expected 0 telegram sends on bad payload, got %d", len(mock.calls))
	}
}

func TestAlertHandler_EmptyAlerts(t *testing.T) {
	mock := &mockSender{}
	h := New(mock, 123456)

	payload := alertmanagerPayload{Status: "firing", Alerts: []alert{}}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/alert", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Alert(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if len(mock.calls) != 0 {
		t.Errorf("expected 0 telegram sends for empty alerts, got %d", len(mock.calls))
	}
}

func TestAlertHandler_ResponseIsJSON(t *testing.T) {
	mock := &mockSender{}
	h := New(mock, 1)

	payload := alertmanagerPayload{Status: "resolved"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/alert", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Alert(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	var decoded alertmanagerPayload
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Errorf("response body is not valid JSON: %v", err)
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !bytes.Contains([]byte(s), []byte(substr)) {
		t.Errorf("expected %q to contain %q", s, substr)
	}
}

func assertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if bytes.Contains([]byte(s), []byte(substr)) {
		t.Errorf("expected %q NOT to contain %q", s, substr)
	}
}
