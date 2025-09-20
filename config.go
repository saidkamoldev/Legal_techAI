package main

import (
    "fmt"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    TelegramBotToken string
    GeminiAPIKey     string // O'zgartirildi
}
// 
func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        return nil, fmt.Errorf(".env fayli topilmadi: %w", err)
    }

    cfg := &Config{
        TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
        GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"), // O'zgartirildi
    }

    if cfg.TelegramBotToken == "" {
        return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN muhit o'zgaruvchisi o'rnatilmagan")
    }
    if cfg.GeminiAPIKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY muhit o'zgaruvchisi o'rnatilmagan")
    }

    return cfg, nil
}