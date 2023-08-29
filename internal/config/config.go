package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env string `env:"ENV" env-required:"true"`
	Database
	HttpServer
}

type Database struct {
	DbUrl string `env:"DB_URL" env-required:"true"`
}

type HttpServer struct {
	Address string `env:"HTTP_SERVER_URL" env-default:"localhost:8080"`
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
