package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/gitgud/internal/config"
	"github.com/user/gitgud/internal/git"
)

var (
	cfgManager *config.Manager
	gitWrapper *git.Wrapper
)

var rootCmd = &cobra.Command{
	Use:   "gg",
	Short: "GitGud - A smart Git wrapper with AI assistance",
	Long: `GitGud is a Git wrapper that provides AI-powered assistance for commit messages,
interactive file selection, and enhanced Git workflows.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		// Handle direct git command passthrough
		if err := gitWrapper.ExecuteCommand(args[0], args[1:]...); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Initialize components
	cfgManager = config.NewManager()
	gitWrapper = git.NewWrapper()

	// Add all subcommands
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newCommitCmd())
	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newLogCmd())
	rootCmd.AddCommand(newDiffCmd())
	rootCmd.AddCommand(newAutoCommitCmd())
	rootCmd.AddCommand(newAutoCommitPerFileCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newLastCmd())
	rootCmd.AddCommand(newInitCmd())

	// Standard git commands
	rootCmd.AddCommand(newBranchCmd())
	rootCmd.AddCommand(newCheckoutCmd())
	rootCmd.AddCommand(newPushCmd())
	rootCmd.AddCommand(newPullCmd())
	rootCmd.AddCommand(newFetchCmd())
	rootCmd.AddCommand(newMergeCmd())
	rootCmd.AddCommand(newRebaseCmd())
	rootCmd.AddCommand(newStashCmd())
	rootCmd.AddCommand(newRemoteCmd())
	rootCmd.AddCommand(newTagCmd())
}

func initConfig() {
	// Initialize configuration if needed
}
