package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"os"
)

type DB struct {
	DbHost     string `yaml:"db_host" env-required:"true"`
	DbPort     string `yaml:"db_port"`
	DbUser     string `yaml:"db_user" env-required:"true"`
	DbPassword string `yaml:"db_password" env-required:"true"`
	DbName     string `yaml:"db_name" env-required:"true"`
	DbSSLMode  string `yaml:"db_sslmode" env-required:"true"`
}

type Server struct {
	ServerHost string `yaml:"server_host" env-required:"true"`
	ServerPort string `yaml:"server_port" env-required:"true"`
}

type Config struct {
	Env    string `yaml:"env"` // local, dev, prod
	DB     DB     `yaml:"auth_db"`
	Server Server `yaml:"auth_server"`
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./authApp/config.yaml"
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
