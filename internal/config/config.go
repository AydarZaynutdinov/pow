package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		App               App          `yaml:"App"`
		AppServer         Server       `yaml:"Server"`
		MaintenanceServer Server       `yaml:"MaintenanceServer"`
		Logger            LoggerConfig `yaml:"Logger"`
		Cache             Cache        `yaml:"Cache"`
		Challenge         Challenge    `yaml:"Challenge"`
		QuotesList        []string     `yaml:"QuotesList"`
	}

	App struct {
		Mode string `yaml:"Mode" validate:"oneof=development prod"`
	}

	Server struct {
		Host             string        `yaml:"Host"`
		Port             int           `yaml:"Port"`
		ReadTimeout      time.Duration `yaml:"ReadTimeout"`
		WriteTimeout     time.Duration `yaml:"WriteTimeout"`
		IdleTimeout      time.Duration `yaml:"IdleTimeout"`
		ShutdownDuration time.Duration `yaml:"ShutdownDuration"`
	}

	Cache struct {
		Address  string `yaml:"Address"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
		PoolSize int    `yaml:"PoolSize"`
	}

	Challenge struct {
		Len        int           `yaml:"Len"`
		TTL        time.Duration `yaml:"TTL"`
		Difficulty int           `yaml:"Difficulty"`
	}
)

func New(configPath string, validator *validator.Validate) (*Config, error) {
	return parse(configPath, validator)
}

// Parses file by received parameter (filePath) to create app config
func parse(filePath string, validator *validator.Validate) (*Config, error) {
	filename, _ := filepath.Abs(filePath)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	err = validator.Struct(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %v", err)
	}

	return &config, nil
}
