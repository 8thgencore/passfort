package http

import (
	"log/slog"
	"strings"

	"github.com/8thgencore/passfort/internal/config"
	"github.com/8thgencore/passfort/internal/delivery/http/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHander.Register)
			user.POST("/login", authHandler.Login)
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
