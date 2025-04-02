package utils

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var Store *sessions.CookieStore

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file:", err)
    }
    
    storeKey := os.Getenv("STORE_KEY")
    if storeKey == "" {
        log.Println("Warning: STORE_KEY not found in environment, using fallback key")
        storeKey = "fallback-secret-key-for-development-only"
    }
    
    Store = sessions.NewCookieStore([]byte(storeKey))
    
    Store.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   300, // 5m
        HttpOnly: true,
        Secure:   false,     // Set to true in production with HTTPS
    }
}