package handlers

import (
	"net/http"

	"github.com/wizzyszn/go_bank/db"
	"github.com/wizzyszn/go_bank/utils"
)

type HealthHandler struct {
	db *db.DB
}

func NewHealthHandler(database *db.DB) *HealthHandler {
	return &HealthHandler{
		db: database,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {

	if err := h.db.Health(); err != nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database Unhealthy")
		return
	}

	stats := h.db.Stats()

	response := map[string]any{
		"status": "healthy",
		"database": map[string]any{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
		},
	}

	utils.WriteSuccess(w, response)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {

	if err := h.db.Ping(); err != nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Not ready")
		return
	}
	utils.WriteSuccess(w, map[string]string{
		"status": "ready",
	})
}

// Live handles GET /live (for Kubernetes liveness probes)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{
		"status": "alive",
	})
}
