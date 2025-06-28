package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AnthropicProvider struct {
	apiKey string
	model  string
	client *http.Client
}

type anthropicRequest struct {
	Model       string                    `json:"model"`
	Messages    []anthropicMessage        `json:"messages"`
	MaxTokens   int                       `json:"max_tokens"`
	Temperature float64                   `json:"temperature"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	ID      string                   `json:"id"`
	Type    string                   `json:"type"`
	Role    string                   `json:"role"`
	Content []anthropicContent       `json:"content"`
	Error   *anthropicError          `json:"error,omitempty"`
}

type anthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *AnthropicProvider) GenerateQuestions(ctx context.Context, prompt string) ([]GeneratedQuestion, error) {
	reqBody := anthropicRequest{
		Model: p.model,
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp anthropicError
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("API error: %s - %s", errorResp.Type, errorResp.Message)
		}
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var anthropicResp anthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if anthropicResp.Error != nil {
		return nil, fmt.Errorf("API error: %s - %s", anthropicResp.Error.Type, anthropicResp.Error.Message)
	}

	// Extract the JSON from the response
	if len(anthropicResp.Content) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	var questions []GeneratedQuestion
	if err := json.Unmarshal([]byte(anthropicResp.Content[0].Text), &questions); err != nil {
		return nil, fmt.Errorf("failed to parse generated questions: %w", err)
	}

	return questions, nil
}