package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gitgud <command> [<args>]")
		fmt.Println("Available commands:")
		fmt.Println("  init")
		fmt.Println("  add <file>")
		fmt.Println("  commit -m <message>")
		fmt.Println("  status")
		fmt.Println("  log")
		fmt.Println("  diff")
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
			fmt.Println("Usage: gitgud <command> [<args>]")
			fmt.Println("\nAvailable commands:")
			fmt.Println("  init                    Initialize a new repository")
			fmt.Println("  add <file>              Add file contents to the index")
			fmt.Println("  commit -m <message>     Record changes to the repository")
			fmt.Println("  status                  Show the working tree status")
			fmt.Println("  log                     Show commit logs")
			fmt.Println("  diff                    Show changes between commits, commit and working tree, etc")
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
			fmt.Println("Run 'gitgud help' for usage.")
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
