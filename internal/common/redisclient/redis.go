package redisclient

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	rdb *redis.Client
)

func Setup(v *viper.Viper) {
	redisURL := v.GetString("redis.url")
	if redisURL == "" {
		log.Fatal().Msg("Failed to connect to Redis")
	}
	if err := InitRedisFromURL(redisURL); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
}

func InitRedisFromURL(redisURL string) error {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("parse redis url failed: %w", err)
	}

	rdb = redis.NewClient(opt)
	return nil
}

func GetClient() *redis.Client {
	return rdb
}
