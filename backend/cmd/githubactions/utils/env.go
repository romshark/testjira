package utils

import (
	"log"
	"os"
)

// GetEnv returns the value of given environment variable or panics if not set
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal("env var '" + key + "' must be set")
	}
	return value
}
