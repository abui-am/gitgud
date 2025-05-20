![Group 1 (2)](https://github.com/user-attachments/assets/a642c270-1e2a-499e-880e-a4feec914445)

```bash
Generating commit message with AI...
Generated commit message:
feat(auto): allow customizable autocommit rules via .autocommit.md
Do you want to commit with this message? (y/n):
```

## Features

![download](https://github.com/user-attachments/assets/edbe19aa-129f-49c5-bd31-53bf04b73a44)

- Executes all standard Git commands
- Provides cleaner success messages
- Passes all arguments to the underlying Git command
- Fall-through behavior for any Git command not explicitly listed
- AI-powered autocommit feature to generate commit messages following Conventional Commits format
- Customizable commit message rules via .autocommit.md file
- Flexible configuration options for OpenAI API key
- API key validation and troubleshooting

## Installation

```
# The default build command would create gitgud.exe (based on the directory name)
# To create gg.exe instead, use:
go build -o gg.exe
```

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

```
./gg init                        # Initialize a new repository
./gg add <file>                  # Add file to staging area
./gg commit -m "commit message"  # Commit staged changes
./gg log                         # View commit history
./gg status                      # Check status of working directory
./gg diff                        # View differences
./gg autocommit                  # Auto-add all changes and generate AI commit message
./gg ac                          # Alias for autocommit
./gg branch                      # List, create, or delete branches
./gg checkout <branch>           # Switch branches
./gg push                        # Push to remote repository
./gg pull                        # Pull from remote repository
```

GitGud passes all arguments directly to Git, so any valid Git command and options will work.

## Using Autocommit

The `autocommit` command (or its shorter alias `ac`):

1. Automatically detects all changes in your repository
2. Includes the current branch name for context-aware commit messages
3. Sends the diff to OpenAI to generate a meaningful commit message following the Conventional Commits format
4. Shows you the suggested commit message and asks for confirmation
5. If you confirm, stages all changes and commits them with the AI-generated message

**Important**: You can customize the commit message format by creating or editing the `.autocommit.md` file. Since this file is in `.gitignore`, you'll need to create it in each repository where you use this tool.

The AI assistant takes into account:

- The current branch name (e.g., `feature/user-auth`, `fix/login-bug`)
- The changes in your working directory
- Your custom commit message rules from `.autocommit.md`

This helps generate more contextually relevant commit messages that align with your branch's purpose.

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
- `test`: Adding or fixing tests
- `chore`: Changes to the build process or auxiliary tools

You can customize the commit message format by creating or editing the `.autocommit.md` file.
