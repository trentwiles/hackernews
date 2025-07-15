package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        //log.Fatalf("Error loading .env file: %v\n", err)
        log.Printf("Error loading .env file: %v\n", err)
    }
    log.Printf("Loaded .env file!\n")
}

func GetEnv(key string) string {
    val := os.Getenv(key)
    if val == "" {
        log.Fatalf("[FATAL] Environment variable %s is empty or not set\n", key)
        //log.Printf("[WARN] Environment variable %s is empty or not set\n", key)
    }
    
    return val
}
