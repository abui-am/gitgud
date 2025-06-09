package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func ParseFileSelection(input string, files []string) ([]string, error) {
	if input == "" {
		return []string{}, nil
	}

	var selectedFiles []string
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Try to parse as number
		var fileNum int
		_, err := fmt.Sscanf(part, "%d", &fileNum)
		if err != nil {
			return nil, fmt.Errorf("invalid file number: %s", part)
		}

		if fileNum < 1 || fileNum > len(files) {
			return nil, fmt.Errorf("file number %d is out of range (1-%d)", fileNum, len(files))
		}

		selectedFiles = append(selectedFiles, files[fileNum-1])
	}

	return selectedFiles, nil
}

func SelectFilesWithArrows(files []string) ([]string, error) {
	if len(files) == 0 {
		return []string{}, nil
	}

	// Create options for the selection prompt
	choices := make([]string, len(files)+2)
	choices[0] = "✅ All files"
	choices[1] = "❌ Exit"
	for i, file := range files {
		choices[i+2] = fmt.Sprintf("📄 %s", file)
	}

	// Create the selection prompt using promptui
	prompt := promptui.Select{
		Label: "Select files to commit",
		Items: choices,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▶ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "{{ . | red | cyan }}",
		},
		Size: 10,
	}

	selectedIndex, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("error running selection prompt: %v", err)
	}

	// Handle the selection
	switch selectedIndex {
	case 0: // All files
		return files, nil
	case 1: // Exit
		return nil, fmt.Errorf("user chose to exit")
	default: // Specific file
		fileIndex := selectedIndex - 2
		if fileIndex >= 0 && fileIndex < len(files) {
			return []string{files[fileIndex]}, nil
		}
		return []string{}, nil
	}
}

func SelectMultipleFilesWithArrows(files []string) ([]string, error) {
	if len(files) == 0 {
		return []string{}, nil
	}

	var selectedFiles []string
	remaining := make([]string, len(files))
	copy(remaining, files)

	fmt.Println("=== Interactive File Selection ===")
	fmt.Println("Select files one by one. You can repeat this process until you're done.")
	fmt.Println()

	for len(remaining) > 0 {
		fmt.Printf("Files selected so far: %d\n", len(selectedFiles))
		if len(selectedFiles) > 0 {
			fmt.Println("Selected files:")
			for _, f := range selectedFiles {
				fmt.Printf("  ✅ %s\n", f)
			}
			fmt.Println()
		}

		fmt.Printf("Remaining files: %d\n", len(remaining))

		// Create options
		choices := make([]string, len(remaining)+3)
		choices[0] = "✅ Proceed with selected files"
		choices[1] = "📄 Select all remaining files"
		choices[2] = "❌ Exit"
		for i, file := range remaining {
			choices[i+3] = fmt.Sprintf("📄 %s", file)
		}

		// Create selection prompt using promptui
		prompt := promptui.Select{
			Label: "Choose an action",
			Items: choices,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}:",
				Active:   "▶ {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "{{ . | red | cyan }}",
			},
			Size: 15,
		}

		selectedIndex, _, err := prompt.Run()
		if err != nil {
			return nil, fmt.Errorf("error running selection prompt: %v", err)
		}

		switch selectedIndex {
		case 0: // Proceed with selected files
			return selectedFiles, nil
		case 1: // Select all remaining
			selectedFiles = append(selectedFiles, remaining...)
			return selectedFiles, nil
		case 2: // Exit
			return nil, fmt.Errorf("user chose to exit")
		default: // Select specific file
			fileIndex := selectedIndex - 3
			if fileIndex >= 0 && fileIndex < len(remaining) {
				// Add to selected files
				selectedFiles = append(selectedFiles, remaining[fileIndex])
				// Remove from remaining
				remaining = append(remaining[:fileIndex], remaining[fileIndex+1:]...)
				fmt.Printf("\n✅ Added: %s\n\n", selectedFiles[len(selectedFiles)-1])
			}
		}
	}

	return selectedFiles, nil
}
