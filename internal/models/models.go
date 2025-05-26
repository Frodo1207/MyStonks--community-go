package models

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup(v *viper.Viper) {
	var (
		mysqlHost     = v.GetString("mysql.host")
		mysqlPort     = v.GetInt("mysql.port")
		mysqlUser     = v.GetString("mysql.user")
		mysqlPassword = v.GetString("mysql.password")
		mysqlDatabase = v.GetString("mysql.database")
		err           error
	)
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlUser,
		mysqlPassword,
		mysqlHost,
		mysqlPort,
		mysqlDatabase)))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MySQL")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get SQL DB")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	// db.AutoMigrate(&User{})
}
