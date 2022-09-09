package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen struct {
		BindIP string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port   string `env:"PORT" env-default:"10000"`
	}
	AppConfig struct {
		LogLevel  string `env:"LOG_LEVEL" env-default:"trace"`
		GinMode   string `env:"GIN_MODE" env-default:"debug"`
		AdminUser struct {
			Username string `env:"ADMIN_USERNAME" env-default:"admin"`
			Password string `env:"ADMIN_PWD" env-default:"admin"`
		}
	}
	PostgreSQL struct {
		Username string `env:"POSTGRES_USER" env-required:"true"`
		Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
		Host     string `env:"POSTGRES_HOST" env-required:"true"`
		Port     string `env:"POSTGRES_PORT" env-required:"true"`
		Database string `env:"POSTGRES_DB" env-required:"true"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("gather config")

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "The Art of Development - Monolith Notes System"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
