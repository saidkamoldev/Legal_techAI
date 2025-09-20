package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "strings"
)

// AI dan keladigan JSON javobni qabul qilish uchun tuzilma
type DocumentAnalysis struct {
	Summary     string            `json:"summary"`
	Obligations map[string]string `json:"obligations"`
	Risks       map[string]string `json:"risks"`
}

// Gemini API javobini qabul qilish uchun asosiy tuzilma
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Gemini API ga yuboriladigan so'rovning tuzilmasi
type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

// analyzeDocumentAI hujjat matnini to'g'ridan-to'g'ri Gemini AI ga tahlil qilish uchun yuboradi
func analyzeDocumentAI(text string, apiKey string) (*DocumentAnalysis, error) {
	prompt := fmt.Sprintf(`
	You are a legal assistant. Analyze the following document and provide a summary in JSON format. The JSON should contain three keys: "summary" for a brief overview, "obligations" for the parties' responsibilities as a JSON object, and "risks" for potential legal risks as a JSON object. All output must be in Russian. Do not include any text before or after the JSON, and do not wrap the JSON in any markdown code blocks.
	
	Document content:
	%s
	`, text)

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
						Text: prompt,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("so'rov tanasini tuzishda xato: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("HTTP so'rovini yaratishda xato: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP so'rovini yuborishda xato: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("javob tanasini o'qishda xato: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API xatosi: %s", string(body))
	}

	var geminiResponse GeminiResponse
	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return nil, fmt.Errorf("javobni tahlil qilishda xato: %w", err)
	}

	if len(geminiResponse.Candidates) == 0 || len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("AI javobida ma'lumot topilmadi")
	}

	jsonString := geminiResponse.Candidates[0].Content.Parts[0].Text

	var analysis DocumentAnalysis
	if err := json.Unmarshal([]byte(jsonString), &analysis); err != nil {
		return nil, fmt.Errorf("AI javobini JSON ga aylantirishda xato: %w. Javob: %s", err, jsonString)
	}

	return &analysis, nil
}