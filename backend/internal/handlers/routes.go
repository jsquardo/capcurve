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
	players.GET("/:id/projection", h.GetProjection)

	api.GET("/players/:id/career-arc", h.GetCareerArc)
	api.GET("/leaderboards", h.GetLeaderboards)
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
