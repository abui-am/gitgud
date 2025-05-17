package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

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

func handleAutoCommit() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file. Please make sure you have a .env file with OPENAI_API_KEY")
		os.Exit(1)
	}

	// Get OpenAI API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY is not set in .env file")
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
	// Check if .autocommit.md exists in the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}

	rulesPath := filepath.Join(currentDir, ".autocommit.md")
	content, err := os.ReadFile(rulesPath)
	if err != nil {
		// If not found in current directory, check executable directory
		exePath, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("error getting executable path: %v", err)
		}

		exeDir := filepath.Dir(exePath)
		rulesPath = filepath.Join(exeDir, ".autocommit.md")
		content, err = os.ReadFile(rulesPath)
		if err != nil {
			// Default rules if .autocommit.md is not found
			return "Please follow the Conventional Commits format: <type>(<scope>): <description>", nil
		}
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
