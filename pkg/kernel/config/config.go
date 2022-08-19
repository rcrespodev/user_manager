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
	Rabbit struct {
		User     string `envconfig:"default=my_user"`
		Password string `envconfig:"default=my_password"`
		Host     string `envconfig:"default=rabbitmq"`
		Port     string `envconfig:"default=5672"`
	}
	Jwt struct {
		Key struct {
			Private string `envconfig:"default=/app/cert/id_rsa"`
			Public  string `envconfig:"default=/app/cert/id_rsa.pub"`
		}
		ExpirationTime string `envconfig:"default=15"`
	}
	Log struct {
		File struct {
			Path string `envconfig:"default=/app/logs"`
		}
	}
	Smtp struct {
		Host     string `envconfig:"default=smtp.gmail.com"`
		Port     string `envconfig:"default=587"`
		Username string `envconfig:"default=user.manager.email.info@gmail.com"`
		Password string `envconfig:"default=kkgaoclbngmnzuhb"`
		Welcome  struct {
			Template string `envconfig:"default=/app/pkg/app/emailSender/domain/welcomeTemplate.txt"`
		}
	}
}

func Setup() *Config {
	if err := envconfig.Init(&Conf); err != nil {
		log.Fatal(err)
	}
	return Conf
}
