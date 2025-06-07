package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// SelectMultipleFiles allows interactive selection of multiple files using arrow keys
func SelectMultipleFiles(files []string) ([]string, error) {
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
			if len(selectedFiles) == 0 {
				fmt.Println("No files selected. Please select at least one file.")
				continue
			}
			return selectedFiles, nil
		case 1: // Select all remaining files
			selectedFiles = append(selectedFiles, remaining...)
			return selectedFiles, nil
		case 2: // Exit
			return nil, fmt.Errorf("user chose to exit")
		default: // Specific file
			fileIndex := selectedIndex - 3
			if fileIndex >= 0 && fileIndex < len(remaining) {
				selectedFile := remaining[fileIndex]
				selectedFiles = append(selectedFiles, selectedFile)

				// Remove the selected file from remaining
				newRemaining := make([]string, 0, len(remaining)-1)
				for i, file := range remaining {
					if i != fileIndex {
						newRemaining = append(newRemaining, file)
					}
				}
				remaining = newRemaining

				fmt.Printf("Added: %s\n\n", selectedFile)
			}
		}
	}

	return selectedFiles, nil
}

// SelectSingleFile allows selection of a single file from a list
func SelectSingleFile(files []string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("no files available for selection")
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
		Label: "Select a file",
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
		return "", fmt.Errorf("error running selection prompt: %v", err)
	}

	// Handle the selection
	switch selectedIndex {
	case 0: // All files (return first file as representative)
		return "ALL_FILES", nil
	case 1: // Exit
		return "", fmt.Errorf("user chose to exit")
	default: // Specific file
		fileIndex := selectedIndex - 2
		if fileIndex >= 0 && fileIndex < len(files) {
			return files[fileIndex], nil
		}
		return "", fmt.Errorf("invalid selection")
	}
}
