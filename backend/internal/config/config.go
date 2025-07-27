package config

import (
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Env           string               `yaml:"env"`
	ExploitRunner *ExploitRunnerConfig `yaml:"exploit_runner"`
	DB            *DBConfig            `yaml:"db"`
	Rabbit        *RabbitMQConfig      `yaml:"rabbitmq"`
}

type ExploitRunnerConfig struct {
	TeamIPs               []string      `yaml:"team_ips"`
	FlagFormat            string        `yaml:"flag_format"`
	RunDuration           time.Duration `yaml:"run_duration"`
	MaxConcurrentExploits int           `yaml:"max_concurrent_exploits"`
	ExploitDirectory      string        `yaml:"exploit_directory"`
	ExploitMaxWorkingTime time.Duration `yaml:"exploit_max_working_time"`
}

type DBConfig struct {
	MigrationsPath string `yaml:"migrations_path"`
}

type RabbitMQConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

const defaultConfigFilepath = "config/init_config.yml"

func MustParseConfig() *Config {
	cfg := Config{}

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = defaultConfigFilepath
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		panic("error reading config: " + err.Error())
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		panic("error parsing config: " + err.Error())
	}

	return &cfg
}
