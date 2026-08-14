package config

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Database DB     `yaml:"db"`
	Redis    Redis  `yaml:"redis"`
	Server   Server `yaml:"server"`
	Kafka    Kafka  `yaml:"kafka"`
	Jwt      Jwt    `yaml:"jwt"`
}

type DB struct {
	Name     string `yaml:"name" required:"true"`
	Host     string `yaml:"host" required:"true"`
	Port     int    `yaml:"port" required:"true"`
	Username string `yaml:"username" required:"true"`
	Password string `env:"DB_PASS" required:"true"`
	SslMode  string `yaml:"ssl_mode" required:"true"`
}

type Redis struct {
	Addr string `yaml:"addr" required:"true"`
}

type Server struct {
	Host string `yaml:"host" required:"true" env-default:":8080"`
	Port int    `yaml:"port" required:"true"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers" required:"true"`
	Topics  Topics   `yaml:"topics" required:"true"`
	GroupId string   `yaml:"group_id" required:"true"`
}

type Topics struct {
	OrderCreated     string `yaml:"order_created" required:"true"`
	PaymentSucceeded string `yaml:"payment_succeeded" required:"true"`
}

type Jwt struct {
	Secret     string        `env:"JWT_SECRET" required:"true"`
	AccessTtl  time.Duration `yaml:"access_ttl" required:"true"`
	RefreshTtl time.Duration `yaml:"refresh_ttl" required:"true"`
}

// MustLoad - загружает .env file,
// читает конфиг и подставляет значения
func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Info("Error loading .env file")
	}

	cfg := Config{}

	err = cleanenv.ReadConfig(os.Getenv("CONFIG_PATH"), &cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return &cfg
}
