package routers

import (
	_ "MyStonks-go/docs"
	"MyStonks-go/internal/common/middleware"
	v1 "MyStonks-go/internal/routers/api/v1"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Router struct {
	Eng     *gin.Engine
	taskApi *v1.TaskApi
	userApi *v1.UserApi
}

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

func InitRouter(taskApi *v1.TaskApi, userApi *v1.UserApi) *Router {

	router := &Router{}
	router.taskApi = taskApi
	router.userApi = userApi

	r := gin.New()
	r.Use(GinLogger())
	r.Use(GinRecovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 可设置具体前端地址，例如 http://localhost:3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := r.Group("/api/v1")
	apiV1Protected := apiV1.Group("", middleware.AuthMiddleware())
	{
		apiV1Auth := apiV1.Group("/auth")
		apiV1Auth.GET("/nonce", router.userApi.GetNonce) // 生成nonce
		apiV1Auth.POST("/login", router.userApi.Login)   // 登录
		apiV1Auth.POST("/logout", router.userApi.Logout) // 登出
		apiV1Auth.POST("/refresh", router.userApi.RefreshToken)
	}
	{
		// 任务相关接口（部分需要登录）
		apiV1.GET("/tasks", router.taskApi.GetTasksByCategory)           // category=daily|newbie|other=
		apiV1.GET("/task/complete", router.taskApi.CheckCompleteTask)    // ?user_id=1&task_id=101
		apiV1.POST("/task/progress", router.taskApi.UpdateTaskProgress)  // ?user_id=1&task_id=201&progress=50
		apiV1.GET("/task/stonks/trade", router.taskApi.CheckStonksTrade) // ?task_id=201
		apiV1.GET("/leaderboard", router.taskApi.GetLeaderboard)
		apiV1Protected.GET("/user/task", router.taskApi.GetUserInfoTask)
		apiV1Protected.GET("/task/finish", router.taskApi.FinishTask)
		apiV1.GET("/user/rank", router.taskApi.GetUserRank)
	}
	router.Eng = r
	return router
}
