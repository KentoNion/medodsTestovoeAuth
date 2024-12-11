package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	dbHost     string `yaml:"db_host" env-required:"true"`
	dbPort     string `yaml:"db_port"`
	dbUser     string `yaml:"db_user" env-required:"true"`
	DbPass     string `yaml:"db_password" env-required:"true"`
	dbName     string `yaml:"db_name" env-required:"true"`
	dbSSL      string `yaml:"db_sslmode" env-required:"true"`
	serverHost string `yaml:"server_host" env-required:"true"`
	ServerPort string `yaml:"server_port" env-required:"true"`
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	//проверка существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "config file not found at %s", configPath)
	}

	var cfg *Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read config file")
	}

	return cfg, nil
}
