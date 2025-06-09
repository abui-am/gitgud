package main

import (
	"os"

	"github.com/user/gitgud/internal/commands"
	"github.com/user/gitgud/internal/help"
)

func main() {
	if len(os.Args) < 2 {
		help.ShowShortUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	commands.HandleCommand(command, args)
}
