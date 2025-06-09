package autocommit

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/user/gitgud/internal/config"
	"github.com/user/gitgud/internal/git"
	"github.com/user/gitgud/internal/ui"
)

type AutocommitRules struct {
	Rules  string
	Source string
	Path   string
}

func HandleAutoCommit() {
	// Get OpenAI API key using our new function
	apiKey, err := config.GetOpenAIAPIKey()
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
	valid, err := config.ValidateAPIKey(apiKey)
	if !valid {
		fmt.Printf("Error: The API key is invalid: %v\n", err)
		fmt.Println("Please run 'gg config reset' to update your API key")
		os.Exit(1)
	}

	// Get current branch name
	branchName, err := git.GetCurrentBranch()
	if err != nil {
		fmt.Printf("Warning: Could not get current branch name: %v\n", err)
		branchName = "unknown"
	}

	// Get autocommit rules
	rules, err := getAutocommitRules()
	if err != nil {
		fmt.Printf("Warning: Could not load autocommit rules: %v\n", err)
		rules = AutocommitRules{
			Rules:  "Please follow the Conventional Commits format: <type>(<scope>): <description>",
			Source: "root",
			Path:   "built-in",
		}
	}

	// Check if user has a custom .autocommit.md file
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Warning: Could not get current directory: %v\n", err)
	} else {
		userRulesPath := filepath.Join(currentDir, ".autocommit.md")
		if _, err := os.Stat(userRulesPath); os.IsNotExist(err) {
			// Only show the note if no custom .autocommit.md exists
			fmt.Println("Note: You can customize the commit message format by creating or editing the .autocommit.md file.")
			fmt.Println("      This file is not tracked by Git (it's in .gitignore).")
		}
	}

	// Print configuration information
	fmt.Println("\nCommit Message Configuration:")
	fmt.Println("===========================")
	fmt.Printf("Using %s rules from: %s\n", rules.Source, rules.Path)
	fmt.Println()

	// Print branch information
	fmt.Println("Current Branch Information:")
	fmt.Println("=========================")
	fmt.Printf("Branch: %s\n", branchName)
	fmt.Println()

	// Print last commit information
	fmt.Println("Last Commit Information:")
	fmt.Println("=======================")
	lastCommitInfo, err := git.GetLastCommitMetadata()
	if err != nil {
		if strings.Contains(err.Error(), "fatal: bad default revision") {
			fmt.Println("No previous commits found.")
		} else {
			fmt.Printf("Warning: Could not get last commit metadata: %v\n", err)
		}
	} else {
		fmt.Println(lastCommitInfo)
	}
	fmt.Println()

	// Check if there are changes to commit
	if !git.HasChangesToCommit() {
		fmt.Println("No changes to commit. Working tree clean.")
		os.Exit(0)
	}

	// Get the diff of changes
	diff, err := git.GetGitDiff()
	if err != nil {
		fmt.Printf("Error getting diff: %v\n", err)
		os.Exit(1)
	}

	if diff == "" {
		fmt.Println("No changes detected in tracked files.")
		fmt.Println("You may need to run 'gg add .' first to stage new files.")
		os.Exit(0)
	}

	// Prompt for custom context
	fmt.Println("\nEnter additional context for the commit message (press Enter to finish):")
	fmt.Println("(This context will help generate a more relevant commit message)")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	customContext := strings.TrimSpace(line)

	// Generate commit message using OpenAI
	fmt.Println("\nGenerating commit message with AI...")
	commitMsg, err := generateCommitMessage(apiKey, diff, customContext)
	if err != nil {
		fmt.Printf("Error generating commit message: %v\n", err)
		fmt.Println("This could be due to an invalid or expired API key.")
		fmt.Println("Please run 'gg config reset' to update your API key")
		os.Exit(1)
	}

	// Display the commit message and ask for confirmation
	fmt.Printf("\nGenerated commit message:\n\n%s\n\n", commitMsg)
	fmt.Print("Do you want to commit with this message? (y/n): ")

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
		if err := git.ExecuteGitCommand("commit", "-m", commitMsg); err != nil {
			fmt.Printf("Error committing changes: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Commit canceled.")
	}
}

