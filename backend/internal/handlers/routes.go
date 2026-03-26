package handlers

import (
	appmiddleware "github.com/jsquardo/capcurve/internal/middleware"
	"github.com/jsquardo/capcurve/internal/syncjob"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, syncStatus *syncjob.StatusStore, adminSecret string) {
	h := &Handler{db: db, syncStatus: syncStatus, adminSecret: adminSecret}

	e.GET("/health", h.HealthCheck)

	api := e.Group("/api/v1", appmiddleware.CORS())
	api.GET("/admin/dashboard", h.AdminDashboard)

	players := api.Group("/players")
	players.GET("", h.ListPlayers)
	players.GET("/:id", h.GetPlayer)

	api.GET("/players/:id/career-arc", h.GetCareerArc)
	api.GET("/players/:id/contracts", h.GetPlayerContracts)
	api.GET("/contracts/:id", h.GetContract)

	leaderboards := api.Group("/leaderboards")
	leaderboards.GET("/most-overpaid", h.MostOverpaid)
	leaderboards.GET("/best-value", h.BestValue)
	leaderboards.GET("/peak-arcs", h.PeakArcs)
}

type Handler struct {
	db          *gorm.DB
	syncStatus  SyncStatusSource
	adminSecret string
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status":  "ok",
		"service": "capcurve-api",
	})
}
