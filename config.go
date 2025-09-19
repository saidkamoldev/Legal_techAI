package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config - dasturning asosiy konfiguratsiya ma'lumotlari uchun tuzilma.
type Config struct {
	TelegramBotToken string
	OpenRouterAPIKey string
}

// LoadConfig - `.env` faylidan konfiguratsiya ma'lumotlarini yuklaydi.
func LoadConfig() (*Config, error) {
	// godotenv.Load() `.env` faylini yuklaydi.
	if err := godotenv.Load(); err != nil {
		// Agar `.env` fayli topilmasa, xato qaytaradi.
		return nil, fmt.Errorf(".env fayli topilmadi: %w", err)
	}

	// Ma'lumotlarni muhit o'zgaruvchilari (environment variables)dan o'qish.
	cfg := &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
	}

	// Agar kerakli kalitlar yo'q bo'lsa, xato qaytaradi.
	if cfg.TelegramBotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN muhit o'zgaruvchisi o'rnatilmagan")
	}
	if cfg.OpenRouterAPIKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY muhit o'zgaruvchisi o'rnatilmagan")
	}

	return cfg, nil
}