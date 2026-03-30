package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jsquardo/capcurve/internal/models"
	"github.com/jsquardo/capcurve/internal/projection"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func newProjectionData(player models.Player, projection careerArcProjection) projectionData {
	return projectionData{
		Player:     newCareerArcPlayerItem(player),
		Projection: projection,
	}
}

func newProjectionPayload(result projection.Result) careerArcProjection {
	payload := careerArcProjection{
		Status:         result.Status,
		Eligible:       result.Eligible,
		Reason:         result.Reason,
		Points:         make([]careerArcProjectionPoint, 0, len(result.Points)),
		ConfidenceBand: make([]careerArcConfidenceBand, 0, len(result.ConfidenceBand)),
		Comparables:    make([]careerArcComparablePlayer, 0, len(result.Comparables)),
	}

	for _, point := range result.Points {
		payload.Points = append(payload.Points, careerArcProjectionPoint{
			Year:         point.Year,
			Age:          point.Age,
			ValueScore:   point.ValueScore,
			IsProjection: true,
		})
	}

	for _, band := range result.ConfidenceBand {
		payload.ConfidenceBand = append(payload.ConfidenceBand, careerArcConfidenceBand{
			Year:  band.Year,
			Lower: band.Lower,
			Upper: band.Upper,
		})
	}

	for _, comparable := range result.Comparables {
		payload.Comparables = append(payload.Comparables, careerArcComparablePlayer{
			PlayerID: comparable.PlayerID,
			MLBID:    comparable.MLBID,
			FullName: comparable.FullName,
			Position: comparable.Position,
		})
	}

	return payload
}

func (h *Handler) buildPlayerProjectionPayload(player models.Player, history []models.SeasonStat) (careerArcProjection, error) {
	service := projection.NewService()

	if !player.Active || len(history) == 0 {
		return newProjectionPayload(service.Build(player, history, nil, nil)), nil
	}

	candidates, candidateStats, err := h.loadProjectionComparableCandidates(player.ID)
	if err != nil {
		return careerArcProjection{}, err
	}

	return newProjectionPayload(service.Build(player, history, candidates, candidateStats)), nil
}

func (h *Handler) loadProjectionComparableCandidates(playerID uint) ([]models.Player, []models.SeasonStat, error) {
	var candidates []models.Player
	if err := h.db.
		Where("id <> ? AND active = ?", playerID, false).
		Find(&candidates).Error; err != nil {
		return nil, nil, err
	}

	if len(candidates) == 0 {
		return []models.Player{}, []models.SeasonStat{}, nil
	}

	candidateIDs := make([]uint, 0, len(candidates))
	for _, candidate := range candidates {
		candidateIDs = append(candidateIDs, candidate.ID)
	}

	var candidateStats []models.SeasonStat
	if err := h.db.
		Where("player_id IN ?", candidateIDs).
		Order("player_id ASC, year ASC, id ASC").
		Find(&candidateStats).Error; err != nil {
		return nil, nil, err
	}

	return candidates, candidateStats, nil
}

func (h *Handler) GetProjection(c echo.Context) error {
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

	var history []models.SeasonStat
	if err := h.db.Where("player_id = ?", player.ID).Order("year ASC, id ASC").Find(&history).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	projectionPayload, err := h.buildPlayerProjectionPayload(player, history)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, projectionResponse{
		Data: newProjectionData(player, projectionPayload),
	})
}
