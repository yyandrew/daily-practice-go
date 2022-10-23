package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database   Database
	Auth       Auth
	ServerPort int `envconfig:"SERVER_PORT" default:"9000"`
}

type Auth struct {
	JWT_TOKEN string `envconfig:"JWT_TOKEN" rquired:"true"`
}

type Database struct {
	Host string `envconfig:"DATABASE_HOST" required:"true"`
	Port int    `envconfig:"DATABASE_PORT" required:"true"`
	Name string `envconfig:"DATABASE_NAME" required:"true"`
}

func NewParsedConfig() (Config, error) {
	_ = godotenv.Load(".env")
	cnf := Config{}
	err := envconfig.Process("", &cnf)

	return cnf, err
}
