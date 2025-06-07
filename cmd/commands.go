package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/gitgud/internal/ai"
	"github.com/user/gitgud/internal/autocommit"
)

// Basic Git commands
func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Git repository",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("init", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("GitGud repository initialized successfully!")
		},
	}
}

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [files...]",
		Short: "Add file contents to the index",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("add", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Record changes to the repository",
		Run: func(cmd *cobra.Command, args []string) {
			message, _ := cmd.Flags().GetString("message")
			if message == "" {
				fmt.Println("Error: Commit message is required")
				os.Exit(1)
			}

			if err := gitWrapper.ExecuteCommand("commit", "-m", message); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Changes committed successfully!")
		},
	}
	cmd.Flags().StringP("message", "m", "", "Commit message")
	cmd.MarkFlagRequired("message")
	return cmd
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the working tree status",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("status", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newLogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "log",
		Short: "Show commit logs",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("log", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Show changes between commits, commit and working tree, etc",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("diff", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

// AI-powered commands
func newAutoCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "autocommit",
		Aliases: []string{"ac"},
		Short:   "Automatically add all changes and generate commit message using AI",
		Run: func(cmd *cobra.Command, args []string) {
			aiClient, err := ai.NewClient(cfgManager)
			if err != nil {
				fmt.Printf("Error initializing AI client: %v\n", err)
				os.Exit(1)
			}

			autoCommitter := autocommit.NewService(gitWrapper, aiClient)
			if err := autoCommitter.AutoCommit(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newAutoCommitPerFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "autocommit-per-file",
		Aliases: []string{"acpf"},
		Short:   "Interactively select and batch commit files with AI-generated messages",
		Run: func(cmd *cobra.Command, args []string) {
			aiClient, err := ai.NewClient(cfgManager)
			if err != nil {
				fmt.Printf("Error initializing AI client: %v\n", err)
				os.Exit(1)
			}

			autoCommitter := autocommit.NewService(gitWrapper, aiClient)
			if err := autoCommitter.AutoCommitPerFile(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "View or update your configuration settings",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cfgManager.HandleConfigCommand(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newLastCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "last",
		Short: "Show information about the last commit",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ShowLastCommit(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

// Standard Git commands
func newBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "branch",
		Short: "List, create, or delete branches",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("branch", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newCheckoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkout",
		Short: "Switch branches or restore working tree files",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("checkout", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Update remote refs along with associated objects",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("push", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Fetch from and integrate with another repository or a local branch",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("pull", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newFetchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Download objects and refs from another repository",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("fetch", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newMergeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "merge",
		Short: "Join two or more development histories together",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("merge", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newRebaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rebase",
		Short: "Reapply commits on top of another base tip",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("rebase", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newStashCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stash",
		Short: "Stash the changes in a dirty working directory away",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("stash", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newRemoteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remote",
		Short: "Manage set of tracked repositories",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("remote", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func newTagCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tag",
		Short: "Create, list, delete or verify a tag object signed with GPG",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitWrapper.ExecuteCommand("tag", args...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
