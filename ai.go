package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	
)




// AIdan kelgan javobni qabul qilish uchun tuzilma
type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AIga yuboriladigan so'rovning tuzilmasi
type OpenRouterRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// analyzeDocumentAI hujjat matnini Gemini Flash 2.0 ga tahlil qilish uchun yuboradi
func analyzeDocumentAI(text string, apiKey string) (string, error) {
	// Konfiguratsiya ma'lumotlarini `.env` faylidan o'qish
	// apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY topilmadi")
	}

	// AIga yuboriladigan so'rovni yaratish
	requestBody, err := json.Marshal(OpenRouterRequest{
		Model: "google/gemini-flash-1.5", // Yoki "google/gemini-pro"
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: fmt.Sprintf(`Quyidagi hujjatni tahlil qilib, rus tilida quyidagi bandlar bo'yicha qisqacha hisobot ber:
				1. Основные пункты документа (asosiy punktlari)
				2. Обязанности сторон (tomomlar majburiyatlari)
				3. Возможные риски (ehtimoliy xavflar)
				
				Hujjat matni:
				%s`, text),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("so'rov tanasini tuzishda xato: %w", err)
	}

	// HTTP so'rovni yaratish
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("HTTP so'rovini yaratishda xato: %w", err)
	}

	// Sarlavhalarni (headers) qo'shish
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HTTP-Referer", "sizning-saytingiz-yoki-loyihangiz-nomi") // OpenRouter talabi

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
	var response OpenRouterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("javobni tahlil qilishda xato: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("AI javobida ma'lumot topilmadi")
	}

	return response.Choices[0].Message.Content, nil
}