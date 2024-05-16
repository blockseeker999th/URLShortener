package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type Config struct {
	Env      string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	*HTTPServer
}

func MustLoad() *Config {
	err := godotenv.Load("./../../.env")
	if err != nil {
		fmt.Println("Error loading .env variables", err)
	}

	return &Config{
		Env:      getEnv("ENV", "local"),
		Host:     "db",
		Port:     getEnv("PG_PORT", "5432"),
		User:     getEnv("PG_USER", "default"),
		Password: getEnv("PG_PASSWORD", "default"),
		DBName:   getEnv("PG_DBNAME", "urlshortener"),
		HTTPServer: &HTTPServer{
			Address:     getEnv("SERVER_ADDRESS", ":3002"),
			Timeout:     parseTime("SERVER_TIMEOUT", "4s"),
			IdleTimeout: parseTime("SERVER_IDLE_TIMEOUT", "4s"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func parseTime(key, fallback string) time.Duration {
	value := getEnv(key, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		fmt.Printf("Error parsing %s duration: %s\n", key, err)
		duration, _ = time.ParseDuration(fallback)
	}
	return duration
}
