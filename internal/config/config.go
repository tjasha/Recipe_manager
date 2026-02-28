package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL      string
	PORT        string
	TEST_DB_URL string
}

func LoadConfig() *Config {

	// Try to load .env file from the current directory or parent directories.
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Could not load .env file, from root directory: %v", err)
		// This makes it work when running tests from subdirectories.
		if err := godotenv.Load("../../.env"); err != nil {
			log.Printf("Could not load .env file, from testing directory: %v", err)
		}
	}

	config := &Config{
		DB_URL:      os.Getenv("DB_URL"),
		PORT:        os.Getenv("PORT"),
		TEST_DB_URL: os.Getenv("TEST_DB_URL"),
	}

	if config.DB_URL == "" {
		log.Fatal("DB_URL is not set!")
	}

	return config
}
