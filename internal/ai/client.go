package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/xian/xsh/internal/config"
)

type Client struct {
	config *config.Config
}

type Provider interface {
	Query(ctx context.Context, prompt string) (string, error)
	GetAvailableModels() ([]string, error)
}

func New(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
	}
}

func (c *Client) Query(prompt string) (string, error) {
	// 确保选择了正确的模型
	if !c.config.HasAnthropicKey() && !c.config.HasGoogleKey() && !c.config.HasOpenAIKey() {
		return "", fmt.Errorf("no valid API key found. Please set one of: OPENAI_API_KEY, ANTHROPIC_API_KEY, or GOOGLE_API_KEY")
	}

	// 优先级：Anthropic > Google > OpenAI
	if c.config.HasAnthropicKey() {
		c.config.SetCurrentModel("claude")
	} else if c.config.HasGoogleKey() {
		c.config.SetCurrentModel("gemini")
	} else if c.config.HasOpenAIKey() {
		c.config.SetCurrentModel("openai")
	}

	modelConfig, exists := c.config.GetCurrentModel()
	if !exists {
		return "", fmt.Errorf("no AI model configured. Please set one of: OPENAI_API_KEY, ANTHROPIC_API_KEY, or GOOGLE_API_KEY")
	}

	ctx := context.Background()

	var provider Provider
	var err error

	switch modelConfig.Provider {
	case "anthropic":
		provider, err = NewAnthropicProvider(modelConfig)
	case "google":
		provider, err = NewGoogleProvider(modelConfig)
	case "openai":
		provider, err = NewOpenAIProvider(modelConfig)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s", modelConfig.Provider)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create AI provider: %w", err)
	}

	// 构建完整的提示词
	systemPrompt := config.GetSystemPrompt()
	fullPrompt := fmt.Sprintf("%s\n\nUser: %s", systemPrompt, prompt)

	response, err := provider.Query(ctx, fullPrompt)

	// 如果是模型不可用错误，尝试使用备用模型
	if err != nil && modelConfig.Provider == "openai" &&
		(strings.Contains(err.Error(), "model_not_found") ||
			strings.Contains(err.Error(), "does not exist")) {

		// 尝试使用 gpt-3.5-turbo 作为备用
		backupConfig := modelConfig
		backupConfig.Model = "gpt-3.5-turbo"

		backupProvider, backupErr := NewOpenAIProvider(backupConfig)
		if backupErr == nil {
			response, err = backupProvider.Query(ctx, fullPrompt)
			if err == nil {
				// 成功使用备用模型，更新配置
				modelConfig.Model = "gpt-3.5-turbo"
			}
		}
	}

	return response, err
}

func (c *Client) SwitchModel(modelName string) error {
	if c.config.SetCurrentModel(modelName) {
		return nil
	}
	return fmt.Errorf("model %s is not available or not configured", modelName)
}

func (c *Client) SwitchModelByDisplayName(displayName, provider string) error {
	if c.config.SetCurrentModelByDisplayName(displayName, provider) {
		return nil
	}
	return fmt.Errorf("model %s (%s) is not available or not configured", displayName, provider)
}

func (c *Client) GetCurrentModel() string {
	modelInfo, exists := c.config.GetCurrentModelInfo()
	if !exists {
		return "unknown"
	}
	return modelInfo.DisplayName
}

func (c *Client) GetAvailableModels() []string {
	modelInfos := c.config.GetAvailableModelInfos()
	var models []string
	for _, info := range modelInfos {
		models = append(models, info.DisplayName)
	}
	return models
}

func (c *Client) GetAvailableModelInfos() []config.ModelInfo {
	var allModels []config.ModelInfo

	// 为每个配置的提供商获取实时模型列表
	modelConfig, exists := c.config.GetCurrentModel()
	if !exists {
		return allModels
	}
	var provider Provider
	var err error

	switch modelConfig.Provider {
	case "openai":
		provider, err = NewOpenAIProvider(modelConfig)
	case "anthropic":
		provider, err = NewAnthropicProvider(modelConfig)
	case "google":
		provider, err = NewGoogleProvider(modelConfig)
	default:
		return allModels
	}

	if err != nil {
		// 如果创建提供商失败，使用默认模型
		allModels = append(allModels, config.ModelInfo{
			Key:         modelConfig.Model,
			DisplayName: modelConfig.Model,
			Provider:    modelConfig.Provider,
		})
		return allModels
	}

	// 获取实时模型列表
	models, err := provider.GetAvailableModels()
	if err != nil {
		// 如果获取失败，使用默认模型
		allModels = append(allModels, config.ModelInfo{
			Key:         modelConfig.Model,
			DisplayName: modelConfig.Model,
			Provider:    modelConfig.Provider,
		})
		return allModels
	}

	// 为每个模型创建 ModelInfo
	for _, model := range models {
		allModels = append(allModels, config.ModelInfo{
			Key:         modelConfig.Provider + "-" + model, // 使用组合键
			DisplayName: model,
			Provider:    modelConfig.Provider,
		})
	}
	return allModels
}
