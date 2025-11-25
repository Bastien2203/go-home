package config

import (
	"context"

	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func LoadFromEnv(ctx context.Context) *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := &Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return cfg
}
