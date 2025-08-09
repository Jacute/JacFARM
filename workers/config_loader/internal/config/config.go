package config

import (
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ExploitRunner *ExploitRunnerConfig `yaml:"exploit_runner"`
	FlagSender    *FlagSenderConfig    `yaml:"flag_sender"`
	DB            *DBConfig
}

type ExploitRunnerConfig struct {
	TeamIPs               []string      `yaml:"team_ips"`
	FlagFormat            string        `yaml:"flag_format"`
	RunDuration           time.Duration `yaml:"run_duration"`
	MaxConcurrentExploits int           `yaml:"max_concurrent_exploits"`
	ExploitDirectory      string        `yaml:"exploit_directory"`
	ExploitMaxWorkingTime time.Duration `yaml:"exploit_max_working_time"`
}

type FlagSenderConfig struct {
	Plugin        string        `yaml:"plugin"`
	PluginDir     string        `yaml:"plugin_directory"`
	JuryFlagURL   string        `yaml:"jury_flag_url_or_host"`
	Token         string        `yaml:"token"`
	FlagTTL       time.Duration `yaml:"flag_ttl"`
	SubmitTimeout time.Duration `yaml:"submit_timeout"`
	SubmitPeriod  time.Duration `yaml:"submit_period"`
	SubmitLimit   int           `yaml:"submit_limit"`
}

type DBConfig struct {
	Username string `envconfig:"PG_USERNAME"`
	Password string `envconfig:"PG_PASSWORD"`
	Host     string `envconfig:"PG_HOST"`
	Port     int    `envconfig:"PG_PORT"`
	DBName   string `envconfig:"PG_DB_NAME"`
}

const defaultConfigFilepath = "./config.yml"

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

	var dbCfg DBConfig
	if err := envconfig.Process("", &dbCfg); err != nil {
		panic("error loading env: " + err.Error())
	}
	cfg.DB = &dbCfg

	return &cfg
}
