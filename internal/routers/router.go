package routers

import (
	apiV1 "MyStonks-go/internal/routers/api/v1"
	"net/http"
	"time"

	_ "MyStonks-go/docs"

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

	apiv1 := r.Group("/api/v1")

	// user相关API
	apiv1.POST("/users", apiV1.CreateUser)
	apiv1.GET("/users/:id", apiV1.GetUser)

	return r
}