func HandleAutoCommitPerFile() {
	// Get OpenAI API key using our existing function
	apiKey, err := config.GetOpenAIAPIKey()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("You can reset your configuration by running 'gg config reset'")
		os.Exit(1)
	}

	if apiKey == "" {
		fmt.Println("Error: OpenAI API key is required for autocommit per file")
		fmt.Println("Please run 'gg config reset' to set up your API key")
		os.Exit(1)
	}

	// Try to validate the key again just to be sure
	valid, err := config.ValidateAPIKey(apiKey)
	if !valid {
		fmt.Printf("Error: The API key is invalid: %v\n", err)
		fmt.Println("Please run 'gg config reset' to update your API key")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Autocommit Per File ===")
	fmt.Println("This will help you commit files individually with AI-generated commit messages.")
	fmt.Println("Use arrow keys to navigate and select files. You can select multiple files interactively.")
	fmt.Println()

	for {
		// Get list of changed files
		changedFiles, err := git.GetChangedFiles()
		if err != nil {
			fmt.Printf("Error getting changed files: %v\n", err)
			os.Exit(1)
		}

		if len(changedFiles) == 0 {
			fmt.Println("No changes to commit. Working tree clean.")
			break
		}

		// Display changed files
		fmt.Println("Changed files:")
		for i, file := range changedFiles {
			fmt.Printf("  %d. %s\n", i+1, file)
		}
		fmt.Println()

		// Use arrow key selection for file selection
		selectedFiles, err := ui.SelectMultipleFilesWithArrows(changedFiles)
		if err != nil {
			if strings.Contains(err.Error(), "user chose to exit") {
				fmt.Println("Exiting autocommit per file.")
				break
			}
			fmt.Printf("Error in file selection: %v\n", err)
			continue
		}

		if len(selectedFiles) == 0 {
			fmt.Println("No files selected.")
			continue
		}

		// Process selected files as a batch
		fmt.Printf("\n--- Processing %d selected file(s) ---\n", len(selectedFiles))
		for _, file := range selectedFiles {
			fmt.Printf("  - %s\n", file)
		}
		fmt.Println()

		// Get combined diff for all selected files
		var combinedDiff strings.Builder
		var validFiles []string

		for _, file := range selectedFiles {
			fileDiff, err := git.GetFileDiff(file)
			if err != nil {
				fmt.Printf("Warning: Could not get diff for %s: %v\n", file, err)
				continue
			}
			if fileDiff == "" {
				fmt.Printf("Warning: No changes detected in %s, skipping.\n", file)
				continue
			}

			combinedDiff.WriteString(fmt.Sprintf("--- %s ---\n", file))
			combinedDiff.WriteString(fileDiff)
			combinedDiff.WriteString("\n")
			validFiles = append(validFiles, file)
		}

		if len(validFiles) == 0 {
			fmt.Println("No valid files to commit, skipping batch.")
			continue
		}

		// Ask for custom context for the batch
		fmt.Printf("Enter additional context for these %d file(s) (press Enter to skip): ", len(validFiles))
		contextLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		customContext := strings.TrimSpace(contextLine)

		// Generate commit message for the batch
		fmt.Printf("Generating commit message for %d file(s)...\n", len(validFiles))
		commitMsg, err := generateBatchCommitMessage(apiKey, validFiles, combinedDiff.String(), customContext)
		if err != nil {
			fmt.Printf("Error generating commit message for batch: %v\n", err)
			continue
		}

		// Display the commit message and ask for confirmation
		fmt.Printf("\nGenerated commit message for batch:\n\n%s\n\n", commitMsg)
		fmt.Print("Do you want to commit these files with this message? (y/n/exit): ")

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		response = strings.TrimSpace(response)

		if strings.ToLower(response) == "exit" {
			fmt.Println("Exiting autocommit per file.")
			return
		}

		if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			// Add all selected files
			for _, file := range validFiles {
				addCmd := exec.Command("git", "add", file)
				addCmd.Stdout = os.Stdout
				addCmd.Stderr = os.Stderr
				if err := addCmd.Run(); err != nil {
					fmt.Printf("Error adding %s: %v\n", file, err)
					continue
				}
			}

			// Commit all files with one message
			if err := git.ExecuteGitCommand("commit", "-m", commitMsg); err != nil {
				fmt.Printf("Error committing batch: %v\n", err)
				continue
			}
			fmt.Printf("Successfully committed %d file(s) in one commit\n", len(validFiles))
		} else {
			fmt.Printf("Skipped committing batch of %d file(s)\n", len(validFiles))
		}

		fmt.Println("\n--- Processing complete ---")
		fmt.Print("Continue with remaining files? (y/n): ")
		continueResponse, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			break
		}
		continueResponse = strings.TrimSpace(continueResponse)
		if strings.ToLower(continueResponse) != "y" && strings.ToLower(continueResponse) != "yes" {
			fmt.Println("Exiting autocommit per file.")
			break
		}
	}
}

