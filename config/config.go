package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"gamification-app/pkg/database"
)

type Config struct {
	Port        string `split_words:"true" default:"9393"`
	Environment string `split_words:"true" default:"development"`
	Debug       bool   `split_words:"true" default:"true"`

	DBRead  database.DBConfig `split_words:"true"`
	DBWrite database.DBConfig `split_words:"true"`

	RedisAddr     string `split_words:"true" default:"localhost:6381"`
	RedisPassword string `split_words:"true" default:""`
	RedisDB       int    `split_words:"true" default:"0"`

	JWTSecret      string `split_words:"true" default:"change-me-in-production"`
	JWTExpireHours int    `split_words:"true" default:"24"`
}

func Read() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warning: .env file not found, using system env")
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("error processing env config: %w", err)
	}

	return &cfg, nil
}
