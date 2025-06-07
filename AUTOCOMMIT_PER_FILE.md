# Autocommit Per File Feature

## Overview

The `autocommit-per-file` (or `acpf`) feature allows you to commit files individually with AI-generated commit messages. This is particularly useful when you have multiple changes that should be committed separately with different commit messages.

## Usage

```bash
gg autocommit-per-file
# or
gg acpf
```

## How it Works

1. **List Changed Files**: The tool displays all files with changes (modified, added, or untracked)
2. **Interactive File Selection**: You can select which files to commit using arrow key navigation:
   - Use ↑/↓ arrow keys to navigate
   - Press Enter to select an option
   - Select individual files one by one
   - Select all remaining files at once
   - Proceed with currently selected files
   - Exit at any time
3. **File Processing**: For each selected file:
   - Shows the diff for that specific file
   - Asks for additional context (optional)
   - Generates an AI-powered commit message
   - Asks for confirmation before committing
4. **Loop Process**: After processing all selected files, you can continue with remaining files or exit

## Features

- ✅ Interactive file selection with arrow key navigation
- ✅ AI-generated commit messages per file
- ✅ Custom context for each file
- ✅ Loop through the process multiple times
- ✅ Exit at any time using the menu
- ✅ Handles new, modified, and staged files
- ✅ Uses your existing autocommit rules from `.autocommit.md`
- ✅ Beautiful terminal UI with colors and icons

## Example Workflow

```
=== Autocommit Per File ===
This will help you commit files individually with AI-generated commit messages.
Use arrow keys to navigate and select files. You can select multiple files interactively.

=== Interactive File Selection ===
Select files one by one. You can repeat this process until you're done.

Files selected so far: 0
Remaining files: 3

Choose an action:
▶ ✅ Proceed with selected files
  📄 Select all remaining files
  ❌ Exit
  📄 main.go
  📄 README.md
  📄 test_feature.txt

[User navigates with arrows and selects main.go, then test_feature.txt]

--- Processing file: main.go ---
Enter additional context for main.go (press Enter to skip): Added new autocommit per file feature
Generating commit message for main.go...

Generated commit message for main.go:

feat: add autocommit per file functionality for selective commits

Do you want to commit this file with this message? (y/n/exit): y
Successfully committed main.go

--- Processing file: test_feature.txt ---
Enter additional context for test_feature.txt (press Enter to skip):
Generating commit message for test_feature.txt...

Generated commit message for test_feature.txt:

docs: add test file demonstrating autocommit per file feature

Do you want to commit this file with this message? (y/n/exit): y
Successfully committed test_feature.txt

--- Batch complete ---
Continue with remaining files? (y/n): y
```

## Requirements

- OpenAI API key configured (same as regular autocommit)
- Git repository initialized
- Changed files to commit

## Tips

- Use specific context to get better commit messages
- You can skip files by answering 'n' when asked to commit
- The feature respects your `.autocommit.md` rules
- Each commit is made separately, so you get granular commit history
- Use arrow keys for easy navigation - no need to type numbers
- You can select files incrementally and see your progress
- The interface shows clear visual feedback with colors and icons
