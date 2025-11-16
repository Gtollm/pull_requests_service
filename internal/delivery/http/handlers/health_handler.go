package handlers

import (
	"encoding/json"
	"net/http"

	"pull-request-review/internal/infrastructure/database"
)

type HealthHandler struct {
	db *database.Database
}

func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Check handles GET /health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(r.Context()); err != nil {
		response := map[string]string{
			"status": "unhealthy",
			"error":  "database connection failed",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
		return
	}

	response := map[string]string{
		"status": "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}