package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// parseDocument - fayldan matnni formatiga qarab ajratib oladi.
func parseDocument(filePath string) (string, error) {
	parts := strings.Split(filePath, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("fayl kengaytmasi topilmadi")
	}
	extension := strings.ToLower(parts[len(parts)-1])

	switch extension {
	case "txt":
		return parseTXT(filePath)
	case "pdf":
		return parsePDF(filePath)
	case "docx", "doc":
		// DOCX va DOC formatlari uchun qo'shimcha kutubxonalar kerak.
		return "", fmt.Errorf("docx/doc formatini tahlil qilish qo'shimcha kutubxonani o'rnatishni talab qiladi")
	default:
		return "", fmt.Errorf("qo'llab-quvvatlanmagan fayl formati: %s", extension)
	}
}

// parseTXT - .txt fayldan matnni o'qiydi.
func parseTXT(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("txt faylini o'qishda xato: %w", err)
	}
	return string(data), nil
}

// parsePDF - .pdf fayldan matnni pdftotext orqali ajratib oladi.
// Bu funksiya ishlashi uchun tizimga 'poppler-utils' o'rnatilgan bo'lishi kerak.
func parsePDF(filePath string) (string, error) {
	// pdftotext buyrug'ini yaratamiz
	// "-enc UTF-8" parametrlari matnni to'g'ri kodirovka qilish uchun.
	cmd := exec.Command("pdftotext", "-enc", "UTF-8", filePath, "-")

	// Buyruqni bajaramiz va natijani (matnni) olamiz
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pdftotext buyrug'ini bajarishda xato: %w", err)
	}

	return string(output), nil
}