func generateBatchCommitMessage(apiKey string, filenames []string, combinedDiff, customContext string) (string, error) {
	// Initialize OpenAI client
	client := openai.NewClient(apiKey)

	// Get current branch name
	branchName, err := git.GetCurrentBranch()
	if err != nil {
		fmt.Printf("Warning: Could not get current branch name: %v\n", err)
		branchName = "unknown"
	}

	// Truncate diff if it's too large (OpenAI has token limits)
	maxDiffLength := 4000
	diffContent := combinedDiff
	if len(combinedDiff) > maxDiffLength {
		diffContent = combinedDiff[:maxDiffLength] + "\n...(diff truncated due to size)"
	}

	// Get autocommit rules
	rules, err := getAutocommitRules()
	if err != nil {
		fmt.Printf("Warning: Could not load autocommit rules: %v\n", err)
		rules = AutocommitRules{
			Rules:  "Please follow the Conventional Commits format: <type>(<scope>): <description>",
			Source: "root",
			Path:   "built-in",
		}
	}

	// Create file list string
	fileListStr := strings.Join(filenames, ", ")

	// Create prompt for OpenAI focused on the batch of files
	prompt := fmt.Sprintf(
		"Generate a commit message for changes to these %d files: %s\n\n"+
			"Combined git diff for these files:\n%s\n\n"+
			"Current branch: %s\n\n"+
			"Additional context provided by the user:\n%s\n\n"+
			"Must follow these rules for the commit message:\n%s\n\n"+
			"Create a unified commit message that summarizes the changes across all these files. "+
			"Reply with ONLY the commit message, nothing else.",
		len(filenames),
		fileListStr,
		diffContent,
		branchName,
		customContext,
		rules.Rules,
	)

	// Create chat completion request
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Dot1Nano,
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

func generateFileCommitMessage(apiKey, filename, diff, customContext string) (string, error) {
	// Initialize OpenAI client
	client := openai.NewClient(apiKey)

	// Get current branch name
	branchName, err := git.GetCurrentBranch()
	if err != nil {
		fmt.Printf("Warning: Could not get current branch name: %v\n", err)
		branchName = "unknown"
	}

	// Truncate diff if it's too large (OpenAI has token limits)
	maxDiffLength := 3000
	diffContent := diff
	if len(diff) > maxDiffLength {
		diffContent = diff[:maxDiffLength] + "\n...(diff truncated due to size)"
	}

	// Get autocommit rules
	rules, err := getAutocommitRules()
	if err != nil {
		fmt.Printf("Warning: Could not load autocommit rules: %v\n", err)
		rules = AutocommitRules{
			Rules:  "Please follow the Conventional Commits format: <type>(<scope>): <description>",
			Source: "root",
			Path:   "built-in",
		}
	}

	// Create prompt for OpenAI focused on the specific file
	prompt := fmt.Sprintf(
		"Generate a commit message for changes to this specific file: %s\n\n"+
			"Git diff for this file:\n%s\n\n"+
			"Current branch: %s\n\n"+
			"Additional context provided by the user:\n%s\n\n"+
			"Must follow these rules for the commit message:\n%s\n\n"+
			"Focus the commit message on what changed in this specific file. "+
			"Reply with ONLY the commit message, nothing else.",
		filename,
		diffContent,
		branchName,
		customContext,
		rules.Rules,
	)

	// Create chat completion request
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Dot1Nano,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 200,
		},
	)

	if err != nil {
		return "", fmt.Errorf("chat completion error: %v", err)
	}

	// Extract the commit message from the response
	commitMessage := resp.Choices[0].Message.Content
	return strings.TrimSpace(commitMessage), nil
}

func getAutocommitRules() (AutocommitRules, error) {
	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return AutocommitRules{}, fmt.Errorf("error getting current directory: %v", err)
	}

	// First, check for user's .autocommit.md in project root
	userRulesPath := filepath.Join(currentDir, ".autocommit.md")
	content, err := os.ReadFile(userRulesPath)
	if err == nil {
		return AutocommitRules{
			Rules:  string(content),
			Source: "project",
			Path:   userRulesPath,
		}, nil
	}

	// If not found in project root, check executable directory for default rules
	exePath, err := os.Executable()
	if err != nil {
		return AutocommitRules{}, fmt.Errorf("error getting executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	defaultRulesPath := filepath.Join(exeDir, ".autocommit.md")
	content, err = os.ReadFile(defaultRulesPath)
	if err == nil {
		return AutocommitRules{
			Rules:  string(content),
			Source: "default",
			Path:   defaultRulesPath,
		}, nil
	}

	// Default rules if no .autocommit.md is found anywhere
	return AutocommitRules{
		Rules:  "Please follow the Conventional Commits format: <type>(<scope>): <description>",
		Source: "root",
		Path:   "built-in",
	}, nil
}

func generateCommitMessage(apiKey, diff string, customContext string) (string, error) {
	// Initialize OpenAI client
	client := openai.NewClient(apiKey)

	// Get current branch name
	branchName, err := git.GetCurrentBranch()
	if err != nil {
		fmt.Printf("Warning: Could not get current branch name: %v\n", err)
		branchName = "unknown"
	}

	// Get last commit metadata
	lastCommitInfo, err := git.GetLastCommitMetadata()
	if err != nil {
		fmt.Printf("Warning: Could not get last commit metadata: %v\n", err)
		lastCommitInfo = ""
	}

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
		rules = AutocommitRules{
			Rules:  "Please follow the Conventional Commits format: <type>(<scope>): <description>",
			Source: "root",
			Path:   "built-in",
		}
	}

	// Create prompt for OpenAI
	prompt := fmt.Sprintf(
		"Generate a commit message for the following git diff:\n\n%s\n\n"+
			"Current branch: %s\n"+
			"%s\n\n"+
			"Additional context provided by the user:\n%s\n\n"+
			"Must follow these rules for the commit message:\n%s\n\n"+
			"Reply with ONLY the commit message, nothing else.",
		diffContent,
		branchName,
		lastCommitInfo,
		customContext,
		rules.Rules,
	)

	// Create chat completion request
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Dot1Nano,
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
