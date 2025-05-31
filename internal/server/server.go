package server

import (
	"MyStonks-go/internal/common/redisclient"
	"MyStonks-go/internal/models"
	"MyStonks-go/internal/routers"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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
	{
		models.Setup(v)
		redisclient.Setup(v)
	}
	gin.SetMode(mode)
	routersInit := routers.InitRouter()
	if env != "prod" {
		routersInit.Use(cors.Default())
	}
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	log.Info().Msgf("start server at %s", endPoint)
	if err := server.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("server start failed")
	}
}
