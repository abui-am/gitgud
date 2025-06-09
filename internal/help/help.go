package help

import "fmt"

// ShowUsage displays the main usage information
func ShowUsage() {
	fmt.Println("GitGud - A wrapper around Git")
	fmt.Println("Usage: gg <command> [<args>]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  init                    Initialize a new repository")
	fmt.Println("  add <file>              Add file contents to the index")
	fmt.Println("  commit -m <message>     Record changes to the repository")
	fmt.Println("  status                  Show the working tree status")
	fmt.Println("  log                     Show commit logs")
	fmt.Println("  diff                    Show changes between commits, commit and working tree, etc")
	fmt.Println("  autocommit (or ac)      Automatically add all changes and generate commit message using AI")
	fmt.Println("  autocommit-per-file (or acpf)  Interactively select and batch commit files with AI-generated messages")
	fmt.Println("  config                  View or update your configuration settings")
	fmt.Println("  branch                  List, create, or delete branches")
	fmt.Println("  checkout                Switch branches or restore working tree files")
	fmt.Println("  push                    Update remote refs along with associated objects")
	fmt.Println("  pull                    Fetch from and integrate with another repository or a local branch")
	fmt.Println("  fetch                   Download objects and refs from another repository")
	fmt.Println("  merge                   Join two or more development histories together")
	fmt.Println("  rebase                  Reapply commits on top of another base tip")
	fmt.Println("  stash                   Stash the changes in a dirty working directory away")
	fmt.Println("  remote                  Manage set of tracked repositories")
	fmt.Println("  tag                     Create, list, delete or verify a tag object signed with GPG")
	fmt.Println("  help                    Display help information")
}

// ShowShortUsage displays the short usage information
func ShowShortUsage() {
	fmt.Println("Usage: gg <command> [<args>]")
	fmt.Println("Available commands:")
	fmt.Println("  init")
	fmt.Println("  add <file>")
	fmt.Println("  commit -m <message>")
	fmt.Println("  status")
	fmt.Println("  log")
	fmt.Println("  diff")
	fmt.Println("  autocommit (or ac)")
	fmt.Println("  autocommit-per-file (or acpf)")
	fmt.Println("  config")
	fmt.Println("  last")
}
