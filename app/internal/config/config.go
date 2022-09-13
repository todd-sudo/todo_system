package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen struct {
		BindIP string `env:"BIND_IP" env-default:"127.0.0.1"`
		Port   string `env:"SERVER_PORT" env-default:"8000"`
	}
	AppConfig struct {
		GinMode   string `env:"GIN_MODE" env-default:"debug"`
		AdminUser struct {
			Username string `env:"ADMIN_USERNAME" env-default:"admin"`
			Password string `env:"ADMIN_PWD" env-default:"admin"`
		}
		Auth struct {
			PasswordHashSalt string `env:"PASSWORD_HASH_SALT" env-required:"true"`
		}
		JWTToken struct {
			AccessTokenPrivateKey  string        `env:"ACCESS_TOKEN_PRIVATE_KEY" env-required:"true"`
			AccessTokenPublicKey   string        `env:"ACCESS_TOKEN_PUBLIC_KEY" env-required:"true"`
			RefreshTokenPrivateKey string        `env:"REFRESH_TOKEN_PRIVATE_KEY" env-required:"true"`
			RefreshTokenPublicKey  string        `env:"REFRESH_TOKEN_PUBLIC_KEY" env-required:"true"`
			AccessTokenExpiresIn   time.Duration `env:"ACCESS_TOKEN_EXPIRED_IN" env-required:"true"`
			RefreshTokenExpiresIn  time.Duration `env:"REFRESH_TOKEN_EXPIRED_IN" env-required:"true"`
			AccessTokenMaxAge      int           `env:"ACCESS_TOKEN_MAXAGE" env-required:"true"`
			RefreshTokenMaxAge     int           `env:"REFRESH_TOKEN_MAXAGE" env-required:"true"`
		}
	}
	PostgreSQL struct {
		Username string `env:"POSTGRES_USER" env-required:"true"`
		Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
		Host     string `env:"POSTGRES_HOST" env-required:"true"`
		Port     string `env:"POSTGRES_PORT" env-required:"true"`
		Database string `env:"POSTGRES_DB" env-required:"true"`
	}
	Redis struct {
		Host string `env:"REDIS_HOST" env-required:"true"`
		Port string `env:"REDIS_PORT" env-required:"true"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("gather config")

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "Go-rshok todo system"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
