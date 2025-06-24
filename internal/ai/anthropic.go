package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xian/xsh/internal/config"
)

type AnthropicProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicResponse struct {
	Content []AnthropicContent `json:"content"`
	Error   *AnthropicError    `json:"error,omitempty"`
}

type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AnthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewAnthropicProvider(cfg config.ModelConfig) (*AnthropicProvider, error) {
	return &AnthropicProvider{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		client:  &http.Client{},
	}, nil
}

func (p *AnthropicProvider) Query(ctx context.Context, prompt string) (string, error) {
	requestBody := AnthropicRequest{
		Model:     p.model,
		MaxTokens: 1000,
		Messages: []AnthropicMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response AnthropicResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Content) == 0 {
		return "", nil
	}

	return response.Content[0].Text, nil
}

func (p *AnthropicProvider) GetAvailableModels() ([]string, error) {
	// Anthropic 不提供公开的模型列表 API，所以使用预定义的 Claude 模型
	models := []string{
		"claude-3-5-sonnet-20241022",
		"claude-3-5-sonnet-20240620",
		"claude-3-sonnet-20240229",
		"claude-3-opus-20240229",
		"claude-3-haiku-20240307",
		"claude-2.1",
		"claude-2.0",
		"claude-instant-1.2",
	}

	return models, nil
}
