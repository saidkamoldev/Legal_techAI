package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

// handleDocument - foydalanuvchi yuborgan faylni qayta ishlaydi.
func handleDocument(c telebot.Context, cfg *Config) error {
	doc := c.Message().Document
    
    // Fayl kengaytmasini kichik harflarga o'tkazib, tekshiramiz
    fileNameLower := strings.ToLower(doc.FileName)
    if !strings.HasSuffix(fileNameLower, ".pdf") && !strings.HasSuffix(fileNameLower, ".docx") && !strings.HasSuffix(fileNameLower, ".doc") && !strings.HasSuffix(fileNameLower, ".txt") {
        return c.Send("К сожалению, я принимаю только файлы в форматах .pdf, .docx, .doc и .txt.")
    }

	// Fayl hajmini tekshirish (10 MB).
	if doc.FileSize > 10*1024*1024 {
		return c.Send("Размер файла не должен превышать 10 МБ.")
	}

	// Faylni vaqtinchalik papkaga yuklab olish.
	tempDir := "temp_files"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return c.Send("Ошибка при создании временной папки: " + err.Error())
	}
	filePath := fmt.Sprintf("%s/%d_%s", tempDir, time.Now().UnixNano(), doc.FileName)
	err := c.Bot().Download(&doc.File, filePath)
	if err != nil {
		return c.Send("Ошибка при загрузке файла: " + err.Error())
	}
	defer os.Remove(filePath) // Fayl tahlil qilingandan so'ng o'chiriladi.

	if err := c.Send("Файл принят, извлекаю текст..."); err != nil {
		return err
	}
    
	// Fayldan matnni ajratish.
	documentText, err := parseDocument(filePath)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при анализе документа: %v", err))
	}

	if len(documentText) < 50 { // Matn hajmi juda qisqa bo'lsa
		return c.Send("Текст документа слишком короткий. Пожалуйста, отправьте более полный документ.")
	}

	if err := c.Send("Отправляю текст на анализ ИИ... Это может занять некоторое время."); err != nil {
		return err
	}
	
	// AIga tahlil uchun yuborish.
	analysis, err := analyzeDocumentAI(documentText, cfg.GeminiAPIKey)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка анализа ИИ: %v", err))
	}

	// Natijalarni chiroyli formatlash
	var builder strings.Builder
	builder.WriteString("✅ **Результат анализа:**\n\n")

	builder.WriteString("**Основные пункты:**\n")
	builder.WriteString(analysis.Summary)
	builder.WriteString("\n\n")

	builder.WriteString("**Обязанности сторон:**\n")
	for party, obligation := range analysis.Obligations {
		builder.WriteString(fmt.Sprintf("- **%s:** %s\n", party, obligation))
	}
	builder.WriteString("\n")

	builder.WriteString("**Возможные риски:**\n")
	for party, risk := range analysis.Risks {
		builder.WriteString(fmt.Sprintf("- **%s:** %s\n", party, risk))
	}

	return c.Send(builder.String(), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}