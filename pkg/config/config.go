package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	APIKey            string  `mapstructure:"api_key"`
	Model             string  `mapstructure:"model"`
	MaxTokens         int     `mapstructure:"max_tokens"`
	Temperature       float64 `mapstructure:"temperature"`
	TopP              float64 `mapstructure:"top_p"`
	SearchContextSize string  `mapstructure:"search_context_size"`
	SearchMode        string  `mapstructure:"search_mode"`
	ReasoningEffort   string  `mapstructure:"reasoning_effort"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Model:             "sonar",
		MaxTokens:         0, // 0 means use model default
		Temperature:       0.2,
		TopP:              0.9,
		SearchContextSize: "low",
		SearchMode:        "web",
		ReasoningEffort:   "medium",
	}
}

// Load loads configuration from environment variables and config file
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Set up viper
	viper.SetEnvPrefix("PPLX")
	viper.AutomaticEnv()

	// Set up config file path
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".pplx")
	configFile := filepath.Join(configDir, "config.yaml")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set config file
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("model", cfg.Model)
	viper.SetDefault("max_tokens", cfg.MaxTokens)
	viper.SetDefault("temperature", cfg.Temperature)
	viper.SetDefault("top_p", cfg.TopP)
	viper.SetDefault("search_context_size", cfg.SearchContextSize)
	viper.SetDefault("search_mode", cfg.SearchMode)
	viper.SetDefault("reasoning_effort", cfg.ReasoningEffort)

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// API key is required - check environment variable first
	apiKey := os.Getenv("PPLX_API_KEY")
	if apiKey == "" {
		apiKey = viper.GetString("api_key")
	}
	cfg.APIKey = apiKey

	return cfg, nil
}

// LoadWithFile loads configuration from a specific config file
func LoadWithFile(configFile string) (*Config, error) {
	cfg := DefaultConfig()

	viper.SetEnvPrefix("PPLX")
	viper.AutomaticEnv()

	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// API key from environment takes precedence
	apiKey := os.Getenv("PPLX_API_KEY")
	if apiKey != "" {
		cfg.APIKey = apiKey
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required. Set PPLX_API_KEY environment variable or add api_key to ~/.pplx/config.yaml")
	}

	// Validate model
	validModels := []string{"sonar", "sonar-pro", "sonar-deep-research", "sonar-reasoning", "sonar-reasoning-pro"}
	if !contains(validModels, c.Model) {
		// Allow any model name (Perplexity may add new ones), but warn
		// We'll accept it anyway
	}

	// Validate temperature
	if c.Temperature < 0 || c.Temperature >= 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	// Validate top_p
	if c.TopP < 0 || c.TopP > 1 {
		return fmt.Errorf("top_p must be between 0 and 1")
	}

	return nil
}

// Save saves the configuration to the config file
func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".pplx")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Don't save API key to config file for security
	// Just save other settings
	viper.Set("model", c.Model)
	viper.Set("max_tokens", c.MaxTokens)
	viper.Set("temperature", c.Temperature)
	viper.Set("top_p", c.TopP)
	viper.Set("search_context_size", c.SearchContextSize)
	viper.Set("search_mode", c.SearchMode)
	viper.Set("reasoning_effort", c.ReasoningEffort)

	configFile := filepath.Join(configDir, "config.yaml")
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pplx")
}

// GetConfigFilePath returns the full path to the config file
func GetConfigFilePath() string {
	return filepath.Join(GetConfigDir(), "config.yaml")
}

// ConfigExists checks if a config file exists
func ConfigExists() bool {
	_, err := os.Stat(GetConfigFilePath())
	return err == nil
}

// CreateDefaultConfig creates a default config file if it doesn't exist
func CreateDefaultConfig() error {
	if ConfigExists() {
		return nil // Don't overwrite existing config
	}

	cfg := DefaultConfig()
	return cfg.Save()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
