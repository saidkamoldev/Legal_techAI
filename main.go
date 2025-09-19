package main

import (
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Konfiguratsiya yuklanmadi: %v", err)
	}

	pref := telebot.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Telegram bot yaratishda xato: %v", err)
	}

	// /start buyrug'ini qayta ishlash
	b.Handle("/start", func(c telebot.Context) error {
		welcomeMessage := "Assalomu alaykum! Men yuridik hujjatlaringizni tahlil qilib beruvchi botman. Menga `.pdf`, `.docx`, `.doc` yoki `.txt` formatidagi hujjatni yuboring, men uning asosiy punktlari, tomonlar majburiyatlari va ehtimoliy xavflarini topib beraman."
		return c.Send(welcomeMessage)
	})

	// Fayllarni qayta ishlash
	b.Handle(telebot.OnDocument, func(c telebot.Context) error {
		return handleDocument(c, cfg)
	})

	log.Println("Bot ishga tushirildi...")
	b.Start()
}