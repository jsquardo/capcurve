package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jsquardo/capcurve/internal/config"
	"github.com/jsquardo/capcurve/internal/database"
	"github.com/jsquardo/capcurve/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
)

func main() {
	// Logger
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	// Config
	cfg := config.Load()

	// Database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	// Routes
	e := echo.New()
	handlers.RegisterRoutes(e, db)

	slog.Info("server starting", "port", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
