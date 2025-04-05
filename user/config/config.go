package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	KeyPath string
}

func New(envPath string) *Config {
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("%v", err)
		panic("failed to load service config")
	}

	return &Config{
		KeyPath: os.Getenv("KEYS_PATH"),
	}
}
