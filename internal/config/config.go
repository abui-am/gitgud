package config

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// Constant for config directory and file names
const (
	ConfigDirName  = ".gg"
	ConfigFileName = "config.json"
)

// Config structure to store the application configuration
type Config struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

func GetOpenAIAPIKey() (string, error) {
	// Try multiple sources for the API key in order of priority

	// 1. Check environment variable first
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		// Validate the key
		valid, err := ValidateAPIKey(apiKey)
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
			valid, err := ValidateAPIKey(apiKey)
			if valid {
				return apiKey, nil
			}
			if err != nil {
				fmt.Printf("Warning: API key in .env file is invalid: %v\n", err)
			}
		}
	}

	// 3. Check user's home directory for config
	homeConfig, err := getUserHomeConfig()
	if err == nil && homeConfig.OpenAIAPIKey != "" {
		// Validate the key
		valid, err := ValidateAPIKey(homeConfig.OpenAIAPIKey)
		if valid {
			return homeConfig.OpenAIAPIKey, nil
		}
		if err != nil {
			fmt.Printf("Warning: API key in home config is invalid: %v\n", err)
		}
	}

	// 4. Check executable directory for .env or config
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)

		// Try .env in exe dir
		_ = godotenv.Load(filepath.Join(exeDir, ".env"))
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey != "" {
			// Validate the key
			valid, err := ValidateAPIKey(apiKey)
			if valid {
				return apiKey, nil
			}
			if err != nil {
				fmt.Printf("Warning: API key in executable directory .env file is invalid: %v\n", err)
			}
		}

		// Try config.json in exe dir
		exeConfig, err := loadConfig(exeDir)
		if err == nil && exeConfig.OpenAIAPIKey != "" {
			// Validate the key
			valid, err := ValidateAPIKey(exeConfig.OpenAIAPIKey)
			if valid {
				return exeConfig.OpenAIAPIKey, nil
			}
			if err != nil {
				fmt.Printf("Warning: API key in executable directory config is invalid: %v\n", err)
			}
		}
	}

	fmt.Println("No valid OpenAI API key found.")
	fmt.Println("You can:\n1. Run 'gg config' to set or update your API key\n2. Provide a key for this session")

	// If we reach here, prompt user to set up config
	return setupConfigInteractively()
}

// ValidateAPIKey checks if the provided API key is valid by making a small request to OpenAI
func ValidateAPIKey(apiKey string) (bool, error) {
	if apiKey == "" {
		return false, fmt.Errorf("API key is empty")
	}

	// Create a client with a short timeout
	client := openai.NewClient(apiKey)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make a minimal request to validate the key
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Test",
				},
			},
			MaxTokens: 5,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "401") ||
			strings.Contains(err.Error(), "invalid_api_key") ||
			strings.Contains(err.Error(), "Incorrect API key") {
			return false, fmt.Errorf("invalid API key")
		}
		// Could be a network error, but the key might still be valid
		return false, fmt.Errorf("could not validate: %v", err)
	}

	// If we get a response, the key is valid
	if len(resp.Choices) > 0 {
		return true, nil
	}

	return false, fmt.Errorf("unexpected response from API")
}

func getUserHomeConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	return loadConfig(filepath.Join(homeDir, ConfigDirName))
}

func loadConfig(dirPath string) (Config, error) {
	configPath := filepath.Join(dirPath, ConfigFileName)

	data, err := os.ReadFile(configPath)
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

func saveConfig(config Config, dirPath string) error {
	// Ensure directory exists
	err := os.MkdirAll(dirPath, 0700) // Restrict to user only
	if err != nil {
		return err
	}

	configPath := filepath.Join(dirPath, ConfigFileName)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600) // Restrict to user only
}

func setupConfigInteractively() (string, error) {
	fmt.Println("OpenAI API key not found. Please enter your OpenAI API key:")
	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", fmt.Errorf("API key cannot be empty")
	}

	fmt.Println("\nWhere would you like to save your API key?")
	fmt.Println("1. User home directory (recommended)")
	fmt.Println("2. Current directory")
	fmt.Println("3. Don't save (use only for this session)")

	choice, err := reader.ReadString('\n')
	if err != nil {
		return apiKey, nil // Return the key but don't save
	}

	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error accessing home directory: %v\n", err)
			return apiKey, nil
		}

		configDir := filepath.Join(homeDir, ConfigDirName)
		config := Config{OpenAIAPIKey: apiKey}

		err = saveConfig(config, configDir)
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
		} else {
			fmt.Printf("API key saved to %s\n", filepath.Join(configDir, ConfigFileName))
		}

	case "2":
		// Save to .env in current directory
		err = os.WriteFile(".env", []byte(fmt.Sprintf("OPENAI_API_KEY=%s", apiKey)), 0600)
		if err != nil {
			fmt.Printf("Error saving .env file: %v\n", err)
		} else {
			fmt.Println("API key saved to .env in current directory")
		}

	default:
		fmt.Println("API key will be used for this session only")
	}

	return apiKey, nil
}

