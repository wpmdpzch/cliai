package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置
type Config struct {
	AI     AIConfig     `yaml:"ai"`
	Exec   ExecConfig   `yaml:"exec"`
	UI     UIConfig     `yaml:"ui"`
	Pkg    PackageConfig `yaml:"packages"`
}

// AIConfig AI 配置
type AIConfig struct {
	Provider  string `yaml:"provider"`
	APIKey    string `yaml:"api_key"`
	BaseURL   string `yaml:"base_url"`
	Model     string `yaml:"model"`
	Temp      float64 `yaml:"temperature"`
	MaxTokens int    `yaml:"max_tokens"`
}

// ExecConfig 执行配置
type ExecConfig struct {
	AutoExec        bool `yaml:"auto_exec"`
	ConfirmDangerous bool `yaml:"confirm_dangerous"`
	Timeout         int  `yaml:"timeout"`
}

// UIConfig UI 配置
type UIConfig struct {
	ModeIndicator bool   `yaml:"mode_indicator"`
	DefaultMode   string `yaml:"default_mode"`
}

// PackageConfig 包配置
type PackageConfig struct {
	Local  string `yaml:"local"`
	System string `yaml:"system"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		AI: AIConfig{
			Provider:  "openai",
			BaseURL:   "https://api.openai.com/v1",
			Model:     "gpt-4o-mini",
			Temp:      0.7,
			MaxTokens: 2048,
		},
		Exec: ExecConfig{
			AutoExec:        false,
			ConfirmDangerous: true,
			Timeout:         30,
		},
		UI: UIConfig{
			ModeIndicator: true,
			DefaultMode:   "cli",
		},
		Pkg: PackageConfig{
			Local:  "~/.cliai/packages",
			System: "/usr/local/cliai/packages",
		},
	}
}

// Load 加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save 保存配置
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
