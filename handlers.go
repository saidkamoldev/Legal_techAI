package main

import (
	"fmt"
	"os"
	"time"
	"strings"

	"gopkg.in/telebot.v3"
)

// handleDocument - foydalanuvchi yuborgan faylni qayta ishlaydi.
func handleDocument(c telebot.Context, cfg *Config) error {
	doc := c.Message().Document
    
    // Fayl nomini kichik harflarga o'tkazib tekshiramiz
    fileNameLower := strings.ToLower(doc.FileName)
    if !strings.HasSuffix(fileNameLower, ".pdf") && !strings.HasSuffix(fileNameLower, ".docx") && !strings.HasSuffix(fileNameLower, ".doc") && !strings.HasSuffix(fileNameLower, ".txt") {
        return c.Send("Kechirasiz, faqat .pdf, .docx, .doc va .txt formatidagi fayllarni qabul qilaman.")
    }

	// Fayl hajmini tekshirish (10 MB).
	if doc.FileSize > 10*1024*1024 {
		return c.Send("Fayl hajmi 10 MB dan oshmasligi kerak.")
	}

	// Faylni vaqtinchalik papkaga yuklab olish.
	tempDir := "temp_files"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return c.Send("Vaqtinchalik papka yaratishda xato: " + err.Error())
	}
	filePath := fmt.Sprintf("%s/%d_%s", tempDir, time.Now().UnixNano(), doc.FileName)
	err := c.Bot().Download(&doc.File, filePath)
	if err != nil {
		return c.Send("Faylni yuklab olishda xato: " + err.Error())
	}
	defer os.Remove(filePath) // Fayl tahlil qilingandan so'ng o'chiriladi.

	if err := c.Send("Fayl qabul qilindi, matnni ajratib olyapman..."); err != nil {
		return err
	}
    
	// Fayldan matnni ajratish.
	documentText, err := parseDocument(filePath)
	if err != nil {
		return c.Send(fmt.Sprintf("Hujjatni tahlil qilishda xato: %v", err))
	}

	if len(documentText) < 50 { // Matn hajmi juda qisqa bo'lsa
		return c.Send("Hujjat matni juda qisqa. Iltimos, to'liqroq hujjat yuboring.")
	}

	if err := c.Send("Matnni AIga tahlil uchun yuboryapman... Bu jarayon biroz vaqt olishi mumkin."); err != nil {
		return err
	}
	
	// AIga tahlil uchun yuborish.
	analysis, err := analyzeDocumentAI(documentText, cfg.OpenRouterAPIKey)
	if err != nil {
		return c.Send(fmt.Sprintf("AI tahlilida xato: %v", err))
	}

	// Natijani foydalanuvchiga yuborish.
	return c.Send("✅ **Tahlil natijasi:**\n" + analysis, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}