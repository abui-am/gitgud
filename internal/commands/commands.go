package commands

import (
	"fmt"
	"os"

	"github.com/user/gitgud/internal/git"
)

// HandleGitCommand handles execution of Git commands with arguments
func HandleGitCommand(command string, args []string) {
	// Special handling for commands that need validation
	switch command {
	case "add":
		handleAddCommand(args)
	case "commit":
		handleCommitCommand(args)
	default:
		// Pass through to git
		err := git.ExecuteGitCommand(command, args...)
		if err != nil {
			fmt.Printf("Error executing git %s: %v\n", command, err)
			os.Exit(1)
		}
	}
}

func handleAddCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing file path")
		fmt.Println("Usage: gg add <file>")
		os.Exit(1)
	}
	git.ExecuteGitCommand("add", args...)
}

func handleCommitCommand(args []string) {
	// Check if -m flag is present
	messageProvided := false
	for i, arg := range args {
		if arg == "-m" && i+1 < len(args) {
			messageProvided = true
			break
		}
	}

	if !messageProvided {
		fmt.Println("Error: Commit message is required")
		fmt.Println("Usage: gg commit -m \"your message\"")
		fmt.Println("Or use: gg autocommit  # for AI-generated messages")
		os.Exit(1)
	}

	git.ExecuteGitCommand("commit", args...)
}
