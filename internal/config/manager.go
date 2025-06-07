package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

const (
	ConfigDirName  = ".gg"
	ConfigFileName = "config.json"
)

// Config represents the application configuration
type Config struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

// Manager handles configuration operations
type Manager struct {
	config Config
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{}
}

// GetOpenAIAPIKey retrieves the OpenAI API key from various sources
func (m *Manager) GetOpenAIAPIKey() (string, error) {
	// Try multiple sources for the API key in order of priority

	// 1. Check environment variable first
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		// Validate the key
		valid, err := m.validateAPIKey(apiKey)
		if valid {
			return apiKey, nil
		}
		// If environment variable contains invalid key, report it but continue searching
		if err != nil {
			fmt.Printf("Warning: Environment variable OPENAI_API_KEY is invalid: %v\n", err)
		}
	}

	// 2. Check .env file in current directory
	err := godotenv.Load()
	if err == nil {
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey != "" {
			// Validate the key
			valid, err := m.validateAPIKey(apiKey)
			if valid {
				return apiKey, nil
			}
			if err != nil {
				fmt.Printf("Warning: .env file OPENAI_API_KEY is invalid: %v\n", err)
			}
		}
	}

	// 3. Check user's home directory config
	homeConfig, err := m.getUserHomeConfig()
	if err == nil && homeConfig.OpenAIAPIKey != "" {
		// Validate the key
		valid, err := m.validateAPIKey(homeConfig.OpenAIAPIKey)
		if valid {
			return homeConfig.OpenAIAPIKey, nil
		}
		if err != nil {
			fmt.Printf("Warning: Home config OPENAI_API_KEY is invalid: %v\n", err)
		}
	}

	// 4. Check current directory config
	currentConfig, err := m.loadConfig(".")
	if err == nil && currentConfig.OpenAIAPIKey != "" {
		// Validate the key
		valid, err := m.validateAPIKey(currentConfig.OpenAIAPIKey)
		if valid {
			return currentConfig.OpenAIAPIKey, nil
		}
		if err != nil {
			fmt.Printf("Warning: Current directory config OPENAI_API_KEY is invalid: %v\n", err)
		}
	}

	// If no valid API key found, prompt user to set it up
	fmt.Println("No valid OpenAI API key found.")
	return m.setupConfigInteractively()
}

// validateAPIKey checks if the API key is valid by making a test request
func (m *Manager) validateAPIKey(apiKey string) (bool, error) {
	if apiKey == "" {
		return false, fmt.Errorf("API key is empty")
	}

	// Basic format validation
	if !strings.HasPrefix(apiKey, "sk-") {
		return false, fmt.Errorf("API key should start with 'sk-'")
	}

	if len(apiKey) < 20 {
		return false, fmt.Errorf("API key appears to be too short")
	}

	// For now, we'll just do basic validation
	// In a production environment, you might want to make a test API call
	return true, nil
}

// getUserHomeConfig loads configuration from user's home directory
func (m *Manager) getUserHomeConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	return m.loadConfig(homeDir)
}

// loadConfig loads configuration from the specified directory
func (m *Manager) loadConfig(dirPath string) (Config, error) {
	configDir := filepath.Join(dirPath, ConfigDirName)
	configFile := filepath.Join(configDir, ConfigFileName)

	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// saveConfig saves configuration to the specified directory
func (m *Manager) saveConfig(config Config, dirPath string) error {
	configDir := filepath.Join(dirPath, ConfigDirName)
	configFile := filepath.Join(configDir, ConfigFileName)

	// Create config directory if it doesn't exist
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

// setupConfigInteractively prompts the user to set up their API key
func (m *Manager) setupConfigInteractively() (string, error) {
	fmt.Println("\n=== GitGud Configuration Setup ===")
	fmt.Println("You need to configure your OpenAI API key to use AI features.")
	fmt.Println("You can get your API key from: https://platform.openai.com/api-keys")
	fmt.Println()

	prompt := promptui.Prompt{
		Label: "Enter your OpenAI API Key",
		Mask:  '*',
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("API key cannot be empty")
			}
			if !strings.HasPrefix(input, "sk-") {
				return fmt.Errorf("API key should start with 'sk-'")
			}
			if len(input) < 20 {
				return fmt.Errorf("API key appears to be too short")
			}
			return nil
		},
	}

	apiKey, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get API key input: %v", err)
	}

	// Ask where to save the configuration
	choices := []string{
		"Save to home directory (recommended)",
		"Save to current project directory",
		"Don't save (use only for this session)",
	}

	selectPrompt := promptui.Select{
		Label: "Where would you like to save this configuration?",
		Items: choices,
	}

	choice, _, err := selectPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get save location choice: %v", err)
	}

	config := Config{OpenAIAPIKey: apiKey}

	switch choice {
	case 0: // Home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %v", err)
		}
		if err := m.saveConfig(config, homeDir); err != nil {
			return "", fmt.Errorf("failed to save config to home directory: %v", err)
		}
		fmt.Println("✅ Configuration saved to home directory")
	case 1: // Current directory
		if err := m.saveConfig(config, "."); err != nil {
			return "", fmt.Errorf("failed to save config to current directory: %v", err)
		}
		fmt.Println("✅ Configuration saved to current project directory")
	case 2: // Don't save
		fmt.Println("⚠️  Configuration will only be used for this session")
	}

	return apiKey, nil
}

