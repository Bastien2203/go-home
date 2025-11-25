package config

type Config struct {
	BrokerUrl        string `env:"BROKER_URL, required"`
	SqliteDbPath     string `env:"SQLITE_DB_PATH,required"`
	ApiPort          int    `env:"API_PORT,default=8080"`
	PluginFolderPath string `env:"PLUGIN_FOLDER_PATH,required"`
}
