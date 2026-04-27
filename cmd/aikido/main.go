package main

import (
	"github.com/xandervr/aikido-cli/internal/cli"
	"github.com/xandervr/aikido-cli/internal/cli/commands"
)

func main() {
	root, g := cli.NewRoot()
	root.AddCommand(
		commands.NewAuth(g),
		commands.NewWorkspace(g),
		commands.NewRepos(g),
		commands.NewIssues(g),
	)
	if err := root.Execute(); err != nil {
		cli.Exit(err)
	}
}
