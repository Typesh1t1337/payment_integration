package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvDev = "dev"
	EnvProd = "prod"
)

type Config struct {
	Env string `env:"ENV" env-default:"dev"`
	DBURL string `env:"DB_URL" env-required:"true"`
	Port string `env:"PORT" env-required:"true"`
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS" env-default:"*"`
	Jwt JwtConfig
}

type JwtConfig struct {
	PrivateKey string `env:"JWT_PRIVATE_KEY_PATH" env-required:"true"`
	PublicKey string `env:"JWT_PUBLIC_KEY_PATH" env-required:"true"`
	AccessTTL time.Duration `env:"JWT_ACCESS_TTL" env-default:"1h"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" env-default:"24h"`
}

func MustLoad() *Config {
    cfg := &Config{}
    
    if err := cleanenv.ReadEnv(cfg); err != nil {
        log.Fatal("Config error: ", err)
    }
	if cfg.Env != EnvDev && cfg.Env != EnvProd {
		log.Fatal("Invalid environment")
	}
    
    return cfg
}