package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/youhey/kendo-server/internal/db"
)

func TestSampleFlow(t *testing.T) {
	store, err := db.Open(filepath.Join(t.TempDir(), "kendo.sqlite3"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer store.Close()

	handler := NewHandler(store)
	mux := http.NewServeMux()
	handler.Register(mux)
	app := AuthMiddleware("secret", mux)

	resp := performRequest(app, http.MethodGet, "/healthz", nil, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want %d", resp.Code, http.StatusOK)
	}

	payload := []byte(`{
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
}`)
	resp = performRequest(app, http.MethodPost, "/api/v1/samples", bytes.NewReader(payload), "secret")
	if resp.Code != http.StatusOK {
		t.Fatalf("post status = %d, want %d", resp.Code, http.StatusOK)
	}

	resp = performRequest(app, http.MethodGet, "/api/v1/samples/recent?node_id=ceiling-01&limit=10", nil, "secret")
	if resp.Code != http.StatusOK {
		t.Fatalf("recent status = %d, want %d", resp.Code, http.StatusOK)
	}

	var body struct {
		OK      bool `json:"ok"`
		Samples []struct {
			NodeID     string `json:"node_id"`
			MeasuredAt string `json:"measured_at"`
			Piezo      struct {
				Peak int64 `json:"peak"`
			} `json:"piezo"`
		} `json:"samples"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode recent: %v", err)
	}
	if !body.OK {
		t.Fatal("recent ok = false, want true")
	}
	if len(body.Samples) != 1 {
		t.Fatalf("recent sample count = %d, want 1", len(body.Samples))
	}
	if body.Samples[0].NodeID != "ceiling-01" {
		t.Fatalf("node_id = %q, want ceiling-01", body.Samples[0].NodeID)
	}
	if body.Samples[0].MeasuredAt != "2026-06-23T15:00:00+09:00" {
		t.Fatalf("measured_at = %q", body.Samples[0].MeasuredAt)
	}
	if body.Samples[0].Piezo.Peak != 920 {
		t.Fatalf("piezo peak = %d, want 920", body.Samples[0].Piezo.Peak)
	}
}

func TestAuthRequired(t *testing.T) {
	store, err := db.Open(filepath.Join(t.TempDir(), "kendo.sqlite3"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer store.Close()

	handler := NewHandler(store)
	mux := http.NewServeMux()
	handler.Register(mux)
	app := AuthMiddleware("secret", mux)

	resp := performRequest(app, http.MethodGet, "/api/v1/samples/recent", nil, "")
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnauthorized)
	}
}

func TestInvalidSample(t *testing.T) {
	store, err := db.Open(filepath.Join(t.TempDir(), "kendo.sqlite3"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer store.Close()

	handler := NewHandler(store)
	mux := http.NewServeMux()
	handler.Register(mux)
	app := AuthMiddleware("secret", mux)

	resp := performRequest(app, http.MethodPost, "/api/v1/samples", bytes.NewBufferString(`{"node_id":"ceiling-01","measured_at":"invalid"}`), "secret")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusBadRequest)
	}
}

func performRequest(handler http.Handler, method, target string, body io.Reader, token string) *httptest.ResponseRecorder {
	if body == nil {
		body = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, target, body)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	return resp
}
