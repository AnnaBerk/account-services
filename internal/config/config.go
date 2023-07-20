package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string `yaml:"env" env-required:"true"`
	StoragePath      string `yaml:"storage_path" env-required:"true"`
	ConnectionString string `yaml:"-"`
	Postgres         `yaml:"db"`
	HTTPServer       `yaml:"http_server"`

	Hasher struct {
		Salt string `env-required:"true" env:"HASHER_SALT"`
	}

	JWT struct {
		SignKey  string        `env-required:"true"                  env:"JWT_SIGN_KEY"`
		TokenTTL time.Duration `env-required:"true" yaml:"token_ttl" env:"JWT_TOKEN_TTL"`
	}
}

type Postgres struct {
	Host        string        `yaml:"host" env:"DB_HOST" env-required:"true"`
	Port        int           `yaml:"port" env:"DB_PORT" env-required:"true"`
	User        string        `env:"DB_USER" env-required:"true"`
	Password    string        `env:"DB_PASSWORD" env-required:"true"`
	DBName      string        `env:"DB_NAME" env-required:"true"`
	SSLMode     string        `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
	Timeout     time.Duration `yaml:"timeout" env:"DB_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"DB_IDLE_TIMEOUT" env-default:"60s"`
	MaxPoolSize int           `yaml:"max_pool_size" env-required:"true"  env:"PG_MAX_POOL_SIZE"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config: %s", err)
	}
	cfg.ConnectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	return &cfg
}
