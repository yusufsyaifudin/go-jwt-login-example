package server

import (
	"context"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/yusufsyaifudin/go-jwt-login-example/apidoc"
	"github.com/yusufsyaifudin/go-jwt-login-example/internal/app/user"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/auth"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/db"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/http"
)

var logger = log.With().Str("pkg", "server").Logger()
var stopped = false

type Config struct {
	ListenAddress   string
	ServerSecretKey string
	DB              db.Query
	Auth            auth.Auth
}

// Run will run the server and return error if error occurred.
func (config *Config) Run() error {
	parentCtx := context.Background()
	defer parentCtx.Done()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(Logger())

	// api documentation
	router.Use(static.Serve("/", apidoc.Static()))

	// to gracefully shutdown the server
	router.Use(func(ctx *gin.Context) {
		// if it's the case then don't receive anymore requests
		if stopped {
			ctx.Status(503)
			return
		}

		ctx.Next()
	})

	router.NoRoute(func(ctx *gin.Context) {
		response := map[string]interface{}{
			"error": map[string]interface{}{
				"message": "route not found",
			},
		}

		ctx.JSON(404, response)
		ctx.Abort()
	})

	router.NoMethod(func(ctx *gin.Context) {
		response := map[string]interface{}{
			"error": map[string]interface{}{
				"message": "method for this route not found",
			},
		}

		ctx.JSON(404, response)
		ctx.Abort()
	})

	userHandler := user.NewUserHandler(config.ServerSecretKey, config.DB, config.Auth)
	protectedMiddleware := http.ChainMiddleware(userHandler.MiddlewareAuthTokenCheck)

	userGroup := router.Group("/api/v1/user")
	userGroup.POST("/login", http.WrapGin(parentCtx, userHandler.LoginUserHandler))
	userGroup.POST("/register", http.WrapGin(parentCtx, userHandler.RegisterUserHandler))
	userGroup.GET("/profile", http.WrapGin(parentCtx, protectedMiddleware(userHandler.ProfileUserHandler)))

	// for debugging purpose
	for _, routeInfo := range router.Routes() {
		logger.Debug().
			Str("path", routeInfo.Path).
			Str("handler", routeInfo.Handler).
			Str("method", routeInfo.Method).
			Msg("registered routes")
	}

	return router.Run(config.ListenAddress)
}

// Shutdown this package
func (config *Config) Shutdown() {
	logger.Info().Msg("not receiving requests anymore")
	stopped = true
}
