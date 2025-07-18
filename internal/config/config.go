package config

import (
    "log"
    "os"
    "path/filepath"

    "github.com/joho/godotenv"
)

func Init() {
    //
}

// first trys WORKING_DIR/.env, then trys ../../.env
func LoadEnv() {
    loadAny([]string{
        ".env",
        filepath.Join("..", "..", ".env"),
    })

    loadAny([]string{
        "frontend/.env",
        filepath.Join("..", "..", "frontend", ".env"),
    })
}

// try the two possible locations for the backend .env
// THEN, try the two possible locations for the frontend .env
func loadAny(paths []string) {
    for _, path := range paths {
        err := godotenv.Load(path)
        if err == nil {
            log.Printf("Loaded .env file from %s\n", path)
            return
        }
    }
    log.Printf("Failed to load .env file(s) from paths: %v\n", paths)
}

func GetEnv(key string) string {
    val := os.Getenv(key)
    if val == "" {
        log.Fatalf("[FATAL] Environment variable %s is empty or not set\n", key)
        //log.Printf("[WARN] Environment variable %s is empty or not set\n", key)
    }
    
    return val
}
