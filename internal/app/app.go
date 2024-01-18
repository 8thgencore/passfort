package app

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/8thgencore/passfort/internal/config"
	"github.com/8thgencore/passfort/internal/database"
	"github.com/8thgencore/passfort/pkg/logger/slogpretty"
)

func Run(configPath string) {
	// Load configuration
	cfg, err := config.NewConfig("./config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Set up logger based on configuration
	log := newSlogLogger(cfg.Log.Slog)

	// Log information about the start of the application
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init database
	ctx := context.Background()
	db, err := database.New(ctx, &cfg.Database)
	if err != nil {
		log.Error("Error initializing database connection", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	log.Info("Successfully connected to the database", "db", cfg.Database.Connection)

	// Migrate database
	err = db.Migrate()
	if err != nil {
		log.Error("Error migrating database", "err", err.Error())
		os.Exit(1)
	}

	log.Info("Successfully migrated the database")
}

func newSlogLogger(c config.Slog) *slog.Logger {
	o := &slog.HandlerOptions{Level: c.Level, AddSource: c.AddSource}
	w := os.Stdout
	var h slog.Handler

	switch c.Format {
	case "pretty":
		h = slogpretty.NewHandler().
			WithAddSource(c.AddSource).
			WithLevel(c.Level).
			WithLevelEmoji(c.Pretty.Emoji).
			WithFieldsFormat(c.Pretty.FieldsFormat)
	case "json":
		h = slog.NewJSONHandler(w, o)
	case "text":
		h = slog.NewTextHandler(w, o)
	}
	return slog.New(h)
}
