package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/gitgud/internal/autocommit"
	"github.com/user/gitgud/internal/commands"
	"github.com/user/gitgud/internal/config"
	"github.com/user/gitgud/internal/git"
)

var rootCmd = &cobra.Command{
	Use:   "gg",
	Short: "GitGud - A smart Git wrapper with AI-powered commit messages",
	Long: `GitGud is a Git wrapper that enhances your Git workflow with AI-powered features.
It supports all standard Git commands while adding intelligent autocommit functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

var autocommitCmd = &cobra.Command{
	Use:     "autocommit",
	Aliases: []string{"ac"},
	Short:   "Generate AI-powered commit messages for all changes",
	Long: `Autocommit analyzes your changes and generates intelligent commit messages
using OpenAI. It follows Conventional Commits format and considers your branch
name and previous commit context.`,
	Run: func(cmd *cobra.Command, args []string) {
		autocommit.HandleAutoCommit()
	},
}

var acpfCmd = &cobra.Command{
	Use:     "autocommit-per-file",
	Aliases: []string{"acpf"},
	Short:   "Commit files individually or in batches with AI-generated messages",
	Long: `Autocommit per file allows you to select specific files and commit them
individually or in batches. Each selection gets its own AI-generated commit message
with retry functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		autocommit.HandleAutoCommitPerFile()
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage GitGud configuration",
	Long:  `View and manage your GitGud configuration including OpenAI API key settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ShowConfigStatus()
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset and update your API key configuration",
	Long:  `Reset your OpenAI API key configuration and set up a new one.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.HandleConfigReset()
	},
}

var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Show detailed information about the last commit",
	Long:  `Display comprehensive information about the most recent commit including metadata and changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		git.HandleLastCommit()
	},
}

// Git passthrough command
var gitCmd = &cobra.Command{
	Use:                "git",
	Short:              "Execute standard Git commands",
	Long:               `Execute any standard Git command. All arguments are passed directly to Git.`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a Git command")
			os.Exit(1)
		}
		commands.HandleGitCommand(args[0], args[1:])
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add config subcommands
	configCmd.AddCommand(configResetCmd)

	// Add all commands to root
	rootCmd.AddCommand(autocommitCmd)
	rootCmd.AddCommand(acpfCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(lastCmd)
	rootCmd.AddCommand(gitCmd)

	// Add standard Git commands as direct subcommands for convenience
	addGitCommand("init", "Initialize a new Git repository")
	addGitCommand("add", "Add file contents to the index")
	addGitCommand("commit", "Record changes to the repository")
	addGitCommand("status", "Show the working tree status")
	addGitCommand("log", "Show commit logs")
	addGitCommand("diff", "Show changes between commits, commit and working tree, etc")
	addGitCommand("push", "Update remote refs along with associated objects")
	addGitCommand("pull", "Fetch from and integrate with another repository or a local branch")
	addGitCommand("branch", "List, create, or delete branches")
	addGitCommand("checkout", "Switch branches or restore working tree files")
	addGitCommand("merge", "Join two or more development histories together")
	addGitCommand("clone", "Clone a repository into a new directory")
	addGitCommand("fetch", "Download objects and refs from another repository")
	addGitCommand("reset", "Reset current HEAD to the specified state")
	addGitCommand("tag", "Create, list, delete or verify a tag object signed with GPG")
	addGitCommand("stash", "Stash the changes in a dirty working directory away")
}

func addGitCommand(name, description string) {
	cmd := &cobra.Command{
		Use:                name,
		Short:              description,
		Long:               fmt.Sprintf("%s - passes all arguments to git %s", description, name),
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			commands.HandleGitCommand(name, args)
		},
	}
	rootCmd.AddCommand(cmd)
}
