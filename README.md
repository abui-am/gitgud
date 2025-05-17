![Group 1 (2)](https://github.com/user-attachments/assets/a642c270-1e2a-499e-880e-a4feec914445)

## Features
![download](https://github.com/user-attachments/assets/edbe19aa-129f-49c5-bd31-53bf04b73a44)
- Executes all standard Git commands
- Provides cleaner success messages
- Passes all arguments to the underlying Git command
- Fall-through behavior for any Git command not explicitly listed
- AI-powered autocommit feature to generate commit messages following Conventional Commits format
- Customizable commit message rules via .autocommit.md file

## Installation

```
# The default build command would create gitgud.exe (based on the directory name)
# To create gg.exe instead, use:
go build -o gg.exe
```

## Setup for AI features

Create a `.env` file in the same directory as the executable with your OpenAI API key:

```
OPENAI_API_KEY=your_openai_api_key_here
```

### Customizing Autocommit Rules

You can modify the commit message format by creating or editing the `.autocommit.md` file. This file contains the rules that will be sent to the AI when generating commit messages.

Note: The `.autocommit.md` file is listed in `.gitignore`, so you'll need to create it manually in each repository where you use this tool. This ensures your commit message customizations don't get committed to your repository.

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
2. Sends the diff to OpenAI to generate a meaningful commit message following the Conventional Commits format
3. Shows you the suggested commit message and asks for confirmation
4. If you confirm, stages all changes and commits them with the AI-generated message

**Important**: You can customize the commit message format by creating or editing the `.autocommit.md` file. Since this file is in `.gitignore`, you'll need to create it in each repository where you use this tool.

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

You can customize the commit message format by editing the `.autocommit.md` file.
