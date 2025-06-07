package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Wrapper provides a convenient interface for Git operations
type Wrapper struct{}

// NewWrapper creates a new Git wrapper
func NewWrapper() *Wrapper {
	return &Wrapper{}
}

// ExecuteCommand executes a Git command with the given arguments
func (w *Wrapper) ExecuteCommand(command string, args ...string) error {
	cmd := exec.Command("git", append([]string{command}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Print custom success messages for certain commands
	if err == nil {
		switch command {
		case "init":
			fmt.Println("GitGud repository initialized successfully!")
		case "commit":
			fmt.Println("Changes committed successfully!")
		}
	}

	return err
}

// GetStatus returns the current Git status
func (w *Wrapper) GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git status: %v", err)
	}
	return string(output), nil
}

// GetChangedFiles returns a list of changed files
func (w *Wrapper) GetChangedFiles() ([]string, error) {
	output, err := w.GetStatus()
	if err != nil {
		return nil, err
	}

	var files []string
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Git status --porcelain format: XY filename
		// We want to extract just the filename
		if len(line) >= 3 {
			filename := strings.TrimSpace(line[2:])
			files = append(files, filename)
		}
	}

	return files, nil
}

// GetDiff returns the Git diff for all changes
func (w *Wrapper) GetDiff() (string, error) {
	// First, add all changes to staging area
	addCmd := exec.Command("git", "add", ".")
	if err := addCmd.Run(); err != nil {
		return "", fmt.Errorf("error adding changes: %v", err)
	}

	// Get diff of staged changes
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}

	return string(output), nil
}

// GetFileDiff returns the diff for a specific file
func (w *Wrapper) GetFileDiff(filename string) (string, error) {
	// Check if file exists in working directory or staging area
	var cmd *exec.Cmd

	// Try to get diff from working directory first
	cmd = exec.Command("git", "diff", filename)
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return string(output), nil
	}

	// If no working directory changes, try staged changes
	cmd = exec.Command("git", "diff", "--cached", filename)
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		return string(output), nil
	}

	// If file is untracked, show the entire file content as new
	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard", filename)
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		// File is untracked, read its content
		content, err := os.ReadFile(filename)
		if err != nil {
			return "", fmt.Errorf("error reading untracked file %s: %v", filename, err)
		}

		// Format as a diff showing the entire file as new
		diff := fmt.Sprintf("diff --git a/%s b/%s\nnew file mode 100644\nindex 0000000..1234567\n--- /dev/null\n+++ b/%s\n", filename, filename, filename)
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			diff += "+" + line + "\n"
		}
		return diff, nil
	}

	return "", fmt.Errorf("no changes found for file: %s", filename)
}

// HasChangesToCommit checks if there are any changes to commit
func (w *Wrapper) HasChangesToCommit() bool {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

// GetCurrentBranch returns the current Git branch
func (w *Wrapper) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting current branch: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetLastCommitMetadata returns metadata about the last commit
func (w *Wrapper) GetLastCommitMetadata() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:Commit: %H%nAuthor: %an <%ae>%nDate: %ad%nMessage: %s", "--date=local")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting last commit metadata: %v", err)
	}
	return string(output), nil
}

// ShowLastCommit displays information about the last commit
func (w *Wrapper) ShowLastCommit() error {
	metadata, err := w.GetLastCommitMetadata()
	if err != nil {
		return err
	}

	fmt.Println("=== Last Commit Information ===")
	fmt.Println(metadata)
	fmt.Println()

	// Show the diff of the last commit
	cmd := exec.Command("git", "show", "--stat")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// AddFile adds a specific file to the staging area
func (w *Wrapper) AddFile(filename string) error {
	cmd := exec.Command("git", "add", filename)
	return cmd.Run()
}

// CommitWithMessage commits staged changes with the given message
func (w *Wrapper) CommitWithMessage(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
