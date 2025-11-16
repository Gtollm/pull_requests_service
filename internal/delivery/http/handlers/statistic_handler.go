package handlers

import (
	"encoding/json"
	"net/http"

	"pull-request-review/internal/domain/ports/service"
)

type StatisticsHandler struct {
	statisticsService service.StatisticsService
}

func NewStatisticsHandler(
	statisticsService service.StatisticsService,
) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

// GetStatistics handles GET /statistics
func (h *StatisticsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	response, err := h.statisticsService.GetStatistics(r.Context())
	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}