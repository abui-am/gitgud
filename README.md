# GitGud

A simple Git wrapper CLI app written in Go.

## Features

- Executes all standard Git commands
- Provides cleaner success messages
- Passes all arguments to the underlying Git command
- Fall-through behavior for any Git command not explicitly listed

## Installation

```
go build -o gitgud.exe
```

## Usage

```
./gitgud init                        # Initialize a new repository
./gitgud add <file>                  # Add file to staging area
./gitgud commit -m "commit message"  # Commit staged changes
./gitgud log                         # View commit history
./gitgud status                      # Check status of working directory
./gitgud diff                        # View differences
./gitgud branch                      # List, create, or delete branches
./gitgud checkout <branch>           # Switch branches
./gitgud push                        # Push to remote repository
./gitgud pull                        # Pull from remote repository
```

GitGud passes all arguments directly to Git, so any valid Git command and options will work.
