package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Username string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     uint16 `env:"DB_PORT"`
	Database string `env:"DB_NAME"`
	SslMode  string `env:"DB_SSLMODE"`
}

type JWTConfig struct {
	Secret string `env:"JWT_SECRET"`
}

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
}

func Get() (*Config, error) {
	cfg := &Config{}

	godotenv.Load(".env")
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
