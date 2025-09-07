package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Config struct {
	Env        string `envconfig:"ENV"`
	ApiKey     string `envconfig:"JACFARM_API_KEY"`
	ExploitDir string `envconfig:"EXPLOIT_DIR"`
	DB         *DBConfig
	Rabbit     *RabbitMQConfig
	HTTP       *HTTPConfig
}

type HTTPConfig struct {
	Host         string        `envconfig:"HTTP_HOST"`
	Port         int           `envconfig:"HTTP_PORT"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `envconfig:"HTTP_IDLE_TIMEOUT"`
	CORS         *CORSConfig
}

type CORSConfig struct {
	AllowedOrigins []string `envconfig:"HTTP_CORS_ALLOWED_ORIGINS"`
}

type DBConfig struct {
	Username string `envconfig:"PG_USERNAME"`
	Password string `envconfig:"PG_PASSWORD"`
	Host     string `envconfig:"PG_HOST"`
	Port     int    `envconfig:"PG_PORT"`
	DBName   string `envconfig:"PG_DB_NAME"`
}

type RabbitMQConfig struct {
	Host           string `envconfig:"RABBITMQ_HOST"`
	Port           int    `envconfig:"RABBITMQ_PORT"`
	Username       string `envconfig:"RABBITMQ_USERNAME"`
	Password       string `envconfig:"RABBITMQ_PASSWORD"`
	ManagementHost string `envconfig:"RABBITMQ_MANAGEMENT_HOST"`
	ManagementPort int    `envconfig:"RABBITMQ_MANAGEMENT_PORT"`
}

func MustParseConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic("error loading env: " + err.Error())
	}
	if cfg.Env == "" {
		cfg.Env = "prod"
	}

	return &cfg
}