func HandleConfig() {
	// If no arguments, show current configuration
	if len(os.Args) == 2 {
		showCurrentConfig()
		return
	}

	if len(os.Args) >= 3 {
		switch os.Args[2] {
		case "reset":
			resetConfig()
		default:
			fmt.Println("Unknown config command. Available commands:")
			fmt.Println("  gg config           - Show current configuration")
			fmt.Println("  gg config reset     - Reset and update your OpenAI API key")
		}
	}
}

func showCurrentConfig() {
	fmt.Println("Current Configuration:")

	// Check all possible locations for API keys

	// Environment variable
	envKey := os.Getenv("OPENAI_API_KEY")
	if envKey != "" {
		// Don't show the full key for security
		maskedKey := maskAPIKey(envKey)
		valid, _ := ValidateAPIKey(envKey)
		status := "valid"
		if !valid {
			status = "invalid"
		}
		fmt.Printf("- Environment variable OPENAI_API_KEY: %s (%s)\n", maskedKey, status)
	} else {
		fmt.Println("- Environment variable OPENAI_API_KEY: not set")
	}

	// Local .env
	err := godotenv.Load()
	if err == nil {
		dotEnvKey := os.Getenv("OPENAI_API_KEY")
		if dotEnvKey != "" && dotEnvKey != envKey {
			maskedKey := maskAPIKey(dotEnvKey)
			valid, _ := ValidateAPIKey(dotEnvKey)
			status := "valid"
			if !valid {
				status = "invalid"
			}
			fmt.Printf("- Local .env file: %s (%s)\n", maskedKey, status)
		} else if dotEnvKey == envKey {
			fmt.Println("- Local .env file: same as environment variable")
		} else {
			fmt.Println("- Local .env file: exists but no API key")
		}
	} else {
		fmt.Println("- Local .env file: not found")
	}

	// Home directory config
	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeConfig, err := loadConfig(filepath.Join(homeDir, ConfigDirName))
		if err == nil && homeConfig.OpenAIAPIKey != "" {
			maskedKey := maskAPIKey(homeConfig.OpenAIAPIKey)
			valid, _ := ValidateAPIKey(homeConfig.OpenAIAPIKey)
			status := "valid"
			if !valid {
				status = "invalid"
			}
			fmt.Printf("- Home directory config (~/%s/%s): %s (%s)\n", ConfigDirName, ConfigFileName, maskedKey, status)
		} else {
			fmt.Printf("- Home directory config (~/%s/%s): not found or no API key\n", ConfigDirName, ConfigFileName)
		}
	}

	// Executable directory
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)

		// Check .env in exe dir
		err := godotenv.Load(filepath.Join(exeDir, ".env"))
		exeDotEnvKey := os.Getenv("OPENAI_API_KEY")
		if err == nil && exeDotEnvKey != "" && exeDotEnvKey != envKey {
			maskedKey := maskAPIKey(exeDotEnvKey)
			valid, _ := ValidateAPIKey(exeDotEnvKey)
			status := "valid"
			if !valid {
				status = "invalid"
			}
			fmt.Printf("- Executable directory .env: %s (%s)\n", maskedKey, status)
		} else {
			fmt.Println("- Executable directory .env: not found or no API key")
		}

		// Check config in exe dir
		exeConfig, err := loadConfig(exeDir)
		if err == nil && exeConfig.OpenAIAPIKey != "" {
			maskedKey := maskAPIKey(exeConfig.OpenAIAPIKey)
			valid, _ := ValidateAPIKey(exeConfig.OpenAIAPIKey)
			status := "valid"
			if !valid {
				status = "invalid"
			}
			fmt.Printf("- Executable directory config (%s): %s (%s)\n", ConfigFileName, maskedKey, status)
		} else {
			fmt.Printf("- Executable directory config (%s): not found or no API key\n", ConfigFileName)
		}
	}

	fmt.Println("\nYou can reset your configuration by running 'gg config reset'")
}

func maskAPIKey(key string) string {
	if len(key) < 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func resetConfig() {
	fmt.Println("Resetting your OpenAI API configuration...")
	_, err := setupConfigInteractively()
	if err != nil {
		fmt.Printf("Error setting up configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Configuration updated successfully!")
}
