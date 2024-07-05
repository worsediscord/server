package api

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// handleHealth returns a health status
//
//	@Summary	Checks server health
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	HealthResponse
//	@Failure	500
//	@Router		/health [get]
func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{Status: "ok"}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
