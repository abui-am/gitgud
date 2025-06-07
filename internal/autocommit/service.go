package autocommit

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/user/gitgud/internal/ai"
	"github.com/user/gitgud/internal/git"
	"github.com/user/gitgud/internal/ui"
)

// Service handles autocommit operations
type Service struct {
	git *git.Wrapper
	ai  *ai.Client
}

// NewService creates a new autocommit service
func NewService(gitWrapper *git.Wrapper, aiClient *ai.Client) *Service {
	return &Service{
		git: gitWrapper,
		ai:  aiClient,
	}
}

// AutoCommit performs automatic commit of all changes
func (s *Service) AutoCommit() error {
	// Check if there are any changes to commit
	diff, err := s.git.GetDiff()
	if err != nil {
		return fmt.Errorf("error getting git diff: %v", err)
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("No changes to commit.")
		return nil
	}

	fmt.Println("=== GitGud AutoCommit ===")
	fmt.Println("Analyzing changes and generating commit message...")

	// Generate commit message using AI
	commitMessage, err := s.ai.GenerateCommitMessage(diff, "")
	if err != nil {
		return fmt.Errorf("error generating commit message: %v", err)
	}

	// Display the generated message and ask for confirmation
	fmt.Printf("\n📝 Generated commit message: %s\n\n", commitMessage)

	// Show a preview of what will be committed
	files, err := s.git.GetChangedFiles()
	if err != nil {
		return fmt.Errorf("error getting changed files: %v", err)
	}

	fmt.Println("Files to be committed:")
	for _, file := range files {
		fmt.Printf("  ✓ %s\n", file)
	}

	// Ask for confirmation
	confirm := promptui.Prompt{
		Label:     "Proceed with this commit? (y/n)",
		IsConfirm: true,
	}

	result, err := confirm.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			fmt.Println("Autocommit cancelled.")
			return nil
		}
		return fmt.Errorf("error getting confirmation: %v", err)
	}

	if strings.ToLower(result) == "y" {
		// Commit the changes
		if err := s.git.CommitWithMessage(commitMessage); err != nil {
			return fmt.Errorf("error committing changes: %v", err)
		}
		fmt.Println("✅ Changes committed successfully!")
	} else {
		fmt.Println("Autocommit cancelled.")
	}

	return nil
}

