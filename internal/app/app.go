package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/8thgencore/passfort/docs"
	"github.com/8thgencore/passfort/internal/config"
	"github.com/8thgencore/passfort/internal/database"
	"github.com/8thgencore/passfort/internal/delivery/http"
	"github.com/8thgencore/passfort/internal/delivery/http/handler"
	"github.com/8thgencore/passfort/internal/repository/cache/redis"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres"
	authService "github.com/8thgencore/passfort/internal/service/auth"
	"github.com/8thgencore/passfort/internal/service/token/paseto"
	userService "github.com/8thgencore/passfort/internal/service/user"
	"github.com/8thgencore/passfort/pkg/logger/slogpretty"
)

// @title						PassFort API
// @version					1.0
// @description				This is a simple RESTful Password Manager Service API written in Go using Gin web framework, PostgreSQL database, and Redis cache.
//
// @contact.name				Tom Jerry
// @contact.url				https://github.com/8thgencore/passfort
// @contact.email				test@gmail.com
//
// @license.name				MIT
// @license.url				https://opensource.org/licenses/MIT
//
// @host						api.example.com
// @BasePath					/v1
// @schemes					http https
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func Run(configPath string) {
	// Load configuration
	cfg, err := config.NewConfig("./config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Set up logger based on configuration
	log := newSlogLogger(cfg.Log.Slog)

	// Log information about the start of the application
	log.Info("starting url-shortener", slog.String("env", string(cfg.Env)))
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
		log.Error("Error migrating database", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Successfully migrated the database")

	// Init cache service
	cache, err := redis.New(ctx, &cfg.Cache)
	if err != nil {
		log.Error("Error initializing cache connection", "error", err.Error())
		os.Exit(1)
	}
	defer cache.Close()

	log.Info("Successfully connected to the cache server")

	// Init token service
	token, _ := paseto.New(&cfg.Token)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	// Dependency injection
	// User
	userRepo := postgres.NewUserRepository(db)
	userService := userService.NewUserService(log, userRepo, cache)
	userHandler := handler.NewUserHandler(userService)

	// Auth
	authService := authService.NewAuthService(log, userRepo, cache, token)
	authHandler := handler.NewAuthHandler(authService)

	// Init router
	router, err := http.NewRouter(log, cfg, token, *userHandler, *authHandler)
	if err != nil {
		log.Error("Error initializing router", "error", err.Error())
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	log.Info("Starting the HTTP server", "listen_address", listenAddr)

	err = router.Serve(listenAddr)
	if err != nil {
		log.Error("Error starting the HTTP server", "error", err.Error())
		os.Exit(1)
	}
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
			WithTimeLayout(c.Pretty.TimeLayout).
			WithFieldsFormat(c.Pretty.FieldsFormat)
	case "json":
		h = slog.NewJSONHandler(w, o)
	case "text":
		h = slog.NewTextHandler(w, o)
	}
	return slog.New(h)
}
