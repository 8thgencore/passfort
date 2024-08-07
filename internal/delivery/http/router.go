package http

import (
	"log/slog"
	"strings"

	"github.com/8thgencore/passfort/internal/config"
	"github.com/8thgencore/passfort/internal/delivery/http/handler"
	"github.com/8thgencore/passfort/internal/delivery/http/helper"
	"github.com/8thgencore/passfort/internal/delivery/http/middleware"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

// NewRouter creates a new router with the given parameters
func NewRouter(
	log *slog.Logger,
	cfg *config.Config,
	tokenService service.TokenService,
	masterPasswordService service.MasterPasswordService,
	userHander handler.UserHandler,
	authHandler handler.AuthHandler,
	collectionHandler handler.CollectionHandler,
	secretHandler handler.SecretHandler,
	masterPasswordHandler handler.MasterPasswordHandler,
) (*Router, error) {
	// Disable debug mode in production
	if cfg.Env == config.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.Config{
		AllowOrigins: strings.Split(cfg.HTTP.AllowOrigins, ","),
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}

	// Init router
	router := gin.New()
	router.Use(sloggin.New(log), gin.Recovery(), cors.New(ginConfig))

	// Custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("user_role", helper.UserRoleValidator); err != nil {
			return nil, err
		}
		if err := v.RegisterValidation("secret_type", helper.SecretTypeValidator); err != nil {
			return nil, err
		}
	}

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware
	authMiddleware := middleware.AuthMiddleware(tokenService)
	adminMiddleware := middleware.AdminMiddleware()
	masterPasswordMiddleware := middleware.MasterPasswordMiddleware(masterPasswordService)

	// Endpoints
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Authentication Routes
			auth := v1.Group("/auth")
			{
				auth.POST("/login", authHandler.Login)
				auth.POST("/register", authHandler.Register)
				auth.POST("/register/confirm", authHandler.ConfirmRegistration)
				auth.POST("/register/resend-otp", authHandler.ResendOTPCode)
				auth.POST("/forgot-password", authHandler.ForgotPassword)
				auth.POST("/reset-password", authHandler.ResetPassword)
				auth.POST("/refresh-token", authHandler.RefreshToken)

				authUser := auth.Use(authMiddleware)
				{
					authUser.POST("/logout", authHandler.Logout)
					authUser.PUT("/change-password", authHandler.ChangePassword)
				}
			}

			// Master Password Routes
			masterPassword := v1.Group("/master-password").Use(authMiddleware)
			{
				masterPassword.POST("", masterPasswordHandler.CreateMasterPassword)
				masterPassword.PUT("", masterPasswordHandler.ChangeMasterPassword)
				masterPassword.POST("/activate", masterPasswordHandler.ActivateMasterPassword)
			}

			// User Routes
			usersGroup := v1.Group("/users")
			{
				users := usersGroup.Use(authMiddleware)
				{
					users.GET("/me", userHander.GetUserMe)
					users.GET("/:id", userHander.GetUser)
				}

				admin := users.Use(adminMiddleware)
				{
					admin.GET("", userHander.ListUsers)
					admin.PUT("/:id", userHander.UpdateUser)
					admin.DELETE("/:id", userHander.DeleteUser)
				}
			}

			// Collection Routes
			collectionsGroup := v1.Group("/collections")
			{
				collections := collectionsGroup.Use(authMiddleware).Use(masterPasswordMiddleware)
				{
					collections.GET("/me", collectionHandler.ListMeCollections)
					collections.POST("", collectionHandler.CreateCollection)
					collections.GET("/:collection_id", collectionHandler.GetCollection)
					collections.PUT("/:collection_id", collectionHandler.UpdateCollection)
					collections.DELETE("/:collection_id", collectionHandler.DeleteCollection)
				}

				// Nest the /secrets routes under /collections/:id
				secrets := collectionsGroup.Group("/:collection_id/secrets").Use(authMiddleware).Use(masterPasswordMiddleware)
				{
					secrets.GET("", secretHandler.ListMeSecrets)
					secrets.POST("", secretHandler.CreateSecret)
					secrets.GET("/:secret_id", secretHandler.GetSecret)
					secrets.PUT("/:secret_id", secretHandler.UpdateSecret)
					secrets.DELETE("/:secret_id", secretHandler.DeleteSecret)
				}
			}
		}
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
