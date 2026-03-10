package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/jsquardo/capcurve/internal/config"
	"github.com/jsquardo/capcurve/internal/database"
	"github.com/jsquardo/capcurve/internal/ingestion"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/ingest <mlb_player_id>")
		os.Exit(1)
	}

	playerID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid mlb player id %q\n", os.Args[1])
		os.Exit(1)
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	// The CLI is intentionally thin: all ingestion behavior lives in internal/ingestion.
	service := ingestion.NewService(db, nil, nil)
	player, err := service.SyncPlayer(context.Background(), playerID)
	if err != nil {
		slog.Error("ingestion failed", "err", err, "mlb_id", playerID)
		os.Exit(1)
	}

	slog.Info("player ingested", "mlb_id", player.MLBID, "db_id", player.ID, "season_count", len(player.SeasonStats))
}
