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

	// Maxsus klaviatura yaratish
	var (
		menu = &telebot.ReplyMarkup{ResizeKeyboard: true}
		btnUpload = menu.Text("Загрузить документ")
	)
	menu.Reply(menu.Row(btnUpload))

	// /start buyrug'ini qayta ishlash
	b.Handle("/start", func(c telebot.Context) error {
		welcomeMessage := "Здравствуйте! Я — бот для анализа юридических документов. Отправьте мне файл в формате .pdf, .docx, .doc или .txt, и я подготовлю для вас краткий анализ.\n\nНажмите кнопку ниже, чтобы загрузить документ."
		return c.Send(welcomeMessage, menu)
	})

	// "Загрузить документ" tugmasini bosganda ishlaydigan handler
	b.Handle(&btnUpload, func(c telebot.Context) error {
		return c.Send("Пожалуйста, отправьте ваш документ.")
	})

	// Fayllarni qayta ishlash
	b.Handle(telebot.OnDocument, func(c telebot.Context) error {
		return handleDocument(c, cfg)
	})

	log.Println("Bot ishga tushirildi...")
	b.Start()
}