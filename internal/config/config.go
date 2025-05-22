package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App      AppConfig
	Logger   LoggerConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port int
}

type LoggerConfig struct {
	Level  string
	Format string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	log.Printf("Configuration loaded successfully from %s", path)
	return &cfg, nil
}
