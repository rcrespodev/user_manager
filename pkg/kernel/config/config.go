package config

import (
	"github.com/vrischmann/envconfig"
	"log"
)

var Conf *Config

type Config struct {
	App struct {
		Host string `default:"0.0.0.0"`
		Port string `default:"8080"`
	}
	Redis struct {
		Host string `default:"redis"`
		Port string `default:"6379"`
	}
	MySql struct {
		User     string `default:"root"`
		Password string `default:"my_secret"`
		Database string `default:"user_manager"`
		Host     string `default:"mysql"`
		Port     string `default:"3306"`
	}
}

func Setup() *Config {
	if err := envconfig.Init(&Conf); err != nil {
		log.Fatal(err)
	}
	return Conf
}
