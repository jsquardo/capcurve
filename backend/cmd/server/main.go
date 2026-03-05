package main

import (
	"log"

	"github.com/jsquardo/capcurve/internal/config"
	"github.com/jsquardo/capcurve/internal/database"
	"github.com/jsquardo/capcurve/internal/handlers"
	"github.com/jsquardo/capcurve/internal/middleware"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	e := echo.New()
	e.HideBanner = false

	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(middleware.CORS())

	handlers.RegisterRoutes(e, db)

	log.Printf("🚀 CapCurve API starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
