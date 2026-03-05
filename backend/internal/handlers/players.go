package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/jsquardo/capcurve/internal/models"
)

func (h *Handler) ListPlayers(c echo.Context) error {
	var players []models.Player
	query := h.db.Limit(50)

	if active := c.QueryParam("active"); active != "" {
		query = query.Where("active = ?", active == "true")
	}
	if position := c.QueryParam("position"); position != "" {
		query = query.Where("position = ?", position)
	}

	if err := query.Find(&players).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, players)
}

func (h *Handler) GetPlayer(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var player models.Player
	if err := h.db.Preload("SeasonStats").Preload("Contracts").Preload("CareerArc").First(&player, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "player not found"})
	}
	return c.JSON(http.StatusOK, player)
}

func (h *Handler) SearchPlayers(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "search query required"})
	}

	var players []models.Player
	searchTerm := "%" + query + "%"
	if err := h.db.Where(
		"first_name ILIKE ? OR last_name ILIKE ? OR (first_name || ' ' || last_name) ILIKE ?",
		searchTerm, searchTerm, searchTerm,
	).Limit(20).Find(&players).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, players)
}

func (h *Handler) GetCareerArc(c echo.Context) error {
	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var arc models.CareerArc
	if err := h.db.Where("player_id = ?", playerID).First(&arc).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "career arc not found"})
	}

	var stats []models.SeasonStat
	h.db.Where("player_id = ?", playerID).Order("year ASC").Find(&stats)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"arc":          arc,
		"season_stats": stats,
	})
}

func (h *Handler) GetPlayerContracts(c echo.Context) error {
	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid player id"})
	}

	var contracts []models.Contract
	if err := h.db.Preload("ContractSeasons").Where("player_id = ?", playerID).Find(&contracts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, contracts)
}

func (h *Handler) GetContract(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid contract id"})
	}

	var contract models.Contract
	if err := h.db.Preload("ContractSeasons").First(&contract, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "contract not found"})
	}
	return c.JSON(http.StatusOK, contract)
}

func (h *Handler) MostOverpaid(c echo.Context) error {
	var contracts []models.Contract
	h.db.Where("is_active = ? AND overall_value_score < ?", true, -20).
		Order("overall_value_score ASC").Limit(25).Find(&contracts)
	return c.JSON(http.StatusOK, contracts)
}

func (h *Handler) BestValue(c echo.Context) error {
	var contracts []models.Contract
	h.db.Where("is_active = ? AND overall_value_score > ?", true, 20).
		Order("overall_value_score DESC").Limit(25).Find(&contracts)
	return c.JSON(http.StatusOK, contracts)
}

func (h *Handler) PeakArcs(c echo.Context) error {
	var arcs []models.CareerArc
	h.db.Order("peak_war DESC").Limit(25).Find(&arcs)
	return c.JSON(http.StatusOK, arcs)
}
