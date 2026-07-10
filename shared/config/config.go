package config

import "strings"

type Config struct {
	BrokerUrl      string `env:"BROKER_URL,required"`
	SqliteDbPath   string `env:"SQLITE_DB_PATH,required"`
	ApiPort        int    `env:"API_PORT,default=8080"`
	SessionSecret  string `env:"SESSION_SECRET,required"`
	AppEnv         AppEnv `env:"ENV,default=dev"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS,default="`
}

func (c *Config) ParsedAllowedOrigins() []string {
	if c.AllowedOrigins == "" {
		return nil
	}
	parts := strings.Split(c.AllowedOrigins, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

type PluginConfig struct {
	BrokerUrl string `env:"BROKER_URL,required"`
	AppEnv    AppEnv `env:"ENV,default=dev"`
}

type AppEnv string

const (
	Production AppEnv = "production"
	Dev        AppEnv = "dev"
)
