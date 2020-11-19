package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// AppConfig defines the config from application yaml
type AppConfig struct {
	Logging struct {
		Debug bool `yaml:"debug" env:"MASTERDATA_LOG_LEVEL"`
	} `yaml:"logging"`
	Database struct {
		URL                   string        `yaml:"url" env:"DATASOURCE_URL"`
		MaxIdleConnections    int           `yaml:"idleConnections" env:"MAX_IDLE_CONNECTIONS"`
		MaxOpenConnections    int           `yaml:"openConnections" env:"MAX_OPEN_CONNECTIONS"`
		MaxConnectionLifetime time.Duration `yaml:"connectionLifetime" env:"MAX_CONNECTION_LIFETIME"`
	} `yaml:"database"`
	Cache struct {
		Size int           `yaml:"size" env:"CACHE_SIZE"`
		TTL  time.Duration `yaml:"ttl" env:"CACHE_TTL"`
	} `yaml:"cache"`
}

// LoadAppConfig builds config for database and returns a DbConfig struct
func LoadAppConfig() AppConfig {

	var conf AppConfig
	err := cleanenv.ReadConfig("application.yaml", &conf)
	if err != nil {
		panic(err)
	}
	return conf
}
