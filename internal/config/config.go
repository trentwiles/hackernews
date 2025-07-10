package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v\n", err)
    }
}

func GetEnv(key string) string {
    val := os.Getenv(key)
    if val == "" {
        log.Printf("[WARN] Environment variable %s is empty or not set\n", key)
    }
    
    return val
}
