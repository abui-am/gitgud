# GitGud

A simple Git wrapper CLI app written in Go.

## Features

- Executes all standard Git commands
- Provides cleaner success messages
- Passes all arguments to the underlying Git command
- Fall-through behavior for any Git command not explicitly listed
- AI-powered autocommit feature to generate commit messages following Conventional Commits format

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