// AutoCommitPerFile performs interactive per-file commits
func (s *Service) AutoCommitPerFile() error {
	// Get list of changed files
	files, err := s.git.GetChangedFiles()
	if err != nil {
		return fmt.Errorf("error getting changed files: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("No changed files found.")
		return nil
	}

	fmt.Println("=== GitGud AutoCommit Per File ===")
	fmt.Printf("Found %d changed file(s)\n\n", len(files))

	reader := bufio.NewReader(os.Stdin)

	for {
		if len(files) == 0 {
			fmt.Println("All files have been processed.")
			break
		}

		fmt.Printf("Remaining files (%d):\n", len(files))
		for i, file := range files {
			fmt.Printf("%d. %s\n", i+1, file)
		}
		fmt.Println()

		// Ask user how they want to select files
		fmt.Println("Options:")
		fmt.Println("1. Select files by number (e.g., '1,3,5' or '1-3')")
		fmt.Println("2. Interactive selection with arrow keys")
		fmt.Println("3. Process all remaining files")
		fmt.Println("4. Exit")
		fmt.Print("Choose option (1-4): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}
		input = strings.TrimSpace(input)

		var selectedFiles []string

		switch input {
		case "1":
			selectedFiles, err = s.selectFilesByNumber(files, reader)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
		case "2":
			selectedFiles, err = ui.SelectMultipleFiles(files)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
		case "3":
			selectedFiles = files
		case "4":
			fmt.Println("Exiting autocommit per file.")
			return nil
		default:
			fmt.Println("Invalid option. Please choose 1-4.")
			continue
		}

		if len(selectedFiles) == 0 {
			continue
		}

		// Process selected files
		if err := s.processSelectedFiles(selectedFiles); err != nil {
			return err
		}

		// Remove processed files from the list
		files = s.removeProcessedFiles(files, selectedFiles)

		if len(files) > 0 {
			fmt.Print("\nContinue with remaining files? (y/n): ")
			continueResponse, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading input: %v", err)
			}
			continueResponse = strings.TrimSpace(continueResponse)
			if strings.ToLower(continueResponse) != "y" {
				fmt.Println("Exiting autocommit per file.")
				break
			}
		}
	}

	return nil
}

// selectFilesByNumber allows user to select files by entering numbers
func (s *Service) selectFilesByNumber(files []string, reader *bufio.Reader) ([]string, error) {
	fmt.Print("Enter file numbers (comma-separated, e.g., '1,3,5' or ranges like '1-3'): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return []string{}, nil
	}

	return parseFileSelection(input, files)
}

// parseFileSelection parses user input for file selection
func parseFileSelection(input string, files []string) ([]string, error) {
	var selectedFiles []string
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check if it's a range (e.g., "1-3")
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start number in range: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end number in range: %s", rangeParts[1])
			}

			if start < 1 || end > len(files) || start > end {
				return nil, fmt.Errorf("invalid range %d-%d (valid range: 1-%d)", start, end, len(files))
			}

			for i := start; i <= end; i++ {
				selectedFiles = append(selectedFiles, files[i-1])
			}
		} else {
			// Single number
			fileNum, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid file number: %s", part)
			}

			if fileNum < 1 || fileNum > len(files) {
				return nil, fmt.Errorf("file number %d is out of range (1-%d)", fileNum, len(files))
			}

			selectedFiles = append(selectedFiles, files[fileNum-1])
		}
	}

	return selectedFiles, nil
}

// processSelectedFiles processes the selected files for commit
func (s *Service) processSelectedFiles(selectedFiles []string) error {
	for _, file := range selectedFiles {
		fmt.Printf("\n--- Processing: %s ---\n", file)

		// Get file diff
		diff, err := s.git.GetFileDiff(file)
		if err != nil {
			fmt.Printf("Error getting diff for %s: %v\n", file, err)
			continue
		}

		// Generate commit message
		fmt.Println("Generating commit message...")
		commitMsg, err := s.ai.GenerateFileCommitMessage(file, diff, "")
		if err != nil {
			fmt.Printf("Error generating commit message for %s: %v\n", file, err)
			continue
		}

		fmt.Printf("📝 Generated message: %s\n", commitMsg)

		// Ask for confirmation
		confirm := promptui.Prompt{
			Label:     fmt.Sprintf("Commit %s with this message? (y/n)", file),
			IsConfirm: true,
		}

		result, err := confirm.Run()
		if err != nil {
			if err == promptui.ErrAbort {
				fmt.Printf("Skipped committing %s\n", file)
				continue
			}
			return fmt.Errorf("error getting confirmation: %v", err)
		}

		if strings.ToLower(result) == "y" {
			// Add and commit the file
			if err := s.git.AddFile(file); err != nil {
				fmt.Printf("Error adding %s: %v\n", file, err)
				continue
			}

			if err := s.git.CommitWithMessage(commitMsg); err != nil {
				fmt.Printf("Error committing %s: %v\n", file, err)
				continue
			}

			fmt.Printf("✅ Successfully committed %s\n", file)
		} else {
			fmt.Printf("Skipped committing %s\n", file)
		}
	}

	return nil
}

// removeProcessedFiles removes processed files from the file list
func (s *Service) removeProcessedFiles(files, processedFiles []string) []string {
	processedMap := make(map[string]bool)
	for _, file := range processedFiles {
		processedMap[file] = true
	}

	var remaining []string
	for _, file := range files {
		if !processedMap[file] {
			remaining = append(remaining, file)
		}
	}

	return remaining
}
