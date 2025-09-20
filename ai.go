package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Gemini API javobini qabul qilish uchun tuzilma
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Gemini API'ga yuboriladigan so'rovning tuzilmasi
type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

// analyzeDocumentAI hujjat matnini to'g'ridan-to'g'ri Gemini AI ga tahlil qilish uchun yuboradi
func analyzeDocumentAI(text string, apiKey string) (string, error) {
	// AIga yuboriladigan so'rovni yaratish
	requestBody, err := json.Marshal(GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{
						Text: fmt.Sprintf(`Quyidagi hujjatni tahlil qilib, rus tilida quyidagi bandlar bo'yicha qisqacha hisobot ber:
						1. Основные пункты документа (asosiy punktlari)
						2. Обязанности сторон (tomomlar majburiyatlari)
						3. Возможные риски (ehtimoliy xavflar)
						
						Hujjat matni:
						%s`, text),
					},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("so'rov tanasini tuzishda xato: %w", err)
	}

	// HTTP so'rovni yaratish
url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("HTTP so'rovini yaratishda xato: %w", err)
	}

	// Sarlavhalarni (headers) qo'shish
	req.Header.Set("Content-Type", "application/json")

	// HTTP so'rovini yuborish
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP so'rovini yuborishda xato: %w", err)
	}
	defer resp.Body.Close()

	// Javobni o'qish
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("javob tanasini o'qishda xato: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API xatosi: %s", string(body))
	}

	// Javobni tahlil qilish
	var response GeminiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("javobni tahlil qilishda xato: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("AI javobida ma'lumot topilmadi")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}