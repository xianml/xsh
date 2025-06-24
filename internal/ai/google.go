package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/xian/xsh/internal/config"
)

type GoogleProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

type GoogleRequest struct {
	Contents []GoogleContent `json:"contents"`
}

type GoogleContent struct {
	Parts []GooglePart `json:"parts"`
}

type GooglePart struct {
	Text string `json:"text"`
}

type GoogleResponse struct {
	Candidates []GoogleCandidate `json:"candidates"`
	Error      *GoogleError      `json:"error,omitempty"`
}

type GoogleCandidate struct {
	Content GoogleContent `json:"content"`
}

type GoogleError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GoogleModelsResponse struct {
	Models []GoogleModelInfo `json:"models"`
}

type GoogleModelInfo struct {
	Name                string   `json:"name"`
	BaseModelId         string   `json:"baseModelId"`
	Version             string   `json:"version"`
	DisplayName         string   `json:"displayName"`
	Description         string   `json:"description"`
	SupportedGeneration []string `json:"supportedGenerationMethods"`
}

func NewGoogleProvider(cfg config.ModelConfig) (*GoogleProvider, error) {
	return &GoogleProvider{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		client:  &http.Client{},
	}, nil
}

func (p *GoogleProvider) Query(ctx context.Context, prompt string) (string, error) {
	requestBody := GoogleRequest{
		Contents: []GoogleContent{
			{
				Parts: []GooglePart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", p.baseURL, p.model, p.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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

	var response GoogleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", nil
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

func (p *GoogleProvider) GetAvailableModels() ([]string, error) {
	url := fmt.Sprintf("%s/v1beta/models?key=%s", p.baseURL, p.apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response GoogleModelsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var models []string
	for _, model := range response.Models {
		// 只包含支持文本生成的 Gemini 模型
		for _, method := range model.SupportedGeneration {
			if method == "generateContent" && strings.Contains(model.Name, "gemini") {
				// 提取模型名称（去掉路径前缀）
				parts := strings.Split(model.Name, "/")
				if len(parts) > 0 {
					modelName := parts[len(parts)-1]
					models = append(models, modelName)
				}
				break
			}
		}
	}

	return models, nil
}
