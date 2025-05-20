package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gg <command> [<args>]")
		fmt.Println("Available commands:")
		fmt.Println("  init")
		fmt.Println("  add <file>")
		fmt.Println("  commit -m <message>")
		fmt.Println("  status")
		fmt.Println("  log")
		fmt.Println("  diff")
		fmt.Println("  autocommit (or ac)")
		fmt.Println("  config")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		executeGitCommand("init")
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing file path")
			os.Exit(1)
		}
		executeGitCommand("add", os.Args[2:]...)
	case "commit":
		commitFlag := flag.NewFlagSet("commit", flag.ExitOnError)
		message := commitFlag.String("m", "", "Commit message")
		commitFlag.Parse(os.Args[2:])

		if *message == "" {
			fmt.Println("Error: Commit message is required")
			os.Exit(1)
		}

		executeGitCommand("commit", "-m", *message)
	case "status":
		executeGitCommand("status")
	case "log":
		executeGitCommand("log")
	case "diff":
		executeGitCommand("diff")
	case "autocommit", "ac":
		handleAutoCommit()
	case "config":
		handleConfig()
	case "branch":
		executeGitCommand("branch", os.Args[2:]...)
	case "checkout":
		executeGitCommand("checkout", os.Args[2:]...)
	case "push":
		executeGitCommand("push", os.Args[2:]...)
	case "pull":
		executeGitCommand("pull", os.Args[2:]...)
	case "fetch":
		executeGitCommand("fetch", os.Args[2:]...)
	case "merge":
		executeGitCommand("merge", os.Args[2:]...)
	case "rebase":
		executeGitCommand("rebase", os.Args[2:]...)
	case "stash":
		executeGitCommand("stash", os.Args[2:]...)
	case "remote":
		executeGitCommand("remote", os.Args[2:]...)
	case "tag":
		executeGitCommand("tag", os.Args[2:]...)
	case "help":
		if len(os.Args) > 2 {
			executeGitCommand("help", os.Args[2])
		} else {
			fmt.Println("GitGud - A wrapper around Git")
			fmt.Println("Usage: gg <command> [<args>]")
			fmt.Println("\nAvailable commands:")
			fmt.Println("  init                    Initialize a new repository")
			fmt.Println("  add <file>              Add file contents to the index")
			fmt.Println("  commit -m <message>     Record changes to the repository")
			fmt.Println("  status                  Show the working tree status")
			fmt.Println("  log                     Show commit logs")
			fmt.Println("  diff                    Show changes between commits, commit and working tree, etc")
			fmt.Println("  autocommit (or ac)      Automatically add all changes and generate commit message using AI")
			fmt.Println("  config                  View or update your configuration settings")
			fmt.Println("  branch                  List, create, or delete branches")
			fmt.Println("  checkout                Switch branches or restore working tree files")
			fmt.Println("  push                    Update remote refs along with associated objects")
			fmt.Println("  pull                    Fetch from and integrate with another repository or a local branch")
			fmt.Println("  fetch                   Download objects and refs from another repository")
			fmt.Println("  merge                   Join two or more development histories together")
			fmt.Println("  rebase                  Reapply commits on top of another base tip")
			fmt.Println("  stash                   Stash the changes in a dirty working directory away")
			fmt.Println("  remote                  Manage set of tracked repositories")
			fmt.Println("  tag                     Create, list, delete or verify a tag object signed with GPG")
			fmt.Println("  help                    Display help information")
		}
	default:
		// Try to execute as a direct git command
		err := executeGitCommand(command, os.Args[2:]...)
		if err != nil {
			fmt.Printf("Error: Unknown command '%s'\n", command)
			fmt.Println("Run 'gg help' for usage.")
			os.Exit(1)
		}
	}
}

func executeGitCommand(command string, args ...string) error {
	// Just directly use the args passed in
	cmd := exec.Command("git", append([]string{command}, args...)...)

	// Set output and error to be displayed directly
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()

	// Print a custom message for certain commands
	switch command {
	case "init":
		if err == nil {
			fmt.Println("GitGud repository initialized successfully!")
		}
	case "commit":
		if err == nil {
			fmt.Println("Changes committed successfully!")
		}
	}

	return err
}

