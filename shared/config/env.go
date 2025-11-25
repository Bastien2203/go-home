package config

import (
	"context"
	"os"

	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func LoadFromEnv(ctx context.Context) *Config {
	env := os.Getenv("ENV")

	// No .env file for productions
	if env != "PROD" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	cfg := &Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return cfg
}

func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
