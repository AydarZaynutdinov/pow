package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	serverConfigFile = "SERVER_CONFIG_FILE"
	clientConfigFile = "CLIENT_CONFIG_FILE"
)

type ServerConfig struct {
	Server    Server    `yaml:"server" env-prefix:"SERVER_"`
	PoW       PoWConfig `yaml:"pow" env-prefix:"POW_"`
	QuoteList []string  `yaml:"quote_list" env:"QUOTE_LIST"`
}

type Server struct {
	Address         string        `yaml:"address" env:"SERVER_ADDRESS" env-default:":8080"`
	LogLevel        string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	ShutDownTimeout time.Duration `yaml:"shut_down_timeout" env:"SHUT_DOWN_TIMEOUT" env-default:"10s"`
}

type PoWConfig struct {
	WorkersCount   int           `yaml:"workers_count" env:"POW_WORKERS_COUNT" env-default:"10"`
	HandlerTimeout time.Duration `yaml:"handler_timeout" env:"POW_HANDLER_TIMEOUT" env-default:"5s"`
	Complexity     uint8         `yaml:"complexity" env:"POW_COMPLEXITY" env-default:"2"`
}

func ParseServer() (*ServerConfig, error) {
	if path, ok := os.LookupEnv(serverConfigFile); ok {
		return parseServerByFile(path)
	}
	return parseServerByEnv()
}

func parseServerByFile(path string) (*ServerConfig, error) {
	cfg := &ServerConfig{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseServerByEnv() (*ServerConfig, error) {
	cfg := &ServerConfig{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type ClientConfig struct {
	Address  string `yaml:"address" env:"ADDRESS" env-default:"127.0.0.1:8080"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	RPS      int    `yaml:"rps" env:"RPS" env-default:"1"`
	Total    int    `yaml:"total" env:"TOTAL" env-default:"100"`
}

func ParseClient() (*ClientConfig, error) {
	if path, ok := os.LookupEnv(clientConfigFile); ok {
		return parseClientByFile(path)
	}
	return parseClientByEnv()
}

func parseClientByFile(path string) (*ClientConfig, error) {
	cfg := &ClientConfig{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseClientByEnv() (*ClientConfig, error) {
	cfg := &ClientConfig{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
