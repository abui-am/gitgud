![Group 1 (2)](https://github.com/user-attachments/assets/a642c270-1e2a-499e-880e-a4feec914445)

```bash
Generating commit message with AI...
Generated commit message:
feat(auto): allow customizable autocommit rules via .autocommit.md
Do you want to commit with this message? (y/n/r=retry):
```

## Features

![download](https://github.com/user-attachments/assets/edbe19aa-129f-49c5-bd31-53bf04b73a44)

- Executes all standard Git commands
- Provides cleaner success messages
- Passes all arguments to the underlying Git command
- Fall-through behavior for any Git command not explicitly listed
- AI-powered autocommit feature to generate commit messages following Conventional Commits format
- **Cobra CLI Framework**: Professional command-line interface with automatic help generation
- **Retry functionality**: Regenerate commit messages until you get one you like
- **Batch commit mode**: Select multiple files and commit them together with one unified message
- **Individual file selection**: Choose specific files from directories instead of entire folders
- Customizable commit message rules via .autocommit.md file
- Flexible configuration options for OpenAI API key
- API key validation and troubleshooting

## Installation

```bash
# Build the application (now using Cobra CLI framework)
go build -o gg.exe
```

**Note**: The application now uses the Cobra CLI framework for better command structure and help system.

### Adding to System PATH

For convenience, you can add the directory containing `gg.exe` to your system PATH to run it from any location:

**On Windows:**

1. Open System Properties → Advanced → Environment Variables
2. Edit the PATH variable and add the directory containing gg.exe
3. Open a new command prompt for the changes to take effect

**On Linux/macOS:**

```
export PATH="$PATH:/path/to/directory/containing/gg"
```

## Configuration Management

### Viewing and Managing Configuration

The `gg config` command helps you manage your OpenAI API key and other settings:

```
gg config             # View current configuration status
gg config reset       # Reset and update your API key
```

When viewing your configuration, the app will show all locations where it looks for your API key, whether each exists, and whether the keys are valid or invalid.

### Handling Invalid API Keys

If your API key is invalid or expired, you'll receive a specific error message when using features that require the OpenAI API. You can:

1. Run `gg config reset` to update your API key
2. The app will prompt you to enter a new key and choose where to save it
3. The new key will be validated immediately to ensure it works

## Setup for OpenAI API Key

The application will look for your OpenAI API key in several locations, in the following order:

1. System environment variable (`OPENAI_API_KEY`)
2. `.env` file in the current working directory
3. Config file in your home directory (`~/.gg/config.json`)
4. `.env` file or config in the same directory as the executable
5. If no API key is found, it will prompt you to enter one interactively

### Option 1: Environment Variable

Set the `OPENAI_API_KEY` environment variable:

**On Windows:**

```
set OPENAI_API_KEY=your_api_key_here
```

**On Linux/macOS:**

```
export OPENAI_API_KEY=your_api_key_here
```

### Option 2: Local .env File

Create a `.env` file in your project directory:

```
OPENAI_API_KEY=your_openai_api_key_here
```

### Option 3: User Home Configuration (Recommended)

The recommended approach is to store your API key in your user home directory:

```
~/.gg/config.json
```

With content:

```json
{
  "openai_api_key": "your_api_key_here"
}
```

This file will be created automatically if you run the autocommit command and choose to save the API key to your home directory when prompted.

### Customizing Autocommit Rules

You can customize the commit message format by creating or editing the `.autocommit.md` file. This file contains the rules that will be sent to the AI when generating commit messages.

The application follows this hierarchy for finding autocommit rules:

1. Project root `.autocommit.md` (highest priority)

   - Create this file in your project's root directory to override default rules
   - This file is listed in `.gitignore`, so it won't be committed to your repository
   - Perfect for project-specific commit message conventions

2. Default `.autocommit.md` (fallback)

   - Located in the same directory as the `gg` executable
   - Used when no project-specific rules are found
   - Provides sensible defaults for all users

3. Built-in default (lowest priority)
   - Used only if no `.autocommit.md` files are found
   - Follows the Conventional Commits format

Example `.autocommit.md` content:

```markdown
# Autocommit Rules

Please follow the Conventional Commits 1.0.0 specification for the commit message.

<type>[optional scope]: <description>

[optional body]

[optional footer(s)]

Common types include:

- feat: A new feature
- fix: A bug fix
- docs: Documentation changes
- style: Code style changes (formatting, etc.)
- refactor: Code refactoring
- test: Testing changes
- chore: Maintenance tasks
```

## Usage

### Getting Help

```bash
gg --help                        # Show all available commands
gg <command> --help             # Show help for specific command
gg autocommit --help            # Show help for autocommit
gg config --help                # Show help for config management
```

### Basic Commands

```bash
gg init                         # Initialize a new repository
gg add <file>                   # Add file to staging area
gg commit -m "commit message"   # Commit staged changes
gg log                          # View commit history
gg status                       # Check status of working directory
gg diff                         # View differences
gg branch                       # List, create, or delete branches
gg checkout <branch>            # Switch branches
gg push                         # Push to remote repository
gg pull                         # Pull from remote repository
```

### AI-Powered Commands

```bash
gg autocommit                   # Auto-add all changes and generate AI commit message
gg ac                           # Alias for autocommit
gg autocommit-per-file          # Interactive per-file commits with AI messages
gg acpf                         # Alias for autocommit-per-file
```

### Configuration

```bash
gg config                       # Show current configuration
gg config reset                 # Reset and update API key
gg last                         # Show detailed information about the last commit
```

GitGud passes all arguments directly to Git, so any valid Git command and options will work. The Cobra framework provides comprehensive help for all commands.

### Viewing Last Commit Information

The `last` command provides detailed information about the most recent commit:

```
./gg last
```

This will show:

- Commit hash
- Author name
- Commit date
- Commit message
- Detailed changes (files modified, insertions, deletions)

## Using Autocommit

The `autocommit` command (or its shorter alias `ac`):

1. Automatically detects all changes in your repository
2. Includes the current branch name for context-aware commit messages
3. Provides context from the last commit (hash, author, date, and message)
4. Sends the diff to OpenAI to generate a meaningful commit message following the Conventional Commits format
5. Shows you the suggested commit message with retry options
6. **NEW**: Press `r` or `retry` to regenerate the message until you're satisfied
7. If you confirm, stages all changes and commits them with the AI-generated message

**Important**: You can customize the commit message format by creating or editing the `.autocommit.md` file. Since this file is in `.gitignore`, you'll need to create it in each repository where you use this tool.

The AI assistant takes into account:

- The current branch name (e.g., `feature/user-auth`, `fix/login-bug`)
- Last commit information (e.g., "Last commit: a1b2c3d by John Doe on Mon Jan 1 12:00:00 2024 - feat: initial implementation")
- The changes in your working directory
- Your custom commit message rules from `.autocommit.md`

This helps generate more contextually relevant commit messages that align with your branch's purpose and development history.

Example:

```
./gg autocommit
# or
./gg ac
```

### Conventional Commits Format

The autocommit command generates commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Common types include:

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that don't affect code functionality (formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features

## Retry Functionality

### What's New

Both `autocommit` and `autocommit-per-file` now support **retry functionality**. If you don't like the AI-generated commit message, you can regenerate it without starting over.

### How to Use

When prompted with a commit message, you now have these options:

- `y` or `yes` - Commit with the current message
- `n` or `no` - Cancel/skip the commit
- `r` or `retry` - **Generate a new message** 🔄
- `exit` - Exit the program (acpf only)

### Example

```bash
gg autocommit

Generated commit message:

fix: update configuration settings

Do you want to commit with this message? (y/n/r=retry): r

Regenerating commit message...

Generated commit message:

feat(config): implement dynamic configuration management

Do you want to commit with this message? (y/n/r=retry): y
Changes committed successfully!
```

### Why Use Retry?

- **AI Variability**: Each generation can produce different results
- **Better Quality**: Keep trying until you get a message you like
- **No Wasted Work**: Same changes, just different AI interpretation
- **Learn Patterns**: See how AI interprets your changes differently

## Autocommit Per File Feature

### Overview

The `autocommit-per-file` (or `acpf`) feature allows you to commit files individually or in batches with AI-generated commit messages. This feature has been significantly enhanced with:

- **Individual file selection**: Choose specific files instead of entire directories
- **Batch commit mode**: Select multiple files and commit them together with one unified message
- **Retry functionality**: Regenerate commit messages until you're satisfied
- **Improved file listing**: See individual files instead of directory names

### Usage

```bash
gg autocommit-per-file
# or
gg acpf
```

### How it Works

1. **List Changed Files**: The tool displays all individual files with changes (no more directory grouping!)
2. **Interactive File Selection**: You can select which files to commit using arrow key navigation:
   - Use ↑/↓ arrow keys to navigate
   - Press Enter to select an option
   - Select individual files one by one
   - Select all remaining files at once
   - Proceed with currently selected files
   - Exit at any time
3. **Batch Processing**: Selected files are processed together:
   - **NEW**: All selected files are committed together as one batch
   - Shows combined diff for all selected files
   - Asks for context for the entire batch (not per file)
   - Generates one unified AI-powered commit message
   - **NEW**: Option to retry message generation with `r` or `retry`
   - Commits all selected files with one message
4. **Loop Process**: After processing the batch, you can continue with remaining files or exit

### Features

- ✅ Interactive file selection with arrow key navigation
- ✅ **NEW**: Individual file listing (no more directory grouping)
- ✅ **NEW**: Batch commit mode - commit multiple files with one message
- ✅ **NEW**: Retry functionality - regenerate messages until satisfied
- ✅ AI-generated commit messages for file batches
- ✅ Custom context for each batch
- ✅ Loop through the process multiple times
- ✅ Exit at any time using the menu
- ✅ Handles new, modified, staged, and renamed files
- ✅ Uses your existing autocommit rules from `.autocommit.md`
- ✅ Beautiful terminal UI with colors and icons

### Example Workflow

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

--- Processing 2 selected file(s) ---
  - main.go
  - test_feature.txt

Enter additional context for these 2 file(s) (press Enter to skip): Added new autocommit per file feature with batch processing
Generating commit message for 2 file(s)...

Generated commit message for batch:

feat: add autocommit per file functionality with batch processing

Do you want to commit these files with this message? (y/n/r=retry/exit): r

Regenerating commit message for 2 file(s)...

Generated commit message for batch:

feat(autocommit): implement batch processing for selective file commits

Do you want to commit these files with this message? (y/n/r=retry/exit): y
Successfully committed 2 file(s) in one commit

--- Processing complete ---
Continue with remaining files? (y/n): y
```

### Requirements

- OpenAI API key configured (same as regular autocommit)
- Git repository initialized
- Changed files to commit

### Tips

- **NEW**: Use the retry feature (`r`) to get better commit messages
- Use specific context to get better commit messages for the batch
- You can skip batches by answering 'n' when asked to commit
- The feature respects your `.autocommit.md` rules
- **NEW**: Selected files are committed together in one commit (not separately)
- Use arrow keys for easy navigation - no need to type numbers
- You can select files incrementally and see your progress
- The interface shows clear visual feedback with colors and icons
- **NEW**: Individual files are shown instead of directory names (e.g., `internal/config/config.go` instead of `internal/`)

## Custom Context

### What is Custom Context?

Custom context is a powerful feature in GitGud that helps you create more meaningful commit messages. It allows you to provide additional information about your changes, which the AI assistant uses along with your branch name and last commit details to generate better commit messages.

### How to Use Custom Context

1. **Setup Your Rules**
   Create a `.autocommit.md` file in your project root:

   ```markdown
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description>

   - #time: The time spent on the task (use the #time tag in the custom context to add the time spent on the task)

   Example:
   feat(auth): #1h add login form

   Types:

   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   ```

2. **Run Autocommit**

   ```bash
   gg autocommit
   ```

3. **Add Your Context**
   When prompted, enter your time spent:

   ```
   #1h
   ```

4. **Review Generated Message**
   The AI will create a commit message using your context:
   ```
   feat(auth): #1h add login form
   ```

### How It Works

GitGud combines three sources of information to generate your commit message:

1. **Your Custom Context**

   - Time spent on the task
   - Any additional tags you provide

2. **Last Commit Information**

   - Commit hash
   - Author name
   - Commit date
   - Previous commit message

3. **Branch Name**
   - `feature/user-auth` → Authentication-related changes
   - `fix/login-bug` → Bug fixes in login
   - `docs/api-update` → Documentation updates

### Examples of Using Custom Context

1. **Time Tracking**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description>
   - #time: The time spent on the task (use the #time tag in the custom context)

   Example:
   feat(auth): #1h add login form

   Types:
   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   ```

   ```bash
   # Create feature branch
   gg checkout -b feature/user-auth

   # Make changes to auth system...

   # Commit with time context
   gg autocommit
   # Enter: #2h

   # Result:
   feat(auth): #2h implement OAuth2 authentication
   ```

2. **Multiple Context Tags**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description> <#complexity> <#impact>
   - #time: Time spent on the task
   - #complexity: Complexity level (low/medium/high)
   - #impact: Impact level (low/medium/high/critical)

   Example:
   fix(auth): #1h resolve validation #complexity:high #impact:critical

   Types:
   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   ```

   ```bash
   # Create fix branch
   gg checkout -b fix/login-validation

   # Fix validation issues...

   # Commit with multiple context tags
   gg autocommit
   # Enter: #1h #complexity:high #impact:critical

   # Result:
   fix(auth): #1h resolve login validation issues #complexity:high #impact:critical
   ```

### Branch Name Context

1. **Feature Branch**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description>
   - #time: The time spent on the task
   - Scope should match branch name prefix (feature/fix/docs)

   Example:
   feat(payment): #2h add payment processing

   Types:
   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   ```

   ```bash
   # Branch name indicates feature type
   gg checkout -b feature/payment-gateway

   # Implement payment system...

   # Commit with time context
   gg autocommit
   # Enter: #4h

   # Result:
   feat(payment): #4h implement Stripe payment integration
   ```

2. **Bug Fix Branch**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description>
   - #time: The time spent on the task
   - Scope should match branch name prefix (feature/fix/docs)

   Example:
   fix(api): #1h resolve timeout issues

   Types:
   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   ```

   ```bash
   # Branch name indicates bug fix
   gg checkout -b fix/api-timeout

   # Fix timeout issues...

   # Commit with time context
   gg autocommit
   # Enter: #1h

   # Result:
   fix(api): #1h resolve request timeout issues
   ```

### Last Commit Context

1. **Related Feature**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description>
   - #time: The time spent on the task
   - Consider last commit context for related changes

   Example:
   feat(profile): #2h add user settings

   Types:
   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   ```

   ```bash
   # Previous commit: feat(auth): #2h add login form
   # Current branch: feature/user-profile

   # Add profile functionality...

   # Commit with time context
   gg autocommit
   # Enter: #3h

   # Result:
   feat(profile): #3h implement user profile management
   ```

2. **Bug Fix Follow-up**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description>
   - #time: The time spent on the task
   - Consider last commit context for related fixes

   Example:
   fix(auth): #1h handle edge cases

   Types:
   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   ```

   ```bash
   # Previous commit: fix(auth): #1h resolve login validation
   # Current branch: fix/auth-edge-cases

   # Fix additional edge cases...

   # Commit with time context
   gg autocommit
   # Enter: #2h

   # Result:
   fix(auth): #2h handle additional login edge cases
   ```

### Combined Context Examples

1. **Feature with Related Fix**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <branch-name> <#time> <description>
   - #time: The time spent on the task
   - branch-name: The name of the branch (will be the card name from JIRA)
   - Consider both branch name and last commit context

   Example:
   fix(payment): #1h improve validation

   Types:
   - feat: JIRA-1234: A new feature
   - fix: JIRA-1235: A bug fix
   - docs: JIRA-1236: Documentation changes
   ```

   ```bash
   # Previous commit: feat(payment): #4h add payment processing
   # Current branch: JIRA-1235

   # Fix payment validation...

   # Commit with time context
   gg autocommit
   # Enter: #1h

   # Result:
   fix(payment): JIRA-1235 #1h improve payment validation rules
   ```

2. **Documentation Update**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description>
   - #time: The time spent on the task
   - Consider last commit for related documentation

   Example:
   docs(api): #1h update endpoint docs

   Types:
   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   ```

   ```bash
   # Previous commit: feat(api): #3h add new endpoints
   # Current branch: docs/api-update

   # Update API documentation...

   # Commit with time context
   gg autocommit
   # Enter: #2h

   # Result:
   docs(api): #2h document new API endpoints
   ```

3. **Refactoring with Context**

   ```bash
   # .autocommit.md configuration
   Please follow the Conventional Commits format: <type>(<scope>): <#time> <description>
   - #time: The time spent on the task
   - Consider last commit for related refactoring

   Example:
   refactor(ui): #2h optimize components

   Types:
   - feat: A new feature
   - fix: A bug fix
   - docs: Documentation changes
   - refactor: Code restructuring
   ```

   ```bash
   # Previous commit: feat(ui): #5h implement dashboard
   # Current branch: refactor/optimize-ui

   # Optimize UI components...

   # Commit with time context
   gg autocommit
   # Enter: #3h

   # Result:
   refactor(ui): #3h optimize dashboard performance
   ```

### Troubleshooting

If your custom context isn't working as expected:

1. **Check File Format**

   - Ensure your `.autocommit.md` file is in the correct location
   - Verify the file format matches the example

2. **Verify Tag Usage**

   - Make sure you're using the correct tag format
   - Check that tags are properly formatted (e.g., `#1h` not `1h`)

3. **Check Input**

   - Confirm your custom context is entered correctly
   - Ensure there are no extra spaces or special characters

4. **Review Output**
   - Check if the AI is following your custom rules
   - Verify the generated commit message format
