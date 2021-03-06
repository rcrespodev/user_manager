package config

import (
	"github.com/vrischmann/envconfig"
	"log"
)

var Conf *Config

type Config struct {
	App struct {
		Host string `envconfig:"default=0.0.0.0"`
		Port string `envconfig:"default=8080"`
	}
	Redis struct {
		Host string `envconfig:"default=redis"`
		Port string `envconfig:"default=6379"`
	}
	MySql struct {
		User     string `envconfig:"default=root"`
		Password string `envconfig:"default=my_secret"`
		Database string `envconfig:"default=user_manager"`
		Host     string `envconfig:"default=mysql"`
		Port     string `envconfig:"default=3306"`
	}
}

func Setup() *Config {
	if err := envconfig.Init(&Conf); err != nil {
		log.Fatal(err)
	}
	return Conf
}
