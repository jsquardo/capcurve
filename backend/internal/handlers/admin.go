package handlers

import (
	"net/http"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/jsquardo/capcurve/internal/syncjob"
	"github.com/labstack/echo/v4"
)

type SyncStatusSource interface {
	Snapshot() syncjob.StatusSnapshot
}

type AdminDashboardResponse struct {
	TotalPlayers  int64                  `json:"total_players"`
	ActivePlayers int64                  `json:"active_players"`
	SyncStatus    syncjob.StatusSnapshot `json:"sync_status"`
}

func (h *Handler) AdminDashboard(c echo.Context) error {
	var totalPlayers int64
	if err := h.db.Model(&models.Player{}).Count(&totalPlayers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var activePlayers int64
	if err := h.db.Model(&models.Player{}).Where("active = ?", true).Count(&activePlayers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, AdminDashboardResponse{
		TotalPlayers:  totalPlayers,
		ActivePlayers: activePlayers,
		SyncStatus:    h.syncStatus.Snapshot(),
	})
}
