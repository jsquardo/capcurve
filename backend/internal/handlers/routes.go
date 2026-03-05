package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	h := &Handler{db: db}

	e.GET("/health", h.HealthCheck)

	api := e.Group("/api/v1")

	players := api.Group("/players")
	players.GET("", h.ListPlayers)
	players.GET("/:id", h.GetPlayer)
	players.GET("/search", h.SearchPlayers)

	api.GET("/players/:id/career-arc", h.GetCareerArc)
	api.GET("/players/:id/contracts", h.GetPlayerContracts)
	api.GET("/contracts/:id", h.GetContract)

	leaderboards := api.Group("/leaderboards")
	leaderboards.GET("/most-overpaid", h.MostOverpaid)
	leaderboards.GET("/best-value", h.BestValue)
	leaderboards.GET("/peak-arcs", h.PeakArcs)
}

type Handler struct {
	db *gorm.DB
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status":  "ok",
		"service": "capcurve-api",
	})
}
