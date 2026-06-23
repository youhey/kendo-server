package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/youhey/kendo-server/internal/db"
	"github.com/youhey/kendo-server/internal/model"
)

type Handler struct {
	store *db.Store
}

func NewHandler(store *db.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.healthz)
	mux.HandleFunc("POST /api/v1/samples", h.createSample)
	mux.HandleFunc("GET /api/v1/samples/recent", h.recentSamples)
}

func (h *Handler) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) createSample(w http.ResponseWriter, r *http.Request) {
	var sample model.Sample
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&sample); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := validateSample(sample); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.SaveSample(r.Context(), sample); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save sample")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) recentSamples(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, err := parseLimit(query.Get("limit"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	samples, err := h.store.RecentSamples(r.Context(), query.Get("node_id"), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load samples")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"samples": samples,
	})
}

func validateSample(sample model.Sample) error {
	if sample.NodeID == "" {
		return errors.New("node_id is required")
	}
	if sample.MeasuredAt == "" {
		return errors.New("measured_at is required")
	}
	if _, err := time.Parse(time.RFC3339, sample.MeasuredAt); err != nil {
		return errors.New("measured_at must be RFC3339")
	}

	return nil
}

func parseLimit(value string) (int, error) {
	if value == "" {
		return 100, nil
	}

	limit, err := strconv.Atoi(value)
	if err != nil || limit < 1 {
		return 0, errors.New("limit must be a positive integer")
	}
	if limit > 1000 {
		return 1000, nil
	}

	return limit, nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"ok":    false,
		"error": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
