package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbHost     string
	Env        string
}

func getenv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	return &Config{
		Port:       getenv("PORT", "9090"),
		DbPort:     getenv("DB_PORT", "5432"),
		DbUser:     getenv("DB_USER", "postgres"),
		DbPassword: getenv("DB_PASSWORD", "password"),
		DbName:     getenv("DB_NAME", "shortener"),
		DbHost:     getenv("DB_HOST", "localhost"),
		Env:        getenv("ENV", "dev"),
	}
}

func DSNbuilder(cfg *Config) string {
	return "postgres://" + cfg.DbUser + ":" + cfg.DbPassword + "@" + cfg.DbHost + ":" + cfg.DbPort + "/" + cfg.DbName + "?sslmode=disable"
}
