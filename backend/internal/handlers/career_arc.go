package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) GetCareerArc(c echo.Context) error {
	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var player models.Player
	if err := h.db.First(&player, playerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "player not found"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var stats []models.SeasonStat
	if err := h.db.Where("player_id = ?", player.ID).Order("year ASC, id ASC").Find(&stats).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var arc *models.CareerArc
	var arcRecord models.CareerArc
	if err := h.db.Where("player_id = ?", player.ID).First(&arcRecord).Error; err != nil {
		// Missing arc metadata should not block the chart endpoint from returning
		// the player's historical timeline.
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	} else {
		arc = &arcRecord
	}

	projectionPayload, err := h.buildPlayerProjectionPayload(player, stats)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, careerArcResponse{
		Data: newCareerArcData(player, stats, arc, projectionPayload),
	})
}
