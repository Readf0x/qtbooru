//go:build debug

package main

import (
	"log"
	"os"
	"qtbooru/config"

	"github.com/joho/godotenv"
)

func loadEnv() *config.ApiConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return &config.ApiConfig{
		Username: os.Getenv("API_USER"),
		Key:      os.Getenv("API_KEY"),
	}
}
