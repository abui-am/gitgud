package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitgud/internal/autocommit"
	"github.com/user/gitgud/internal/config"
	"github.com/user/gitgud/internal/git"
	"github.com/user/gitgud/internal/help"
)

// HandleCommand handles the routing and execution of commands
func HandleCommand(command string, args []string) {
	switch command {
	case "init":
		git.ExecuteGitCommand("init")
	case "add":
		handleAddCommand(args)
	case "commit":
		handleCommitCommand(args)
	case "status":
		git.ExecuteGitCommand("status")
	case "log":
		git.ExecuteGitCommand("log")
	case "diff":
		git.ExecuteGitCommand("diff")
	case "autocommit", "ac":
		autocommit.HandleAutoCommit()
	case "autocommit-per-file", "acpf":
		autocommit.HandleAutoCommitPerFile()
	case "config":
		config.HandleConfig()
	case "last":
		git.HandleLastCommit()
	case "branch":
		git.ExecuteGitCommand("branch", args...)
	case "checkout":
		git.ExecuteGitCommand("checkout", args...)
	case "push":
		git.ExecuteGitCommand("push", args...)
	case "pull":
		git.ExecuteGitCommand("pull", args...)
	case "fetch":
		git.ExecuteGitCommand("fetch", args...)
	case "merge":
		git.ExecuteGitCommand("merge", args...)
	case "rebase":
		git.ExecuteGitCommand("rebase", args...)
	case "stash":
		git.ExecuteGitCommand("stash", args...)
	case "remote":
		git.ExecuteGitCommand("remote", args...)
	case "tag":
		git.ExecuteGitCommand("tag", args...)
	case "help":
		handleHelpCommand(args)
	default:
		handleUnknownCommand(command, args)
	}
}

func handleAddCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing file path")
		os.Exit(1)
	}
	git.ExecuteGitCommand("add", args...)
}

func handleCommitCommand(args []string) {
	commitFlag := flag.NewFlagSet("commit", flag.ExitOnError)
	message := commitFlag.String("m", "", "Commit message")
	commitFlag.Parse(args)

	if *message == "" {
		fmt.Println("Error: Commit message is required")
		os.Exit(1)
	}

	git.ExecuteGitCommand("commit", "-m", *message)
}

func handleHelpCommand(args []string) {
	if len(args) > 0 {
		git.ExecuteGitCommand("help", args[0])
	} else {
		help.ShowUsage()
	}
}

func handleUnknownCommand(command string, args []string) {
	// Try to execute as a direct git command
	err := git.ExecuteGitCommand(command, args...)
	if err != nil {
		fmt.Printf("Error: Unknown command '%s'\n", command)
		fmt.Println("Run 'gg help' for usage.")
		os.Exit(1)
	}
}
