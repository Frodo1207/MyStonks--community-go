package server

import (
	"MyStonks-go/internal/common/redisclient"
	"MyStonks-go/internal/models"
	"MyStonks-go/internal/routers"
	v1 "MyStonks-go/internal/routers/api/v1"
	"MyStonks-go/internal/service"
	"MyStonks-go/internal/store"
	"fmt"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
}

func StartServer(v *viper.Viper) {
	var (
		env            = v.GetString("common.env")
		mode           = v.GetString("server.RunMode")
		endPoint       = fmt.Sprintf(":%d", v.GetInt("server.HttpPort"))
		readTimeout    = time.Duration(v.GetInt("server.ReadTimeout")) * time.Second
		writeTimeout   = time.Duration(v.GetInt("server.WriteTimeout")) * time.Second
		maxHeaderBytes = 1 << 20
	)

	gin.SetMode(mode)
	db := models.Setup(v)
	redisclient.Setup(v)

	userStore := store.NewUserStore(db)
	srv := service.NewUserSrv(userStore)
	userApi := v1.NewUserApi(srv)

	taskStore := store.NewTaskStore(db)
	taskService := service.NewTaskService(taskStore, userStore)
	taskApi := v1.NewTaskApi(taskService)

	eventStore := store.NewEventStore(db)
	eventService := service.NewEventSrv(eventStore)
	eventApi := v1.NewEventApi(eventService)

	router_ := routers.InitRouter(taskApi, userApi, eventApi)

	if env != "prod" {
		router_.Eng.Use(cors.Default())
	}
	server := &http.Server{
		Addr:           endPoint,
		Handler:        router_.Eng,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	log.Info().Msgf("start server at %s", endPoint)
	if err := server.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("server start failed")
	}
}
