package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jsquardo/capcurve/internal/config"
	"github.com/jsquardo/capcurve/internal/database"
	"github.com/jsquardo/capcurve/internal/handlers"
	"github.com/jsquardo/capcurve/internal/ingestion"
	"github.com/jsquardo/capcurve/internal/syncjob"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	syncStatus := syncjob.NewStatusStore(cfg.SyncEnabled)

	if cfg.SyncEnabled {
		location, err := time.LoadLocation(cfg.SyncTimeZone)
		if err != nil {
			slog.Error("failed to load sync timezone", "timezone", cfg.SyncTimeZone, "err", err)
			os.Exit(1)
		}

		syncService := syncjob.NewService(db, ingestion.NewService(db, nil, nil), syncjob.Options{
			Logger: slog.Default(),
			Schedule: syncjob.Schedule{
				Hour:    cfg.SyncHour,
				Minute:  cfg.SyncMinute,
				Weekday: time.Weekday(cfg.SyncWeekday),
			},
			Location: location,
			Status:   syncStatus,
		})

		go syncService.Start(ctx)
	}

	// Routes
	e := echo.New()
	handlers.RegisterRoutes(e, db, syncStatus)

	slog.Info("server starting", "port", cfg.Port)

	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "err", err)
		os.Exit(1)
	}
}