// HandleConfigCommand handles the config command
func (m *Manager) HandleConfigCommand() error {
	choices := []string{
		"Show current configuration",
		"Update OpenAI API Key",
		"Reset configuration",
	}

	prompt := promptui.Select{
		Label: "Configuration Options",
		Items: choices,
	}

	choice, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("failed to get configuration choice: %v", err)
	}

	switch choice {
	case 0:
		return m.showCurrentConfig()
	case 1:
		_, err := m.setupConfigInteractively()
		return err
	case 2:
		return m.resetConfig()
	}

	return nil
}

// showCurrentConfig displays the current configuration
func (m *Manager) showCurrentConfig() error {
	fmt.Println("\n=== Current Configuration ===")

	// Check environment variable
	envKey := os.Getenv("OPENAI_API_KEY")
	if envKey != "" {
		fmt.Printf("Environment Variable: %s\n", m.maskAPIKey(envKey))
	} else {
		fmt.Println("Environment Variable: Not set")
	}

	// Check .env file
	_ = godotenv.Load()
	envFileKey := os.Getenv("OPENAI_API_KEY")
	if envFileKey != "" && envFileKey != envKey {
		fmt.Printf(".env file: %s\n", m.maskAPIKey(envFileKey))
	} else {
		fmt.Println(".env file: Not found or same as environment")
	}

	// Check home directory config
	homeConfig, err := m.getUserHomeConfig()
	if err == nil && homeConfig.OpenAIAPIKey != "" {
		fmt.Printf("Home directory config: %s\n", m.maskAPIKey(homeConfig.OpenAIAPIKey))
	} else {
		fmt.Println("Home directory config: Not found")
	}

	// Check current directory config
	currentConfig, err := m.loadConfig(".")
	if err == nil && currentConfig.OpenAIAPIKey != "" {
		fmt.Printf("Current directory config: %s\n", m.maskAPIKey(currentConfig.OpenAIAPIKey))
	} else {
		fmt.Println("Current directory config: Not found")
	}

	// Show which one would be used
	fmt.Println("\n--- Priority Order ---")
	fmt.Println("1. Environment Variable")
	fmt.Println("2. .env file")
	fmt.Println("3. Home directory config")
	fmt.Println("4. Current directory config")

	return nil
}

// maskAPIKey masks most of the API key for display
func (m *Manager) maskAPIKey(key string) string {
	if len(key) <= 10 {
		return strings.Repeat("*", len(key))
	}
	return key[:7] + strings.Repeat("*", len(key)-10) + key[len(key)-3:]
}

// resetConfig removes configuration files
func (m *Manager) resetConfig() error {
	choices := []string{
		"Reset home directory config",
		"Reset current directory config",
		"Reset both",
		"Cancel",
	}

	prompt := promptui.Select{
		Label: "What would you like to reset?",
		Items: choices,
	}

	choice, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("failed to get reset choice: %v", err)
	}

	switch choice {
	case 0: // Home directory
		homeDir, _ := os.UserHomeDir()
		configDir := filepath.Join(homeDir, ConfigDirName)
		if err := os.RemoveAll(configDir); err != nil {
			return fmt.Errorf("failed to reset home config: %v", err)
		}
		fmt.Println("✅ Home directory configuration reset")
	case 1: // Current directory
		configDir := filepath.Join(".", ConfigDirName)
		if err := os.RemoveAll(configDir); err != nil {
			return fmt.Errorf("failed to reset current config: %v", err)
		}
		fmt.Println("✅ Current directory configuration reset")
	case 2: // Both
		homeDir, _ := os.UserHomeDir()
		homeConfigDir := filepath.Join(homeDir, ConfigDirName)
		currentConfigDir := filepath.Join(".", ConfigDirName)

		if err := os.RemoveAll(homeConfigDir); err != nil {
			return fmt.Errorf("failed to reset home config: %v", err)
		}
		if err := os.RemoveAll(currentConfigDir); err != nil {
			return fmt.Errorf("failed to reset current config: %v", err)
		}
		fmt.Println("✅ Both configurations reset")
	case 3: // Cancel
		fmt.Println("Configuration reset cancelled")
	}

	return nil
}
