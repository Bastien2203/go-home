package config

type Config struct {
	BrokerUrl     string `env:"BROKER_URL, required"`
	SqliteDbPath  string `env:"SQLITE_DB_PATH,required"`
	ApiPort       int    `env:"API_PORT,default=8080"`
	SessionSecret string `env:"SESSION_SECRET"`
	AppEnv        AppEnv `env:"APP_ENV,default=dev"`
}

type AppEnv string

const (
	Production AppEnv = "production"
	Dev        AppEnv = "dev"
)
