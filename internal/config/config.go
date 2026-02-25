package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL string
	PORT   string
}

func LoadConfig() *Config {

	// Read  from .env file
	_ = godotenv.Load()

	config := &Config{
		DB_URL: os.Getenv("DB_URL"),
		PORT:   os.Getenv("PORT"),
	}

	if config.DB_URL == "" {
		log.Fatal("DB_URL is not set!")
	}

	return config
}
