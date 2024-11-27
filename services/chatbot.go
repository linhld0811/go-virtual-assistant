package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultGroqURL   = "https://api.groq.com/openai/v1/chat/completions"
	defaultAuthToken = "Bearer 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"
)

type ChatbotConfig struct {
	URL           string
	Authorization string
}

// NewDefaultConfig returns a default configuration
func NewDefaultConfig() *ChatbotConfig {
	return &ChatbotConfig{
		URL:           defaultGroqURL,
		Authorization: defaultAuthToken,
	}
}

type RequestBody struct {
	Model       string  `json:"model"`
	MaxToken    int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	Messages    []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func ChatbotService(query string, config *ChatbotConfig) (string, error) {
	if config == nil {
		config = NewDefaultConfig()
	}
	contentType := "application/json"
	authorization := "Bearer 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"

	requestBody := RequestBody{
		Model:       "llama3-8b-8192",
		MaxToken:    1000,
		Temperature: 0.2,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: query},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from chatbot")
	}

	return response.Choices[0].Message.Content, nil
}
