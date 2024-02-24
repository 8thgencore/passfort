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

func NewRouter(
	log *slog.Logger,
	cfg *config.Config,
	token service.TokenService,
	userHander handler.UserHandler,
	authHandler handler.AuthHandler,
	collectionHandler handler.CollectionHandler,
	secretHandler handler.SecretHandler,
) (*Router, error) {
	// Disable debug mode in production
	if cfg.Env == config.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	allowOrigins := cfg.HTTP.AllowOrigins
	originsList := strings.Split(allowOrigins, ",")
	ginConfig.AllowOrigins = originsList

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
	authMiddleware := middleware.AuthMiddleware(token)
	adminMiddleware := middleware.AdminMiddleware()

	// Endpoints
	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)

			authUser := auth.Group("/").Use(authMiddleware)
			{
				authUser.PUT("/change-password", authHandler.ChangePassword)
			}
		}
		user := v1.Group("/users")
		{
			user.POST("/register", userHander.Register)

			authUser := user.Group("/").Use(authMiddleware)
			{
				authUser.GET("/me", userHander.GetUserMe)
				authUser.GET("/:id", userHander.GetUser)

				admin := authUser.Use(adminMiddleware)
				{
					admin.GET("/", userHander.ListUsers)
					admin.PUT("/:id", userHander.UpdateUser)
					admin.DELETE("/:id", userHander.DeleteUser)
				}
			}
		}
		collection := v1.Group("/collections")
		{
			authCollection := collection.Use(authMiddleware)
			{
				authCollection.GET("/me", collectionHandler.ListMeCollections)
				authCollection.POST("/", collectionHandler.CreateCollection)
				authCollection.GET("/:collection_id", collectionHandler.GetCollection)
				authCollection.PUT("/:collection_id", collectionHandler.UpdateCollection)
				authCollection.DELETE("/:collection_id", collectionHandler.DeleteCollection)

				// Nest the /secrets routes under /collections/:id
				authSecret := collection.Group("/:collection_id/secrets").Use(authMiddleware)
				{
					authSecret.GET("/", secretHandler.ListMeSecrets)
					authSecret.POST("/", secretHandler.CreateSecret)
					authSecret.GET("/:secret_id", secretHandler.GetSecret)
					// authSecret.PUT("/:secret_id", secretHandler.UpdateSecret)
					authSecret.DELETE("/:secret_id", secretHandler.DeleteSecret)
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