func getOpenAIAPIKey() (string, error) {
	// Try multiple sources for the API key in order of priority

	// 1. Check environment variable first
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		// Validate the key
		valid, err := validateAPIKey(apiKey)
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
			valid, err := validateAPIKey(apiKey)
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
		valid, err := validateAPIKey(homeConfig.OpenAIAPIKey)
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
			valid, err := validateAPIKey(apiKey)
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
			valid, err := validateAPIKey(exeConfig.OpenAIAPIKey)
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

// validateAPIKey checks if the provided API key is valid by making a small request to OpenAI
func validateAPIKey(apiKey string) (bool, error) {
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

func handleConfig() {
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
		valid, _ := validateAPIKey(envKey)
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
			valid, _ := validateAPIKey(dotEnvKey)
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
			valid, _ := validateAPIKey(homeConfig.OpenAIAPIKey)
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
			valid, _ := validateAPIKey(exeDotEnvKey)
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
			valid, _ := validateAPIKey(exeConfig.OpenAIAPIKey)
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

func handleAutoCommit() {
	// Get OpenAI API key using our new function
	apiKey, err := getOpenAIAPIKey()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("You can reset your configuration by running 'gg config reset'")
		os.Exit(1)
	}

	if apiKey == "" {
		fmt.Println("Error: OpenAI API key is required for autocommit")
		fmt.Println("Please run 'gg config reset' to set up your API key")
		os.Exit(1)
	}

	// Try to validate the key again just to be sure
	valid, err := validateAPIKey(apiKey)
	if !valid {
		fmt.Printf("Error: The API key is invalid: %v\n", err)
		fmt.Println("Please run 'gg config reset' to update your API key")
		os.Exit(1)
	}

	// Inform users they can modify the rules
	fmt.Println("Note: You can customize the commit message format by creating or editing the .autocommit.md file.")
	fmt.Println("      This file is not tracked by Git (it's in .gitignore).")

	// Check if there are changes to commit
	if !hasChangesToCommit() {
		fmt.Println("No changes to commit. Working tree clean.")
		os.Exit(0)
	}

	// Get the diff of changes
	diff, err := getGitDiff()
	if err != nil {
		fmt.Printf("Error getting diff: %v\n", err)
		os.Exit(1)
	}

	if diff == "" {
		fmt.Println("No changes detected in tracked files.")
		fmt.Println("You may need to run 'gg add .' first to stage new files.")
		os.Exit(0)
	}

	// Generate commit message using OpenAI
	fmt.Println("Generating commit message with AI...")
	commitMsg, err := generateCommitMessage(apiKey, diff)
	if err != nil {
		fmt.Printf("Error generating commit message: %v\n", err)
		fmt.Println("This could be due to an invalid or expired API key.")
		fmt.Println("Please run 'gg config reset' to update your API key")
		os.Exit(1)
	}

	// Display the commit message and ask for confirmation
	fmt.Printf("\nGenerated commit message:\n\n%s\n\n", commitMsg)
	fmt.Print("Do you want to commit with this message? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	response = strings.TrimSpace(response)
	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		// Add all changes
		addCmd := exec.Command("git", "add", ".")
		addCmd.Stdout = os.Stdout
		addCmd.Stderr = os.Stderr
		if err := addCmd.Run(); err != nil {
			fmt.Printf("Error adding changes: %v\n", err)
			os.Exit(1)
		}

		// Commit changes
		if err := executeGitCommand("commit", "-m", commitMsg); err != nil {
			fmt.Printf("Error committing changes: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Commit canceled.")
	}
}

func hasChangesToCommit() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error checking git status: %v\n", err)
		os.Exit(1)
	}

	return len(output) > 0
}

func getGitDiff() (string, error) {
	// Get staged changes
	stagedCmd := exec.Command("git", "diff", "--staged")
	stagedOutput, err := stagedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting staged diff: %v", err)
	}

	// Get unstaged changes
	unstagedCmd := exec.Command("git", "diff")
	unstagedOutput, err := unstagedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting unstaged diff: %v", err)
	}

	// Combine both outputs
	combinedDiff := string(stagedOutput) + string(unstagedOutput)

	// Get untracked files
	untrackedCmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	untrackedOutput, err := untrackedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting untracked files: %v", err)
	}

	// If there are untracked files, add them to the diff summary
	if len(untrackedOutput) > 0 {
		untrackedFiles := strings.Split(strings.TrimSpace(string(untrackedOutput)), "\n")
		combinedDiff += "\n\nUntracked files:\n"
		for _, file := range untrackedFiles {
			combinedDiff += "  " + file + "\n"
		}
	}

	return combinedDiff, nil
}

func getAutocommitRules() (string, error) {
	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}

	// First, check for user's .autocommit.md in project root
	userRulesPath := filepath.Join(currentDir, ".autocommit.md")
	content, err := os.ReadFile(userRulesPath)
	if err == nil {
		return string(content), nil
	}

	// If not found in project root, check executable directory for default rules
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	defaultRulesPath := filepath.Join(exeDir, ".autocommit.md")
	content, err = os.ReadFile(defaultRulesPath)
	if err != nil {
		// Default rules if no .autocommit.md is found anywhere
		return "Please follow the Conventional Commits format: <type>(<scope>): <description>", nil
	}

	return string(content), nil
}

func generateCommitMessage(apiKey, diff string) (string, error) {
	// Initialize OpenAI client
	client := openai.NewClient(apiKey)

	// Truncate diff if it's too large (OpenAI has token limits)
	maxDiffLength := 4000
	diffContent := diff
	if len(diff) > maxDiffLength {
		diffContent = diff[:maxDiffLength] + "\n...(diff truncated due to size)"
	}

	// Get autocommit rules
	rules, err := getAutocommitRules()

	if err != nil {
		fmt.Printf("Warning: Could not load autocommit rules: %v\n", err)
		rules = "Please follow the Conventional Commits format: <type>(<scope>): <description>"
	}

	// Create prompt for OpenAI
	prompt := fmt.Sprintf(
		"Generate a commit message for the following git diff:\n\n%s\n\n"+
			"Must follow these rules for the commit message:\n%s\n\n"+
			"Reply with ONLY the commit message, nothing else.",
		diffContent,
		rules,
	)

	// Create chat completion request
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 250,
		},
	)

	if err != nil {
		return "", fmt.Errorf("chat completion error: %v", err)
	}

	// Extract the commit message from the response
	commitMessage := resp.Choices[0].Message.Content
	return strings.TrimSpace(commitMessage), nil
}
