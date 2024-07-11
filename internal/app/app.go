package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/8thgencore/passfort/docs"
	mailGrpc "github.com/8thgencore/passfort/internal/clients/mail/grpc"
	"github.com/8thgencore/passfort/internal/config"
	"github.com/8thgencore/passfort/internal/database"
	"github.com/8thgencore/passfort/internal/delivery/http"
	"github.com/8thgencore/passfort/internal/delivery/http/handler"
	"github.com/8thgencore/passfort/internal/repository/cache/redis"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres"
	authSvc "github.com/8thgencore/passfort/internal/service/auth"
	collectionSvc "github.com/8thgencore/passfort/internal/service/collection"
	masterPasswordSvc "github.com/8thgencore/passfort/internal/service/master_password"
	otpSvc "github.com/8thgencore/passfort/internal/service/otp"
	secretSvc "github.com/8thgencore/passfort/internal/service/secret"
	tokenSvc "github.com/8thgencore/passfort/internal/service/token"
	userSvc "github.com/8thgencore/passfort/internal/service/user"
	"github.com/8thgencore/passfort/pkg/logger/sl"
	"github.com/8thgencore/passfort/pkg/logger/slogpretty"
	"github.com/hibiken/asynq"
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
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Set up logger based on configuration
	log := newSlogLogger(cfg.Log.Slog)

	// Log information about the start of the application
	log.Info("starting passfort", slog.String("env", string(cfg.Env)))
	log.Debug("debug messages are enabled")

	// Init database
	ctx := context.Background()
	db, err := database.New(ctx, &cfg.Database)
	if err != nil {
		log.Error("Error initializing database connection", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	log.Info("Successfully connected to the database", "db", cfg.Database.Connection)

	// Migrate database
	err = db.Migrate()
	if err != nil {
		log.Error("Error migrating database", sl.Err(err))
		os.Exit(1)
	}

	log.Info("Successfully migrated the database")

	// Init cache service
	cache, err := redis.New(ctx, &cfg.Cache)
	if err != nil {
		log.Error("Error initializing cache connection", sl.Err(err))
		os.Exit(1)
	}
	defer cache.Close()

	log.Info("Successfully connected to the cache server")

	// Init token service
	tokenService := tokenSvc.NewTokenService(log, cfg.Token.SigningKey, cfg.Token.AccessTokenTTL, cfg.Token.RefreshTokenTTL, cache)

	// Otp service
	otpService := otpSvc.NewOtpService(log, cache)

	// Register external microservices
	mailClient, err := mailGrpc.New(ctx,
		log,
		cfg.Clients.Mail.Address,
		cfg.Clients.Mail.Timeout,
		cfg.Clients.Mail.RetriesCount,
	)
	if err != nil {
		log.Error("Error initializing mail client", "error", err)
		os.Exit(1)
	}

	log.Info("Successfully initializing the mail client")

	// Init asynq client and server and register task handlers
	asynqCfg := asynq.RedisClientOpt{Addr: cfg.Cache.Addr, Password: cfg.Cache.Password}
	asynqClient := asynq.NewClient(asynqCfg)
	asynqServer := asynq.NewServer(asynqCfg, asynq.Config{Concurrency: 10})

	// Dependency injection
	// User
	userRepo := postgres.NewUserRepository(db)
	userService := userSvc.NewUserService(log, userRepo, cache)
	userHandler := handler.NewUserHandler(userService)

	// Auth
	authService := authSvc.NewAuthService(log, userRepo, cache, tokenService, otpService, mailClient)
	authHandler := handler.NewAuthHandler(authService)

	// Collection
	collectionRepo := postgres.NewCollectionRepository(db)
	collectionService := collectionSvc.NewCollectionService(log, collectionRepo)
	collectionHandler := handler.NewCollectionHandler(collectionService)

	// Secret
	secretRepo := postgres.NewSecretRepository(db)
	secretService := secretSvc.NewSecretService(log, secretRepo, collectionRepo, cache, asynqClient)
	secretHandler := handler.NewSecretHandler(secretService)

	// MasterPassword
	masterPasswordService := masterPasswordSvc.NewMasterPasswordService(log, userRepo, cache, *secretService, cfg.MasterPassword.MasterPasswordTTL)
	masterPasswordHandler := handler.NewMasterPasswordHandler(masterPasswordService)

	mux := asynq.NewServeMux()
	mux.HandleFunc(secretSvc.TypeReencryptSecrets, secretService.HandleReencryptSecretsTask)

	go func() {
		log.Info("Starting asynq server")
		if err := asynqServer.Run(mux); err != nil {
			log.Error("Could not run asynq server:", "error", err)
			os.Exit(1)
		}
	}()

	// Init router
	router, err := http.NewRouter(
		log,
		cfg,
		tokenService,
		masterPasswordService,
		*userHandler,
		*authHandler,
		*collectionHandler,
		*secretHandler,
		*masterPasswordHandler,
	)
	if err != nil {
		log.Error("Error initializing router", sl.Err(err))
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	log.Info("Starting the HTTP server", "listen_address", listenAddr)

	err = router.Serve(listenAddr)
	if err != nil {
		log.Error("Error starting the HTTP server", sl.Err(err))
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
