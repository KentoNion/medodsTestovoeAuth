package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"os"
)

type DB struct {
	DbHost     string `yaml:"host" env-required:"true"`
	DbPort     string `yaml:"port"`
	DbUser     string `yaml:"user" env-required:"true"`
	DbPassword string `yaml:"password" env-required:"true"`
	DbName     string `yaml:"name" env-required:"true"`
	DbSSLMode  string `yaml:"sslmode" env-required:"true"`
}

type Server struct {
	ServerHost string `yaml:"host" env-required:"true"`
	ServerPort string `yaml:"port"`
}

type Config struct {
	Env    string `yaml:"env"` // local, dev, prod
	DB     DB     `yaml:"auth_db"`
	Server Server `yaml:"auth_server"`
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	//проверка существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "config file not found at %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read config file")
	}

	return &cfg, nil
}
