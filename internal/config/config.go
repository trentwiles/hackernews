package config

import (
    "log"
    "os"
    "fmt"

    "github.com/joho/godotenv"
)

func LoadEnv() {
    err := godotenv.Load("../../.env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

func GetEnv(key string) string {
    val := os.Getenv(key)
    if val == "" {
        fmt.Printf("WARNING: Environment variable %s is empty or not set\n", key)
    }
    
    return val
}
