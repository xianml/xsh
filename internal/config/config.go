package config

import (
	"os"
)

type Config struct {
	CurrentModel    string
	Models          map[string]ModelConfig
	AnthropicAPIKey string
	GoogleAPIKey    string
	OpenAIAPIKey    string
}

type ModelConfig struct {
	Provider string
	APIKey   string
	BaseURL  string
	Model    string
}

// Load 从环境变量加载配置
func Load() *Config {
	config := &Config{
		CurrentModel: getEnv("XSH_MODEL", "openai"),
		Models:       make(map[string]ModelConfig),
	}

	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		// Anthropic Claude 配置
		config.Models["claude"] = ModelConfig{
			Provider: "anthropic",
			APIKey:   apiKey,
			BaseURL:  getEnv("ANTHROPIC_BASE_URL", "https://api.anthropic.com"),
			Model:    getEnv("ANTHROPIC_MODEL", "claude-3-sonnet-20240229"),
		}
		config.AnthropicAPIKey = apiKey
		config.CurrentModel = "claude"
	} else if apiKey := os.Getenv("GOOGLE_API_KEY"); apiKey != "" {
		// Google Gemini 配置
		config.Models["gemini"] = ModelConfig{
			Provider: "google",
			APIKey:   apiKey,
			BaseURL:  getEnv("GOOGLE_BASE_URL", "https://generativelanguage.googleapis.com"),
			Model:    getEnv("GOOGLE_MODEL", "gemini-pro"),
		}
		config.GoogleAPIKey = apiKey
		config.CurrentModel = "gemini"
	} else if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		// OpenAI 配置
		config.Models["openai"] = ModelConfig{
			Provider: "openai",
			APIKey:   apiKey,
			BaseURL:  getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
			Model:    getEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
		}
		config.OpenAIAPIKey = apiKey
		config.CurrentModel = "openai"
	}

	// 如果当前选择的模型不可用，选择第一个可用的模型
	if _, exists := config.Models[config.CurrentModel]; !exists {
		for name := range config.Models {
			config.CurrentModel = name
			break
		}
	}

	return config
}

// GetCurrentModel 获取当前选择的模型配置
func (c *Config) GetCurrentModel() (ModelConfig, bool) {
	model, exists := c.Models[c.CurrentModel]
	return model, exists
}

// SetCurrentModel 设置当前模型
func (c *Config) SetCurrentModel(name string) bool {
	if _, exists := c.Models[name]; exists {
		c.CurrentModel = name
		return true
	}
	return false
}

// GetAvailableModels 获取所有可用模型列表
func (c *Config) GetAvailableModels() []string {
	var models []string
	for name := range c.Models {
		models = append(models, name)
	}
	return models
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetSystemPrompt 获取系统提示词
func GetSystemPrompt() string {
	return `You are an AI assistant helping with shell commands and system administration tasks.

When a user asks for help with a command or describes what they want to do, provide:
1. The exact shell command(s) they need
2. A brief explanation of what the command does
3. Any important warnings or considerations

Format your response with shell commands clearly marked, for example:
$ ls -la
$ cd /path/to/directory

Be concise but helpful. Focus on practical, safe commands. If the user's request is unclear, ask for clarification.
Current shell: zsh
OS: ` + getOS()
}

func getOS() string {
	if os := os.Getenv("OS"); os != "" {
		return os
	}
	if runtime := os.Getenv("GOOS"); runtime != "" {
		return runtime
	}
	return "unix-like"
}

// HasAnthropicKey checks if the Anthropic API key is set
func (c *Config) HasAnthropicKey() bool {
	return c.AnthropicAPIKey != ""
}

// HasGoogleKey checks if the Google API key is set
func (c *Config) HasGoogleKey() bool {
	return c.GoogleAPIKey != ""
}

// HasOpenAIKey checks if the OpenAI API key is set
func (c *Config) HasOpenAIKey() bool {
	return c.OpenAIAPIKey != ""
}

// ModelInfo 包含模型的显示信息
type ModelInfo struct {
	Key         string // 内部键名，如 "openai", "claude"
	DisplayName string // 显示名称，如 "gpt-3.5-turbo", "claude-3-sonnet-20240229"
	Provider    string // 提供商名称，如 "openai", "anthropic"
}

// GetAvailableModelInfos 获取所有可用模型的详细信息
func (c *Config) GetAvailableModelInfos() []ModelInfo {
	var models []ModelInfo
	for key, modelConfig := range c.Models {
		models = append(models, ModelInfo{
			Key:         key,
			DisplayName: modelConfig.Model,
			Provider:    modelConfig.Provider,
		})
	}
	return models
}

// GetCurrentModelInfo 获取当前模型的详细信息
func (c *Config) GetCurrentModelInfo() (ModelInfo, bool) {
	modelConfig, exists := c.Models[c.CurrentModel]
	if !exists {
		return ModelInfo{}, false
	}

	return ModelInfo{
		Key:         c.CurrentModel,
		DisplayName: modelConfig.Model,
		Provider:    modelConfig.Provider,
	}, true
}

// SetCurrentModelByDisplayName 通过显示名称和提供商设置当前模型
func (c *Config) SetCurrentModelByDisplayName(displayName, provider string) bool {
	for key, modelConfig := range c.Models {
		if modelConfig.Provider == provider {
			// 更新模型配置中的具体模型名称
			modelConfig.Model = displayName
			c.Models[key] = modelConfig
			c.CurrentModel = key
			return true
		}
	}
	return false
}
