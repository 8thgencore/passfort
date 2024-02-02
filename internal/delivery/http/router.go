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

	router := gin.New()
	router.Use(sloggin.New(log), gin.Recovery(), cors.New(ginConfig))

	// Custom validator
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		if err := v.RegisterValidation("user_role", helper.UserRoleValidator); err != nil {
			return nil, err
		}
	}

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			authUser := auth.Group("/").Use(middleware.AuthMiddleware(token))
			{
				authUser.PUT("/change-password", authHandler.ChangePassword)
			}
		}
		user := v1.Group("/users")
		{
			authUser := user.Group("/").Use(middleware.AuthMiddleware(token))
			{
				authUser.GET("/", userHander.ListUsers)
				authUser.GET("/:id", userHander.GetUser)

				admin := authUser.Use(middleware.AdminMiddleware())
				{
					admin.PUT("/:id", userHander.UpdateUser)
					admin.DELETE("/:id", userHander.DeleteUser)
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
