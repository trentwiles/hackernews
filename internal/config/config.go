package config

import (
    "log"
    "os"
    "path/filepath"

    "github.com/joho/godotenv"
)

// first trys WORKING_DIR/.env, then trys ../../.env
func LoadEnv() {
    paths := []string{
        ".env",
        filepath.Join("..", "..", ".env"),
    }

    var err error
    for _, path := range paths {
        err = godotenv.Load(path)
        if err == nil {
            log.Printf("Loaded .env file from %s\n", path)
            return
        }
    }

    log.Printf("Failed to load .env file from known paths: %v\n", err)
}

func GetEnv(key string) string {
    val := os.Getenv(key)
    if val == "" {
        log.Fatalf("[FATAL] Environment variable %s is empty or not set\n", key)
        //log.Printf("[WARN] Environment variable %s is empty or not set\n", key)
    }
    
    return val
}
