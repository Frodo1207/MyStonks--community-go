package routers

import (
	_ "MyStonks-go/docs"
	"MyStonks-go/internal/common/middleware"
	v1 "MyStonks-go/internal/routers/api/v1"
	"net/http"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		evt := log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP())

		if raw != "" {
			evt.Str("query", raw)
		}

		evt.Msg("Incoming request")
	}
}

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Interface("error", err).
					Str("path", c.Request.URL.Path).
					Msg("Panic recovered")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(GinLogger())
	r.Use(GinRecovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := r.Group("/api/v1")
	apiV1Protected := apiV1.Group("", middleware.AuthMiddleware())
	{
		apiV1Auth := apiV1.Group("/auth")
		apiV1Auth.GET("/nonce", v1.GetNonce)
		apiV1Auth.POST("/login", v1.Login)
		apiV1Auth.POST("/logout", v1.Logout)
		apiV1Auth.POST("/refresh", v1.RefreshToken)
	}
	{
		apiV1Protected.GET("/user", v1.GetUserInfo)
	}
	return r
}
