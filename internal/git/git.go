package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ExecuteGitCommand(command string, args ...string) error {
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

func HasChangesToCommit() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error checking git status: %v\n", err)
		os.Exit(1)
	}

	return len(output) > 0
}

func GetGitDiff() (string, error) {
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

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting current branch: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func GetLastCommitMetadata() (string, error) {
	// Get the last commit's metadata using git log
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%h|%an|%ad|%s")
	output, err := cmd.Output()
	if err != nil {
		// If there's no previous commit, return empty string
		if strings.Contains(err.Error(), "fatal: bad default revision") {
			return "", nil
		}
		return "", fmt.Errorf("error getting last commit metadata: %v", err)
	}

	// Parse the output
	parts := strings.Split(string(output), "|")
	if len(parts) != 4 {
		return "", fmt.Errorf("unexpected commit metadata format")
	}

	// Format: commit hash, author name, date, and subject
	return fmt.Sprintf("Last commit: %s by %s on %s - %s",
		parts[0], // hash
		parts[1], // author
		parts[2], // date
		parts[3], // subject
	), nil
}

func GetChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting git status: %v", err)
	}

	var files []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Git status --porcelain format: XY filename
		// We want to extract just the filename
		if len(line) >= 3 {
			filename := strings.TrimSpace(line[2:])

			// Handle renamed files: "R  oldfile -> newfile"
			if line[0] == 'R' && strings.Contains(filename, " -> ") {
				// For renamed files, we want the new filename
				parts := strings.Split(filename, " -> ")
				if len(parts) == 2 {
					filename = strings.TrimSpace(parts[1])
				}
			}

			// Check if this is an untracked directory (ends with /)
			if strings.HasSuffix(filename, "/") {
				// Expand directory to individual files
				dirFiles, err := getFilesInDirectory(filename)
				if err != nil {
					// If we can't expand, just add the directory
					files = append(files, filename)
				} else {
					files = append(files, dirFiles...)
				}
			} else {
				files = append(files, filename)
			}
		}
	}

	return files, nil
}

// Helper function to get all files in a directory recursively
func getFilesInDirectory(dirPath string) ([]string, error) {
	var files []string

	// Use git ls-files to get untracked files in the directory
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard", dirPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(string(output)) == "" {
		return nil, fmt.Errorf("no files found in directory")
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}

func GetFileDiff(filename string) (string, error) {
	// Check if file is staged
	stagedCmd := exec.Command("git", "diff", "--staged", "--", filename)
	stagedOutput, err := stagedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting staged diff for %s: %v", filename, err)
	}

	// Check if file has unstaged changes
	unstagedCmd := exec.Command("git", "diff", "--", filename)
	unstagedOutput, err := unstagedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting unstaged diff for %s: %v", filename, err)
	}

	// Check if it's an untracked file
	statusCmd := exec.Command("git", "status", "--porcelain", "--", filename)
	statusOutput, err := statusCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting status for %s: %v", filename, err)
	}

	statusLine := strings.TrimSpace(string(statusOutput))

	// Combine outputs
	combinedDiff := string(stagedOutput) + string(unstagedOutput)

	// If it's an untracked file, show that it's new
	if len(statusLine) > 0 && statusLine[0] == '?' {
		combinedDiff += fmt.Sprintf("\nNew file: %s", filename)

		// Try to show the content of new file (if it's text and not too large)
		if fileContent, err := os.ReadFile(filename); err == nil && len(fileContent) < 2000 {
			combinedDiff += fmt.Sprintf("\nFile content:\n%s", string(fileContent))
		}
	}

	return combinedDiff, nil
}

func HandleLastCommit() {
	// Get last commit metadata
	lastCommitInfo, err := GetLastCommitMetadata()
	if err != nil {
		if strings.Contains(err.Error(), "fatal: bad default revision") {
			fmt.Println("No commits found in the repository.")
			os.Exit(0)
		}
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Print the information
	fmt.Println("\nLast Commit Information:")
	fmt.Println("=======================")
	fmt.Println(lastCommitInfo)
}
